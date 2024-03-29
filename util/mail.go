package util

import (
	"fmt"
	"os"
	"os/exec"

	"crdx.org/hereduck"
)

func SendMail(name, email, subject, body string) error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	template := hereduck.D(`
		To: %s <%s>
		Subject: %s: %s
		Content-Type: text/plain

		%s
	`)

	payload := fmt.Sprintf(
		template,
		name,
		email,
		hostname,
		subject,
		body,
	)

	command := exec.Command("sendmail", "-i", email)

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
