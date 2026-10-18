package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetApiUrl(t *testing.T) {
	oldProtocol, oldHost := ApiProtocol, ApiHost
	defer func() {
		ApiProtocol = oldProtocol
		ApiHost = oldHost
	}()

	ApiProtocol = "https"
	ApiHost = "example.com:1234"

	if got := GetApiUrl(); got != "https://example.com:1234" {
		t.Errorf("GetApiUrl() = %q, want %q", got, "https://example.com:1234")
	}
}

func TestReadConfigFile(t *testing.T) {
	dir := t.TempDir()
	content := `application:
  vaultIdentifier: my-vault
  token: secret-token
  recipient: alice@example.com
`
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(content), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	oldBaseDir, oldConf := ConfigBaseDir, Conf
	defer func() {
		ConfigBaseDir = oldBaseDir
		Conf = oldConf
	}()

	ConfigBaseDir = dir + "/"

	if err := ReadConfigFile(); err != nil {
		t.Fatalf("ReadConfigFile: %v", err)
	}

	if Conf.Application.VaultIdentifier != "my-vault" {
		t.Errorf("VaultIdentifier = %q, want %q", Conf.Application.VaultIdentifier, "my-vault")
	}
	if Conf.Application.Token != "secret-token" {
		t.Errorf("Token = %q, want %q", Conf.Application.Token, "secret-token")
	}
	if Conf.Application.Recipient != "alice@example.com" {
		t.Errorf("Recipient = %q, want %q", Conf.Application.Recipient, "alice@example.com")
	}
}

func TestReadConfigFileMissing(t *testing.T) {
	oldBaseDir := ConfigBaseDir
	defer func() { ConfigBaseDir = oldBaseDir }()

	ConfigBaseDir = t.TempDir() + "/"

	if err := ReadConfigFile(); err == nil {
		t.Fatal("expected error for missing config file, got nil")
	}
}
