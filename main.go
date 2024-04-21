package main

import (
	"log"
	"os"

	"git.rickiekarp.net/rickie/fileguardian/config"
	"git.rickiekarp.net/rickie/fileguardian/modules/filestorage"
	"git.rickiekarp.net/rickie/fileguardian/storage"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	storage.Init()

	files := []filestorage.File{}
	if *config.ShouldAddToStorage {
		// load files from disk
		if *config.DataPath == "" {
			log.Println("No path given!")
			os.Exit(1)
		}
		files = filestorage.LoadDataFromDisk(*config.DataPath)
	}

	// open connection to database and add data sets
	if *config.ShouldAddToStorage || *config.ShouldListStorage {
		storage.OpenDatabase()
		defer storage.StoragePtr.Close()
	}

	// write files to storage
	if *config.ShouldAddToStorage {
		filestorage.AddToStorage(storage.StoragePtr, files)
	}

	if *config.ShouldListStorage {
		filestorage.DisplayEntries(storage.StoragePtr)
	}
}
