package spotify

import (
	"fmt"
	"github.com/crdx/mission/args"
	"github.com/crdx/mission/logger"
	"github.com/crdx/mission/util"
	"path/filepath"

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

func Run(args args.Args, logger *logger.Logger) error {
	credentials, err := getCredentials(args)
	if err != nil {
		return err
	}

	saveDir, err := getSaveDir(args.SyncFilesDir)
	if err != nil {
		return err
	}

	token, err := loadToken(args.StoreDir)
	if err != nil {
		return err
	}

	client := newAuthenticator(credentials.clientId, credentials.clientSecret).NewClient(&token)

	defer logger.HandleError(saveToken(&client, args.StoreDir), "spotify")
	return savePlaylists(saveDir, &client, logger)
}
