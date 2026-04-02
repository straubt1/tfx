// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

// Package hclconfig reads and writes TFx profile blocks in ~/.tfx.hcl.
//
// Profile block format:
//
//	profile "default" {
//	  hostname     = "app.terraform.io"
//	  organization = "my-org"
//	  token        = "abc123..."
//	}
//
// The block label is the profile name (a user-editable alias — not the
// hostname). hostname is an optional key inside the block; it defaults to
// DefaultHostname ("app.terraform.io") when omitted.
package hclconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

const (
	// DefaultProfileName is the profile name used when none is specified.
	DefaultProfileName = "default"
	// DefaultHostname is the TFE/HCP Terraform hostname used when none is specified.
	DefaultHostname = "app.terraform.io"
)

// Profile holds configuration for one TFx profile.
type Profile struct {
	Name         string // block label — user-editable alias
	Hostname     string // hostname value; defaults to DefaultHostname if omitted
	Organization string // organization value; may be empty
	Token        string // token value
}

var (
	reProfileStart = regexp.MustCompile(`^profile\s+"([^"]+)"\s*\{`)
	reKeyValue     = regexp.MustCompile(`^\s+(hostname|organization|token)\s*=\s*"([^"]*)"`)
	reBlockEnd     = regexp.MustCompile(`^\}`)
)

// DefaultConfigPath returns the canonical path to ~/.tfx.hcl.
func DefaultConfigPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".tfx.hcl"), nil
}

// ListProfiles parses path and returns all profile blocks in file order.
//
//   - Each profile block becomes one Profile. Name = block label;
//     Hostname = hostname inside the block, or DefaultHostname when absent.
//   - File not found: returns nil, nil.
//   - Files without profile blocks: returns nil, nil.
func ListProfiles(path string) ([]Profile, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	var profiles []Profile
	var current *Profile
	for _, line := range lines {
		if m := reProfileStart.FindStringSubmatch(line); m != nil {
			current = &Profile{Name: m[1]}
			continue
		}
		if current != nil {
			if m := reKeyValue.FindStringSubmatch(line); m != nil {
				switch m[1] {
				case "hostname":
					current.Hostname = m[2]
				case "organization":
					current.Organization = m[2]
				case "token":
					current.Token = m[2]
				}
				continue
			}
			if reBlockEnd.MatchString(line) {
				if current.Hostname == "" {
					current.Hostname = DefaultHostname
				}
				profiles = append(profiles, *current)
				current = nil
			}
		}
	}
	return profiles, nil
}

// WriteProfile adds or replaces the named profile block in the file at path.
// If path does not exist it is created. All other profiles and file content
// are preserved. The file is written with 0600 permissions because it
// contains an API token.
//
// name defaults to DefaultProfileName when empty.
// hostname defaults to DefaultHostname when empty.
// organization may be empty; a commented placeholder is written instead.
func WriteProfile(path, name, hostname, organization, token string) error {
	if name == "" {
		name = DefaultProfileName
	}
	if hostname == "" {
		hostname = DefaultHostname
	}

	existing := ""
	if data, err := os.ReadFile(path); err == nil {
		existing = string(data)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("reading config file: %w", err)
	}

	stripped := removeProfileBlock(existing, name)

	var orgLine string
	if organization != "" {
		orgLine = fmt.Sprintf("  organization = %q", organization)
	} else {
		orgLine = `  # organization = "" # set this to your organization name`
	}
	block := fmt.Sprintf(
		"\nprofile %q {\n  hostname     = %q\n%s\n  token        = %q\n}\n",
		name, hostname, orgLine, token,
	)

	content := strings.TrimRight(stripped, "\n\t ")
	if content != "" {
		content += "\n"
	}
	content += block

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}
	return os.WriteFile(path, []byte(content), 0600)
}

// removeProfileBlock strips the profile block with the given name (block label)
// from content, returning everything else.
func removeProfileBlock(content, name string) string {
	target := fmt.Sprintf("profile %q", name)
	lines := strings.Split(content, "\n")
	var out []string
	skip := false
	depth := 0
	for _, line := range lines {
		if !skip && strings.Contains(line, target) && reProfileStart.MatchString(line) {
			skip = true
			depth = 1
			continue
		}
		if skip {
			depth += strings.Count(line, "{") - strings.Count(line, "}")
			if depth <= 0 {
				skip = false
			}
			continue
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}
