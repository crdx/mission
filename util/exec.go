package util

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
)

func ExecCommand(cmd ...string) (string, error) {
	bytes, err := exec.Command(cmd[0], cmd[1:]...).Output()
	return strings.TrimSpace(string(bytes)), err
}

type ExecContext struct {
	workDir string
	env     []string
	writer  func(string)
}

func NewExecContext(workDir string, writer func(string), env []string) ExecContext {
	return ExecContext{
		workDir: workDir,
		env:     env,
		writer:  writer,
	}
}

func (self ExecContext) NewCommand(cmd ...string) *exec.Cmd {
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Dir = self.workDir

	if len(self.env) > 0 {
		command.Env = append(os.Environ(), self.env...)
	}

	return command
}

func (self ExecContext) GetStdout(cmd ...string) (string, error) {
	command := self.NewCommand(cmd...)
	bytes, err := command.Output()
	return string(bytes), err
}

func (self ExecContext) Run(cmd ...string) error {
	command := self.NewCommand(cmd...)

	stdout, err := command.StdoutPipe()
	if err != nil {
		return err
	}

	command.Stderr = command.Stdout

	if err = command.Start(); err != nil {
		return err
	}

	for scanner := bufio.NewScanner(stdout); scanner.Scan(); {
		self.writer(scanner.Text())
	}

	return command.Wait()
}
