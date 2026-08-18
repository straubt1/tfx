// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	view "github.com/straubt1/tfx/cmd/views"
	"github.com/straubt1/tfx/output"
)

const defaultPublicRegistryURL = "https://registry.terraform.io"

// PublicProviderDownloadConfig is the input for DownloadPublicProvider.
type PublicProviderDownloadConfig struct {
	// RegistryBaseURL is the public registry origin, e.g. https://registry.terraform.io.
	// Empty uses the public Terraform Registry.
	RegistryBaseURL string
	Namespace       string
	Name            string
	Version         string
	// Directory is the staging base directory. Files are written to
	// <Directory>/<Namespace>/<Name>/<Version>/. Created if missing.
	Directory string
	// Platforms are os_arch values such as linux_amd64. Ignored when AllPlatforms is true.
	Platforms    []string
	AllPlatforms bool
}

type wellKnownDocument struct {
	ProvidersV1 string `json:"providers.v1"`
}

type publicProviderVersionsResponse struct {
	Versions []publicProviderVersion `json:"versions"`
}

type publicProviderVersion struct {
	Version   string                   `json:"version"`
	Platforms []publicProviderPlatform `json:"platforms"`
}

type publicProviderPlatform struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

type publicProviderPackage struct {
	OS                  string `json:"os"`
	Arch                string `json:"arch"`
	Filename            string `json:"filename"`
	DownloadURL         string `json:"download_url"`
	ShasumsURL          string `json:"shasums_url"`
	ShasumsSignatureURL string `json:"shasums_signature_url"`
	Shasum              string `json:"shasum"`
	SigningKeys         struct {
		GPGPublicKeys []struct {
			KeyID      string `json:"key_id"`
			AsciiArmor string `json:"ascii_armor"`
		} `json:"gpg_public_keys"`
	} `json:"signing_keys"`
}

// ParsePlatform splits an os_arch string (e.g. linux_amd64) into os and arch.
func ParsePlatform(s string) (osName, arch string, err error) {
	s = strings.TrimSpace(s)
	osName, arch, ok := strings.Cut(s, "_")
	if !ok || osName == "" || arch == "" {
		return "", "", errors.Errorf("invalid platform %q, expected os_arch (e.g. linux_amd64)", s)
	}
	return osName, arch, nil
}

func platformKey(osName, arch string) string {
	return osName + "_" + arch
}

