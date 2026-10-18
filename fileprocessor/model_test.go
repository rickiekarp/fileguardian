package fileprocessor

import (
	"os"
	"path/filepath"
	"testing"

	"git.rickiekarp.net/rickie/fileguardian/config"
	"git.rickiekarp.net/rickie/yubase"
)

func TestGetProcessMode(t *testing.T) {
	oldEncrypt, oldDecrypt := *config.FlagEncrypt, *config.FlagDecrypt
	defer func() {
		*config.FlagEncrypt = oldEncrypt
		*config.FlagDecrypt = oldDecrypt
	}()

	*config.FlagEncrypt = false
	*config.FlagDecrypt = false
	if got := getProcessMode(); got != Print {
		t.Errorf("getProcessMode() = %d, want Print", got)
	}

	*config.FlagEncrypt = true
	if got := getProcessMode(); got != Encrypt {
		t.Errorf("getProcessMode() = %d, want Encrypt", got)
	}

	*config.FlagEncrypt = false
	*config.FlagDecrypt = true
	if got := getProcessMode(); got != Decrypt {
		t.Errorf("getProcessMode() = %d, want Decrypt", got)
	}
}

func TestGetQuantumLockMode(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "a.txt")
	if err := os.WriteFile(file, []byte("data"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	if got, exists := getQuantumLockMode(file); !exists || got != yubase.File {
		t.Errorf("getQuantumLockMode(file) = %v, %v; want File, true", got, exists)
	}
	if got, exists := getQuantumLockMode(dir); !exists || got != yubase.Dir {
		t.Errorf("getQuantumLockMode(dir) = %v, %v; want Dir, true", got, exists)
	}
	if got, exists := getQuantumLockMode(filepath.Join(dir, "missing")); exists || got != yubase.String {
		t.Errorf("getQuantumLockMode(missing) = %v, %v; want String, false", got, exists)
	}
}
