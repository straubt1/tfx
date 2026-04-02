// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package hclconfig

import (
	"os"
	"path/filepath"
	"testing"
)

// helper writes content to a temp file and returns its path.
func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), ".tfx.hcl")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestListProfiles_FileNotFound(t *testing.T) {
	profiles, err := ListProfiles("/tmp/does-not-exist-hclconfig-test.hcl")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if profiles != nil {
		t.Fatalf("expected nil profiles, got %v", profiles)
	}
}

func TestListProfiles_EmptyFile(t *testing.T) {
	path := writeTempConfig(t, "")
	profiles, err := ListProfiles(path)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if profiles != nil {
		t.Fatalf("expected nil profiles, got %v", profiles)
	}
}

// --- Flat format (legacy) ---

func TestListProfiles_FlatOldKeys(t *testing.T) {
	path := writeTempConfig(t, `
tfeHostname     = "tfe.example.com"
tfeOrganization = "my-org"
tfeToken        = "tok123"
`)
	profiles, err := ListProfiles(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(profiles))
	}
	p := profiles[0]
	assertProfile(t, p, "default", "tfe.example.com", "my-org", "tok123")
}

func TestListProfiles_FlatNewKeys(t *testing.T) {
	path := writeTempConfig(t, `
hostname            = "tfe.example.com"
defaultOrganization = "my-org"
token               = "tok456"
`)
	profiles, err := ListProfiles(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(profiles))
	}
	p := profiles[0]
	assertProfile(t, p, "default", "tfe.example.com", "my-org", "tok456")
}

func TestListProfiles_FlatNewKeysShortOrg(t *testing.T) {
	path := writeTempConfig(t, `
hostname     = "tfe.example.com"
organization = "short-org"
token        = "tok789"
`)
	profiles, err := ListProfiles(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(profiles))
	}
	p := profiles[0]
	assertProfile(t, p, "default", "tfe.example.com", "short-org", "tok789")
}

func TestListProfiles_FlatTokenOnly(t *testing.T) {
	path := writeTempConfig(t, `token = "tok-only"`)
	profiles, err := ListProfiles(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(profiles))
	}
	p := profiles[0]
	assertProfile(t, p, "default", DefaultHostname, "", "tok-only")
}

func TestListProfiles_FlatNoToken(t *testing.T) {
	path := writeTempConfig(t, `hostname = "tfe.example.com"`)
	profiles, err := ListProfiles(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if profiles != nil {
		t.Fatalf("expected nil profiles when no token, got %v", profiles)
	}
}

// --- Profile block format ---

func TestListProfiles_BlockNewKeys(t *testing.T) {
	path := writeTempConfig(t, `
profile "default" {
  hostname            = "tfe.example.com"
  defaultOrganization = "block-org"
  token               = "block-tok"
}
`)
	profiles, err := ListProfiles(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(profiles))
	}
	assertProfile(t, profiles[0], "default", "tfe.example.com", "block-org", "block-tok")
}

func TestListProfiles_BlockOldKeys(t *testing.T) {
	path := writeTempConfig(t, `
profile "legacy" {
  tfeHostname     = "tfe.old.com"
  tfeOrganization = "old-org"
  tfeToken        = "old-tok"
}
`)
	profiles, err := ListProfiles(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(profiles))
	}
	assertProfile(t, profiles[0], "legacy", "tfe.old.com", "old-org", "old-tok")
}

func TestListProfiles_BlockHostnameDefaultsWhenOmitted(t *testing.T) {
	path := writeTempConfig(t, `
profile "nope" {
  token = "just-tok"
}
`)
	profiles, err := ListProfiles(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(profiles))
	}
	assertProfile(t, profiles[0], "nope", DefaultHostname, "", "just-tok")
}

func TestListProfiles_BlockLabelWithDotUsedAsHostname(t *testing.T) {
	path := writeTempConfig(t, `
profile "tfe.myco.internal" {
  token = "dot-tok"
}
`)
	profiles, err := ListProfiles(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(profiles))
	}
	assertProfile(t, profiles[0], "tfe.myco.internal", "tfe.myco.internal", "", "dot-tok")
}

func TestListProfiles_MultipleBlocks(t *testing.T) {
	path := writeTempConfig(t, `
profile "default" {
  hostname = "app.terraform.io"
  token    = "tok-1"
}

profile "staging" {
  hostname            = "tfe.staging.com"
  defaultOrganization = "staging-org"
  token               = "tok-2"
}
`)
	profiles, err := ListProfiles(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(profiles))
	}
	assertProfile(t, profiles[0], "default", "app.terraform.io", "", "tok-1")
	assertProfile(t, profiles[1], "staging", "tfe.staging.com", "staging-org", "tok-2")
}

// --- WriteProfile ---

func TestWriteProfile_CreatesNewFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", ".tfx.hcl")
	err := WriteProfile(path, "default", "app.terraform.io", "my-org", "write-tok")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	profiles, err := ListProfiles(path)
	if err != nil {
		t.Fatalf("unexpected error reading back: %v", err)
	}
	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(profiles))
	}
	assertProfile(t, profiles[0], "default", "app.terraform.io", "my-org", "write-tok")
}

func TestWriteProfile_DefaultsNameAndHostname(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".tfx.hcl")
	err := WriteProfile(path, "", "", "org", "tok")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	profiles, err := ListProfiles(path)
	if err != nil {
		t.Fatalf("unexpected error reading back: %v", err)
	}
	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(profiles))
	}
	assertProfile(t, profiles[0], DefaultProfileName, DefaultHostname, "org", "tok")
}

func TestWriteProfile_EmptyOrgWritesComment(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".tfx.hcl")
	err := WriteProfile(path, "test", "host.com", "", "tok")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("unexpected error reading file: %v", err)
	}
	content := string(data)
	if !contains(content, "# defaultOrganization") {
		t.Errorf("expected commented org placeholder, got:\n%s", content)
	}
}

func TestWriteProfile_ReplacesExisting(t *testing.T) {
	path := writeTempConfig(t, `
profile "default" {
  hostname = "old.com"
  token    = "old-tok"
}
`)
	err := WriteProfile(path, "default", "new.com", "new-org", "new-tok")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	profiles, err := ListProfiles(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(profiles))
	}
	assertProfile(t, profiles[0], "default", "new.com", "new-org", "new-tok")
}

func assertProfile(t *testing.T, p Profile, name, hostname, org, token string) {
	t.Helper()
	if p.Name != name {
		t.Errorf("Name: got %q, want %q", p.Name, name)
	}
	if p.Hostname != hostname {
		t.Errorf("Hostname: got %q, want %q", p.Hostname, hostname)
	}
	if p.Organization != org {
		t.Errorf("Organization: got %q, want %q", p.Organization, org)
	}
	if p.Token != token {
		t.Errorf("Token: got %q, want %q", p.Token, token)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
