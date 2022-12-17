package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/crdx/mission/util"
)

func validate(config Config) error {
	if len(config.Tasks) == 0 {
		return fmt.Errorf("missing tasks")
	}

	for _, task := range config.Tasks {
		if task.External {
			if !util.IsExecutable(task.EntryPoint) {
				return fmt.Errorf("entrypoint for tasks.%s (%s) does not exist or is not executable", task.Slug, task.EntryPoint)
			}

			bytes, err := os.ReadFile(task.EntryPoint)
			if err != nil {
				return fmt.Errorf("unable to read entrypoint for tasks.%s", task.Slug)
			}

			str := string(bytes)

			if strings.HasPrefix(str, "#!/bin/bash\n") {
				if !strings.HasPrefix(str, "#!/bin/bash\nset -euo pipefail\n") {
					return fmt.Errorf("external bash scripts must start with set -euo pipefail")
				}
			}
		}

		if !task.External {
			if task.GetBuiltIn() == nil {
				return fmt.Errorf("built in task %s has no implementor", task.Slug)
			}
		}
	}

	if config.User.Name == "" {
		return fmt.Errorf("missing user")
	}

	userInfo, err := util.GetUserInfo(config.User.Name)
	if err != nil {
		return fmt.Errorf("unable to find user: %w", err)
	}
	if userInfo.UserId == 0 {
		return fmt.Errorf("unable to find valid user (uid is %d)", userInfo.UserId)
	}

	if config.PassBin == "" {
		return fmt.Errorf("missing PassBin")
	}

	if !util.IsExecutable(config.PassBin) {
		return fmt.Errorf("PassBin (%s) is not an executable file", config.PassBin)
	}

	for _, key := range []string{"sync", "local", "logs", "helpers", "tasks"} {
		dir := config.Storage[key]

		if dir.Path == "" {
			return fmt.Errorf("missing storage.%s.path", key)
		}

		if !util.IsDirectory(dir.Path) {
			return fmt.Errorf("dirs.%s.path (%s) is not a directory", key, dir.Path)
		}

		if dir.Commit && !util.IsGitRepository(dir.Path) {
			return fmt.Errorf("dirs.%s.commit is true but dirs.%s.path (%s) is not a git repository", key, key, dir.Path)
		}
	}

	if config.Ping.Enabled {
		if config.Ping.Endpoint == "" {
			return fmt.Errorf("missing ping.endpoint")
		}
	}

	if config.Notify.Enabled {
		if !util.IsExecutable("/bin/notify-send") {
			return fmt.Errorf("missing notify-send dependency required by notify")
		}
	}

	if config.Mail.Enabled {
		if config.Mail.Type != "sendmail" {
			return fmt.Errorf("sendmail is the only valid value for mail.type")
		}
	}

	for _, str := range config.Filters {
		_, err = regexp.Compile(str)
		if err != nil {
			return fmt.Errorf("filter is invalid: %w", err)
		}
	}

	return nil
}
