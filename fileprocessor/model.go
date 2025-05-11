package fileprocessor

import (
	"git.rickiekarp.net/rickie/fileguardian/config"
	"git.rickiekarp.net/rickie/goutilkit"
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

//
// FileType
//

type FileType string

const (
	File FileType = "file"
	Dir  FileType = "dir"
	None FileType = ""
)

func getFileMode(filePath string) FileType {
	fileType := File

	// check if given file path exists locally
	pathExists, fileInfo := goutilkit.PathExists(filePath)
	if pathExists {
		if (*fileInfo).IsDir() {
			fileType = Dir
		}
	} else {
		fileType = None
	}

	return fileType
}
