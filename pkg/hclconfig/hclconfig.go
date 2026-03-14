// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

// Package hclconfig reads and writes TFx profile blocks in ~/.tfx.hcl.
//
// Profile block format:
//
//	profile "default" {
//	  tfeHostname     = "app.terraform.io"
//	  tfeOrganization = "my-org"
//	  tfeToken        = "abc123..."
//	}
//
// The block label is the profile name (user-editable alias). tfeHostname is a
// key inside the block. The name defaults to the hostname when written by
// `tfx login`.
//
// Old (flat) format is still supported for reading; it is treated as a single
// profile whose name and hostname both equal the value of tfeHostname.
package hclconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

// Profile holds configuration for one TFx profile.
type Profile struct {
	Name         string // block label — user-editable alias
	Hostname     string // tfeHostname value inside the block
	Organization string
	Token        string
}

var (
	reProfileStart = regexp.MustCompile(`^profile\s+"([^"]+)"\s*\{`)
	reKeyValue     = regexp.MustCompile(`^\s+(tfeHostname|tfeOrganization|tfeToken)\s*=\s*"([^"]*)"`)
	reFlatHostname = regexp.MustCompile(`^tfeHostname\s*=\s*"([^"]*)"`)
	reFlatOrg      = regexp.MustCompile(`^tfeOrganization\s*=\s*"([^"]*)"`)
	reFlatToken    = regexp.MustCompile(`^tfeToken\s*=\s*"([^"]*)"`)
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
//   - New format: each profile block becomes one entry (Name = block label,
//     Hostname = tfeHostname inside the block; falls back to Name if omitted).
//   - Old format (flat tfeHostname/tfeOrganization/tfeToken keys): returns a
//     single Profile with Name and Hostname both equal to the tfeHostname value.
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

	// Detect old (flat) format: look for a top-level tfeHostname line.
	for _, line := range lines {
		if m := reFlatHostname.FindStringSubmatch(line); m != nil {
			p := Profile{Name: m[1], Hostname: m[1]}
			for _, l := range lines {
				if mo := reFlatOrg.FindStringSubmatch(l); mo != nil {
					p.Organization = mo[1]
				}
				if mt := reFlatToken.FindStringSubmatch(l); mt != nil {
					p.Token = mt[1]
				}
			}
			return []Profile{p}, nil
		}
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
				case "tfeHostname":
					current.Hostname = m[2]
				case "tfeOrganization":
					current.Organization = m[2]
				case "tfeToken":
					current.Token = m[2]
				}
				continue
			}
			if reBlockEnd.MatchString(line) {
				// Backward compat: if tfeHostname was absent, use the block label.
				if current.Hostname == "" {
					current.Hostname = current.Name
				}
				profiles = append(profiles, *current)
				current = nil
			}
		}
	}
	return profiles, nil
}

// WriteProfile adds or replaces the named profile block in the file at path.
// If path does not exist it is created. All other profiles and content
// (flat keys, comments) are preserved unchanged. The file is written with 0600
// permissions because it contains an API token.
//
// name is the block label (profile alias). hostname is written as tfeHostname
// inside the block.
func WriteProfile(path, name, hostname, organization, token string) error {
	existing := ""
	if data, err := os.ReadFile(path); err == nil {
		existing = string(data)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("reading config file: %w", err)
	}

	stripped := removeProfileBlock(existing, name)

	var orgLine string
	if organization != "" {
		orgLine = fmt.Sprintf("  tfeOrganization = %q", organization)
	} else {
		orgLine = `  # tfeOrganization = "" # set this to your organization name`
	}
	block := fmt.Sprintf(
		"\nprofile %q {\n  tfeHostname     = %q\n%s\n  tfeToken        = %q\n}\n",
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
