// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

type RegistryModuleVersionListView struct{ *BaseView }
type RegistryModuleVersionCreateView struct{ *BaseView }
type RegistryModuleVersionDeleteView struct{ *BaseView }

func NewRegistryModuleVersionListView() *RegistryModuleVersionListView {
	return &RegistryModuleVersionListView{NewBaseView()}
}
func NewRegistryModuleVersionCreateView() *RegistryModuleVersionCreateView {
	return &RegistryModuleVersionCreateView{NewBaseView()}
}
func NewRegistryModuleVersionDeleteView() *RegistryModuleVersionDeleteView {
	return &RegistryModuleVersionDeleteView{NewBaseView()}
}

func (v *RegistryModuleVersionListView) Render(module *tfe.RegistryModule) error {
	if v.IsJSON() {
		return v.renderer.RenderJSON(module.VersionStatuses)
	}
	headers := []string{"Version", "Status"}
	rows := make([][]interface{}, len(module.VersionStatuses))
	for i, vs := range module.VersionStatuses {
		rows[i] = []interface{}{vs.Version, vs.Status}
	}
	return v.renderer.RenderTable(headers, rows)
}

func (v *RegistryModuleVersionCreateView) Render(mv *tfe.RegistryModuleVersion) error {
	if v.IsJSON() {
		return v.renderer.RenderJSON(mv)
	}
	props := []PropertyPair{
		{Key: "ID", Value: mv.ID},
		{Key: "Module", Value: mv.RegistryModule.Name},
		{Key: "Created", Value: mv.CreatedAt},
	}
	return v.renderer.RenderProperties(props)
}

func (v *RegistryModuleVersionDeleteView) Render(name string) error {
	if v.IsJSON() {
		return v.renderer.RenderJSON(map[string]interface{}{"status": "Success", "name": name})
	}
	return v.renderer.RenderProperties([]PropertyPair{{Key: "Status", Value: "Success"}, {Key: "Name", Value: name}})
}
