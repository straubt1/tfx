// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParsePlatform(t *testing.T) {
	osName, arch, err := ParsePlatform("linux_amd64")
	if err != nil {
		t.Fatal(err)
	}
	if osName != "linux" || arch != "amd64" {
		t.Fatalf("got %s %s", osName, arch)
	}

	if _, _, err := ParsePlatform("linux"); err == nil {
		t.Fatal("expected error for missing arch")
	}
}

func TestDownloadPublicProvider(t *testing.T) {
	zip := []byte("fake-azurerm-linux-amd64-zip")
	sum := sha256.Sum256(zip)
	shasum := hex.EncodeToString(sum[:])
	filename := "terraform-provider-azurerm_5.0.0_linux_amd64.zip"
	shasumsBody := shasum + "  " + filename + "\n"
	armor := "-----BEGIN PGP PUBLIC KEY BLOCK-----\n\nhashicorp-key\n-----END PGP PUBLIC KEY BLOCK-----"
	armorJSON, err := json.Marshal(armor)
	if err != nil {
		t.Fatal(err)
	}

	var server *httptest.Server
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/terraform.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"providers.v1":"/v1/providers/"}`))
	})
	mux.HandleFunc("/v1/providers/hashicorp/azurerm/versions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"versions":[{"version":"5.0.0","platforms":[{"os":"linux","arch":"amd64"},{"os":"darwin","arch":"arm64"}]}]}`))
	})
	mux.HandleFunc("/v1/providers/hashicorp/azurerm/5.0.0/download/linux/amd64", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{
			"os":"linux","arch":"amd64","filename":"%s",
			"download_url":"%s/files/%s",
			"shasums_url":"%s/files/terraform-provider-azurerm_5.0.0_SHA256SUMS",
			"shasums_signature_url":"%s/files/terraform-provider-azurerm_5.0.0_SHA256SUMS.72D7468F.sig",
			"shasum":"%s",
			"signing_keys":{"gpg_public_keys":[{"key_id":"34365D9472D7468F","ascii_armor":%s}]}
		}`, filename, server.URL, filename, server.URL, server.URL, shasum, armorJSON)
	})
	mux.HandleFunc("/files/"+filename, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(zip)
	})
	mux.HandleFunc("/files/terraform-provider-azurerm_5.0.0_SHA256SUMS", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(shasumsBody))
	})
	mux.HandleFunc("/files/terraform-provider-azurerm_5.0.0_SHA256SUMS.72D7468F.sig", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("fake-sig"))
	})
	server = httptest.NewServer(mux)
	defer server.Close()

	dir := t.TempDir()
	result, err := DownloadPublicProvider(PublicProviderDownloadConfig{
		RegistryBaseURL: server.URL,
		Namespace:       "hashicorp",
		Name:            "azurerm",
		Version:         "5.0.0",
		Directory:       dir,
		Platforms:       []string{"linux_amd64", "windows_amd64"},
	})
	if err != nil {
		t.Fatal(err)
	}

	wantDir := filepath.Join(dir, "hashicorp", "azurerm", "5.0.0")
	if result.Directory != wantDir {
		t.Fatalf("directory = %s, want %s", result.Directory, wantDir)
	}
	if result.KeyID != "34365D9472D7468F" {
		t.Fatalf("key id = %s", result.KeyID)
	}
	if len(result.Platforms) != 1 || result.Platforms[0].OS != "linux" {
		t.Fatalf("platforms = %+v", result.Platforms)
	}
	if len(result.SkippedPlatforms) != 1 || result.SkippedPlatforms[0] != "windows_amd64" {
		t.Fatalf("skipped = %+v", result.SkippedPlatforms)
	}

	zipPath := filepath.Join(wantDir, filename)
	if _, err := os.Stat(zipPath); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(result.ShasumsPath); err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(result.ShasumsSigPath, "SHA256SUMS.72D7468F.sig") {
		t.Fatalf("sig path = %s", result.ShasumsSigPath)
	}

	wantKey := filepath.Join(wantDir, "34365D9472D7468F.asc")
	if result.GPGPublicKeyPath != wantKey {
		t.Fatalf("gpg public key path = %s, want %s", result.GPGPublicKeyPath, wantKey)
	}
	gotArmor, err := os.ReadFile(wantKey)
	if err != nil {
		t.Fatal(err)
	}
	if string(gotArmor) != armor+"\n" {
		t.Fatalf("gpg public key contents = %q", gotArmor)
	}

	if _, err := os.Stat(filepath.Join(wantDir, "tfx-provider.json")); !os.IsNotExist(err) {
		t.Fatalf("tfx-provider.json should not be written, stat err = %v", err)
	}
}

func TestDownloadPublicProviderWritesPartnerGPGKey(t *testing.T) {
	zip := []byte("fake-cosign-darwin-arm64-zip")
	sum := sha256.Sum256(zip)
	shasum := hex.EncodeToString(sum[:])
	filename := "terraform-provider-cosign_0.4.16_darwin_arm64.zip"
	shasumsBody := shasum + "  " + filename + "\n"
	armor := "-----BEGIN PGP PUBLIC KEY BLOCK-----\n\nchainguard-key\n-----END PGP PUBLIC KEY BLOCK-----"
	armorJSON, err := json.Marshal(armor)
	if err != nil {
		t.Fatal(err)
	}

	var server *httptest.Server
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/terraform.json", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"providers.v1":"/v1/providers/"}`))
	})
	mux.HandleFunc("/v1/providers/chainguard-dev/cosign/versions", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"versions":[{"version":"0.4.16","platforms":[{"os":"darwin","arch":"arm64"}]}]}`))
	})
	mux.HandleFunc("/v1/providers/chainguard-dev/cosign/0.4.16/download/darwin/arm64", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
			"os":"darwin","arch":"arm64","filename":"%s",
			"download_url":"%s/files/%s",
			"shasums_url":"%s/files/terraform-provider-cosign_0.4.16_SHA256SUMS",
			"shasums_signature_url":"%s/files/terraform-provider-cosign_0.4.16_SHA256SUMS.sig",
			"shasum":"%s",
			"signing_keys":{"gpg_public_keys":[{"key_id":"5bbee08f6bf07616","ascii_armor":%s}]}
		}`, filename, server.URL, filename, server.URL, server.URL, shasum, armorJSON)
	})
	mux.HandleFunc("/files/"+filename, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(zip)
	})
	mux.HandleFunc("/files/terraform-provider-cosign_0.4.16_SHA256SUMS", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(shasumsBody))
	})
	mux.HandleFunc("/files/terraform-provider-cosign_0.4.16_SHA256SUMS.sig", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("fake-sig"))
	})
	server = httptest.NewServer(mux)
	defer server.Close()

	dir := t.TempDir()
	result, err := DownloadPublicProvider(PublicProviderDownloadConfig{
		RegistryBaseURL: server.URL,
		Namespace:       "chainguard-dev",
		Name:            "cosign",
		Version:         "0.4.16",
		Directory:       dir,
		Platforms:       []string{"darwin_arm64"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.KeyID != "5BBEE08F6BF07616" {
		t.Fatalf("key id = %s", result.KeyID)
	}
	wantKey := filepath.Join(dir, "chainguard-dev", "cosign", "0.4.16", "5BBEE08F6BF07616.asc")
	if result.GPGPublicKeyPath != wantKey {
		t.Fatalf("gpg public key path = %s, want %s", result.GPGPublicKeyPath, wantKey)
	}
	gotArmor, err := os.ReadFile(wantKey)
	if err != nil {
		t.Fatal(err)
	}
	if string(gotArmor) != armor+"\n" {
		t.Fatalf("gpg public key contents = %q", gotArmor)
	}

	if _, err := os.Stat(filepath.Join(dir, "chainguard-dev", "cosign", "0.4.16", "tfx-provider.json")); !os.IsNotExist(err) {
		t.Fatalf("tfx-provider.json should not be written, stat err = %v", err)
	}
}

func TestDownloadPublicProviderChecksumMismatch(t *testing.T) {
	zip := []byte("zip-bytes")
	filename := "terraform-provider-azurerm_5.0.0_linux_amd64.zip"
	wrongSum := strings.Repeat("0", 64)

	var server *httptest.Server
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/terraform.json", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"providers.v1":"/v1/providers/"}`))
	})
	mux.HandleFunc("/v1/providers/hashicorp/azurerm/versions", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"versions":[{"version":"5.0.0","platforms":[{"os":"linux","arch":"amd64"}]}]}`))
	})
	mux.HandleFunc("/v1/providers/hashicorp/azurerm/5.0.0/download/linux/amd64", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
			"os":"linux","arch":"amd64","filename":"%s",
			"download_url":"%s/files/%s",
			"shasums_url":"%s/files/SHA256SUMS",
			"shasums_signature_url":"%s/files/SHA256SUMS.sig",
			"shasum":"%s",
			"signing_keys":{"gpg_public_keys":[{"key_id":"ABC"}]}
		}`, filename, server.URL, filename, server.URL, server.URL, wrongSum)
	})
	mux.HandleFunc("/files/"+filename, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(zip)
	})
	mux.HandleFunc("/files/SHA256SUMS", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "%s  %s\n", wrongSum, filename)
	})
	mux.HandleFunc("/files/SHA256SUMS.sig", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("sig"))
	})
	server = httptest.NewServer(mux)
	defer server.Close()

	_, err := DownloadPublicProvider(PublicProviderDownloadConfig{
		RegistryBaseURL: server.URL,
		Namespace:       "hashicorp",
		Name:            "azurerm",
		Version:         "5.0.0",
		Directory:       t.TempDir(),
		Platforms:       []string{"linux_amd64"},
	})
	if err == nil || !strings.Contains(err.Error(), "checksum mismatch") {
		t.Fatalf("expected checksum mismatch, got %v", err)
	}
}

