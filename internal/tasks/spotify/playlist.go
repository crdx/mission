package spotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"crdx.org/mission/internal/logger"
	"github.com/zmb3/spotify"
)

type Playlist struct {
	Id             spotify.ID              `json:"-"`
	Name           string                  `json:"-"`
	Instance       spotify.SimplePlaylist  `json:"playlist"`
	PlaylistTracks []spotify.PlaylistTrack `json:"tracks"`

	client *spotify.Client
}

func newPlaylist(client *spotify.Client, playlist spotify.SimplePlaylist) Playlist {
	return Playlist{
		Id:       playlist.ID,
		Name:     playlist.Name,
		Instance: playlist,

		client: client,
	}
}

func (self Playlist) FileName() string {
	return fmt.Sprintf("%s_%s.json", self.Id, self.SanitisedName())
}

func (self Playlist) SanitisedName() string {
	if self.Name != "" {
		return sanitise(self.Name)
	} else {
		return "unknown-playlist-name"
	}
}

func (self *Playlist) loadTracks(playlistTrackPage *spotify.PlaylistTrackPage) error {
	self.PlaylistTracks = []spotify.PlaylistTrack{}

	for {
		self.PlaylistTracks = append(self.PlaylistTracks, playlistTrackPage.Tracks...)

		err := self.client.NextPage(playlistTrackPage)
		if errors.Is(err, spotify.ErrNoMorePages) {
			return nil
		}

		if err != nil {
			return err
		}
	}
}

func (self *Playlist) SaveTo(dir string) (int, error) {
	playlistTrackPage, err := self.client.GetPlaylistTracks(self.Id)
	if err != nil {
		return 0, err
	}

	if err := self.loadTracks(playlistTrackPage); err != nil {
		return 0, err
	}

	bytes, err := json.MarshalIndent(self, "", "    ")
	if err != nil {
		return 0, err
	}

	return len(bytes), os.WriteFile(path.Join(dir, self.FileName()), bytes, 0666)
}

func getPlaylists(client *spotify.Client, playlistPage *spotify.SimplePlaylistPage) ([]Playlist, error) {
	var playlists []Playlist

	for {
		for _, playlist := range playlistPage.Playlists {
			playlists = append(playlists, newPlaylist(client, playlist))
		}

		err := client.NextPage(playlistPage)
		if errors.Is(err, spotify.ErrNoMorePages) {
			return playlists, nil
		}
		if err != nil {
			return nil, err
		}
	}
}

func readSavedPlaylists(dir string) ([]string, error) {
	var fileNames []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileNames = append(fileNames, entry.Name())
	}

	return fileNames, nil
}

func savePlaylists(dir string, client *spotify.Client, logger *logger.Logger) error {
	logger.Println("Fetching playlists")
	playlistPage, err := client.CurrentUsersPlaylists()
	if err != nil {
		return err
	}

	playlists, err := getPlaylists(client, playlistPage)
	if err != nil {
		return err
	}

	logger.Printf("Found %d playlists\n", len(playlists))

	seen := map[string]bool{}

	for _, playlist := range playlists {
		logger.Printf("Downloading %s... ", playlist.FileName())

		byteCount, err := playlist.SaveTo(dir)
		if err != nil {
			logger.Printf("\n")
			return err
		}

		logger.PrintRawf("%dK\n", byteCount/1000)
		seen[playlist.FileName()] = true
	}

	savedPlaylists, err := readSavedPlaylists(dir)
	if err != nil {
		return err
	}

	for _, fileName := range savedPlaylists {
		if _, ok := seen[fileName]; !ok {
			logger.Println("Deleting", fileName)
			if err := os.Remove(path.Join(dir, fileName)); err != nil {
				return err
			}
		}
	}

	return nil
}
