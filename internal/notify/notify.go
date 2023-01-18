package notify

import (
	"fmt"
	"os/exec"
	"time"

	"crdx.org/mission/internal/util"
)

type Icon = string

const (
	IconInfo  Icon = "dialog-information"
	IconError Icon = "dialog-error"
)

func Start(userName string) error {
	return notify(IconInfo, "Run started", userName)
}

func Finish(userName string, completedIn time.Duration) error {
	message := fmt.Sprintf("Run finished in %s", util.FormatDuration(completedIn))
	return notify(IconInfo, message, userName)
}

func Fail(userName string, completedIn time.Duration) error {
	message := fmt.Sprintf("Run finished with errors or warnings in %s", util.FormatDuration(completedIn))
	return notify(IconError, message, userName)
}

// —————————————————————————————————————————————————————————————————————————————————————————————————

func notify(icon Icon, message string, userName string) error {
	userInfo, err := util.GetUserInfo(userName)
	if err != nil {
		return err
	}

	dbusAddress := fmt.Sprintf("unix:path=/run/user/%d/bus", userInfo.UserId)

	cmd := []string{
		"sudo",
		"-u",
		userName,
		"DISPLAY=:0",
		fmt.Sprintf("DBUS_SESSION_BUS_ADDRESS=%s", dbusAddress),
		"notify-send",
		message,
		"--icon",
		icon,
		"--urgency",
		"critical",
	}

	return exec.Command(cmd[0], cmd[1:]...).Run()
}
