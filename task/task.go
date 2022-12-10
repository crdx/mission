package task

import (
	"fmt"
	"github.com/crdx/mission/tasks/spotify"
	"path"

	"github.com/crdx/col"
)

type Task struct {
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	Scheduled bool   `json:"scheduled"` // Is this task run on a schedule?
	External  bool   `json:"external"`  // Is this task an external script?
	Post      bool   `json:"post"`      // Should this task run after commit & chown?
}

func (self Task) GetScriptPath() string {
	if self.External {
		return path.Join("tasks", self.Slug, "run")
	} else {
		return ""
	}
}

func (self Task) GetBuiltIn() Action {
	tasks := map[string]Action{
		"spotify": spotify.Run,
	}

	return tasks[self.Slug]
}

func (self Task) GetShortString() string {
	return self.Slug
}

func (self Task) GetLongString() string {
	return fmt.Sprintf("%s [%s, %s]", self.Slug, col.Bold(self.Name), self.GetDisplayEnabled())
}

func (self Task) GetDisplayEnabled() string {
	if self.Scheduled {
		return col.Green("auto")
	} else {
		return col.Red("manual")
	}
}
