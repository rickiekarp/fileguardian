package fileprocessor

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"git.rickiekarp.net/rickie/fileguardian/config"
)

func Run(args []string) error {

	if len(args) < 1 {
		return errors.New("no argument provided")
	}

	processMode := getProcessMode()

	for _, arg := range args {

		fileMode := getFileMode(arg)
		baseFile := filepath.Base(arg)
		hasDataExtension := strings.HasSuffix(baseFile, "."+config.DataExtension)

		switch processMode {
		case Print:
			resp, err := sendRequest(baseFile, fileMode, config.StorageContext)
			if err != nil {
				return err
			}

			if hasDataExtension {
				if resp != nil {
					fmt.Println(resp.Source)
				}
			} else {
				// print result depending on the fetched base file
				if strings.EqualFold(baseFile, resp.Source) {
					fmt.Println(resp.Target)
				} else if strings.EqualFold(baseFile, resp.Target) {
					fmt.Println(resp.Source)
				}
			}

		case Decrypt:
			switch fileMode {
			case File:
				if !hasDataExtension {
					return errors.New("wrong file extension found during decryption")
				}

				err := decryptFile(arg, baseFile, fileMode)
				if err != nil {
					return err
				}

			case Dir:
				err := decryptFolderContent(arg)
				if err != nil {
					return err
				}
			}

		case Encrypt:

			switch fileMode {
			case File:
				if hasDataExtension {
					return errors.New("can not encrypt an already encrypted file")
				}

				err := encryptFile(arg, baseFile, fileMode)
				if err != nil {
					return err
				}

			case Dir:
				err := encryptFolderContent(arg)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
