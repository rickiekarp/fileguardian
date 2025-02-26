package config

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ConfigFile struct {
	Application struct {
		VaultIdentifier string `yaml:"vaultIdentifier"`
		Token           string `yaml:"token"`
		Recipient       string `yaml:"recipient"`
	}
}

var Conf ConfigFile

// ReadConfigFile reads the given config file and tries to unmarshal it into the given configStruct
func ReadConfigFile() error {

	// if the ConfigBaseDir has not been set (e.g. from ldflags), set it here
	if ConfigBaseDir == "" {
		binaryDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logrus.Error("binaryDir.Get err: ", err)
			return err
		}
		ConfigBaseDir = binaryDir + "/"
	}

	// read config file
	yamlFile, err := os.ReadFile(ConfigBaseDir + "config.yaml")
	if err != nil {
		logrus.Error("yamlFile.Get err: ", err)
		return err
	}

	// unmarshal config file depending on given configStruct
	err = yaml.Unmarshal(yamlFile, &Conf)
	if err != nil {
		logrus.Error("Unmarshal failed: ", err)
		return err
	}

	return nil
}
