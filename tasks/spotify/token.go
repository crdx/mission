package spotify

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

func getTokenFile(storeDir string) string {
	return filepath.Join(storeDir, "spotify")
}

func loadToken(storeDir string) (oauth2.Token, error) {
	contents, err := os.ReadFile(getTokenFile(storeDir))
	if err != nil {
		return oauth2.Token{}, err
	}

	var token oauth2.Token
	err = json.Unmarshal(contents, &token)
	return token, err
}

func saveToken(client *spotify.Client, storeDir string) error {
	token, err := client.Token()
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(token, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(getTokenFile(storeDir), bytes, 0666)
}
