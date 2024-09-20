package main

import (
	"bytes"
	"encoding/json"
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

			// if the file has a DataExcention, we assume it's already hashed and therefor the target
			if strings.HasSuffix(baseFile, "."+config.DataExtension) {

				// fetch entry
				resp := request(baseFile, fileType, config.StorageContext)
				if resp == nil {
					os.Exit(1)
				}

				fmt.Println(resp.Source)

			} else {

				fileInfo := filestorage.Evaluate(arg)
				if fileInfo != nil {
					if (*fileInfo).IsDir() {
						fileType = "dir"
					}
				} else {
					fileType = ""
				}

				resp := request(baseFile, fileType, config.StorageContext)
				if resp == nil {
					os.Exit(1)
				}

				switch baseFile {
				case resp.Source:
					fmt.Println(resp.Target)
				case resp.Target:
					fmt.Println(resp.Source)
				}
			}
		}
	}
}

func request(fileName string, fileType string, context string) *FileGuardianEventMessage {
	url := config.ApiProtocol + "://" + config.ApiHost + "/fileguardian/v1/fetch"

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
		logrus.Error(err)
		return nil
	}

	// We can set the content type here
	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		logrus.Error(err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		logrus.Error("Could not fetch data")
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	var res FileGuardianEventMessage
	err = json.Unmarshal([]byte(body), &res)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	if res.Source == "" || res.Target == "" {
		logrus.Error("source and target are empty for body: ", string(body))
		return nil
	}

	return &res
}

type FileGuardianEventMessage struct {
	Id         *int64 `json:"id,omitempty"`
	Type       string `json:"type,omitempty"`
	Source     string `json:"source,omitempty"`
	Target     string `json:"target,omitempty"`
	Context    string `json:"context,omitempty"`
	Inserttime *int64 `json:"inserttime,omitempty"`
}
