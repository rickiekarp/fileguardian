package utils

import (
	"fmt"
	"strings"

	"git.rickiekarp.net/rickie/fileguardian/config"
	"git.rickiekarp.net/rickie/gomain"
)

/// Encryption helper functions

func Encrypt(src string, dst string, recipient string) error {
	outDir := *config.FlagOutput
	if len(outDir) > 0 && !strings.HasSuffix(*config.FlagOutput, "/") {
		outDir = outDir + "/"
	}

	opts := gomain.DefaultOptions()
	opts.OutputType = gomain.CommandOutputSilent

	cmdResp := gomain.ExecuteWithOptions(
		opts,
		"gpg",
		"--output", outDir+dst,
		"--encrypt",
		"--recipient", recipient,
		src,
	)

	if cmdResp.ExitCode != 0 || cmdResp.Err != nil {
		return fmt.Errorf("could not encrypt file (exit code: %d)", cmdResp.ExitCode)
	}

	return nil
}

/// Decryption helper functions

func Decrypt(sourceFile string, target string, passphrase string) error {
	outDir := *config.FlagOutput
	if len(outDir) > 0 && !strings.HasSuffix(*config.FlagOutput, "/") {
		outDir = outDir + "/"
	}

	opts := gomain.DefaultOptions()
	opts.OutputType = gomain.CommandOutputSilent

	cmdResp := gomain.ExecuteWithOptions(
		opts,
		"gpg",
		"--pinentry-mode", "loopback",
		"--passphrase", passphrase,
		"--output", outDir+sourceFile,
		"--decrypt", target,
	)
	if cmdResp.ExitCode != 0 || cmdResp.Err != nil {
		return fmt.Errorf("could not decrypt file (exit code: %d)", cmdResp.ExitCode)
	}
	return nil
}
