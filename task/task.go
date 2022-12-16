package task

import (
	"fmt"

	"github.com/crdx/mission/tasks/spotify"

	"github.com/crdx/col"
)

type Task struct {
	Slug       string `json:"slug"`
	Name       string `json:"name"`
	Scheduled  bool   `json:"scheduled"`
	External   bool   `json:"external"`
	Post       bool   `json:"post"`
	EntryPoint string `json:"entrypoint"`
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
	return fmt.Sprintf(
		"%s [%s, %s]",
		self.Slug,
		col.Bold(self.Name),
		self.GetDisplayEnabled(),
	)
}

func (self Task) GetDisplayEnabled() string {
	if self.Scheduled {
		return col.Green("auto")
	} else {
		return col.Red("manual")
	}
}
