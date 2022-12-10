package spotify

import (
	"fmt"
	"github.com/crdx/mission/args"
)

type Credentials struct {
	clientId     string
	clientSecret string
}

func getCredentials(args args.Args) (Credentials, error) {
	clientSecret, err := args.GetPassValue("spotify_api_client_secret")
	if err != nil {
		return Credentials{}, err
	}

	clientId, err := args.GetPassValue("spotify_api_client_id")
	if err != nil {
		return Credentials{}, err
	}

	if len(clientId) == 0 || len(clientSecret) == 0 {
		err = fmt.Errorf("unable to find clientId or clientSecret")
		return Credentials{}, err
	}

	return Credentials{
		clientId:     clientId,
		clientSecret: clientSecret,
	}, err
}
