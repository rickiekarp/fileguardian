package utils

import (
	"fmt"
)

/// Encryption helper functions

func PrintEncryption(src string, dst string, recipient string) {
	fmt.Println("gpg --output " + dst + " --encrypt --recipient " + recipient + " " + src)
}

/// Decryption helper functions

func PrintDecryption(src string, dst string) {
	fmt.Println("gpg --output " + src + " --decrypt " + dst)
}
