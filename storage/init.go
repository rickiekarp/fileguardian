package storage

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"git.rickiekarp.net/rickie/fileguardian/config"
	"git.rickiekarp.net/rickie/fileguardian/modules/filestorage"
)

var Version = "development" // Version set during go build using ldflags

var StoragePtr *sql.DB

func Init() {
	flag.Parse()

	if *config.PrintHelp {
		config.PrintUsage()
		os.Exit(0)
	}

	if *config.PrintVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	if *config.ShouldCreateStorage {
		CreateDatabase()
		OpenDatabase()
		filestorage.CreateTable(StoragePtr)
		StoragePtr.Close()
	}
}

func CreateDatabase() {
	log.Println("Creating database...")
	file, err := os.Create(*config.StorageName + "." + config.Extension)
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("Database created")
}

func OpenDatabase() {
	StoragePtr, _ = sql.Open("sqlite3", "./"+*config.StorageName+"."+config.Extension)
}
