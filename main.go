package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"git.rickiekarp.net/rickie/fileguardian/config"
	"git.rickiekarp.net/rickie/fileguardian/modules/filestorage"
	"github.com/sirupsen/logrus"
)

func main() {

	if *config.PrintHelp {
		config.PrintUsage()
		os.Exit(0)
	}

	if *config.PrintVersion {
		fmt.Println(config.Version)
		os.Exit(0)
	}

	arguments := flag.Args()

	if len(arguments) > 0 {

		fileType := "file"

		for _, arg := range arguments {

			baseFile := filepath.Base(arg)

			// if the file has a DataExtension, we assume it's already hashed and therefor the target
			if strings.HasSuffix(baseFile, "."+config.DataExtension) {

				// fetch entry
				resp, err := sendRequest(baseFile, fileType, config.StorageContext)
				if err != nil {
					logrus.Error(err)
					os.Exit(1)
				}

				if resp != nil {
					fmt.Println(resp.Source)
				} else {
					os.Exit(0)
				}

			} else {

				fileInfo := filestorage.PathExists(arg)
				if fileInfo != nil {
					if (*fileInfo).IsDir() {
						fileType = "dir"
					}
				} else {
					fileType = ""
				}

				resp, err := sendRequest(baseFile, fileType, config.StorageContext)
				if err != nil {
					logrus.Error(err)
					os.Exit(1)
				}

				if resp != nil {
					switch baseFile {
					case resp.Source:
						fmt.Println(resp.Target)
					case resp.Target:
						fmt.Println(resp.Source)
					}
				} else {
					os.Exit(0)
				}
			}
		}
	}
}

func sendRequest(fileName string, fileType string, context string) (*FileGuardianEventMessage, error) {
	url := config.ApiProtocol + "://" + config.ApiHost + "/fileguardian/v1/fetch"

	if *config.FlagCheck {
		url += "?check=true"
	}

	// create post body using an instance of the Person struct
	requestEvent := FileGuardianEventMessage{
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

	var res FileGuardianEventMessage
	err = json.Unmarshal([]byte(body), &res)
	if err != nil {
		return nil, err
	}

	if res.Source == "" || res.Target == "" {
		return nil, errors.New("source and target are empty for body")
	}

	return &res, nil
}

type FileGuardianEventMessage struct {
	Id         *int64 `json:"id,omitempty"`
	Type       string `json:"type,omitempty"`
	Source     string `json:"source,omitempty"`
	Target     string `json:"target,omitempty"`
	Context    string `json:"context,omitempty"`
	Inserttime *int64 `json:"inserttime,omitempty"`
}
