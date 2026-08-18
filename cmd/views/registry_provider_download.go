// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	"fmt"
	"path/filepath"
	"strings"
)

// RegistryProviderDownloadedPlatform is one staged zip from a public registry download.
type RegistryProviderDownloadedPlatform struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Filename string `json:"filename"`
	Path     string `json:"path"`
	Shasum   string `json:"shasum"`
}

// RegistryProviderDownloadResult is the staged public-registry provider download.
type RegistryProviderDownloadResult struct {
	Directory        string                               `json:"directory"`
	Namespace        string                               `json:"namespace"`
	Name             string                               `json:"name"`
	Version          string                               `json:"version"`
	KeyID            string                               `json:"key_id"`
	ShasumsPath      string                               `json:"shasums"`
	ShasumsSigPath   string                               `json:"shasums_sig"`
	Platforms        []RegistryProviderDownloadedPlatform `json:"platforms"`
	SkippedPlatforms []string                             `json:"skipped_platforms,omitempty"`
	Commands         []string                             `json:"commands,omitempty"`
}

type RegistryProviderDownloadView struct{ *BaseView }

func NewRegistryProviderDownloadView() *RegistryProviderDownloadView {
	return &RegistryProviderDownloadView{NewBaseView()}
}

func (v *RegistryProviderDownloadView) Render(result *RegistryProviderDownloadResult) error {
	result.Commands = uploadCommands(result)

	if v.IsJSON() {
		return v.Output().RenderJSON(result)
	}

	props := []PropertyPair{
		{Key: "Namespace", Value: result.Namespace},
		{Key: "Name", Value: result.Name},
		{Key: "Version", Value: result.Version},
		{Key: "Directory", Value: result.Directory},
		{Key: "GPG Key ID", Value: result.KeyID},
		{Key: "SHA256SUMS", Value: result.ShasumsPath},
		{Key: "SHA256SUMS.sig", Value: result.ShasumsSigPath},
	}
	if err := v.Output().RenderProperties(props); err != nil {
		return err
	}

	if len(result.Platforms) > 0 {
		headers := []string{"OS", "Arch", "File"}
		rows := make([][]interface{}, len(result.Platforms))
		for i, p := range result.Platforms {
			rows[i] = []interface{}{p.OS, p.Arch, filepath.Base(p.Path)}
		}
		if err := v.Output().RenderTable(headers, rows); err != nil {
			return err
		}
	}

	if len(result.SkippedPlatforms) > 0 {
		v.Output().Message("Skipped platforms (not published for this version): %s", strings.Join(result.SkippedPlatforms, ", "))
	}

	v.Output().Message("Upload to a private registry with:")
	for _, c := range result.Commands {
		v.Output().Message("  %s", c)
	}
	return nil
}

func uploadCommands(result *RegistryProviderDownloadResult) []string {
	cmds := []string{
		fmt.Sprintf("tfx registry provider create --name %s", result.Name),
		fmt.Sprintf(
			"tfx registry provider version create --name %s --version %s --key-id %s --shasums %s --shasums-sig %s",
			result.Name,
			result.Version,
			result.KeyID,
			result.ShasumsPath,
			result.ShasumsSigPath,
		),
	}
	for _, p := range result.Platforms {
		cmds = append(cmds, fmt.Sprintf(
			"tfx registry provider version platform create --name %s --version %s --os %s --arch %s -f %s",
			result.Name,
			result.Version,
			p.OS,
			p.Arch,
			p.Path,
		))
	}
	return cmds
}
