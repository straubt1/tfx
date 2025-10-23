// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

type RegistryProviderListView struct{ *BaseView }
type RegistryProviderCreateView struct{ *BaseView }
type RegistryProviderShowView struct{ *BaseView }
type RegistryProviderDeleteView struct{ *BaseView }

func NewRegistryProviderListView() *RegistryProviderListView {
	return &RegistryProviderListView{NewBaseView()}
}
func NewRegistryProviderCreateView() *RegistryProviderCreateView {
	return &RegistryProviderCreateView{NewBaseView()}
}
func NewRegistryProviderShowView() *RegistryProviderShowView {
	return &RegistryProviderShowView{NewBaseView()}
}
func NewRegistryProviderDeleteView() *RegistryProviderDeleteView {
	return &RegistryProviderDeleteView{NewBaseView()}
}

func (v *RegistryProviderListView) Render(items []*tfe.RegistryProvider) error {
	if v.IsJSON() {
		return v.renderer.RenderJSON(items)
	}
	headers := []string{"Name", "Registry", "ID", "Published"}
	rows := make([][]interface{}, len(items))
	for i, p := range items {
		rows[i] = []interface{}{p.Name, p.RegistryName, p.ID, p.UpdatedAt}
	}
	return v.renderer.RenderTable(headers, rows)
}

func (v *RegistryProviderCreateView) Render(p *tfe.RegistryProvider) error {
	if v.IsJSON() {
		return v.renderer.RenderJSON(p)
	}
	props := []PropertyPair{
		{Key: "Name", Value: p.Name},
		{Key: "ID", Value: p.ID},
		{Key: "Namespace", Value: p.Namespace},
		{Key: "Created", Value: p.UpdatedAt},
	}
	return v.renderer.RenderProperties(props)
}

func (v *RegistryProviderShowView) Render(p *tfe.RegistryProvider) error {
	if v.IsJSON() {
		return v.renderer.RenderJSON(p)
	}
	props := []PropertyPair{
		{Key: "Name", Value: p.Name},
		{Key: "ID", Value: p.ID},
		{Key: "Namespace", Value: p.Namespace},
		{Key: "Created", Value: p.UpdatedAt},
	}
	return v.renderer.RenderProperties(props)
}

func (v *RegistryProviderDeleteView) Render(name string) error {
	if v.IsJSON() {
		return v.renderer.RenderJSON(map[string]interface{}{"status": "Success", "name": name})
	}
	return v.renderer.RenderProperties([]PropertyPair{{Key: "Status", Value: "Success"}, {Key: "Name", Value: name}})
}
