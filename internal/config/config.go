package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/crdx/mission/internal/jsonc"
)

type Config struct {
	Tasks   []Task             `json:"tasks"`
	User    User               `json:"user"`
	PassBin string             `json:"passBin"`
	Storage map[string]Storage `json:"storage"`
	Ping    Ping               `json:"ping"`
	Notify  Notify             `json:"notify"`
	Mail    Mail               `json:"mail"`
	Filters []string           `json:"filters"`
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Storage struct {
	Path   string `json:"path"`
	Commit bool   `json:"commit"`
	Chown  bool   `json:"chown"`
}

type Ping struct {
	Enabled  bool   `json:"enabled"`
	Endpoint string `json:"endpoint"`
}

type Notify struct {
	Enabled bool `json:"enabled"`
}

type Mail struct {
	Enabled bool   `json:"enabled"`
	Type    string `json:"type"`
}

func (self Ping) GetEndpoint() (string, error) {
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

func (self Config) GetRunnableTasks(slugs []string) []Task {
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
		err = fmt.Errorf("validation failure for %s: %w", path, err)
		return
	}

	return
}

// —————————————————————————————————————————————————————————————————————————————————————————————————

func (self Config) getScheduledTasks() []Task {
	var tasks []Task

	for _, task := range self.Tasks {
		if task.Scheduled {
			tasks = append(tasks, task)
		}
	}

	return tasks
}

func (self Config) getTasksBySlugs(slugs []string) []Task {
	var tasks []Task

	for _, slug := range slugs {
		for _, task := range self.Tasks {
			if task.Slug == slug {
				tasks = append(tasks, task)
			}
		}
	}

	return tasks
}
