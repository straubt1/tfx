// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

// AdminTerraformVersionListView handles rendering for admin terraform-version list command
type AdminTerraformVersionListView struct {
	*BaseView
}

func NewAdminTerraformVersionListView() *AdminTerraformVersionListView {
	return &AdminTerraformVersionListView{
		BaseView: NewBaseView(),
	}
}

// terraformVersionListOutput is a JSON-safe representation of a Terraform version for list views
type terraformVersionListOutput struct {
	Version    string `json:"version"`
	ID         string `json:"id"`
	Enabled    bool   `json:"enabled"`
	Official   bool   `json:"official"`
	Usage      int    `json:"usage"`
	Deprecated bool   `json:"deprecated"`
}

// Render renders Terraform versions list
func (v *AdminTerraformVersionListView) Render(versions []*tfe.AdminTerraformVersion) error {
	if v.IsJSON() {
		output := make([]terraformVersionListOutput, len(versions))
		for i, tfv := range versions {
			output[i] = terraformVersionListOutput{
				Version:    tfv.Version,
				ID:         tfv.ID,
				Enabled:    tfv.Enabled,
				Official:   tfv.Official,
				Usage:      tfv.Usage,
				Deprecated: tfv.Deprecated,
			}
		}
		return v.Output().RenderJSON(output)
	}

	// Terminal mode: render as table
	headers := []string{"Version", "ID", "Enabled", "Official", "Usage", "Deprecated"}
	rows := make([][]interface{}, len(versions))
	for i, tfv := range versions {
		rows[i] = []interface{}{
			tfv.Version,
			tfv.ID,
			tfv.Enabled,
			tfv.Official,
			tfv.Usage,
			tfv.Deprecated,
		}
	}

	return v.Output().RenderTable(headers, rows)
}
