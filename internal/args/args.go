package args

import (
	"fmt"
	"strings"

	"github.com/crdx/mission/internal/util"
)

type Args struct {
	Storage    map[string]string
	TargetUser string
	PassBin    string
}

func New(storage map[string]string, targetUser string, passBin string) *Args {
	return &Args{
		TargetUser: targetUser,
		PassBin:    passBin,
		Storage:    storage,
	}
}

func (self Args) GetPassValue(key string) (string, error) {
	return util.ExecCommand(self.PassBin, key)
}

func (self Args) ToEnvironmentVariables() []string {
	env := []string{
		fmt.Sprintf("TARGET_USER=%s", self.TargetUser),
		fmt.Sprintf("PASS_BIN=%s", self.PassBin),
	}

	for name, path := range self.Storage {
		env = append(env, fmt.Sprintf("%s_DIR=%s", strings.ToUpper(name), path))
	}

	return env
}
