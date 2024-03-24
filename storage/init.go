package storage

import (
	"database/sql"
	"flag"
	"log"
	"os"

	"git.rickiekarp.net/rickie/fileguardian/config"
	"git.rickiekarp.net/rickie/fileguardian/modules/filestorage"
)

var printHelp = flag.Bool("h", false, "print help")

var ShouldCreateStorage = flag.Bool("create", false, "create storage file")
var ShouldListStorage = flag.Bool("list", true, "list storage content")
var ShouldAddToStorage = flag.Bool("add", false, "add content to storage")
var StorageName = flag.String("name", "storage", "storage name")
var DataPath = flag.String("data", "", "data path")

var StoragePtr *sql.DB

func Init() {
	flag.Parse()

	if *printHelp {
		config.PrintUsage()
		os.Exit(0)
	}

	if *ShouldCreateStorage {
		CreateDatabase()
		OpenDatabase()
		filestorage.CreateTable(StoragePtr)
		StoragePtr.Close()
	}
}

func CreateDatabase() {
	log.Println("Creating database...")
	file, err := os.Create(*StorageName + "." + config.Extension)
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("Database created")
}

func OpenDatabase() {
	StoragePtr, _ = sql.Open("sqlite3", "./"+*StorageName+"."+config.Extension)
}
