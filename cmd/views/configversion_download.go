// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

type ConfigVersionDownloadView struct{ *BaseView }

func NewConfigVersionDownloadView() *ConfigVersionDownloadView {
	return &ConfigVersionDownloadView{NewBaseView()}
}

func (v *ConfigVersionDownloadView) Render(directory string) error {
	if v.IsJSON() {
		return v.renderer.RenderJSON(map[string]string{"status": "Success", "directory": directory})
	}
	props := []PropertyPair{
		{Key: "Status", Value: "Success"},
		{Key: "Directory", Value: directory},
	}
	return v.renderer.RenderProperties(props)
}
