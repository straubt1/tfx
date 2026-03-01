// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

// AdminTerraformVersionShowView handles rendering for admin terraform-version show command
type AdminTerraformVersionShowView struct {
	*BaseView
}

func NewAdminTerraformVersionShowView() *AdminTerraformVersionShowView {
	return &AdminTerraformVersionShowView{
		BaseView: NewBaseView(),
	}
}

// terraformVersionShowOutput is a JSON-safe representation of a Terraform version
type terraformVersionShowOutput struct {
	Version string `json:"version"`
	ID      string `json:"id"`
	URL     string `json:"url"`
	SHA     string `json:"sha"`
	Enabled bool   `json:"enabled"`
	Beta    bool   `json:"beta"`
}

// Render renders a single Terraform version's details
func (v *AdminTerraformVersionShowView) Render(tfv *tfe.AdminTerraformVersion) error {
	if v.IsJSON() {
		output := terraformVersionShowOutput{
			Version: tfv.Version,
			ID:      tfv.ID,
			URL:     tfv.URL,
			SHA:     tfv.Sha,
			Enabled: tfv.Enabled,
			Beta:    tfv.Beta,
		}
		return v.Output().RenderJSON(output)
	}

	// Terminal mode: render key fields
	properties := []PropertyPair{
		{Key: "Version", Value: tfv.Version},
		{Key: "ID", Value: tfv.ID},
		{Key: "URL", Value: tfv.URL},
		{Key: "SHA", Value: tfv.Sha},
		{Key: "Enabled", Value: tfv.Enabled},
		{Key: "Beta", Value: tfv.Beta},
	}

	return v.Output().RenderProperties(properties)
}
