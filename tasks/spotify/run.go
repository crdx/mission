package spotify

import (
	"fmt"
	"path/filepath"

	"crdx.org/mission/args"
	"crdx.org/mission/logger"
	"crdx.org/mission/util"

	"github.com/zmb3/spotify"
)

func newAuthenticator(clientId, clientSecret string) spotify.Authenticator {
	authenticator := spotify.NewAuthenticator(
		"http://localhost:5432",
		spotify.ScopeUserLibraryRead,
		spotify.ScopePlaylistReadPrivate,
		spotify.ScopePlaylistReadCollaborative,
	)
	authenticator.SetAuthInfo(clientId, clientSecret)
	return authenticator
}

func getSaveDir(root string) (saveDir string, err error) {
	saveDir = filepath.Join(root, "spotify")

	if !util.IsDirectory(saveDir) {
		err = fmt.Errorf("the Spotify playlists directory does not exist: %s", saveDir)
	}

	return
}

func validate(args *args.Args) error {
	for _, name := range []string{"sync", "sessions"} {
		if _, ok := args.Storage[name]; !ok {
			return fmt.Errorf("%s directory missing", name)
		}
	}

	return nil
}

func Run(args *args.Args, logger *logger.Logger) error {
	err := validate(args)
	if err != nil {
		return err
	}

	credentials, err := getCredentials(args)
	if err != nil {
		return err
	}

	saveDir, err := getSaveDir(args.Storage["sync"])
	if err != nil {
		return err
	}

	token, err := loadToken(args.Storage["sessions"])
	if err != nil {
		return err
	}

	client := newAuthenticator(credentials.clientId, credentials.clientSecret).NewClient(&token)

	defer logger.HandleError(saveToken(&client, args.Storage["sessions"]), "spotify")
	return savePlaylists(saveDir, &client, logger)
}
