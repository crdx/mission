package config

import (
	"fmt"

	"crdx.org/col"
	"crdx.org/mission/internal/args"
	"crdx.org/mission/internal/logger"
	"crdx.org/mission/internal/tasks/spotify"
)

type Action func(args *args.Args, logger *logger.Logger) error

const (
	TaskTypeExec    = "exec"
	TaskTypeBuiltIn = "builtin"
)

var ValidTaskTypes = []string{
	TaskTypeExec,
	TaskTypeBuiltIn,
}

type Task struct {
	Slug       string `json:"slug"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Scheduled  bool   `json:"scheduled"`
	Post       bool   `json:"post"`
	EntryPoint string `json:"entrypoint"`
}

func (self Task) GetBuiltInAction() Action {
	tasks := map[string]Action{
		"spotify": spotify.Run,
	}

	return tasks[self.Slug]
}

func (self Task) GetShortString() string {
	return self.Slug
}

func (self Task) GetLongString() string {
	return fmt.Sprintf(
		"%s [%s, %s]",
		self.Slug,
		col.Bold(self.Name),
		self.GetDisplayEnabled(),
	)
}

func (self Task) GetDisplayEnabled() string {
	if self.Scheduled {
		return col.Green("scheduled")
	} else {
		return col.Red("manual")
	}
}