func TestDownloadPublicProviderZipMissingFromSHA256SUMS(t *testing.T) {
	zip := []byte("zip-bytes")
	sum := sha256.Sum256(zip)
	shasum := hex.EncodeToString(sum[:])
	filename := "terraform-provider-azurerm_5.0.0_linux_amd64.zip"

	var server *httptest.Server
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/terraform.json", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"providers.v1":"/v1/providers/"}`))
	})
	mux.HandleFunc("/v1/providers/hashicorp/azurerm/versions", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"versions":[{"version":"5.0.0","platforms":[{"os":"linux","arch":"amd64"}]}]}`))
	})
	mux.HandleFunc("/v1/providers/hashicorp/azurerm/5.0.0/download/linux/amd64", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
			"os":"linux","arch":"amd64","filename":"%s",
			"download_url":"%s/files/%s",
			"shasums_url":"%s/files/SHA256SUMS",
			"shasums_signature_url":"%s/files/SHA256SUMS.sig",
			"shasum":"%s",
			"signing_keys":{"gpg_public_keys":[{"key_id":"ABC"}]}
		}`, filename, server.URL, filename, server.URL, server.URL, shasum)
	})
	mux.HandleFunc("/files/"+filename, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(zip)
	})
	mux.HandleFunc("/files/SHA256SUMS", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "%s  other-file.zip\n", shasum)
	})
	mux.HandleFunc("/files/SHA256SUMS.sig", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("sig"))
	})
	server = httptest.NewServer(mux)
	defer server.Close()

	_, err := DownloadPublicProvider(PublicProviderDownloadConfig{
		RegistryBaseURL: server.URL,
		Namespace:       "hashicorp",
		Name:            "azurerm",
		Version:         "5.0.0",
		Directory:       t.TempDir(),
		Platforms:       []string{"linux_amd64"},
	})
	if err == nil || !strings.Contains(err.Error(), "not listed in SHA256SUMS") {
		t.Fatalf("expected missing SHA256SUMS entry, got %v", err)
	}
}

