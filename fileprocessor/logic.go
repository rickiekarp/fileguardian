package fileprocessor

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"git.rickiekarp.net/rickie/fileguardian/config"
	"git.rickiekarp.net/rickie/fileguardian/utils"
	"git.rickiekarp.net/rickie/goutilkit"
	"git.rickiekarp.net/rickie/nexuscore"
	"git.rickiekarp.net/rickie/nexusform"
	"github.com/sirupsen/logrus"
)

func encryptFile(arg string, baseFile string, fileType FileType) error {

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

	// fetch entry
	resp, err := sendRequest(baseFile, fileType, config.StorageContext)
	if err != nil {
		return err
	}

	return utils.Encrypt(arg, resp.Target, recipient)
}

func encryptFolderContent(folderPath string) error {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return err
	}

	// make sure the arg argument has a / suffix
	if !strings.HasSuffix(folderPath, "/") {
		folderPath = folderPath + "/"
	}

	for _, entry := range entries {
		var baseFile = entry.Name()
		if strings.HasSuffix(baseFile, config.DataExtension) {
			continue
		}
		err := encryptFile(folderPath+baseFile, baseFile, File)
		if err != nil {
			return err
		}
	}

	return nil
}

func decryptFile(arg string, baseFile string, fileType FileType) error {
	// fetch entry
	resp, err := sendRequest(baseFile, fileType, config.StorageContext)
	if err != nil {
		return err
	}

	if resp == nil {
		return errors.New("file not found for decryption: " + baseFile)
	}

	pathExists, _ := goutilkit.PathExists(arg)
	if !pathExists {
		return errors.New("file not found on disk: " + baseFile)
	}

	// load config file
	err = config.ReadConfigFile()
	if err != nil {
		return errors.New("could not load config")
	}

	vaultEntry, err := fetchPassphrase(config.Conf.Application.VaultIdentifier, config.Conf.Application.Token)
	if err != nil {
		return errors.New("could not fetch data from vault")
	}

	if len(vaultEntry.Content) == 0 {
		return errors.New("invalid content found")
	}
	return utils.Decrypt(resp.Source, arg, string(vaultEntry.Content))
}

func decryptFolderContent(folderPath string) error {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return err
	}

	// make sure the folderPath argument has a / suffix
	if !strings.HasSuffix(folderPath, "/") {
		folderPath = folderPath + "/"
	}

	for _, entry := range entries {
		var baseFile = entry.Name()
		if !strings.HasSuffix(baseFile, config.DataExtension) {
			continue
		}
		err := decryptFile(folderPath+baseFile, baseFile, File)
		if err != nil {
			return err
		}
	}

	return nil
}

func fetchPassphrase(identifier string, token string) (*nexusform.VaultEntry, error) {
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
