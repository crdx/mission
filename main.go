package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/crdx/mission/config"
	"github.com/crdx/mission/logger"
	"github.com/crdx/mission/util"

	"github.com/crdx/col"
	"github.com/crdx/duckopt"
)

func getUsage() string {
	return `
        Usage:
            $0 [options] run [--task SLUG...] [--headless]
            $0 [options] list [--verbose]
            $0 [options] check
            $0 [options] dump

        Commands:
            run     Run all tasks or specific tasks
            list    List all available tasks
            check   Validate configuration
            dump    Dump parsed configuration as JSON

        Options:
            --headless           Run headlessly
            -c, --config PATH    Configuration file
            -t, --task SLUG      One or more tasks to run
            -q, --quiet          Be quiet
            -v, --verbose        Be verbose
            -C, --no-color       Disable colours
            -h, --help           Show help
	`
}

type Opts struct {
	Run   bool `docopt:"run"`
	List  bool `docopt:"list"`
	Check bool `docopt:"check"`
	Dump  bool `docopt:"dump"`

	Config   string   `docopt:"--config"`
	Verbose  bool     `docopt:"--verbose"`
	Headless bool     `docopt:"--headless"`
	Quiet    bool     `docopt:"--quiet"`
	Tasks    []string `docopt:"--task"`
	NoColor  bool     `docopt:"--no-color"`
	Help     bool     `docopt:"--help"`
}

const (
	DefaultConfigFilePath = "mission.config.json"
	LockFilePath          = "/tmp/mission.lock"
)

func die(format string, args ...any) {
	log.Fatalf(col.Red("Error: "+format), args...)
}

func getLogFileName(startedAt time.Time) string {
	return fmt.Sprintf(
		"%s.txt",
		startedAt.Format("2006-01-02-15-04-05"),
	)
}

func createLogger(headless bool, config config.Config, startedAt time.Time) (*logger.Logger, error) {
	logger := logger.New(config.Filters)

	if headless {
		logFile := path.Join(config.Storage["logs"].Path, getLogFileName(startedAt))
		file, err := os.Create(logFile)

		if err != nil {
			return logger, fmt.Errorf("unable to create logfile %s: %w", logFile, err)
		}

		logger.AddFile(file)
	} else {
		logger.AddFile(os.Stdout)
	}

	return logger, nil
}

func lock() bool {
	if util.IsReadableFile(LockFilePath) {
		return false
	}
	_, err := os.Create(LockFilePath)
	return err == nil
}

func unlock() {
	os.Remove(LockFilePath)
}

func getConfigFilePath(path string) string {
	if path != "" {
		return path
	} else {
		return DefaultConfigFilePath
	}
}

func main() {
	var opts Opts
	log.SetFlags(0)
	if err := duckopt.Parse(getUsage(), "$0").Bind(&opts); err != nil {
		die("unable to parse arguments: %s", err)
	}

	col.InitUnless(opts.NoColor || opts.Headless)

	startedAt := time.Now()

	config, err := config.Get(getConfigFilePath(opts.Config))
	if err != nil {
		die("%s", err)
	}

	if opts.Dump {
		j, err := json.MarshalIndent(config, "", "    ")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Println(string(j))
			os.Exit(0)
		}
	}

	if opts.Check {
		os.Exit(0)
	}

	if opts.List {
		for _, task := range config.Tasks {
			if opts.Verbose {
				fmt.Println(task.GetLongString())
			} else {
				fmt.Println(task.GetShortString())
			}
		}
		os.Exit(0)
	}

	if !lock() {
		die("unable to obtain exclusive lock")
	}

	logger, err := createLogger(opts.Headless, config, startedAt)
	if err != nil {
		unlock()
		die("%s", err)
	}

	err = NewRunner(opts.Headless, opts.Quiet, config, logger, startedAt).Run(opts.Tasks)

	logger.Close()
	unlock()

	if err != nil {
		die("%s", err)
	}
}
