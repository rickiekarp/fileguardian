package utils

import (
	"fmt"

	"git.rickiekarp.net/rickie/goutilkit"
)

/// Encryption helper functions

func Encrypt(src string, dst string, recipient string) error {
	command := "gpg --output " + dst + " --encrypt --recipient " + recipient + " " + src
	exitCode, err := goutilkit.ExecuteCmdSilent(command)
	if exitCode != 0 || err != nil {
		return fmt.Errorf("could not encrypt file (exit code: %d)", exitCode)
	}
	return nil
}

/// Decryption helper functions

func Decrypt(src string, dst string, passphrase string) error {
	command := "gpg --pinentry-mode loopback --passphrase " + passphrase + " --output " + "src" + " --decrypt " + dst
	exitCode, err := goutilkit.ExecuteCmdSilent(command)
	if exitCode != 0 || err != nil {
		return fmt.Errorf("could not decrypt file (exit code: %d)", exitCode)
	}
	return nil
}
