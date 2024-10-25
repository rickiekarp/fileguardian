package main

import (
	"flag"
	"fmt"
	"os"

	"git.rickiekarp.net/rickie/fileguardian/config"
	"git.rickiekarp.net/rickie/fileguardian/fileprocessor"
	"git.rickiekarp.net/rickie/goutilkit"
	"github.com/sirupsen/logrus"
)

func main() {

	if *config.PrintHelp {
		goutilkit.PrintUsageAndExit()
	}

	if *config.PrintVersion {
		fmt.Println(config.Version)
		os.Exit(0)
	}

	arguments := flag.Args()

	validateInput()

	err := fileprocessor.Run(arguments)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func validateInput() {
	if *config.FlagEncrypt && *config.FlagDecrypt {
		logrus.Error("can't use -e and -d flag at the same time")
		os.Exit(1)
	}
}
