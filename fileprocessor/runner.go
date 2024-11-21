package fileprocessor

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"git.rickiekarp.net/rickie/fileguardian/config"
	"git.rickiekarp.net/rickie/fileguardian/utils"
	"git.rickiekarp.net/rickie/goutilkit"
	"git.rickiekarp.net/rickie/nexuscore"
	"git.rickiekarp.net/rickie/nexusform"
	"github.com/sirupsen/logrus"
)

func Run(args []string) error {

	// default mode: print
	processMode := Print

	// encrypts or decrypts a given file
	if *config.FlagEncrypt {
		processMode = Encrypt
	} else if *config.FlagDecrypt {
		processMode = Decrypt
	}

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
					return err
				}

				switch processMode {
				case Print:
					// print result if it exists
					if resp != nil {
						fmt.Println(resp.Source)
					}
				case Decrypt:
					if resp == nil {
						return errors.New("file not found for decryption: " + baseFile)
					}

					pathExists, _ := goutilkit.PathExists(arg)
					if !pathExists {
						return errors.New("file not found on disk: " + baseFile)
					}

					// load config file
					err := config.ReadConfigFile()
					if err != nil {
						return errors.New("could not load config")
					}

					vaultEntry, err := Fetch(config.Conf.Application.VaultIdentifier, config.Conf.Application.Token)
					if err != nil {
						return errors.New("could not fetch data from vault")
					}
					vaultContent := string(vaultEntry.Content)

					if len(vaultContent) == 0 {
						return errors.New("invalid content found")
					}
					return utils.Decrypt(resp.Source, arg, vaultContent)
				case Encrypt:
					return errors.New("can't encrypt an already encrypted file")
				}
			} else {

				// check if given file path arg exists locally
				pathExists, fileInfo := goutilkit.PathExists(arg)
				if pathExists {
					if (*fileInfo).IsDir() {
						fileType = "dir"
					}
				} else {
					fileType = ""
				}

				// fetch entry
				resp, err := sendRequest(baseFile, fileType, config.StorageContext)
				if err != nil {
					return err
				}

				switch processMode {
				case Print:
					// print result depending on the fetched base file
					if resp != nil {
						switch baseFile {
						case resp.Source:
							fmt.Println(resp.Target)
						case resp.Target:
							fmt.Println(resp.Source)
						}
					}
				case Decrypt:
					return errors.New("can't decrypt an already decrypted file")
				case Encrypt:
					if fileType != "file" {
						return errors.New("wrong fileType detected for encryption: " + fileType)
					}

					recipient := *config.FlagRecipient

					// make sure pgp recipient is set
					if len(*config.FlagRecipient) == 0 {
						// load config file
						err := config.ReadConfigFile()
						if err != nil {
							return errors.New("could not load config")
						}

						if len(config.Conf.Application.Recipient) == 0 {
							return errors.New("invalid recipient configured")
						}

						recipient = config.Conf.Application.Recipient
					}

					if len(recipient) == 0 {
						return errors.New("invalid recipient found")
					}

					return utils.Encrypt(resp.Source, resp.Target, recipient)
				}
			}
		}
	}

	return nil
}

func Fetch(identifier string, token string) (*nexusform.VaultEntry, error) {
	url := config.GetApiUrl() + nexuscore.ApiVaultFetchKey

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Vault-Identifier", identifier)
	req.Header.Set("X-Vault-Token", token)

	// We can set the content type here
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	if len(body) == 0 {
		return nil, errors.New("no file found")
	}

	res := nexusform.VaultEntry{}
	err = json.Unmarshal([]byte(body), &res)
	if err != nil {
		logrus.Error("Could not unmarshal weather data! ", err)
		return nil, err
	}

	return &res, nil
}
