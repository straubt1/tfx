// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

// AdminTerraformVersionUpdateView handles rendering for admin terraform-version enable/disable commands
type AdminTerraformVersionUpdateView struct {
	*BaseView
}

func NewAdminTerraformVersionUpdateView() *AdminTerraformVersionUpdateView {
	return &AdminTerraformVersionUpdateView{
		BaseView: NewBaseView(),
	}
}

// updateResult represents the result of updating versions
type updateResult struct {
	Results map[string]string `json:"results"`
}

// Render renders the update results
func (v *AdminTerraformVersionUpdateView) Render(results map[string]string) error {
	if v.IsJSON() {
		return v.Output().RenderJSON(updateResult{Results: results})
	}

	// Terminal mode: render results as properties
	properties := make([]PropertyPair, 0, len(results))
	for version, status := range results {
		properties = append(properties, PropertyPair{
			Key:   version,
			Value: status,
		})
	}

	return v.Output().RenderProperties(properties)
}