// DownloadPublicProvider fetches provider artifacts from the public registry
// protocol (service discovery → versions → package metadata → files) and
// writes them under <Directory>/<Namespace>/<Name>/<Version>/.
func DownloadPublicProvider(cfg PublicProviderDownloadConfig) (*view.RegistryProviderDownloadResult, error) {
	log := output.Get().Logger()
	log.Debug("Downloading public provider",
		"namespace", cfg.Namespace,
		"name", cfg.Name,
		"version", cfg.Version,
		"directory", cfg.Directory,
		"allPlatforms", cfg.AllPlatforms,
	)

	if cfg.Namespace == "" {
		return nil, errors.New("namespace is required")
	}
	if cfg.Name == "" {
		return nil, errors.New("name is required")
	}
	if cfg.Version == "" {
		return nil, errors.New("version is required")
	}
	if cfg.Directory == "" {
		return nil, errors.New("directory is required")
	}

	base := strings.TrimRight(cfg.RegistryBaseURL, "/")
	if base == "" {
		base = defaultPublicRegistryURL
	}

	providersBase, err := discoverProvidersV1(base)
	if err != nil {
		return nil, err
	}

	versions, err := listPublicProviderVersions(providersBase, cfg.Namespace, cfg.Name)
	if err != nil {
		return nil, err
	}

	var selected *publicProviderVersion
	for i := range versions {
		if versions[i].Version == cfg.Version {
			selected = &versions[i]
			break
		}
	}
	if selected == nil {
		return nil, errors.Errorf("version %s not found for %s/%s", cfg.Version, cfg.Namespace, cfg.Name)
	}

	available := make(map[string]publicProviderPlatform, len(selected.Platforms))
	for _, p := range selected.Platforms {
		available[platformKey(p.OS, p.Arch)] = p
	}

	wanted, err := resolveWantedPlatforms(cfg, selected.Platforms)
	if err != nil {
		return nil, err
	}

	destDir, err := filepath.Abs(filepath.Join(cfg.Directory, cfg.Namespace, cfg.Name, cfg.Version))
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve staging directory")
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, errors.Wrap(err, "failed to create staging directory")
	}

	result := &view.RegistryProviderDownloadResult{
		Directory: destDir,
		Namespace: cfg.Namespace,
		Name:      cfg.Name,
		Version:   cfg.Version,
	}

	var checksums map[string]string
	var firstPkg *publicProviderPackage

	for _, plat := range wanted {
		key := platformKey(plat.OS, plat.Arch)
		if _, ok := available[key]; !ok {
			log.Warn("Skipping platform not published for this version", "platform", key)
			result.SkippedPlatforms = append(result.SkippedPlatforms, key)
			continue
		}

		pkg, err := findPublicProviderPackage(providersBase, cfg.Namespace, cfg.Name, cfg.Version, plat.OS, plat.Arch)
		if err != nil {
			log.Warn("Skipping platform: package metadata unavailable", "platform", key, "error", err)
			result.SkippedPlatforms = append(result.SkippedPlatforms, key)
			continue
		}

		if firstPkg == nil {
			firstPkg = pkg
			result.KeyID = packageKeyID(pkg)
			shasumsPath, sigPath, sums, err := downloadChecksumFiles(destDir, pkg)
			if err != nil {
				return nil, err
			}
			result.ShasumsPath = shasumsPath
			result.ShasumsSigPath = sigPath
			checksums = sums
			ascPath, err := writeGPGPublicKey(destDir, pkg)
			if err != nil {
				return nil, err
			}
			result.GPGPublicKeyPath = ascPath
		}

		zipPath, err := safeJoin(destDir, pkg.Filename)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid package filename for %s", key)
		}
		log.Debug("Downloading provider zip", "platform", key, "url", pkg.DownloadURL, "path", zipPath)
		if err := DownloadFile(pkg.DownloadURL, zipPath); err != nil {
			return nil, errors.Wrapf(err, "failed to download %s", key)
		}

		sum, err := fileSHA256(zipPath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to hash %s", pkg.Filename)
		}
		if pkg.Shasum != "" && !strings.EqualFold(sum, pkg.Shasum) {
			return nil, errors.Errorf("checksum mismatch for %s: got %s want %s", pkg.Filename, sum, pkg.Shasum)
		}
		expected, ok := checksums[pkg.Filename]
		if !ok {
			expected, ok = checksums[filepath.Base(pkg.Filename)]
		}
		if !ok {
			return nil, errors.Errorf("zip %s is not listed in SHA256SUMS", pkg.Filename)
		}
		if !strings.EqualFold(sum, expected) {
			return nil, errors.Errorf("checksum mismatch for %s vs SHA256SUMS: got %s want %s", pkg.Filename, sum, expected)
		}

		result.Platforms = append(result.Platforms, view.RegistryProviderDownloadedPlatform{
			OS:       plat.OS,
			Arch:     plat.Arch,
			Filename: pkg.Filename,
			Path:     zipPath,
			Shasum:   sum,
		})
	}

	if len(result.Platforms) == 0 {
		return nil, errors.New("no provider platforms were downloaded")
	}

	log.Info("Public provider download complete",
		"name", cfg.Name,
		"version", cfg.Version,
		"directory", destDir,
		"platforms", len(result.Platforms),
		"skipped", len(result.SkippedPlatforms),
	)
	return result, nil
}

