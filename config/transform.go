package config

import (
	"github.com/crdx/mission/util"
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

	if config.PassBin != "" {
		passBin, err := util.GetAbsoluteDir(config.PassBin, config.User.Name)
		if err != nil {
			return err
		}
		config.PassBin = passBin
	}

	return nil
}
