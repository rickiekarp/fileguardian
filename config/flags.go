package config

import (
	"flag"
)

var StorageContext string

var PrintHelp = flag.Bool("h", false, "print help")
var PrintVersion = flag.Bool("v", false, "prints version")

var FlagDecrypt = flag.Bool("d", false, "decrypts a given file")
var FlagEncrypt = flag.Bool("e", false, "encrypts a given file")
var FlagRecipient = flag.String("r", "", "sets the recipient of an encrypted message")
var FlagCheck = flag.Bool("c", false, "only checks if the given file exists on the remote")
var FlagSanitizer = flag.Bool("s", false, "sanitizes all file names in a given directory")

func init() {
	flag.StringVar(&StorageContext, "context", "default", "storage table context")
	flag.Parse()
}
