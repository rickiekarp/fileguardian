package fileprocessor

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"git.rickiekarp.net/rickie/fileguardian/config"
	"git.rickiekarp.net/rickie/fileguardian/utils"
	"git.rickiekarp.net/rickie/goutilkit"
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

					// load config file
					err := config.ReadConfigFile()
					if err != nil {
						return errors.New("could not load config")
					}

					passphrase := utils.ReadStringFromFile(config.Conf.Application.PassphraseFile)
					if len(passphrase) == 0 {
						return errors.New("invalid passphrase given")
					}
					return utils.Decrypt(resp.Source, resp.Target, passphrase)
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
