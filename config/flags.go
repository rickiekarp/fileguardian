package config

import (
	"flag"
)

var StorageContext string

var PrintHelp = flag.Bool("h", false, "print help")
var PrintVersion = flag.Bool("v", false, "prints version")

func init() {
	flag.StringVar(&StorageContext, "context", "default", "storage table context")
	flag.Parse()
}
