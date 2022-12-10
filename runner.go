package main

import (
	"fmt"
	"github.com/crdx/mission/args"
	"github.com/crdx/mission/config"
	"github.com/crdx/mission/logger"
	"github.com/crdx/mission/notify"
	"github.com/crdx/mission/task"
	"github.com/crdx/mission/util"
	"path/filepath"
	"strings"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

type Runner struct {
	headless  bool
	quiet     bool
	config    config.Config
	logger    *logger.Logger
	startedAt time.Time
}

func NewRunner(headless bool, quiet bool, config config.Config, logger *logger.Logger, startedAt time.Time) Runner {
	return Runner{
		headless:  headless,
		quiet:     quiet,
		config:    config,
		logger:    logger,
		startedAt: startedAt,
	}
}

func (self Runner) Run(slugs []string) error {
	tasks := self.config.GetRunnableTasks(slugs)

	if len(tasks) == 0 {
		return fmt.Errorf("no tasks found")
	}

	preTasks, postTasks := self.splitTasks(tasks)

	self.logger.Header("Start")
	self.logger.Printf("StartTime %s\n", self.startedAt.Format(TimeFormat))

	if self.headless {
		self.mailStart()
	}

	if !self.quiet {
		self.notifyStart()
	}

	if len(preTasks) > 0 {
		self.runTasks(preTasks)
		self.commitRepositories()
		self.chownFiles()
	}

	if len(postTasks) > 0 {
		self.runTasks(postTasks)
	}

	self.logger.Header("Finish")
	completedIn := time.Since(self.startedAt).Truncate(time.Second)

	self.logger.Printf("FinishTime %s\n", time.Now().Format(TimeFormat))

	if !self.quiet {
		if self.logger.FoundProblems() {
			self.notifyFail(completedIn)
		} else {
			self.notifyFinish(completedIn)
		}
	}

	warnings := self.logger.PrintWarnings()
	errors := self.logger.PrintErrors()

	if self.headless {
		self.mailFinish(completedIn)
	}

	if self.headless && !self.logger.FoundProblems() {
		self.ping()
	}

	if warnings > 0 || errors > 0 {
		return fmt.Errorf(
			"%d %s and %d %s were emitted during the run",
			errors,
			util.Pluralise(errors, "error", "errors"),
			warnings,
			util.Pluralise(warnings, "warning", "warnings"),
		)
	}

	return nil
}

func (Runner) splitTasks(tasks task.Tasks) (pre task.Tasks, post task.Tasks) {
	for _, task := range tasks {
		if task.Post {
			post = append(post, task)
		} else {
			pre = append(pre, task)
		}
	}
	return
}

func (self Runner) runTasks(tasks task.Tasks) {
	for _, task := range tasks {
		self.logger.HandleError(self.runOne(task), task.Slug)
	}
}

func (self Runner) runOne(task task.Task) error {
	self.logger.Header("Task: " + task.Name)

	if task.External {
		return self.runExternal(self.getArgs(), task.GetScriptPath())
	} else {
		return task.GetBuiltIn()(self.getArgs(), self.logger)
	}
}

func (self Runner) runExternal(args args.Args, entryPoint string) error {
	absoluteEntryPoint, err := filepath.Abs(entryPoint)
	if err != nil {
		return err
	}

	ctx := util.NewExecContext(
		filepath.Dir(absoluteEntryPoint),
		func(str string) { self.logger.Println(str) },
		args.ToEnvironmentVariables(),
	)

	return ctx.Run(absoluteEntryPoint)
}

func (self Runner) getArgs() args.Args {
	return args.New(
		self.config.Storage["sync"].Path,
		self.config.Storage["local"].Path,
		self.config.Storage["helpers"].Path,
		self.config.Storage["logs"].Path,
		self.config.Storage["store"].Path,
		self.config.User.Name,
		self.config.PassBin,
	)
}

func (self Runner) mailStart() {
	self.logger.Println("MailStart")

	subject := "run started"
	body := "Run started."

	self.logger.HandleError(util.SendMail(self.config.User.Name, subject, body), "mailStart")
}

func (self Runner) mailFinish(completedIn time.Duration) {
	self.logger.Header("Mail")
	self.logger.Println("MailFinish")

	subject := fmt.Sprintf("run finished in %s", completedIn)

	if self.logger.FoundProblems() {
		subject = "[FAILURE] " + subject
	}

	body := fmt.Sprintf(
		"Run complete in %s\n\n%s",
		completedIn,
		strings.Join(self.logger.FilteredLines(), ""),
	)

	self.logger.HandleError(util.SendMail(self.config.User.Name, subject, body), "mailFinish")
}

func (self Runner) chownFiles() {
	self.logger.Header("Chown")

	userInfo, err := util.GetUserInfo(self.config.User.Name)
	if err != nil {
		self.logger.HandleWarning(err, "chownFiles")
		return
	}

	for _, dir := range self.config.Storage {
		if !dir.Chown {
			continue
		}

		self.logger.Printf("%s: ", dir.Path)
		count, err := util.ChownDirectory(dir.Path, userInfo.UserId, userInfo.GroupId)
		self.logger.PrintRawf("%d %s\n", count, util.Pluralise(count, "file", "files"))
		self.logger.HandleWarning(err, "chownFiles")
	}
}

func (self Runner) getCommitMessage() string {
	return fmt.Sprintf(
		"%s run complete on %s",
		util.GetRunType(self.headless),
		self.startedAt.Format(TimeFormat),
	)
}

func (self Runner) commitRepositories() {
	self.logger.Header("Commit")

	writer := func(str string) {
		self.logger.PrintRawln(str)
	}

	leadingNewLine := false

	for _, dir := range self.config.Storage {
		if !dir.Commit {
			continue
		}

		if leadingNewLine {
			self.logger.PrintRawf("\n")
		}
		leadingNewLine = true

		self.logger.Printf("RepoDir %s\n", dir.Path)

		ctx := util.NewExecContext(dir.Path, writer, nil)

		self.logger.Println("Status")
		if err := ctx.Run("git", "status", "-sb"); err != nil {
			self.logger.HandleError(err, "commitRepositories")
			continue
		}

		stdout, err := ctx.GetStdout("git", "status", "--short")
		if err != nil {
			self.logger.HandleError(err, "commitRepositories")
			continue
		}

		if len(stdout) == 0 {
			self.logger.Println("NoChanges")
			continue
		}

		self.logger.Println("Add")
		if err := ctx.Run("git", "add", "."); err != nil {
			self.logger.HandleError(err, "commitRepositories")
			continue
		}

		self.logger.Println("Commit")
		if err := ctx.Run("git", "commit", "--message", self.getCommitMessage()); err != nil {
			self.logger.HandleError(err, "commitRepositories")
			continue
		}
	}
}

func (self Runner) ping() {
	if !self.config.Ping.Enabled {
		return
	}

	self.logger.Header("Ping")

	endpoint, err := self.config.Ping.GetEndpoint()
	if err != nil {
		self.logger.HandleWarning(err, "ping")
		return
	}

	body, err := util.HttpGet(endpoint, nil)
	if err != nil {
		self.logger.HandleWarning(err, "ping")
		return
	}

	self.logger.Printf("PingResponse %s\n", body)
}

func (self Runner) notifyStart() {
	if !self.config.Notify.Enabled {
		return
	}

	self.logger.Println("NotifyStart")
	self.logger.HandleWarning(notify.Start(self.config.User.Name), "notifyStart")
}

func (self Runner) notifyFinish(completedIn time.Duration) {
	if !self.config.Notify.Enabled {
		return
	}

	self.logger.Println("NotifyFinish")
	self.logger.HandleWarning(notify.Finish(self.config.User.Name, completedIn), "notifyFinish")
}

func (self Runner) notifyFail(completedIn time.Duration) {
	if !self.config.Notify.Enabled {
		return
	}

	self.logger.Println("NotifyFailure")
	self.logger.HandleWarning(notify.Fail(self.config.User.Name, completedIn), "notifyFail")
}
