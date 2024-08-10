package filestorage

import (
	"fmt"
	"log"
	"strings"

	"git.rickiekarp.net/rickie/fileguardian/config"
)

func DisplayEntries(account Storage) {

	for _, value := range account.Files[config.StorageContext] {
		switch *config.InstructionType {
		case "encrypt":
			if !*config.ShouldSkipCompression {
				PrintCompression(value.Src, value.Dst)
			}
			if *config.EncryptionRecipient == "" {
				log.Fatal("No recipient set, ignoring encryption instructions! Please provide the encryptionRecipient flag")
			} else {
				PrintEncryption(value.Src, value.Dst, *config.EncryptionRecipient)
			}
		case "decrypt":
			PrintDecryption(value.Src, value.Dst)
			if !*config.ShouldSkipCompression {
				PrintExtract(value.Src, value.Dst)
			}
		default:
			log.Println("File: ", value.ProcessId, " ", value.Type, " ", value.Src, " ", value.Dst)
		}
	}
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
