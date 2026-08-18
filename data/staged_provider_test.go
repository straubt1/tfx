// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func writeStagedFiles(t *testing.T, dir, name, version, sigSuffix string, zips map[string][]byte) {
	t.Helper()
	var sums string
	for filename, body := range zips {
		path := filepath.Join(dir, filename)
		if err := os.WriteFile(path, body, 0644); err != nil {
			t.Fatal(err)
		}
		sum := sha256.Sum256(body)
		sums += hex.EncodeToString(sum[:]) + "  " + filename + "\n"
	}
	shasums := "terraform-provider-" + name + "_" + version + "_SHA256SUMS"
	if err := os.WriteFile(filepath.Join(dir, shasums), []byte(sums), 0644); err != nil {
		t.Fatal(err)
	}
	sigName := shasums + ".sig"
	if sigSuffix != "" {
		sigName = shasums + "." + sigSuffix + ".sig"
	}
	if err := os.WriteFile(filepath.Join(dir, sigName), []byte("sig"), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestReadStagedProviderDirectory(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "azurerm", "5.0.0")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	linux := []byte("linux-zip")
	darwin := []byte("darwin-zip")
	writeStagedFiles(t, dir, "azurerm", "5.0.0", "72D7468F", map[string][]byte{
		"terraform-provider-azurerm_5.0.0_linux_amd64.zip":  linux,
		"terraform-provider-azurerm_5.0.0_darwin_arm64.zip": darwin,
	})

	got, err := ReadStagedProviderDirectory(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "azurerm" || got.Version != "5.0.0" {
		t.Fatalf("name/version = %s %s", got.Name, got.Version)
	}
	if got.KeyID != "34365D9472D7468F" {
		t.Fatalf("key id = %s", got.KeyID)
	}
	if len(got.Platforms) != 2 {
		t.Fatalf("platforms = %d", len(got.Platforms))
	}
}

func TestReadStagedProviderDirectoryMissingSHA256SUMS(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "terraform-provider-azurerm_5.0.0_linux_amd64.zip"), []byte("z"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := ReadStagedProviderDirectory(dir)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestReadStagedProviderDirectoryZipNotInSHA256SUMS(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "azurerm", "5.0.0")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	writeStagedFiles(t, dir, "azurerm", "5.0.0", "72D7468F", map[string][]byte{
		"terraform-provider-azurerm_5.0.0_linux_amd64.zip": []byte("linux"),
	})
	extra := "terraform-provider-azurerm_5.0.0_windows_amd64.zip"
	if err := os.WriteFile(filepath.Join(dir, extra), []byte("windows"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := ReadStagedProviderDirectory(dir)
	if err == nil {
		t.Fatal("expected error for zip not in SHA256SUMS")
	}
}

func TestReadStagedProviderDirectoryPathMismatch(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "aws", "4.0.0")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	writeStagedFiles(t, dir, "azurerm", "5.0.0", "72D7468F", map[string][]byte{
		"terraform-provider-azurerm_5.0.0_linux_amd64.zip": []byte("linux"),
	})
	_, err := ReadStagedProviderDirectory(dir)
	if err == nil {
		t.Fatal("expected path mismatch error")
	}
}

func TestInferKeyIDFromSigFilename(t *testing.T) {
	id, err := inferKeyIDFromSigFilename("terraform-provider-azurerm_5.0.0_SHA256SUMS.72D7468F.sig")
	if err != nil {
		t.Fatal(err)
	}
	if id != "34365D9472D7468F" {
		t.Fatalf("got %s", id)
	}

	id, err = inferKeyIDFromSigFilename("terraform-provider-random_3.1.0_SHA256SUMS.348FFC4C.sig")
	if err != nil {
		t.Fatal(err)
	}
	if id != "51852D87348FFC4C" {
		t.Fatalf("got %s", id)
	}

	id, err = inferKeyIDFromSigFilename("terraform-provider-foo_1.0.0_SHA256SUMS.ABCDEF0123456789.sig")
	if err != nil {
		t.Fatal(err)
	}
	if id != "ABCDEF0123456789" {
		t.Fatalf("got %s", id)
	}

	id, err = inferKeyIDFromSigFilename("terraform-provider-foo_1.0.0_SHA256SUMS.sig")
	if err != nil {
		t.Fatal(err)
	}
	if id != "" {
		t.Fatalf("expected empty key id, got %s", id)
	}

	if _, err := inferKeyIDFromSigFilename("terraform-provider-foo_1.0.0_SHA256SUMS.DEADBEEF.sig"); err == nil {
		t.Fatal("expected error for unknown short id")
	}
}
