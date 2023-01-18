package logger

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"crdx.org/mission/internal/util"
)

type Problem struct {
	Err    error
	Source string
}

type Logger struct {
	Warnings []Problem
	Errors   []Problem
	Lines    []string

	filters    []string
	files      []*os.File
	startedAt  time.Time
	hasWritten bool
}

func New(filters []string) *Logger {
	return &Logger{
		startedAt: time.Now(),
		filters:   filters,
	}
}

func (self *Logger) AddFile(file *os.File) {
	self.files = append(self.files, file)
}

func (self *Logger) Header(str string) {
	if self.hasWritten {
		self.write("\n", false)
	}

	self.write("# ———————————————————————————————————————————————\n", false)
	self.write(fmt.Sprintf("# %s\n", str), false)
	self.write("# ———————————————————————————————————————————————\n", false)
	self.write("\n", false)
}

func (self *Logger) FoundProblems() bool {
	return len(self.Errors) > 0 || len(self.Warnings) > 0
}

func (self *Logger) HandleError(err error, source string) {
	if err != nil {
		self.write(fmt.Sprintf("ERROR: %s\n", err), true)
		self.Errors = append(self.Errors, Problem{err, source})
	}
}

func (self *Logger) HandleWarning(err error, source string) {
	if err != nil {
		self.write(fmt.Sprintf("WARN: %s\n", err), true)
		self.Warnings = append(self.Warnings, Problem{err, source})
	}
}

func (self *Logger) PrintWarnings() int {
	return self.printObjects("Warnings", self.Warnings)
}

func (self *Logger) PrintErrors() int {
	return self.printObjects("Errors", self.Errors)
}

func (self *Logger) Printf(format string, args ...any) {
	self.write(fmt.Sprintf(format, args...), true)
}

func (self *Logger) PrintRawf(format string, args ...any) {
	self.write(fmt.Sprintf(format, args...), false)
}

func (self *Logger) Println(str ...string) {
	self.write(strings.Join(str, " ")+"\n", true)
}

func (self *Logger) PrintRawln(str ...string) {
	self.write(strings.Join(str, " ")+"\n", false)
}

func (self *Logger) Close() {
	for _, file := range self.files {
		if file != os.Stdout {
			file.Close()
		}
	}
}

func (self *Logger) FilteredLines() []string {
	skipped := 0
	var filteredLines []string

	for _, line := range self.Lines {
		if self.matchesFilter(line) {
			skipped++
			continue
		} else if skipped > 0 {
			s := fmt.Sprintf(
				"\n<<%d %s filtered>>\n\n",
				skipped,
				util.Pluralise(skipped, "line", "lines"),
			)
			filteredLines = append(filteredLines, s)
			skipped = 0
		}

		filteredLines = append(filteredLines, line)
	}

	return filteredLines
}

// —————————————————————————————————————————————————————————————————————————————————————————————————

func (self *Logger) matchesFilter(line string) bool {
	bytes := []byte(line)

	for _, filter := range self.filters {
		// Guaranteed to already be valid by this point.
		if regexp.MustCompile(filter).Match(bytes) {
			return true
		}
	}

	return false
}

func (self *Logger) printObjects(what string, problems []Problem) int {
	if len(problems) > 0 {
		self.Header(what)
		for _, problem := range problems {
			self.write(fmt.Sprintf("- %s: %s\n", problem.Source, problem.Err), false)
		}
		return len(problems) // 99
	}

	return 0
}

func (self *Logger) write(str string, timestamped bool) {
	if timestamped {
		str = fmt.Sprintf("[%-10s] %s", util.FormatDuration(time.Since(self.startedAt)), str)
	}

	self.hasWritten = true
	self.Lines = append(self.Lines, str)

	for _, file := range self.files {
		_, _ = file.Write([]byte(str))
		_ = file.Sync()
	}
}
