package config

import (
	"path"

	"github.com/crdx/mission/internal/util"
)

func transform(config *Config) error {
	for key, storage := range config.Storage {
		path := storage.Path

		if path == "" {
			continue
		}

		path, err := util.GetAbsoluteDir(path, config.User.Name)
		if err != nil {
			return err
		}

		storage.Path = path
		config.Storage[key] = storage
	}

	for key, task := range config.Tasks {
		if task.EntryPoint != "" {
			entryPoint, err := util.GetAbsoluteDir(task.EntryPoint, config.User.Name)
			if err != nil {
				return err
			}

			task.EntryPoint = entryPoint
		} else {
			task.EntryPoint = path.Join(config.Storage["tasks"].Path, task.Slug, "run")
		}

		config.Tasks[key] = task
	}

	if config.PassBin != "" {
		passBin, err := util.GetAbsoluteDir(config.PassBin, config.User.Name)
		if err != nil {
			return err
		}
		config.PassBin = passBin
	}

	return nil
}
