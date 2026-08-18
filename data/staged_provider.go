// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/pkg/errors"
	"github.com/straubt1/tfx/cmd/flags"
	"github.com/straubt1/tfx/output"
	"github.com/straubt1/tfx/pkg/file"
)

// PublicRegistryHashiCorpNamespace is flags.PublicRegistryHashiCorpNamespace.
const PublicRegistryHashiCorpNamespace = flags.PublicRegistryHashiCorpNamespace

const stagedProviderMetadataFilename = "tfx-provider.json"

// HashiCorp GPG short IDs (8 hex) as used in SHA256SUMS.<id>.sig filenames,
// mapped to the 16-hex long IDs TFE expects for --key-id.
var hashicorpGPGShortIDs = map[string]string{
	"72D7468F": "34365D9472D7468F",
	"348FFC4C": "51852D87348FFC4C",
}

// IsHashiCorpPublicNamespace reports whether namespace is the official
// public-registry publisher "hashicorp" (case-insensitive). Provider name
// and GPG key ID are not part of this check.
func IsHashiCorpPublicNamespace(namespace string) bool {
	return strings.EqualFold(namespace, PublicRegistryHashiCorpNamespace)
}

type stagedProviderMetadata struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	KeyID     string `json:"key_id"`
}

// StagedProviderPlatform is one zip found in a staged provider directory.
type StagedProviderPlatform struct {
	OS       string
	Arch     string
	Filename string
	Path     string
	Shasum   string
}

// StagedProviderDirectory is the inferred contents of a download-staged folder.
type StagedProviderDirectory struct {
	// Namespace is the public Terraform Registry publisher (e.g. hashicorp, chainguard-dev).
	Namespace    string
	Name         string
	Version      string
	KeyID        string
	GPGPublicKey string
	Shasums      string
	ShasumsSig   string
	Platforms    []StagedProviderPlatform
}

// ReadStagedProviderDirectory inspects a directory produced by
// `tfx registry provider download` and infers name, version, GPG key id,
// SHA256SUMS/sig paths, and platform zips present on disk.
func ReadStagedProviderDirectory(dir string) (*StagedProviderDirectory, error) {
	log := output.Get().Logger()
	log.Debug("Reading staged provider directory", "dir", dir)

	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve directory")
	}
	if !file.IsDirectory(abs) {
		return nil, errors.Errorf("directory does not exist: %s", dir)
	}

	shasums, err := globOne(abs, "terraform-provider-*_SHA256SUMS")
	if err != nil {
		return nil, errors.Wrap(err, "SHA256SUMS")
	}
	sig, err := findShasumsSig(abs)
	if err != nil {
		return nil, err
	}

	name, version, err := parseShasumsFilename(filepath.Base(shasums))
	if err != nil {
		return nil, err
	}
	if err := matchPathNameVersion(abs, name, version); err != nil {
		return nil, err
	}

	keyID, err := inferKeyIDFromSigFilename(filepath.Base(sig))
	if err != nil {
		return nil, err
	}

	meta, err := readStagedProviderMetadata(abs)
	if err != nil {
		return nil, err
	}
	namespace := ""
	if meta != nil {
		namespace = strings.TrimSpace(meta.Namespace)
		if meta.KeyID != "" {
			keyID = strings.ToUpper(strings.TrimSpace(meta.KeyID))
		}
		if meta.Name != "" && meta.Name != name {
			return nil, errors.Errorf("tfx-provider.json name %s does not match %s", meta.Name, name)
		}
		if meta.Version != "" && meta.Version != version {
			return nil, errors.Errorf("tfx-provider.json version %s does not match %s", meta.Version, version)
		}
	}

	ascPath, err := findGPGPublicKey(abs, keyID)
	if err != nil {
		return nil, err
	}
	if keyID == "" && ascPath != "" {
		keyID = keyIDFromAscPath(ascPath)
	}
	if namespace == "" && isHashiCorpGPGKeyID(keyID) {
		namespace = PublicRegistryHashiCorpNamespace
	}

	b, err := os.ReadFile(shasums)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read SHA256SUMS")
	}
	checksums := parseSHA256SUMS(string(b))

	zips, err := filepath.Glob(filepath.Join(abs, "terraform-provider-*.zip"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to list zip files")
	}

	var platforms []StagedProviderPlatform
	for _, zipPath := range zips {
		base := filepath.Base(zipPath)
		osName, arch, err := parseProviderZipFilename(name, version, base)
		if err != nil {
			return nil, err
		}
		expected, ok := checksums[base]
		if !ok {
			return nil, errors.Errorf("zip %s is not listed in SHA256SUMS", base)
		}
		sum, err := fileSHA256(zipPath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to hash %s", base)
		}
		if !strings.EqualFold(sum, expected) {
			return nil, errors.Errorf("checksum mismatch for %s: got %s want %s", base, sum, expected)
		}
		platforms = append(platforms, StagedProviderPlatform{
			OS:       osName,
			Arch:     arch,
			Filename: base,
			Path:     zipPath,
			Shasum:   sum,
		})
	}
	if len(platforms) == 0 {
		return nil, errors.New("no provider zip files found in directory")
	}

	log.Debug("Staged provider directory read", "namespace", namespace, "name", name, "version", version, "keyID", keyID, "platforms", len(platforms))
	return &StagedProviderDirectory{
		Namespace:    namespace,
		Name:         name,
		Version:      version,
		KeyID:        keyID,
		GPGPublicKey: ascPath,
		Shasums:      shasums,
		ShasumsSig:   sig,
		Platforms:    platforms,
	}, nil
}