func resolveWantedPlatforms(cfg PublicProviderDownloadConfig, published []publicProviderPlatform) ([]publicProviderPlatform, error) {
	if cfg.AllPlatforms {
		if len(published) == 0 {
			return nil, errors.New("provider version lists no platforms")
		}
		return published, nil
	}

	if len(cfg.Platforms) == 0 {
		return nil, errors.New("at least one platform is required")
	}

	var wanted []publicProviderPlatform
	seen := map[string]bool{}
	for _, raw := range cfg.Platforms {
		osName, arch, err := ParsePlatform(raw)
		if err != nil {
			return nil, err
		}
		key := platformKey(osName, arch)
		if seen[key] {
			continue
		}
		seen[key] = true
		wanted = append(wanted, publicProviderPlatform{OS: osName, Arch: arch})
	}
	return wanted, nil
}

func discoverProvidersV1(registryBase string) (string, error) {
	wellKnown := strings.TrimRight(registryBase, "/") + "/.well-known/terraform.json"
	output.Get().Logger().Debug("Public registry service discovery", "url", wellKnown)

	var doc wellKnownDocument
	if err := getJSON(wellKnown, &doc); err != nil {
		return "", errors.Wrap(err, "failed public registry service discovery")
	}
	if doc.ProvidersV1 == "" {
		return "", errors.New("registry service discovery is missing providers.v1")
	}

	base, err := url.Parse(wellKnown)
	if err != nil {
		return "", errors.Wrap(err, "invalid registry base URL")
	}
	ref, err := url.Parse(doc.ProvidersV1)
	if err != nil {
		return "", errors.Wrap(err, "invalid providers.v1 URL")
	}
	resolved := base.ResolveReference(ref).String()
	if !strings.HasSuffix(resolved, "/") {
		resolved += "/"
	}
	return resolved, nil
}

func listPublicProviderVersions(providersBase, namespace, name string) ([]publicProviderVersion, error) {
	u, err := url.JoinPath(providersBase, namespace, name, "versions")
	if err != nil {
		return nil, err
	}
	output.Get().Logger().Debug("Listing public provider versions", "url", u)

	var resp publicProviderVersionsResponse
	if err := getJSON(u, &resp); err != nil {
		return nil, errors.Wrapf(err, "failed to list versions for %s/%s", namespace, name)
	}
	return resp.Versions, nil
}

func findPublicProviderPackage(providersBase, namespace, name, version, osName, arch string) (*publicProviderPackage, error) {
	u, err := url.JoinPath(providersBase, namespace, name, version, "download", osName, arch)
	if err != nil {
		return nil, err
	}
	output.Get().Logger().Trace("Finding public provider package", "url", u)

	var pkg publicProviderPackage
	if err := getJSON(u, &pkg); err != nil {
		return nil, errors.Wrapf(err, "failed to find package %s/%s %s %s_%s", namespace, name, version, osName, arch)
	}
	if pkg.Filename == "" || pkg.DownloadURL == "" {
		return nil, errors.New("package metadata is missing filename or download_url")
	}
	pkg.DownloadURL = resolveURL(u, pkg.DownloadURL)
	pkg.ShasumsURL = resolveURL(u, pkg.ShasumsURL)
	pkg.ShasumsSignatureURL = resolveURL(u, pkg.ShasumsSignatureURL)
	return &pkg, nil
}

func resolveURL(base, ref string) string {
	if ref == "" {
		return ""
	}
	b, err := url.Parse(base)
	if err != nil {
		return ref
	}
	r, err := url.Parse(ref)
	if err != nil {
		return ref
	}
	return b.ResolveReference(r).String()
}

func filenameFromURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Path == "" {
		return path.Base(raw)
	}
	return path.Base(u.Path)
}

// safeJoin joins dir and name so the result stays under dir. name must be a
// single path component (no separators, not absolute, not "." or "..").
func safeJoin(dir, name string) (string, error) {
	if name == "" {
		return "", errors.New("empty filename")
	}
	cleaned := filepath.Clean(name)
	base := filepath.Base(cleaned)
	if base == "." || base == ".." || base == string(filepath.Separator) {
		return "", errors.Errorf("invalid filename %q", name)
	}
	if filepath.IsAbs(cleaned) || filepath.Dir(cleaned) != "." {
		return "", errors.Errorf("invalid filename %q, must be a basename", name)
	}

	parent, err := filepath.Abs(dir)
	if err != nil {
		return "", errors.Wrap(err, "failed to resolve directory")
	}
	dest := filepath.Join(parent, base)
	rel, err := filepath.Rel(parent, dest)
	if err != nil {
		return "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", errors.Errorf("filename %q escapes staging directory", name)
	}
	return dest, nil
}

