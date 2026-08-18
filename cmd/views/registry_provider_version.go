// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

type RegistryProviderVersionListView struct{ *BaseView }
type RegistryProviderVersionCreateView struct{ *BaseView }
type RegistryProviderVersionShowView struct{ *BaseView }
type RegistryProviderVersionDeleteView struct{ *BaseView }

func NewRegistryProviderVersionListView() *RegistryProviderVersionListView {
	return &RegistryProviderVersionListView{NewBaseView()}
}
func NewRegistryProviderVersionCreateView() *RegistryProviderVersionCreateView {
	return &RegistryProviderVersionCreateView{NewBaseView()}
}
func NewRegistryProviderVersionShowView() *RegistryProviderVersionShowView {
	return &RegistryProviderVersionShowView{NewBaseView()}
}
func NewRegistryProviderVersionDeleteView() *RegistryProviderVersionDeleteView {
	return &RegistryProviderVersionDeleteView{NewBaseView()}
}

func (v *RegistryProviderVersionListView) Render(items []*tfe.RegistryProviderVersion) error {
	if v.IsJSON() {
		return v.Output().RenderJSON(items)
	}
	headers := []string{"Version", "ID", "Published", "SHASUM", "SHASUM Sig"}
	rows := make([][]interface{}, len(items))
	for i, p := range items {
		rows[i] = []interface{}{p.Version, p.ID, p.UpdatedAt, p.ShasumsUploaded, p.ShasumsSigUploaded}
	}
	return v.Output().RenderTable(headers, rows)
}

func (v *RegistryProviderVersionCreateView) Render(p *tfe.RegistryProviderVersion) error {
	if v.IsJSON() {
		return v.Output().RenderJSON(p)
	}
	props := []PropertyPair{{Key: "ID", Value: p.ID}, {Key: "Version", Value: p.Version}, {Key: "Created", Value: p.UpdatedAt}}
	return v.Output().RenderProperties(props)
}

// RegistryProviderVersionCreateFromDirectoryResult is the directory-mode create output.
type RegistryProviderVersionCreateFromDirectoryResult struct {
	Name            string                          `json:"name"`
	Version         string                          `json:"version"`
	KeyID           string                          `json:"key_id"`
	ProviderCreated bool                            `json:"provider_created"`
	GPGKeyCreated   bool                            `json:"gpg_key_created"`
	ProviderVersion *tfe.RegistryProviderVersion    `json:"provider_version"`
	Platforms       []*tfe.RegistryProviderPlatform `json:"platforms"`
}

func (v *RegistryProviderVersionCreateView) RenderFromDirectory(result *RegistryProviderVersionCreateFromDirectoryResult) error {
	if v.IsJSON() {
		return v.Output().RenderJSON(result)
	}
	props := []PropertyPair{
		{Key: "Name", Value: result.Name},
		{Key: "Version", Value: result.Version},
		{Key: "GPG Key ID", Value: result.KeyID},
		{Key: "Provider Created", Value: result.ProviderCreated},
		{Key: "GPG Key Created", Value: result.GPGKeyCreated},
		{Key: "ID", Value: result.ProviderVersion.ID},
		{Key: "Created", Value: result.ProviderVersion.UpdatedAt},
	}
	if err := v.Output().RenderProperties(props); err != nil {
		return err
	}
	if len(result.Platforms) == 0 {
		return nil
	}
	headers := []string{"OS", "Arch", "Filename", "ID"}
	rows := make([][]interface{}, len(result.Platforms))
	for i, p := range result.Platforms {
		rows[i] = []interface{}{p.OS, p.Arch, p.Filename, p.ID}
	}
	return v.Output().RenderTable(headers, rows)
}

func (v *RegistryProviderVersionShowView) Render(p *tfe.RegistryProviderVersion, shasums string) error {
	if v.IsJSON() {
		return v.Output().RenderJSON(map[string]interface{}{"version": p, "shasums": shasums})
	}
	props := []PropertyPair{
		{Key: "Version", Value: p.Version},
		{Key: "ID", Value: p.ID},
		{Key: "Shasums Uploaded", Value: p.ShasumsUploaded},
		{Key: "Shasums Sig Uploaded", Value: p.ShasumsSigUploaded},
	}
	if err := v.Output().RenderProperties(props); err != nil {
		return err
	}
	if shasums != "" {
		v.Output().Message("\nShasums:\n%s", shasums)
	}
	return nil
}

func (v *RegistryProviderVersionDeleteView) Render(name string) error {
	if v.IsJSON() {
		return v.Output().RenderJSON(map[string]interface{}{"status": "Success", "name": name})
	}
	return v.Output().RenderProperties([]PropertyPair{{Key: "Status", Value: "Success"}, {Key: "Name", Value: name}})
}
