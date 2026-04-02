// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package hclconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteProfile_UsesOrganizationKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".tfx.hcl")

	err := WriteProfile(path, "default", "app.terraform.io", "my-org", "tok123")
	if err != nil {
		t.Fatalf("WriteProfile() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	content := string(data)

	if !strings.Contains(content, `organization = "my-org"`) {
		t.Errorf("expected 'organization = \"my-org\"' in output, got:\n%s", content)
	}
}

func TestWriteProfile_EmptyOrganization(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".tfx.hcl")

	err := WriteProfile(path, "default", "app.terraform.io", "", "tok123")
	if err != nil {
		t.Fatalf("WriteProfile() error = %v", err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)

	if !strings.Contains(content, `# organization = ""`) {
		t.Errorf("expected commented organization placeholder, got:\n%s", content)
	}
}

func TestWriteProfile_DefaultName(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".tfx.hcl")

	err := WriteProfile(path, "", "app.terraform.io", "my-org", "tok123")
	if err != nil {
		t.Fatalf("WriteProfile() error = %v", err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)

	if !strings.Contains(content, `profile "default"`) {
		t.Errorf("expected profile name 'default', got:\n%s", content)
	}
}

func TestWriteProfile_DefaultHostname(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".tfx.hcl")

	err := WriteProfile(path, "default", "", "my-org", "tok123")
	if err != nil {
		t.Fatalf("WriteProfile() error = %v", err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)

	if !strings.Contains(content, `hostname     = "app.terraform.io"`) {
		t.Errorf("expected default hostname, got:\n%s", content)
	}
}

func TestListProfiles_SingleProfile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".tfx.hcl")

	config := `profile "default" {
  hostname     = "app.terraform.io"
  organization = "my-org"
  token        = "tok123"
}
`
	os.WriteFile(path, []byte(config), 0600)

	profiles, err := ListProfiles(path)
	if err != nil {
		t.Fatalf("ListProfiles() error = %v", err)
	}
	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(profiles))
	}
	p := profiles[0]
	if p.Name != "default" {
		t.Errorf("expected name 'default', got %q", p.Name)
	}
	if p.Hostname != "app.terraform.io" {
		t.Errorf("expected hostname 'app.terraform.io', got %q", p.Hostname)
	}
	if p.Organization != "my-org" {
		t.Errorf("expected organization 'my-org', got %q", p.Organization)
	}
	if p.Token != "tok123" {
		t.Errorf("expected token 'tok123', got %q", p.Token)
	}
}

func TestListProfiles_MultipleProfiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".tfx.hcl")

	config := `profile "default" {
  hostname     = "app.terraform.io"
  organization = "my-org"
  token        = "tok123"
}

profile "staging" {
  hostname     = "tfe.staging.co"
  organization = "staging-org"
  token        = "tok456"
}
`
	os.WriteFile(path, []byte(config), 0600)

	profiles, err := ListProfiles(path)
	if err != nil {
		t.Fatalf("ListProfiles() error = %v", err)
	}
	if len(profiles) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(profiles))
	}

	if profiles[0].Name != "default" || profiles[0].Organization != "my-org" {
		t.Errorf("profile 0 = %+v", profiles[0])
	}
	if profiles[1].Name != "staging" || profiles[1].Organization != "staging-org" {
		t.Errorf("profile 1 = %+v", profiles[1])
	}
}

func TestListProfiles_DefaultHostname(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".tfx.hcl")

	config := `profile "default" {
  organization = "my-org"
  token        = "tok123"
}
`
	os.WriteFile(path, []byte(config), 0600)

	profiles, err := ListProfiles(path)
	if err != nil {
		t.Fatalf("ListProfiles() error = %v", err)
	}
	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(profiles))
	}
	if profiles[0].Hostname != DefaultHostname {
		t.Errorf("expected default hostname %q, got %q", DefaultHostname, profiles[0].Hostname)
	}
}

func TestListProfiles_FileNotFound(t *testing.T) {
	profiles, err := ListProfiles("/nonexistent/path/.tfx.hcl")
	if err != nil {
		t.Fatalf("ListProfiles() error = %v", err)
	}
	if profiles != nil {
		t.Errorf("expected nil profiles for missing file, got %v", profiles)
	}
}

func TestListProfiles_NoProfileBlocks(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".tfx.hcl")

	// File exists but has no profile blocks
	os.WriteFile(path, []byte("# just a comment\n"), 0600)

	profiles, err := ListProfiles(path)
	if err != nil {
		t.Fatalf("ListProfiles() error = %v", err)
	}
	if profiles != nil {
		t.Errorf("expected nil profiles for file without blocks, got %v", profiles)
	}
}

func TestWriteProfile_PreservesOtherProfiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".tfx.hcl")

	// Write first profile
	err := WriteProfile(path, "default", "app.terraform.io", "org1", "tok1")
	if err != nil {
		t.Fatalf("WriteProfile(default) error = %v", err)
	}

	// Write second profile
	err = WriteProfile(path, "staging", "staging.co", "org2", "tok2")
	if err != nil {
		t.Fatalf("WriteProfile(staging) error = %v", err)
	}

	// Both should exist
	profiles, err := ListProfiles(path)
	if err != nil {
		t.Fatalf("ListProfiles() error = %v", err)
	}
	if len(profiles) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(profiles))
	}
}

func TestWriteProfile_UpdatesExistingProfile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".tfx.hcl")

	// Write initial
	WriteProfile(path, "default", "app.terraform.io", "org1", "tok1")

	// Update same profile
	WriteProfile(path, "default", "new.host", "org2", "tok2")

	profiles, _ := ListProfiles(path)
	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile after update, got %d", len(profiles))
	}
	if profiles[0].Hostname != "new.host" || profiles[0].Token != "tok2" {
		t.Errorf("profile not updated: %+v", profiles[0])
	}
}
