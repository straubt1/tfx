// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

// Package hclconfig reads and writes TFx profile blocks in ~/.tfx.hcl.
//
// New profile block format:
//
//	profile "default" {
//	  hostname            = "app.terraform.io"
//	  defaultOrganization = "my-org"
//	  token               = "abc123..."
//	}
//
// The block label is the profile name (a user-editable alias — not the
// hostname). hostname is an optional key inside the block; it defaults to
// DefaultHostname ("app.terraform.io") when omitted.
//
// Legacy flat format (no profile blocks) is still read; it is treated as a
// single profile named DefaultProfileName with hostname DefaultHostname.
// The old tfe-prefixed keys (tfeHostname, tfeOrganization, tfeToken) are
// still accepted inside profile blocks for backward compatibility.
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
	Organization string // defaultOrganization value; may be empty
	Token        string // token value
}

var (
	reProfileStart = regexp.MustCompile(`^profile\s+"([^"]+)"\s*\{`)
	reKeyValue     = regexp.MustCompile(`^\s+(tfeHostname|tfeOrganization|tfeToken|hostname|defaultOrganization|organization|token)\s*=\s*"([^"]*)"`)
	reFlatHostname = regexp.MustCompile(`^hostname\s*=\s*"([^"]*)"`)
	reFlatOrg      = regexp.MustCompile(`^(?:defaultOrganization|organization)\s*=\s*"([^"]*)"`)
	reFlatToken    = regexp.MustCompile(`^token\s*=\s*"([^"]*)"`)
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

// ListProfiles parses path and returns all profiles in file order.
//
//   - New format (file contains one or more profile blocks): each block
//     becomes one Profile. Name = block label; Hostname = hostname inside
//     the block, or DefaultHostname when the key is absent.
//   - Legacy flat format (no profile blocks): returns a single Profile with
//     Name = DefaultProfileName and Hostname = tfeHostname value or
//     DefaultHostname. Only returned when at least a token is present.
//   - File not found: returns nil, nil.
func ListProfiles(path string) ([]Profile, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	// Decide format: if ANY profile block exists → new format; otherwise legacy.
	hasBlocks := false
	for _, line := range lines {
		if reProfileStart.MatchString(line) {
			hasBlocks = true
			break
		}
	}

	if !hasBlocks {
		// Legacy flat format: parse top-level keys with sensible defaults.
		p := Profile{
			Name:     DefaultProfileName,
			Hostname: DefaultHostname,
		}
		for _, l := range lines {
			if m := reFlatHostname.FindStringSubmatch(l); m != nil {
				p.Hostname = m[1]
			}
			if m := reFlatOrg.FindStringSubmatch(l); m != nil {
				p.Organization = m[1]
			}
			if m := reFlatToken.FindStringSubmatch(l); m != nil {
				p.Token = m[1]
			}
		}
		if p.Token == "" {
			// Empty or comment-only file — no usable profile.
			return nil, nil
		}
		return []Profile{p}, nil
	}

	// New format — scan for profile blocks.
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
				case "tfeHostname", "hostname":
					current.Hostname = m[2]
				case "tfeOrganization", "defaultOrganization", "organization":
					current.Organization = m[2]
				case "tfeToken", "token":
					current.Token = m[2]
				}
				continue
			}
			if reBlockEnd.MatchString(line) {
				// If hostname was absent, apply defaults:
				//   - Block label looks like a hostname (contains ".") → backward compat
				//     for files written before name/hostname were separated.
				//   - Otherwise → DefaultHostname. Never use the profile name as a
				//     hostname for arbitrary aliases like "default" or "prod".
				if current.Hostname == "" {
					if strings.Contains(current.Name, ".") {
						current.Hostname = current.Name
					} else {
						current.Hostname = DefaultHostname
					}
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
// (flat keys, comments) are preserved. The file is written with 0600
// permissions because it contains an API token.
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
		orgLine = fmt.Sprintf("  defaultOrganization = %q", organization)
	} else {
		orgLine = `  # defaultOrganization = "" # set this to your organization name`
	}
	block := fmt.Sprintf(
		"\nprofile %q {\n  hostname            = %q\n%s\n  token               = %q\n}\n",
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
