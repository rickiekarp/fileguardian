package filestorage

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"git.rickiekarp.net/rickie/fileguardian/config"
)

func DisplayEntries(db *sql.DB) {

	var lookupQuery = "SELECT * FROM files ORDER BY id"

	if *LookupSrcFile != "" {
		lookupQuery = "SELECT * FROM files WHERE src = '" + *LookupSrcFile + "'"
	} else if *LookupDstFile != "" {
		lookupQuery = "SELECT * FROM files WHERE dst = '" + *LookupDstFile + "'"
	}

	row, err := db.Query(lookupQuery)
	if err != nil {
		log.Println(err)
		return
	}
	defer row.Close()
	for row.Next() {
		var id int
		var tag string
		var src string
		var dst string
		row.Scan(&id, &tag, &src, &dst)

		switch *InstructionType {
		case "encrypt":
			if !*config.ShouldSkipCompression {
				PrintCompression(src, dst)
			}
			if *EncryptionRecipient == "" {
				log.Fatal("No recipient set, ignoring encryption instructions! Please provide the encryptionRecipient flag")
			} else {
				PrintEncryption(src, dst, *EncryptionRecipient)
			}
		case "decrypt":
			PrintDecryption(src, dst)
			if !*config.ShouldSkipCompression {
				PrintExtract(src, dst)
			}
		default:
			log.Println("File: ", tag, " ", src, " ", dst)
		}
	}
}

func ReadEntry(db *sql.DB, category string, fileName string) (*File, error) {
	var lookupQuery = "SELECT * FROM files WHERE " + category + " = '" + fileName + "'"
	row, err := db.Query(lookupQuery)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if row.Next() {
		var id int
		var tag string
		var src string
		var dst string
		row.Scan(&id, &tag, &src, &dst)

		file := File{Tag: tag, Src: src, Dst: dst}
		return &file, nil
	}

	return nil, nil
}

/// Encryption helper functions

func PrintCompression(src string, dst string) {
	// ignore file that seems to compressed already
	if strings.HasSuffix(src, "."+config.CompressExtension) {
		return
	}

	fmt.Println("tar cJPf " + src + "." + config.CompressExtension + " " + src)
}

func PrintEncryption(src string, dst string, recipient string) {

	if !strings.HasSuffix(src, "."+config.CompressExtension) && !*config.ShouldSkipCompression {
		src += "." + config.CompressExtension
	}

	fmt.Println("gpg --output " + dst + " --encrypt --recipient " + recipient + " " + src)
}

/// Decryption helper functions

func PrintDecryption(src string, dst string) {
	if !strings.HasSuffix(src, "."+config.CompressExtension) && !*config.ShouldSkipCompression {
		fmt.Println("gpg --output " + src + ".tar.xz --decrypt " + dst)
		return
	}

	fmt.Println("gpg --output " + src + " --decrypt " + dst)
}

func PrintExtract(src string, dst string) {
	if !strings.HasSuffix(src, "."+config.CompressExtension) {
		src += "." + config.CompressExtension
	}
	fmt.Println("tar -xf " + src)
}

func PrintData(files []File) {
	for _, elem := range files {
		fmt.Println(elem)
	}
}
