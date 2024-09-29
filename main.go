package main

import (
	"flag"
	"fmt"
	"os"

	"git.rickiekarp.net/rickie/fileguardian/config"
	"git.rickiekarp.net/rickie/fileguardian/fileprocessor"
	"git.rickiekarp.net/rickie/filesanitizer"
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
	fileprocessor.Run(arguments)
}
