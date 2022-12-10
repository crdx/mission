package spotify

import "regexp"

var sanitiseRegex = regexp.MustCompile("[^A-Za-z0-9-_]")

func sanitise(name string) string {
	return string(sanitiseRegex.ReplaceAll([]byte(name), []byte("")))
}
