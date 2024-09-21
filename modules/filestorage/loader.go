package filestorage

import (
	"io/fs"
	"os"
)

func PathExists(path string) *fs.FileInfo {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil
	}
	return &fileInfo
}
