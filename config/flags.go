package config

import (
	"flag"
)

var StorageContext string

var PrintHelp = flag.Bool("h", false, "print help")
var PrintVersion = flag.Bool("v", false, "prints version")
var FlagCheck = flag.Bool("c", false, "only checks if the given file exists on the remote")

func init() {
	flag.StringVar(&StorageContext, "context", "default", "storage table context")
	flag.Parse()
}
