package filestorage

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func LoadDataFromDisk(path string) []File {
	files := []File{}

	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()

	// calculate seconds since midnight
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	timeSinceMidnight := now.Sub(midnight)

	// create a tag that represents the current processing job
	var processTag = now.Format("06") + strconv.Itoa(now.YearDay()) + strconv.Itoa(int(timeSinceMidnight.Seconds()))

	// create slice of files to add to the storage
	for idx, e := range entries {
		if e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}

		bs := []byte(processTag + "-" + strconv.Itoa(idx))
		files = append(files, File{Tag: processTag, Src: e.Name(), Dst: fmt.Sprintf("%x", md5.Sum(bs)) + ".fgd"})
	}

	return files
}
