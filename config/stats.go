package config

import "time"

var Version = "development" // Version set during go build using ldflags

// start time of the app
var StartupTime = time.Now().UTC()