func writeStagedProviderMetadata(dir, namespace, name, version, keyID string) error {
	meta := stagedProviderMetadata{
		Namespace: namespace,
		Name:      name,
		Version:   version,
		KeyID:     keyID,
	}
	b, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to encode tfx-provider.json")
	}
	b = append(b, '\n')
	path := filepath.Join(dir, stagedProviderMetadataFilename)
	if err := os.WriteFile(path, b, 0644); err != nil {
		return errors.Wrap(err, "failed to write tfx-provider.json")
	}
	return nil
}

func readStagedProviderMetadata(dir string) (*stagedProviderMetadata, error) {
	path := filepath.Join(dir, stagedProviderMetadataFilename)
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to read tfx-provider.json")
	}
	var meta stagedProviderMetadata
	if err := json.Unmarshal(b, &meta); err != nil {
		return nil, errors.Wrap(err, "failed to parse tfx-provider.json")
	}
	return &meta, nil
}

func findGPGPublicKey(dir, keyID string) (string, error) {
	if keyID != "" {
		p := filepath.Join(dir, strings.ToUpper(keyID)+".asc")
		if file.IsFile(p) {
			return p, nil
		}
	}
	matches, err := filepath.Glob(filepath.Join(dir, "*.asc"))
	if err != nil {
		return "", errors.Wrap(err, "failed to list GPG public key files")
	}
	var hexMatches []string
	for _, m := range matches {
		base := strings.TrimSuffix(filepath.Base(m), ".asc")
		if len(base) == 16 && isHex(strings.ToUpper(base)) {
			hexMatches = append(hexMatches, m)
		}
	}
	if len(hexMatches) == 0 {
		return "", nil
	}
	if len(hexMatches) > 1 {
		return "", errors.New("multiple GPG public key .asc files found")
	}
	return hexMatches[0], nil
}

func keyIDFromAscPath(path string) string {
	return strings.ToUpper(strings.TrimSuffix(filepath.Base(path), ".asc"))
}

func isHashiCorpGPGKeyID(keyID string) bool {
	id := strings.ToUpper(strings.TrimSpace(keyID))
	for _, v := range hashicorpGPGShortIDs {
		if v == id {
			return true
		}
	}
	return false
}

