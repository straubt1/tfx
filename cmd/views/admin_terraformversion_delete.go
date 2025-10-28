// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

// AdminTerraformVersionDeleteView handles rendering for admin terraform-version delete command
type AdminTerraformVersionDeleteView struct {
	*BaseView
}

func NewAdminTerraformVersionDeleteView() *AdminTerraformVersionDeleteView {
	return &AdminTerraformVersionDeleteView{
		BaseView: NewBaseView(),
	}
}

// Render renders the deletion success message
func (v *AdminTerraformVersionDeleteView) Render() error {
	if v.IsJSON() {
		return v.Output().RenderJSON(deleteResult{Status: "success"})
	}

	// Terminal mode: render success message
	properties := []PropertyPair{
		{Key: "Status", Value: "Success"},
	}

	return v.Output().RenderProperties(properties)
}
