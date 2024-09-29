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
	"git.rickiekarp.net/rickie/fileguardian/fileprocessor"
	"git.rickiekarp.net/rickie/fileguardian/utils"
	"git.rickiekarp.net/rickie/filesanitizer"
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

	// if the -s flag is set, attempt to sanitize the filenames of all files in a given directory
	if *config.FlagSanitizer {
		if len(arguments) > 0 {
			filesanitizer.SanitizeFilesInFolder(arguments[0])
		} else {
			os.Exit(1)
		}
		os.Exit(0)
	}

	// encrypts or decrypts a given file
	if *config.FlagEncrypt || *config.FlagDecrypt {
		os.Exit(0)
	}

	// process the given files
	run(arguments)
}

func run(args []string) {
	// fetch source/hashed file name from the storage
	if len(args) > 0 {

		fileType := "file"

		for _, arg := range args {

			baseFile := filepath.Base(arg)

			// if the file has a DataExtension, we assume it's already hashed and therefor the target
			if strings.HasSuffix(baseFile, "."+config.DataExtension) {

				// fetch entry
				resp, err := sendRequest(baseFile, fileType, config.StorageContext)
				if err != nil {
					logrus.Error(err)
					os.Exit(1)
				}

				// print result
				if resp != nil {
					fmt.Println(resp.Source)
				}

			} else {

				// check if given file arg exists locally
				fileInfo := utils.PathExists(arg)
				if fileInfo != nil {
					if (*fileInfo).IsDir() {
						fileType = "dir"
					}
				} else {
					fileType = ""
				}

				// fetch entry
				resp, err := sendRequest(baseFile, fileType, config.StorageContext)
				if err != nil {
					logrus.Error(err)
					os.Exit(1)
				}

				// print result
				if resp != nil {
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
}

func sendRequest(fileName string, fileType string, context string) (*fileprocessor.FileGuardianEventMessage, error) {
	url := config.ApiProtocol + "://" + config.ApiHost + "/fileguardian/v1/fetch"

	if *config.FlagCheck {
		url += "?check=true"
	}

	// create post body using an instance of the Person struct
	requestEvent := fileprocessor.FileGuardianEventMessage{
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

	var res fileprocessor.FileGuardianEventMessage
	err = json.Unmarshal([]byte(body), &res)
	if err != nil {
		return nil, err
	}

	if res.Source == "" || res.Target == "" {
		return nil, errors.New("source and target are empty for body")
	}

	return &res, nil
}
