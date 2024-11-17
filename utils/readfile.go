package utils

import (
	"io"
	"log"
	"os"
)

func ReadStringFromFile(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Println(err)
		}
	}()
	b, err := io.ReadAll(file)
	if err != nil {
		return ""
	}
	return string(b)
}
