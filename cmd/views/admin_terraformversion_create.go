// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

// AdminTerraformVersionCreateView handles rendering for admin terraform-version create command
type AdminTerraformVersionCreateView struct {
	*BaseView
}

func NewAdminTerraformVersionCreateView() *AdminTerraformVersionCreateView {
	return &AdminTerraformVersionCreateView{
		BaseView: NewBaseView(),
	}
}

// Render renders the created Terraform version
func (v *AdminTerraformVersionCreateView) Render(tfv *tfe.AdminTerraformVersion) error {
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
