// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

// writeConfig writes content to a temp .tfx.hcl and returns its path.
func writeConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".tfx.hcl")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("writeConfig: %v", err)
	}
	return path
}

// resetState resets Viper and package-level state for a clean test.
func resetState(t *testing.T) {
	t.Helper()
	viper.Reset()
	userChangedFlags = make(map[string]bool)
}

func TestResolveProfile_DefaultProfile(t *testing.T) {
	resetState(t)
	path := writeConfig(t, `
profile "staging" {
  hostname     = "staging.co"
  organization = "staging-org"
  token        = "staging-tok"
}
profile "default" {
  hostname     = "app.terraform.io"
  organization = "default-org"
  token        = "default-tok"
}
`)
	viper.SetConfigFile(path)
	viper.ReadInConfig()

	err := resolveProfile()
	if err != nil {
		t.Fatalf("resolveProfile() error = %v", err)
	}

	if got := viper.GetString("organization"); got != "default-org" {
		t.Errorf("expected organization=default-org, got %q", got)
	}
	if got := viper.GetString("profile"); got != "default" {
		t.Errorf("expected profile=default, got %q", got)
	}
}

func TestResolveProfile_NamedProfile(t *testing.T) {
	resetState(t)
	path := writeConfig(t, `
profile "default" {
  hostname     = "app.terraform.io"
  organization = "default-org"
  token        = "default-tok"
}
profile "staging" {
  hostname     = "staging.co"
  organization = "staging-org"
  token        = "staging-tok"
}
`)
	viper.SetConfigFile(path)
	viper.ReadInConfig()
	viper.Set("profile", "staging")

	err := resolveProfile()
	if err != nil {
		t.Fatalf("resolveProfile() error = %v", err)
	}

	if got := viper.GetString("organization"); got != "staging-org" {
		t.Errorf("expected organization=staging-org, got %q", got)
	}
	if got := viper.GetString("hostname"); got != "staging.co" {
		t.Errorf("expected hostname=staging.co, got %q", got)
	}
}

func TestResolveProfile_NoDefaultProfile_Skips(t *testing.T) {
	resetState(t)
	path := writeConfig(t, `
profile "prod" {
  hostname     = "prod.co"
  organization = "prod-org"
  token        = "prod-tok"
}
profile "staging" {
  hostname     = "staging.co"
  organization = "staging-org"
  token        = "staging-tok"
}
`)
	viper.SetConfigFile(path)
	viper.ReadInConfig()

	err := resolveProfile()
	if err != nil {
		t.Fatalf("expected no error (flags/env may provide creds), got: %v", err)
	}
	// No profile loaded — token should be empty
	if got := viper.GetString("token"); got != "" {
		t.Errorf("expected empty token, got %q", got)
	}
}

func TestResolveProfile_UnknownNamedProfile_Error(t *testing.T) {
	resetState(t)
	path := writeConfig(t, `
profile "default" {
  hostname     = "app.terraform.io"
  organization = "my-org"
  token        = "tok"
}
`)
	viper.SetConfigFile(path)
	viper.ReadInConfig()
	viper.Set("profile", "nonexistent")

	err := resolveProfile()
	if err == nil {
		t.Fatal("expected error for unknown profile, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' error, got: %v", err)
	}
}

func TestResolveProfile_NoProfiles_Skips(t *testing.T) {
	resetState(t)
	path := writeConfig(t, `# empty config, no profiles
`)
	viper.SetConfigFile(path)
	viper.ReadInConfig()

	err := resolveProfile()
	if err != nil {
		t.Fatalf("expected no error (flags/env may provide creds), got: %v", err)
	}
}

func TestResolveProfile_CLIFlagsOverrideProfile(t *testing.T) {
	resetState(t)
	path := writeConfig(t, `
profile "default" {
  hostname     = "profile.co"
  organization = "profile-org"
  token        = "profile-tok"
}
`)
	viper.SetConfigFile(path)
	viper.ReadInConfig()

	// Simulate user typed --token on CLI
	userChangedFlags["token"] = true
	viper.Set("token", "flag-tok")

	err := resolveProfile()
	if err != nil {
		t.Fatalf("resolveProfile() error = %v", err)
	}

	// Token should come from flag, not profile
	if got := viper.GetString("token"); got != "flag-tok" {
		t.Errorf("expected token=flag-tok, got %q", got)
	}
	// Organization should come from profile (no flag override)
	if got := viper.GetString("organization"); got != "profile-org" {
		t.Errorf("expected organization=profile-org, got %q", got)
	}
}

func TestResolveProfile_EnvVarsOverrideProfile(t *testing.T) {
	resetState(t)
	path := writeConfig(t, `
profile "default" {
  hostname     = "profile.co"
  organization = "profile-org"
  token        = "profile-tok"
}
`)
	viper.SetConfigFile(path)
	viper.ReadInConfig()

	viper.BindEnv("hostname", "TFE_HOSTNAME")
	viper.BindEnv("organization", "TFE_ORGANIZATION")
	viper.BindEnv("token", "TFE_TOKEN")

	t.Setenv("TFE_TOKEN", "env-tok")

	err := resolveProfile()
	if err != nil {
		t.Fatalf("resolveProfile() error = %v", err)
	}

	// Token should come from env, not profile
	if got := viper.GetString("token"); got != "env-tok" {
		t.Errorf("expected token=env-tok, got %q", got)
	}
	// Organization should come from profile (no env override)
	if got := viper.GetString("organization"); got != "profile-org" {
		t.Errorf("expected organization=profile-org, got %q", got)
	}
}

func TestResolveProfile_CLIFlagsOverrideEnvVars(t *testing.T) {
	resetState(t)

	viper.BindEnv("token", "TFE_TOKEN")
	t.Setenv("TFE_TOKEN", "env-tok")

	userChangedFlags["token"] = true
	viper.Set("token", "flag-tok")

	// No config file — resolveProfile returns nil (no config to load)
	err := resolveProfile()
	if err != nil {
		t.Fatalf("resolveProfile() error = %v", err)
	}

	if got := viper.GetString("token"); got != "flag-tok" {
		t.Errorf("expected token=flag-tok, got %q", got)
	}
}

func TestResolveProfile_NoConfigFile(t *testing.T) {
	resetState(t)
	// No config file set — should return nil (flags/env handle auth)

	err := resolveProfile()
	if err != nil {
		t.Fatalf("resolveProfile() error = %v", err)
	}
}

func TestResolveProfile_ProfileWithFlagOverride(t *testing.T) {
	resetState(t)
	path := writeConfig(t, `
profile "default" {
  hostname     = "default.co"
  organization = "default-org"
  token        = "default-tok"
}
`)
	viper.SetConfigFile(path)
	viper.ReadInConfig()
	// User typed --organization to override profile value
	userChangedFlags["organization"] = true
	viper.Set("organization", "override-org")

	err := resolveProfile()
	if err != nil {
		t.Fatalf("resolveProfile() error = %v", err)
	}

	// Organization from flag override
	if got := viper.GetString("organization"); got != "override-org" {
		t.Errorf("expected organization=override-org, got %q", got)
	}
	// Token from profile (no flag override)
	if got := viper.GetString("token"); got != "default-tok" {
		t.Errorf("expected token=default-tok, got %q", got)
	}
}
