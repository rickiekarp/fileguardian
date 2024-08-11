package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git.rickiekarp.net/rickie/fileguardian/config"
	"git.rickiekarp.net/rickie/fileguardian/modules/filestorage"
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

	var acc *filestorage.Storage
	if _, err := os.Stat(config.StorageFileName); errors.Is(err, os.ErrNotExist) {
		acc, _ = filestorage.Generate()
		filestorage.Persist(*acc)
	} else {
		acc, _ = filestorage.Load()
	}

	if len(arguments) > 0 {

		var storageModified = false
		var foundFile *filestorage.File = nil

		for _, arg := range arguments {

			foundFile = nil
			if strings.HasSuffix(arg, "."+config.DataExtension) {

				for _, storedFile := range acc.Files[config.StorageContext] {
					if storedFile.Dst == arg {
						foundFile = &storedFile
						break
					}
				}

				if foundFile != nil {
					fmt.Println(foundFile.Src)
					continue
				}
			} else {

				// look for the file in the hashed data
				for _, storedFile := range acc.Files[config.StorageContext] {
					if storedFile.Dst == filepath.Base(arg) {
						foundFile = &storedFile
						break
					}
				}

				// if the file still can't be found, look into the Src
				if foundFile == nil {
					for _, storedFile := range acc.Files[config.StorageContext] {
						if storedFile.Src == filepath.Base(arg) {
							foundFile = &storedFile
							break
						}
					}
				}

				if foundFile != nil {
					// if yes -> print decoded
					fmt.Println(foundFile.Dst)
					continue
				} else {
					loadedFiles, err := filestorage.LoadDataFromDisk(arg)
					if err != nil {
						fmt.Println("Could not load files from disk")
						fmt.Println(err)
						os.Exit(1)
					}

					acc.Files[config.StorageContext] = append(acc.Files[config.StorageContext], loadedFiles[0])
					storageModified = true
					fmt.Println(loadedFiles[0].Dst)
					continue
				}
			}
		}

		if storageModified {
			filestorage.Persist(*acc)
		}

	} else {
		filestorage.DisplayEntries(*acc)
	}
}
