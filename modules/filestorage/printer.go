package filestorage

import (
	"database/sql"
	"fmt"
	"log"
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
		log.Fatal(err)
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
			PrintCompression(src, dst)
			if *EncryptionRecipient == "" {
				log.Fatal("No recipient set, ignoring encryption instructions! Please provide the encryptionRecipient flag")
			} else {
				PrintEncryption(src, dst, *EncryptionRecipient)
			}
		case "decrypt":
			PrintDecryption(src, dst)
			PrintExtract(src, dst)
		default:
			log.Println("File: ", id, " ", tag, " ", src, " ", dst)
		}

	}
}

/// Encryption helper functions

func PrintCompression(src string, dst string) {
	fmt.Println("tar cJPf " + src + ".tar.xz " + src)
}

func PrintEncryption(src string, dst string, recipient string) {
	fmt.Println("gpg --output " + dst + " --encrypt --recipient " + recipient + " " + src + ".tar.xz")
}

/// Decryption helper functions

func PrintDecryption(src string, dst string) {
	fmt.Println("gpg --output " + src + ".tar.xz --decrypt " + dst)
}

func PrintExtract(src string, dst string) {
	fmt.Println("tar -xf " + src + ".tar.xz")
}

func PrintData(files []File) {
	for _, elem := range files {
		fmt.Println(elem)
	}
}
