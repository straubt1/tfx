// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

type RegistryModuleListView struct{ *BaseView }
type RegistryModuleCreateView struct{ *BaseView }
type RegistryModuleShowView struct{ *BaseView }
type RegistryModuleDeleteView struct{ *BaseView }

func NewRegistryModuleListView() *RegistryModuleListView {
	return &RegistryModuleListView{NewBaseView()}
}
func NewRegistryModuleCreateView() *RegistryModuleCreateView {
	return &RegistryModuleCreateView{NewBaseView()}
}
func NewRegistryModuleShowView() *RegistryModuleShowView {
	return &RegistryModuleShowView{NewBaseView()}
}
func NewRegistryModuleDeleteView() *RegistryModuleDeleteView {
	return &RegistryModuleDeleteView{NewBaseView()}
}

type registryModuleListOutput struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
	ID       string `json:"id"`
	Status   string `json:"status"`
	Versions int    `json:"versions"`
}

func (v *RegistryModuleListView) Render(items []*tfe.RegistryModule) error {
	if v.IsJSON() {
		out := make([]registryModuleListOutput, len(items))
		for i, m := range items {
			out[i] = registryModuleListOutput{m.Name, m.Provider, m.ID, string(m.Status), len(m.VersionStatuses)}
		}
		return v.Output().RenderJSON(out)
	}
	headers := []string{"Name", "Provider", "ID", "Status", "Versions"}
	rows := make([][]interface{}, len(items))
	for i, m := range items {
		rows[i] = []interface{}{m.Name, m.Provider, m.ID, string(m.Status), len(m.VersionStatuses)}
	}
	return v.Output().RenderTable(headers, rows)
}

func (v *RegistryModuleCreateView) Render(m *tfe.RegistryModule) error {
	if v.IsJSON() {
		return v.Output().RenderJSON(m)
	}
	props := []PropertyPair{
		{Key: "Name", Value: m.Name},
		{Key: "Provider", Value: m.Provider},
		{Key: "ID", Value: m.ID},
		{Key: "Namespace", Value: m.Namespace},
	}
	return v.Output().RenderProperties(props)
}

func (v *RegistryModuleShowView) Render(m *tfe.RegistryModule) error {
	if v.IsJSON() {
		return v.Output().RenderJSON(m)
	}
	props := []PropertyPair{
		{Key: "ID", Value: m.ID},
		{Key: "Status", Value: m.Status},
		{Key: "Created", Value: m.CreatedAt},
		{Key: "Updated", Value: m.UpdatedAt},
		{Key: "Versions", Value: len(m.VersionStatuses)},
	}
	if len(m.VersionStatuses) > 0 {
		props = append(props, PropertyPair{Key: "Latest Version", Value: m.VersionStatuses[0].Version})
	}
	return v.Output().RenderProperties(props)
}

func (v *RegistryModuleDeleteView) Render(name string) error {
	if v.IsJSON() {
		return v.Output().RenderJSON(map[string]interface{}{"status": "Success", "name": name})
	}
	return v.Output().RenderProperties([]PropertyPair{{Key: "Status", Value: "Success"}, {Key: "Name", Value: name}})
}
