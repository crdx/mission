package config

import "github.com/crdx/mission/task"

type Config struct {
	Tasks   task.Tasks               `json:"tasks"`
	User    UserConfig               `json:"user"`
	PassBin string                   `json:"passBin"`
	Storage map[string]StorageConfig `json:"storage"`
	Ping    PingConfig               `json:"ping"`
	Notify  NotifyConfig             `json:"notify"`
	Mail    MailConfig               `json:"mail"`
	Filters []string                 `json:"filters"`
}

type UserConfig struct {
	Name string `json:"name"`
}

type StorageConfig struct {
	Path   string `json:"path"`
	Commit bool   `json:"commit"`
	Chown  bool   `json:"chown"`
}

type PingConfig struct {
	Enabled  bool   `json:"enabled"`
	Endpoint string `json:"endpoint"`
}

type NotifyConfig struct {
	Enabled bool `json:"enabled"`
}

type MailConfig struct {
	Enabled bool   `json:"enabled"`
	Type    string `json:"type"`
}
