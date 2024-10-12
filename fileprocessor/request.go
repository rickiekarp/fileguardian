package fileprocessor

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"git.rickiekarp.net/rickie/fileguardian/config"
	"git.rickiekarp.net/rickie/nexusform"
)

func sendRequest(fileName string, fileType string, context string) (*nexusform.FileGuardianEntry, error) {
	url := config.ApiProtocol + "://" + config.ApiHost + "/fileguardian/v1/fetch"

	if *config.FlagCheck {
		url += "?check=true"
	}

	// create post body using an instance of the Person struct
	requestEvent := nexusform.FileGuardianEntry{
		Type:    fileType,
		Context: context,
	}

	// from a fileguardian perspective our file is always the source
	requestEvent.Source = fileName

	// convert p to JSON data
	jsonData, err := json.Marshal(requestEvent)
	if err != nil {
		return nil, errors.New("could not marshal requestEvent")
	}

	// We can set the content type here
	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonData))
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

	var res nexusform.FileGuardianEntry
	err = json.Unmarshal([]byte(body), &res)
	if err != nil {
		return nil, err
	}

	if res.Source == "" || res.Target == "" {
		return nil, errors.New("source and target are empty for body")
	}

	return &res, nil
}
