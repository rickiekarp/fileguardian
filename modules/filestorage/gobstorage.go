package filestorage

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"
	"time"

	"git.rickiekarp.net/rickie/fileguardian/config"
	"git.rickiekarp.net/rickie/fileguardian/utils"
	"github.com/sirupsen/logrus"
)

var Profile *Storage

func Generate() (*Storage, error) {

	initDate := time.Now().UTC().Unix()

	// Create an instance of the account struct
	account := Storage{
		Id:           utils.RandSeq(32),
		CreationDate: initDate,
		LastModified: initDate,
		Files:        make(map[string][]File),
	}
	return &account, nil
}

func Persist(account Storage) {

	account.LastModified = time.Now().UTC().Unix()

	// Create a new buffer to write the serialized data to
	var b bytes.Buffer

	// Create a new gob encoder and use it to encode the account struct
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(account); err != nil {
		logrus.Println("Error encoding struct:", err)
		return
	}

	// The serialized data can now be found in the buffer
	serializedData := b.Bytes()

	storageFile := config.GetStorageFile()
	err := os.WriteFile(storageFile, serializedData, 0644)
	if err != nil {
		logrus.Println("Could not write file", config.StorageFileName)
	}
}

func Load(storageFile string) (*Storage, error) {

	file, err := os.Open(storageFile)
	if err != nil {
		logrus.Println(err)
		return nil, err
	}
	defer file.Close()

	// The serialized data to be deserialized
	serializedData, err := os.ReadFile(storageFile)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new buffer from the serialized data
	b := bytes.NewBuffer(serializedData)

	// Create a new gob decoder and use it to decode the account struct
	var account Storage
	dec := gob.NewDecoder(b)
	if err := dec.Decode(&account); err != nil {
		logrus.Println("Error decoding struct:", err)
		return nil, err
	}

	// The account struct has now been deserialized
	return &account, nil
}