func downloadChecksumFiles(destDir string, pkg *publicProviderPackage) (shasumsPath, sigPath string, checksums map[string]string, err error) {
	if pkg.ShasumsURL == "" || pkg.ShasumsSignatureURL == "" {
		return "", "", nil, errors.New("package metadata is missing shasums URLs")
	}

	shasumsPath, err = safeJoin(destDir, filenameFromURL(pkg.ShasumsURL))
	if err != nil {
		return "", "", nil, errors.Wrap(err, "invalid SHA256SUMS filename")
	}
	sigPath, err = safeJoin(destDir, filenameFromURL(pkg.ShasumsSignatureURL))
	if err != nil {
		return "", "", nil, errors.Wrap(err, "invalid SHA256SUMS signature filename")
	}

	output.Get().Logger().Debug("Downloading SHA256SUMS", "url", pkg.ShasumsURL, "path", shasumsPath)
	if err := DownloadFile(pkg.ShasumsURL, shasumsPath); err != nil {
		return "", "", nil, errors.Wrap(err, "failed to download SHA256SUMS")
	}
	output.Get().Logger().Debug("Downloading SHA256SUMS signature", "url", pkg.ShasumsSignatureURL, "path", sigPath)
	if err := DownloadFile(pkg.ShasumsSignatureURL, sigPath); err != nil {
		return "", "", nil, errors.Wrap(err, "failed to download SHA256SUMS signature")
	}

	b, err := os.ReadFile(shasumsPath)
	if err != nil {
		return "", "", nil, errors.Wrap(err, "failed to read SHA256SUMS")
	}
	return shasumsPath, sigPath, parseSHA256SUMS(string(b)), nil
}

func parseSHA256SUMS(content string) map[string]string {
	sums := make(map[string]string)
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		filename := strings.TrimPrefix(fields[len(fields)-1], "*")
		sums[filename] = fields[0]
	}
	return sums
}

func packageKeyID(pkg *publicProviderPackage) string {
	if pkg == nil || len(pkg.SigningKeys.GPGPublicKeys) == 0 {
		return ""
	}
	return strings.ToUpper(strings.TrimSpace(pkg.SigningKeys.GPGPublicKeys[0].KeyID))
}

func writeGPGPublicKey(destDir string, pkg *publicProviderPackage) (string, error) {
	if pkg == nil || len(pkg.SigningKeys.GPGPublicKeys) == 0 {
		return "", nil
	}
	keyID := packageKeyID(pkg)
	armor := strings.TrimSpace(pkg.SigningKeys.GPGPublicKeys[0].AsciiArmor)
	if keyID == "" || armor == "" {
		return "", nil
	}
	if !strings.HasSuffix(armor, "\n") {
		armor += "\n"
	}
	ascPath, err := safeJoin(destDir, keyID+".asc")
	if err != nil {
		return "", errors.Wrap(err, "invalid GPG public key filename")
	}
	output.Get().Logger().Debug("Writing GPG public key", "keyID", keyID, "path", ascPath)
	if err := os.WriteFile(ascPath, []byte(armor), 0644); err != nil {
		return "", errors.Wrap(err, "failed to write GPG public key")
	}
	return ascPath, nil
}

func fileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func getJSON(rawURL string, dest interface{}) error {
	resp, err := publicRegistryGet(rawURL)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return errors.Errorf("not found: %s", rawURL)
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("unexpected status %d from %s", resp.StatusCode, rawURL)
	}
	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return errors.Wrap(err, "failed to decode JSON")
	}
	return nil
}
