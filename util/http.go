package util

import (
	"io"
	"net/http"
)

func HttpGet(endpoint string, params map[string]string) ([]byte, error) {
	client := http.Client{}
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)

	if err != nil {
		return nil, err
	}

	query := request.URL.Query()
	for key, value := range params {
		query.Add(key, value)
	}
	request.URL.RawQuery = query.Encode()

	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	return io.ReadAll(response.Body)
}
