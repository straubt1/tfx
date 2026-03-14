// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

// Package hclconfig reads and writes TFx profile blocks in ~/.tfx.hcl.
//
// New format — profile name is the hostname:
//
//	profile "app.terraform.io" {
//	  tfeOrganization = "my-org"
//	  tfeToken        = "abc123..."
//	}
//
// Old (flat) format is still supported for reading; it is treated as a single
// profile whose name equals the value of tfeHostname.
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
// In the new format the Hostname field is the profile block label.
type Profile struct {
	Hostname     string
	Organization string
	Token        string
}

var (
	reProfileStart = regexp.MustCompile(`^profile\s+"([^"]+)"\s*\{`)
	reKeyValue     = regexp.MustCompile(`^\s+(tfeOrganization|tfeToken)\s*=\s*"([^"]*)"`)
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
//   - New format: each profile block becomes one entry (Hostname = block label).
//   - Old format (flat tfeHostname/tfeOrganization/tfeToken keys): returns a
//     single Profile with Hostname = the tfeHostname value.
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
			// Old format — collect flat values.
			p := Profile{Hostname: m[1]}
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
			current = &Profile{Hostname: m[1]}
			continue
		}
		if current != nil {
			if m := reKeyValue.FindStringSubmatch(line); m != nil {
				switch m[1] {
				case "tfeOrganization":
					current.Organization = m[2]
				case "tfeToken":
					current.Token = m[2]
				}
				continue
			}
			if reBlockEnd.MatchString(line) {
				profiles = append(profiles, *current)
				current = nil
			}
		}
	}
	return profiles, nil
}

// WriteProfile adds or replaces the profile block for hostname in the file at
// path. If path does not exist it is created. All other profiles and content
// (flat keys, comments) are preserved unchanged. The file is written with 0600
// permissions because it contains an API token.
func WriteProfile(path, hostname, organization, token string) error {
	// Read existing content (empty string if file doesn't exist yet).
	existing := ""
	if data, err := os.ReadFile(path); err == nil {
		existing = string(data)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("reading config file: %w", err)
	}

	// Remove existing block for this hostname so we don't duplicate it.
	stripped := removeProfileBlock(existing, hostname)

	// Build new block.
	block := fmt.Sprintf(
		"\nprofile %q {\n  tfeOrganization = %q\n  tfeToken        = %q\n}\n",
		hostname, organization, token,
	)

	// Normalise trailing whitespace on the surviving content, then append.
	content := strings.TrimRight(stripped, "\n\t ")
	if content != "" {
		content += "\n"
	}
	content += block

	// Ensure the parent directory exists (e.g. if path is ~/foo/bar.hcl).
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}
	return os.WriteFile(path, []byte(content), 0600)
}

// removeProfileBlock strips the profile block for hostname from content,
// returning everything else. Uses a simple brace-depth counter so it handles
// the closing } correctly even if the block contains nested structures.
func removeProfileBlock(content, hostname string) string {
	target := fmt.Sprintf("profile %q", hostname)
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
