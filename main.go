package main

import (
	"log"
	"os"

	"git.rickiekarp.net/rickie/fileguardian/modules/filestorage"
	"git.rickiekarp.net/rickie/fileguardian/storage"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	storage.Init()

	files := []filestorage.File{}
	if *storage.ShouldAddToStorage {
		// load files from disk
		if *storage.DataPath == "" {
			log.Println("No path given!")
			os.Exit(1)
		}
		files = filestorage.LoadDataFromDisk(*storage.DataPath)
	}

	// open connection to database and add data sets
	if *storage.ShouldAddToStorage || *storage.ShouldListStorage {
		storage.OpenDatabase()
		defer storage.StoragePtr.Close()
	}

	// write files to storage
	if *storage.ShouldAddToStorage {
		filestorage.AddToStorage(storage.StoragePtr, files)
	}

	if *storage.ShouldListStorage {
		filestorage.DisplayEntries(storage.StoragePtr)
	}
}
