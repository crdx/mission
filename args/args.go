package args

import (
	"fmt"
	"github.com/crdx/mission/util"
)

type Args struct {
	SyncFilesDir  string
	LocalFilesDir string
	HelpersDir    string
	LogsDir       string
	StoreDir      string
	TargetUser    string
	PassBin       string
}

func New(syncFilesDir, localFilesDir, helpersDir, logsDir, storeDir, targetUser, passBin string) Args {
	return Args{
		SyncFilesDir:  syncFilesDir,
		LocalFilesDir: localFilesDir,
		HelpersDir:    helpersDir,
		LogsDir:       logsDir,
		StoreDir:      storeDir,
		TargetUser:    targetUser,
		PassBin:       passBin,
	}
}

func (self Args) GetPassValue(key string) (string, error) {
	return util.ExecCommand(self.PassBin, key)
}

func (self Args) ToEnvironmentVariables() []string {
	return []string{
		fmt.Sprintf("SYNC_FILES_DIR=%s", self.SyncFilesDir),
		fmt.Sprintf("LOCAL_FILES_DIR=%s", self.LocalFilesDir),
		fmt.Sprintf("HELPERS_DIR=%s", self.HelpersDir),
		fmt.Sprintf("LOGS_DIR=%s", self.LogsDir),
		fmt.Sprintf("STORE_DIR=%s", self.StoreDir),
		fmt.Sprintf("TARGET_USER=%s", self.TargetUser),
		fmt.Sprintf("PASS_BIN=%s", self.PassBin),
	}
}
