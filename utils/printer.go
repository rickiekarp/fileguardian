package utils

import (
	"fmt"
	"strings"

	"git.rickiekarp.net/rickie/fileguardian/config"
)

/// Encryption helper functions

func PrintCompression(src string, dst string) {
	// ignore file that seems to compressed already
	if strings.HasSuffix(src, "."+config.CompressExtension) {
		return
	}

	fmt.Println("tar cJPf " + src + "." + config.CompressExtension + " " + src)
}

func PrintEncryption(src string, dst string, recipient string) {

	if !strings.HasSuffix(src, "."+config.CompressExtension) {
		src += "." + config.CompressExtension
	}

	fmt.Println("gpg --output " + dst + " --encrypt --recipient " + recipient + " " + src)
}

/// Decryption helper functions

func PrintDecryption(src string, dst string) {
	if !strings.HasSuffix(src, "."+config.CompressExtension) {
		fmt.Println("gpg --output " + src + "." + config.CompressExtension + " --decrypt " + dst)
		return
	}

	fmt.Println("gpg --output " + src + " --decrypt " + dst)
}
