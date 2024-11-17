package config

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ConfigFile struct {
	Application struct {
		PassphraseFile string `yaml:"passphraseFile"`
		Recipient      string `yaml:"recipient"`
	}
}

var Conf ConfigFile

// ReadConfigFile reads the given config file and tries to unmarshal it into the given configStruct
func ReadConfigFile() error {

	yamlFile, err := os.ReadFile("config.yaml")
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
