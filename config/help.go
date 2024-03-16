package config

import (
	"flag"
	"fmt"
	"os"
)

// Function to print basic program usage
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(0)
}
