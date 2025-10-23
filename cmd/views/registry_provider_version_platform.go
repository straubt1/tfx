// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

type RegistryProviderPlatformListView struct{ *BaseView }
type RegistryProviderPlatformCreateView struct{ *BaseView }
type RegistryProviderPlatformShowView struct{ *BaseView }
type RegistryProviderPlatformDeleteView struct{ *BaseView }

func NewRegistryProviderPlatformListView() *RegistryProviderPlatformListView {
	return &RegistryProviderPlatformListView{NewBaseView()}
}
func NewRegistryProviderPlatformCreateView() *RegistryProviderPlatformCreateView {
	return &RegistryProviderPlatformCreateView{NewBaseView()}
}
func NewRegistryProviderPlatformShowView() *RegistryProviderPlatformShowView {
	return &RegistryProviderPlatformShowView{NewBaseView()}
}
func NewRegistryProviderPlatformDeleteView() *RegistryProviderPlatformDeleteView {
	return &RegistryProviderPlatformDeleteView{NewBaseView()}
}

func (v *RegistryProviderPlatformListView) Render(items []*tfe.RegistryProviderPlatform) error {
	if v.IsJSON() {
		return v.renderer.RenderJSON(items)
	}
	headers := []string{"OS", "Arch", "ID", "Filename", "Shasum"}
	rows := make([][]interface{}, len(items))
	for i, p := range items {
		rows[i] = []interface{}{p.OS, p.Arch, p.ID, p.Filename, p.Shasum}
	}
	return v.renderer.RenderTable(headers, rows)
}

func (v *RegistryProviderPlatformCreateView) Render(p *tfe.RegistryProviderPlatform) error {
	if v.IsJSON() {
		return v.renderer.RenderJSON(p)
	}
	props := []PropertyPair{{Key: "ID", Value: p.ID}, {Key: "OS", Value: p.OS}, {Key: "Arch", Value: p.Arch}}
	return v.renderer.RenderProperties(props)
}

func (v *RegistryProviderPlatformShowView) Render(p *tfe.RegistryProviderPlatform) error {
	if v.IsJSON() {
		return v.renderer.RenderJSON(p)
	}
	props := []PropertyPair{
		{Key: "ID", Value: p.ID},
		{Key: "OS", Value: p.OS},
		{Key: "ARCH", Value: p.Arch},
		{Key: "Filename", Value: p.Filename},
		{Key: "Shasum", Value: p.Shasum},
	}
	return v.renderer.RenderProperties(props)
}

func (v *RegistryProviderPlatformDeleteView) Render(name string) error {
	if v.IsJSON() {
		return v.renderer.RenderJSON(map[string]interface{}{"status": "Success", "name": name})
	}
	return v.renderer.RenderProperties([]PropertyPair{{Key: "Status", Value: "Success"}, {Key: "Name", Value: name}})
}
