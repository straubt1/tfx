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
		return v.renderer.RenderJSON(items)
	}
	headers := []string{"Version", "ID", "Published", "SHASUM", "SHASUM Sig"}
	rows := make([][]interface{}, len(items))
	for i, p := range items {
		rows[i] = []interface{}{p.Version, p.ID, p.UpdatedAt, p.ShasumsUploaded, p.ShasumsSigUploaded}
	}
	return v.renderer.RenderTable(headers, rows)
}

func (v *RegistryProviderVersionCreateView) Render(p *tfe.RegistryProviderVersion) error {
	if v.IsJSON() {
		return v.renderer.RenderJSON(p)
	}
	props := []PropertyPair{{Key: "ID", Value: p.ID}, {Key: "Version", Value: p.Version}, {Key: "Created", Value: p.UpdatedAt}}
	return v.renderer.RenderProperties(props)
}

func (v *RegistryProviderVersionShowView) Render(p *tfe.RegistryProviderVersion, shasums string) error {
	if v.IsJSON() {
		return v.renderer.RenderJSON(map[string]interface{}{"version": p, "shasums": shasums})
	}
	props := []PropertyPair{
		{Key: "Version", Value: p.Version},
		{Key: "ID", Value: p.ID},
		{Key: "Shasums Uploaded", Value: p.ShasumsUploaded},
		{Key: "Shasums Sig Uploaded", Value: p.ShasumsSigUploaded},
	}
	if err := v.renderer.RenderProperties(props); err != nil {
		return err
	}
	if shasums != "" {
		v.renderer.Message("\nShasums:\n%s", shasums)
	}
	return nil
}

func (v *RegistryProviderVersionDeleteView) Render(name string) error {
	if v.IsJSON() {
		return v.renderer.RenderJSON(map[string]interface{}{"status": "Success", "name": name})
	}
	return v.renderer.RenderProperties([]PropertyPair{{Key: "Status", Value: "Success"}, {Key: "Name", Value: name}})
}
