package fileprocessor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"git.rickiekarp.net/rickie/fileguardian/config"
	"git.rickiekarp.net/rickie/yubase"
)

func withAPIHost(t *testing.T, server *httptest.Server) {
	t.Helper()
	oldProtocol, oldHost := config.ApiProtocol, config.ApiHost
	config.ApiProtocol = "http"
	config.ApiHost = strings.TrimPrefix(server.URL, "http://")
	t.Cleanup(func() {
		config.ApiProtocol = oldProtocol
		config.ApiHost = oldHost
	})
}

func TestSendRequestSuccess(t *testing.T) {
	var gotRequest yubase.FileGuardianEntry

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/storage/v1/fileguardian/fetch" {
			t.Errorf("path = %q, want /storage/v1/fileguardian/fetch", r.URL.Path)
		}
		_ = json.NewDecoder(r.Body).Decode(&gotRequest)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(yubase.FileGuardianEntry{
			Source: "src.txt",
			Target: "src.txt.fgd",
		})
	}))
	defer server.Close()
	withAPIHost(t, server)

	res, err := sendRequest("src.txt", yubase.File, "default")
	if err != nil {
		t.Fatalf("sendRequest: %v", err)
	}
	if res == nil {
		t.Fatal("sendRequest returned nil result")
	}
	if res.Source != "src.txt" || res.Target != "src.txt.fgd" {
		t.Errorf("result = %+v, want {src.txt src.txt.fgd}", res)
	}

	if gotRequest.Source != "src.txt" {
		t.Errorf("request Source = %q, want %q", gotRequest.Source, "src.txt")
	}
	if gotRequest.Type != string(yubase.File) {
		t.Errorf("request Type = %q, want %q", gotRequest.Type, string(yubase.File))
	}
	if gotRequest.Context != "default" {
		t.Errorf("request Context = %q, want %q", gotRequest.Context, "default")
	}
}

func TestSendRequestCheckQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "check=true" {
			t.Errorf("query = %q, want %q", r.URL.RawQuery, "check=true")
		}
		_ = json.NewEncoder(w).Encode(yubase.FileGuardianEntry{Source: "a", Target: "b"})
	}))
	defer server.Close()
	withAPIHost(t, server)

	oldCheck := *config.FlagCheck
	*config.FlagCheck = true
	t.Cleanup(func() { *config.FlagCheck = oldCheck })

	if _, err := sendRequest("src.txt", yubase.File, "default"); err != nil {
		t.Fatalf("sendRequest: %v", err)
	}
}

func TestSendRequestEmptyBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	withAPIHost(t, server)

	res, err := sendRequest("src.txt", yubase.File, "default")
	if err != nil {
		t.Fatalf("sendRequest: %v", err)
	}
	if res != nil {
		t.Errorf("result = %+v, want nil", res)
	}
}

func TestSendRequestInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("not-json"))
	}))
	defer server.Close()
	withAPIHost(t, server)

	if _, err := sendRequest("src.txt", yubase.File, "default"); err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestSendRequestMissingSourceTarget(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(yubase.FileGuardianEntry{Type: string(yubase.File)})
	}))
	defer server.Close()
	withAPIHost(t, server)

	if _, err := sendRequest("src.txt", yubase.File, "default"); err == nil {
		t.Fatal("expected error for missing source/target, got nil")
	}
}
