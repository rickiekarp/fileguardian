package fileprocessor

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"git.rickiekarp.net/rickie/fileguardian/config"
	"git.rickiekarp.net/rickie/yubase"
)

func Run(args []string) error {

	if len(args) < 1 {
		return errors.New("no argument provided")
	}

	processMode := getProcessMode()

	for _, arg := range args {

		quantumLockMode, pathExists := getQuantumLockMode(arg)
		input := arg
		hasDataExtension := false
		if pathExists {
			input = filepath.Base(arg)
			hasDataExtension = strings.HasSuffix(input, "."+yubase.FileGuardianDataExtension)
		}

		switch processMode {
		case Print:
			resp, err := sendRequest(input, quantumLockMode, config.StorageContext)
			if err != nil {
				return err
			}

			if hasDataExtension {
				if resp != nil {
					fmt.Println(resp.Source)
				}
			} else {
				// print result depending on the fetched base file
				if strings.EqualFold(input, resp.Source) {
					fmt.Println(resp.Target)
				} else if strings.EqualFold(input, resp.Target) {
					fmt.Println(resp.Source)
				}
			}

		case Decrypt:
			switch quantumLockMode {
			case yubase.File:
				if !hasDataExtension {
					return errors.New("wrong file extension found during decryption")
				}

				err := decryptFile(arg, input, quantumLockMode)
				if err != nil {
					return err
				}

			case yubase.Dir:
				err := decryptFolderContent(arg)
				if err != nil {
					return err
				}
			case yubase.String:
				return errors.New("wrong quantum lock mode found during decryption")
			}

		case Encrypt:

			switch quantumLockMode {
			case yubase.File:
				if hasDataExtension {
					return errors.New("can not encrypt an already encrypted file")
				}

				err := encryptFile(arg, input, quantumLockMode)
				if err != nil {
					return err
				}

			case yubase.Dir:
				err := encryptFolderContent(arg)
				if err != nil {
					return err
				}
			case yubase.String:
				return errors.New("wrong quantum lock mode found during encryption")
			}
		}
	}

	return nil
}
