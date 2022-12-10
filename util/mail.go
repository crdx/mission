package util

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/crdx/hereduck"
)

func SendMail(toName, subject, body string) error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	toAddress := "root@" + hostname

	template := hereduck.D(`
		To: %s <%s>
		Subject: %s: %s
		Content-Type: text/plain

		%s
	`)

	payload := fmt.Sprintf(
		template,
		toName,
		toAddress,
		hostname,
		subject,
		body,
	)

	command := exec.Command("sendmail", "-i", "root")

	stdin, err := command.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		if _, err = stdin.Write([]byte(payload)); err != nil {
			panic(err)
		}

		if err = stdin.Close(); err != nil {
			panic(err)
		}
	}()

	return command.Run()
}
