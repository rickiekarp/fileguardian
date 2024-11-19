package config

var ApiProtocol = "http"          // ApiProtocol set during go build using ldflags
var ApiHost = "localhost:12000"   // ApiHost set during go build using ldflags
var ConfigBaseDir = "deployment/" // ConfigBaseDir set during go build using ldflags

func GetApiUrl() string {
	return ApiProtocol + "://" + ApiHost
}
