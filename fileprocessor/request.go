package fileprocessor

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"git.rickiekarp.net/rickie/fileguardian/config"
	"git.rickiekarp.net/rickie/yubase"
)

func sendRequest(fileName string, fileType FileType, context string) (*yubase.FileGuardianEntry, error) {
	url := config.ApiProtocol + "://" + config.ApiHost + "/storage/v1/fileguardian/fetch"

	if *config.FlagCheck {
		url += "?check=true"
	}

	// create post body using an instance of the Person struct
	requestEvent := yubase.FileGuardianEntry{
		Type:    string(fileType),
		Context: context,
	}

	// from a fileguardian perspective our file is always the source
	requestEvent.Source = fileName

	// convert p to JSON data
	jsonData, err := json.Marshal(requestEvent)
	if err != nil {
		return nil, errors.New("could not marshal requestEvent")
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			TLSHandshakeTimeout: 30 * time.Second,
		},
	}

	resp, err := client.Post(url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, nil
	}

	var res yubase.FileGuardianEntry
	err = json.Unmarshal([]byte(body), &res)
	if err != nil {
		return nil, err
	}

	if res.Source == "" || res.Target == "" {
		return nil, errors.New("source and target are empty for body")
	}

	return &res, nil
}
