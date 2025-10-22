// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

type StateVersionDownloadView struct{ *BaseView }

func NewStateVersionDownloadView() *StateVersionDownloadView {
	return &StateVersionDownloadView{NewBaseView()}
}

func (v *StateVersionDownloadView) Render(filename string) error {
	if v.IsJSON() {
		return v.renderer.RenderJSON(map[string]string{"file": filename, "status": "Success"})
	}
	props := []PropertyPair{
		{Key: "Status", Value: "Success"},
		{Key: "File", Value: filename},
	}
	return v.renderer.RenderProperties(props)
}
