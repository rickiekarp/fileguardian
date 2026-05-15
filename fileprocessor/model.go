package fileprocessor

import (
	"git.rickiekarp.net/rickie/fileguardian/config"
	"git.rickiekarp.net/rickie/gomain"
	"git.rickiekarp.net/rickie/yubase"
)

//
// ProcessType
//

type ProcessType int

const (
	Print ProcessType = iota
	Encrypt
	Decrypt
)

func getProcessMode() ProcessType {

	// default mode: print
	processMode := Print

	// encrypts or decrypts a given file
	if *config.FlagEncrypt {
		processMode = Encrypt
	} else if *config.FlagDecrypt {
		processMode = Decrypt
	}

	return processMode
}

func getQuantumLockMode(filePath string) (yubase.QuantumLockType, bool) {
	lockType := yubase.None

	// check if given file path exists locally
	pathExists, fileInfo := gomain.PathExists(filePath)
	if pathExists {
		if (*fileInfo).IsDir() {
			lockType = yubase.Dir
		} else {
			lockType = yubase.File
		}
	} else {
		lockType = yubase.String
	}

	return lockType, pathExists
}
