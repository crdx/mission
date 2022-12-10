package config

import (
	"encoding/json"
	"fmt"
	"github.com/crdx/mission/jsonc"
	"github.com/crdx/mission/task"
	"os"
	"strings"
)

func (self PingConfig) GetEndpoint() (string, error) {
	if strings.Contains(self.Endpoint, "%s") {
		hostname, err := os.Hostname()
		if err != nil {
			return "", err
		}

		return fmt.Sprintf(self.Endpoint, hostname), nil
	} else {
		return self.Endpoint, nil
	}
}

func (self Config) GetRunnableTasks(slugs []string) task.Tasks {
	if len(slugs) > 0 {
		return self.getTasksBySlugs(slugs)
	} else {
		return self.getScheduledTasks()
	}
}

func Get(path string) (config Config, err error) {
	configJsonC, err := os.ReadFile(path)

	if err != nil {
		err = fmt.Errorf("unable to read config file: %w", err)
		return
	}

	configJson, err := jsonc.Decode(configJsonC)

	if err != nil {
		err = fmt.Errorf("unable to decode jsonc: %w", err)
		return
	}

	err = json.Unmarshal(configJson, &config)

	if err != nil {
		err = fmt.Errorf("unable to parse %s: %w", path, err)
		return
	}

	err = transform(&config)

	if err != nil {
		err = fmt.Errorf("unable to transform %s: %w", path, err)
		return
	}

	err = validate(config)

	if err != nil {
		err = fmt.Errorf("unable to validate %s: %w", path, err)
		return
	}

	return
}

// —————————————————————————————————————————————————————————————————————————————————————————————————

func (self Config) getScheduledTasks() task.Tasks {
	var tasks task.Tasks

	for _, task := range self.Tasks {
		if task.Scheduled {
			tasks = append(tasks, task)
		}
	}

	return tasks
}

func (self Config) getTasksBySlugs(slugs []string) task.Tasks {
	var tasks task.Tasks

	for _, slug := range slugs {
		for _, task := range self.Tasks {
			if task.Slug == slug {
				tasks = append(tasks, task)
			}
		}
	}

	return tasks
}
