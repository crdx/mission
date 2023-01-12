package config

import (
	"fmt"
	"net/mail"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/crdx/mission/util"
	"golang.org/x/exp/slices"
)

func validate(config Config) error {
	if len(config.Tasks) == 0 {
		return fmt.Errorf("missing tasks")
	}

	for _, task := range config.Tasks {
		if task.Type == "" {
			return fmt.Errorf("type for tasks.%s is missing", task.Slug)
		}

		if !slices.Contains(ValidTaskTypes, task.Type) {
			return fmt.Errorf("type for tasks.%s (%s) is invalid", task.Slug, task.Type)
		}

		if task.Type == TaskTypeExec {
			if !util.IsExecutable(task.EntryPoint) {
				return fmt.Errorf(
					"entrypoint for tasks.%s (%s) does not exist or is not executable",
					task.Slug,
					task.EntryPoint,
				)
			}

			bytes, err := os.ReadFile(task.EntryPoint)
			if err != nil {
				return fmt.Errorf("unable to read entrypoint for tasks.%s", task.Slug)
			}

			str := string(bytes)

			if strings.HasPrefix(str, "#!/bin/bash\n") {
				if !strings.HasPrefix(str, "#!/bin/bash\nset -euo pipefail\n") {
					return fmt.Errorf("bash scripts must start with set -euo pipefail")
				}
			}
		}

		if task.Type == TaskTypeBuiltIn {
			if task.GetBuiltInAction() == nil {
				return fmt.Errorf("built in task %s has no implementation", task.Slug)
			}
		}
	}

	if config.User.Name == "" {
		return fmt.Errorf("missing username")
	}

	userInfo, err := util.GetUserInfo(config.User.Name)
	if err != nil {
		return fmt.Errorf("unable to find user: %w", err)
	}
	if userInfo.UserId == 0 {
		return fmt.Errorf("unable to find valid user (uid is %d)", userInfo.UserId)
	}

	if config.Mail.Enabled {
		if config.User.Email == "" {
			return fmt.Errorf("missing email address")
		}

		if _, err := mail.ParseAddress(config.User.Email); err != nil {
			return fmt.Errorf("email address is invalid: %w", err)
		}
	}

	if config.PassBin != "" && !util.IsExecutable(config.PassBin) {
		return fmt.Errorf("PassBin (%s) is not an executable file", config.PassBin)
	}

	for _, key := range []string{"tasks", "logs"} {
		dir := config.Storage[key]

		if dir.Path == "" {
			return fmt.Errorf("missing storage.%s.path", key)
		}
	}

	for key, dir := range config.Storage {
		if !util.IsDirectory(dir.Path) {
			return fmt.Errorf("dirs.%s.path (%s) is not a directory", key, dir.Path)
		}

		if dir.Commit && !util.IsGitRepository(dir.Path) {
			return fmt.Errorf(
				"dirs.%s.commit is true but dirs.%s.path (%s) is not a git repository",
				key,
				key, dir.Path,
			)
		}
	}

	if config.Ping.Enabled {
		if config.Ping.Endpoint == "" {
			return fmt.Errorf("missing ping.endpoint")
		}
	}

	if config.Notify.Enabled {
		if _, err := exec.LookPath("notify-send"); err != nil {
			return fmt.Errorf("missing notify-send dependency required by notify")
		}
	}

	if config.Mail.Enabled {
		if config.Mail.Type != "sendmail" {
			return fmt.Errorf("sendmail is the only valid value for mail.type")
		}

		if _, err := exec.LookPath("sendmail"); err != nil {
			return fmt.Errorf("missing sendmail dependency required by mail")
		}
	}

	for _, str := range config.Filters {
		_, err = regexp.Compile(str)
		if err != nil {
			return fmt.Errorf("filter regex is invalid: %w", err)
		}
	}

	return nil
}
