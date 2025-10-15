// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

// OrganizationListView handles rendering for organization list command
type OrganizationListView struct {
	*BaseView
}

func NewOrganizationListView() *OrganizationListView {
	return &OrganizationListView{
		BaseView: NewBaseView(),
	}
}

// organizationListOutput is a JSON-safe representation of an organization for list views
type organizationListOutput struct {
	Name                   string `json:"name"`
	Email                  string `json:"email"`
	ExternalID             string `json:"externalId"`
	CollaboratorAuthPolicy string `json:"collaboratorAuthPolicy"`
	CostEstimationEnabled  bool   `json:"costEstimationEnabled"`
	SAMLEnabled            bool   `json:"samlEnabled"`
	TwoFactorConformant    bool   `json:"twoFactorConformant"`
	DefaultExecutionMode   string `json:"defaultExecutionMode"`
	IsUnified              bool   `json:"isUnified"`
}

// Render renders organizations list
func (v *OrganizationListView) Render(orgs []*tfe.Organization) error {
	if v.IsJSON() {
		// Convert to JSON-safe representation
		output := make([]organizationListOutput, len(orgs))
		for i, org := range orgs {
			output[i] = organizationListOutput{
				Name:                   org.Name,
				Email:                  org.Email,
				ExternalID:             org.ExternalID,
				CollaboratorAuthPolicy: string(org.CollaboratorAuthPolicy),
				CostEstimationEnabled:  org.CostEstimationEnabled,
				SAMLEnabled:            org.SAMLEnabled,
				TwoFactorConformant:    org.TwoFactorConformant,
				DefaultExecutionMode:   org.DefaultExecutionMode,
				IsUnified:              org.IsUnified,
			}
		}
		return v.renderer.RenderJSON(output)
	}

	// Terminal mode: render as table
	headers := []string{"Name", "ID", "Email"}
	rows := make([][]interface{}, len(orgs))

	for i, org := range orgs {
		rows[i] = []interface{}{org.Name, org.ExternalID, org.Email}
	}

	return v.renderer.RenderTable(headers, rows)
}
