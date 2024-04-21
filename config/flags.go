package config

import "flag"

var ShouldCreateStorage = flag.Bool("create", false, "create storage file")
var ShouldListStorage = flag.Bool("list", true, "list storage content")
var ShouldAddToStorage = flag.Bool("add", false, "add content to storage")
var ShouldSkipCompression = flag.Bool("skipTar", false, "whether to skip file compression")
var StorageName = flag.String("name", "storage", "storage name")
var DataPath = flag.String("data", "", "data path")

var PrintHelp = flag.Bool("h", false, "print help")
var PrintVersion = flag.Bool("v", false, "prints version")