func TestSafeJoin(t *testing.T) {
	dir := t.TempDir()

	got, err := safeJoin(dir, "terraform-provider-azurerm_5.0.0_linux_amd64.zip")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(dir, "terraform-provider-azurerm_5.0.0_linux_amd64.zip")
	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	rejects := []string{"", ".", "..", "/etc/passwd", "../evil.zip", "foo/bar.zip"}
	for _, name := range rejects {
		if _, err := safeJoin(dir, name); err == nil {
			t.Errorf("expected error for %q", name)
		}
	}
}

func TestDownloadPublicProviderRejectsPathTraversalFilename(t *testing.T) {
	zip := []byte("zip-bytes")
	sum := sha256.Sum256(zip)
	shasum := hex.EncodeToString(sum[:])
	filename := "../evil.zip"

	var server *httptest.Server
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/terraform.json", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"providers.v1":"/v1/providers/"}`))
	})
	mux.HandleFunc("/v1/providers/hashicorp/azurerm/versions", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"versions":[{"version":"5.0.0","platforms":[{"os":"linux","arch":"amd64"}]}]}`))
	})
	mux.HandleFunc("/v1/providers/hashicorp/azurerm/5.0.0/download/linux/amd64", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
			"os":"linux","arch":"amd64","filename":"%s",
			"download_url":"%s/files/ok.zip",
			"shasums_url":"%s/files/SHA256SUMS",
			"shasums_signature_url":"%s/files/SHA256SUMS.sig",
			"shasum":"%s",
			"signing_keys":{"gpg_public_keys":[{"key_id":"ABC"}]}
		}`, filename, server.URL, server.URL, server.URL, shasum)
	})
	mux.HandleFunc("/files/ok.zip", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(zip)
	})
	mux.HandleFunc("/files/SHA256SUMS", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "%s  ok.zip\n", shasum)
	})
	mux.HandleFunc("/files/SHA256SUMS.sig", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("sig"))
	})
	server = httptest.NewServer(mux)
	defer server.Close()

	_, err := DownloadPublicProvider(PublicProviderDownloadConfig{
		RegistryBaseURL: server.URL,
		Namespace:       "hashicorp",
		Name:            "azurerm",
		Version:         "5.0.0",
		Directory:       t.TempDir(),
		Platforms:       []string{"linux_amd64"},
	})
	if err == nil || !strings.Contains(err.Error(), "invalid") {
		t.Fatalf("expected invalid filename error, got %v", err)
	}
}
