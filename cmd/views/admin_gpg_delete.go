// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

// AdminGPGDeleteView handles rendering for admin gpg delete command
type AdminGPGDeleteView struct {
	*BaseView
}

func NewAdminGPGDeleteView() *AdminGPGDeleteView {
	return &AdminGPGDeleteView{
		BaseView: NewBaseView(),
	}
}

// deleteResult represents a successful deletion
type deleteResult struct {
	Status string `json:"status"`
}

// Render renders the deletion success message
func (v *AdminGPGDeleteView) Render() error {
	if v.IsJSON() {
		return v.Output().RenderJSON(deleteResult{Status: "success"})
	}

	// Terminal mode: render success message
	properties := []PropertyPair{
		{Key: "Status", Value: "Success"},
	}

	return v.Output().RenderProperties(properties)
}
