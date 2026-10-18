package utils

import (
	"io"
	"os"
	"strings"
	"testing"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("close pipe: %v", err)
	}
	os.Stdout = old

	out, _ := io.ReadAll(r)
	return string(out)
}

func TestPrintEncryption(t *testing.T) {
	out := captureStdout(t, func() {
		PrintEncryption("src.txt", "dst.gpg", "alice")
	})

	want := "gpg --output dst.gpg --encrypt --recipient alice src.txt"
	if !strings.Contains(out, want) {
		t.Errorf("output = %q, want it to contain %q", out, want)
	}
}

func TestPrintDecryption(t *testing.T) {
	out := captureStdout(t, func() {
		PrintDecryption("src.txt", "dst.gpg")
	})

	want := "gpg --output src.txt --decrypt dst.gpg"
	if !strings.Contains(out, want) {
		t.Errorf("output = %q, want it to contain %q", out, want)
	}
}
