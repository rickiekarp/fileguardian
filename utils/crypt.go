package utils

import (
	"fmt"
	"strings"

	"git.rickiekarp.net/rickie/fileguardian/config"
	"git.rickiekarp.net/rickie/goutilkit"
)

/// Encryption helper functions

func Encrypt(src string, dst string, recipient string) error {
	outDir := *config.FlagOutput
	if len(outDir) > 0 && !strings.HasSuffix(*config.FlagOutput, "/") {
		outDir = outDir + "/"
	}

	command := "gpg --output " + outDir + dst + " --encrypt --recipient " + recipient + " " + src
	exitCode, err := goutilkit.ExecuteCmdSilent(command)
	if exitCode != 0 || err != nil {
		return fmt.Errorf("could not encrypt file (exit code: %d)", exitCode)
	}
	return nil
}

/// Decryption helper functions

func Decrypt(sourceFile string, target string, passphrase string) error {
	outDir := *config.FlagOutput
	if len(outDir) > 0 && !strings.HasSuffix(*config.FlagOutput, "/") {
		outDir = outDir + "/"
	}

	command := "gpg --pinentry-mode loopback --passphrase " + passphrase + " --output " + outDir + sourceFile + " --decrypt " + target
	exitCode, err := goutilkit.ExecuteCmdSilent(command)
	if exitCode != 0 || err != nil {
		return fmt.Errorf("could not decrypt file (exit code: %d)", exitCode)
	}
	return nil
}