func globOne(dir, pattern string) (string, error) {
	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		return "", errors.Errorf("no file matching %s", pattern)
	}
	if len(matches) > 1 {
		return "", errors.Errorf("multiple files matching %s", pattern)
	}
	return matches[0], nil
}

func findShasumsSig(dir string) (string, error) {
	plain, err := filepath.Glob(filepath.Join(dir, "terraform-provider-*_SHA256SUMS.sig"))
	if err != nil {
		return "", err
	}
	keyed, err := filepath.Glob(filepath.Join(dir, "terraform-provider-*_SHA256SUMS.*.sig"))
	if err != nil {
		return "", err
	}
	matches := append(plain, keyed...)
	if len(matches) == 0 {
		return "", errors.New("no SHA256SUMS signature file found")
	}
	if len(matches) > 1 {
		return "", errors.New("multiple SHA256SUMS signature files found")
	}
	return matches[0], nil
}

func parseShasumsFilename(base string) (name, version string, err error) {
	const prefix = "terraform-provider-"
	const suffix = "_SHA256SUMS"
	if !strings.HasPrefix(base, prefix) || !strings.HasSuffix(base, suffix) {
		return "", "", errors.Errorf("unexpected SHA256SUMS filename %s", base)
	}
	rest := strings.TrimSuffix(strings.TrimPrefix(base, prefix), suffix)
	name, version, ok := strings.Cut(rest, "_")
	if !ok || name == "" || version == "" || strings.Contains(version, "_") {
		return "", "", errors.Errorf("could not parse name and version from %s", base)
	}
	if _, err := semver.NewVersion(version); err != nil {
		return "", "", errors.Wrapf(err, "invalid version in %s", base)
	}
	return name, version, nil
}

func parseProviderZipFilename(name, version, base string) (osName, arch string, err error) {
	prefix := "terraform-provider-" + name + "_" + version + "_"
	if !strings.HasPrefix(base, prefix) || !strings.HasSuffix(base, ".zip") {
		return "", "", errors.Errorf("zip %s does not match terraform-provider-%s_%s_<os>_<arch>.zip", base, name, version)
	}
	rest := strings.TrimSuffix(strings.TrimPrefix(base, prefix), ".zip")
	osName, arch, ok := strings.Cut(rest, "_")
	if !ok || osName == "" || arch == "" {
		return "", "", errors.Errorf("could not parse os/arch from %s", base)
	}
	return osName, arch, nil
}

func matchPathNameVersion(abs, name, version string) error {
	base := filepath.Base(abs)
	parent := filepath.Base(filepath.Dir(abs))
	if _, err := semver.NewVersion(base); err != nil {
		return nil
	}
	if parent == "." || parent == string(filepath.Separator) || parent == "" {
		return nil
	}
	if parent != name || base != version {
		return errors.Errorf("directory %s/%s does not match provider %s version %s", parent, base, name, version)
	}
	return nil
}

func inferKeyIDFromSigFilename(base string) (string, error) {
	const marker = "_SHA256SUMS"
	idx := strings.Index(base, marker)
	if idx < 0 || !strings.HasSuffix(base, ".sig") {
		return "", errors.Errorf("unexpected signature filename %s", base)
	}
	rest := strings.TrimSuffix(base[idx+len(marker):], ".sig")
	if rest == "" {
		return "", nil
	}
	id := strings.TrimPrefix(rest, ".")
	id = strings.ToUpper(id)
	if len(id) == 16 && isHex(id) {
		return id, nil
	}
	if mapped, ok := hashicorpGPGShortIDs[id]; ok {
		return mapped, nil
	}
	if len(id) == 8 && isHex(id) {
		return "", errors.Errorf("could not infer GPG key id from short suffix %s; pass --key-id", id)
	}
	return "", errors.Errorf("could not infer GPG key id from %s; pass --key-id", base)
}

func isHex(s string) bool {
	for _, c := range s {
		if (c < '0' || c > '9') && (c < 'A' || c > 'F') {
			return false
		}
	}
	return len(s) > 0
}
