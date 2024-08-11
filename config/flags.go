package config

import (
	"flag"
)

var StoragePath string
var StorageName string
var StorageContext string

var StorageFileName string

var ShouldSkipCompression = flag.Bool("skipTar", false, "whether to skip file compression")

var InstructionType = flag.String("instructionType", "", "prints compression instructions")

// gpg recipient
var EncryptionRecipient = flag.String("encryptionRecipient", "", "sets the recipient of the encrypted file")

var PrintHelp = flag.Bool("h", false, "print help")
var PrintVersion = flag.Bool("v", false, "prints version")

func init() {
	flag.StringVar(&StoragePath, "dir", "", "storage name")
	flag.StringVar(&StorageName, "name", "storage", "storage name")
	flag.StringVar(&StorageContext, "context", "default", "storage table context")
	flag.Parse()

	StorageFileName = StorageName + "." + Extension
}

func GetStorageFile() string {
	storageFilePath := StorageFileName
	if StoragePath != "" {
		storageFilePath = StoragePath + "/" + StorageFileName
	}
	return storageFilePath
}
