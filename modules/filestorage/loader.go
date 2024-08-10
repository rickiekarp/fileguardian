package filestorage

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"git.rickiekarp.net/rickie/fileguardian/config"
)

func LoadDataFromDisk(path string) ([]File, error) {
	files := []File{}

	// This returns an *os.FileInfo type
	fileInfo, err := os.Stat(path)
	if err != nil {
		return files, err
	}

	now := time.Now().UTC()

	// create a processId that represents the current processing job
	var processId = strconv.Itoa(int(config.StartupTime.Unix())) + "-" + strconv.Itoa(os.Getpid())

	// populate files slice
	if fileInfo.IsDir() {

		files = append(files,
			File{
				ProcessId: processId,
				Type:      "dir",
				Src:       fileInfo.Name(),
				Dst:       fmt.Sprintf("%x", getMd5Sum(processId, time.Now().UTC().Nanosecond())) + "." + config.DataExtension,
				CreatedAt: now.Unix(),
			},
		)

	} else {
		files = append(files,
			File{
				ProcessId: processId,
				Type:      "file",
				Src:       fileInfo.Name(),
				Dst:       fmt.Sprintf("%x", getMd5Sum(processId, time.Now().UTC().Nanosecond())) + "." + config.DataExtension,
				CreatedAt: now.Unix(),
			})
	}

	return files, nil
}

func getMd5Sum(processId string, idx int) [16]byte {
	return md5.Sum([]byte(processId + "-" + strconv.Itoa(idx)))
}

func IsValidPath(path string) (bool, bool) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Println(err)
		return false, false
	}

	if fileInfo.IsDir() {
		return true, false
	}

	return true, true
}
