// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// cacheDir returns (and creates if absent) the base TFx cache directory.
// Path: ~/.tfx/cache/
func cacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home dir: %w", err)
	}
	dir := filepath.Join(home, ".tfx", "cache")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("create cache dir: %w", err)
	}
	return dir, nil
}

// stateJSONPath returns the cache path for a state version's JSON.
// Path: ~/.tfx/cache/state/<svID>.json
func stateJSONPath(svID string) (string, error) {
	base, err := cacheDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, "state")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("create state cache dir: %w", err)
	}
	return filepath.Join(dir, svID+".json"), nil
}

// cvArchivePath returns the cache path for a config version's .tar.gz archive.
// Path: ~/.tfx/cache/cv/<cvID>.tar.gz
func cvArchivePath(cvID string) (string, error) {
	base, err := cacheDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, "cv")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("create cv cache dir: %w", err)
	}
	return filepath.Join(dir, cvID+".tar.gz"), nil
}

// cvExtractDir returns the directory where a config version's archive is extracted.
// Path: ~/.tfx/cache/cv/<cvID>/
func cvExtractDir(cvID string) (string, error) {
	base, err := cacheDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, "cv", cvID)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("create cv extract dir: %w", err)
	}
	return dir, nil
}

// extractTarGz extracts the .tar.gz file at src into destDir.
// All paths are sanitized to prevent directory traversal attacks.
func extractTarGz(src, destDir string) error {
	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open archive: %w", err)
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("gzip reader: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read tar: %w", err)
		}

		// Sanitize path: join with destDir and verify it stays inside.
		target := filepath.Join(destDir, filepath.Clean("/"+hdr.Name))
		if !isSubPath(destDir, target) {
			continue // skip path traversal attempts
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0700); err != nil {
				return fmt.Errorf("mkdir %s: %w", target, err)
			}
		case tar.TypeReg, tar.TypeRegA:
			if err := os.MkdirAll(filepath.Dir(target), 0700); err != nil {
				return fmt.Errorf("mkdir parent %s: %w", target, err)
			}
			out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
			if err != nil {
				return fmt.Errorf("create file %s: %w", target, err)
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return fmt.Errorf("write file %s: %w", target, err)
			}
			out.Close()
		// Symlinks, hard links, etc. are silently skipped for safety.
		}
	}
	return nil
}

// cvExtractDirPath returns the extraction directory for cvID without creating it.
// Use this for display purposes (status bar, hints) where side effects are undesired.
func cvExtractDirPath(cvID string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".tfx", "cache", "cv", cvID)
}

// tildePath replaces the home directory prefix with "~" for compact display.
func tildePath(p string) string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return p
	}
	if strings.HasPrefix(p, home) {
		return "~" + p[len(home):]
	}
	return p
}

// isSubPath reports whether child is inside (or equal to) parent.
func isSubPath(parent, child string) bool {
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}
	// filepath.Rel returns ".." or starts with "../" for paths outside parent.
	return rel != ".." && (len(rel) < 3 || rel[:3] != ".." + string(filepath.Separator))
}
