// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	"fmt"

	tfe "github.com/hashicorp/go-tfe"
)

// OrganizationShowView handles rendering for organization show command
type OrganizationShowView struct {
	*BaseView
}

func NewOrganizationShowView() *OrganizationShowView {
	return &OrganizationShowView{
		BaseView: NewBaseView(),
	}
}

// organizationShowOutput is a JSON-safe representation of an organization
type organizationShowOutput struct {
	Name                   string                         `json:"name"`
	Email                  string                         `json:"email"`
	ExternalID             string                         `json:"externalId"`
	CreatedAt              string                         `json:"createdAt"`
	CollaboratorAuthPolicy string                         `json:"collaboratorAuthPolicy"`
	CostEstimationEnabled  bool                           `json:"costEstimationEnabled"`
	OwnersTeamSAMLRoleID   string                         `json:"ownersTeamSamlRoleId"`
	SAMLEnabled            bool                           `json:"samlEnabled"`
	SessionRemember        int                            `json:"sessionRemember"`
	SessionTimeout         int                            `json:"sessionTimeout"`
	TwoFactorConformant    bool                           `json:"twoFactorConformant"`
	TrialExpiresAt         string                         `json:"trialExpiresAt"`
	DefaultExecutionMode   string                         `json:"defaultExecutionMode"`
	IsUnified              bool                           `json:"isUnified"`
	Permissions            *organizationPermissionsOutput `json:"permissions,omitempty"`
}

type organizationPermissionsOutput struct {
	CanCreateTeam               bool `json:"canCreateTeam"`
	CanCreateWorkspace          bool `json:"canCreateWorkspace"`
	CanCreateWorkspaceMigration bool `json:"canCreateWorkspaceMigration"`
	CanDestroy                  bool `json:"canDestroy"`
	CanManageRunTasks           bool `json:"canManageRunTasks"`
	CanTraverse                 bool `json:"canTraverse"`
	CanUpdate                   bool `json:"canUpdate"`
}

// Render renders a single organization's details
func (v *OrganizationShowView) Render(org *tfe.Organization) error {
	if org == nil {
		return v.RenderError(fmt.Errorf("organization not found"))
	}

	if v.IsJSON() {
		// JSON mode: convert to JSON-safe structure
		output := organizationShowOutput{
			Name:                   org.Name,
			Email:                  org.Email,
			ExternalID:             org.ExternalID,
			CreatedAt:              org.CreatedAt.String(),
			CollaboratorAuthPolicy: string(org.CollaboratorAuthPolicy),
			CostEstimationEnabled:  org.CostEstimationEnabled,
			OwnersTeamSAMLRoleID:   org.OwnersTeamSAMLRoleID,
			SAMLEnabled:            org.SAMLEnabled,
			SessionRemember:        org.SessionRemember,
			SessionTimeout:         org.SessionTimeout,
			TwoFactorConformant:    org.TwoFactorConformant,
			TrialExpiresAt:         org.TrialExpiresAt.String(),
			DefaultExecutionMode:   org.DefaultExecutionMode,
			IsUnified:              org.IsUnified,
		}

		// Add permissions if available
		if org.Permissions != nil {
			output.Permissions = &organizationPermissionsOutput{
				CanCreateTeam:               org.Permissions.CanCreateTeam,
				CanCreateWorkspace:          org.Permissions.CanCreateWorkspace,
				CanCreateWorkspaceMigration: org.Permissions.CanCreateWorkspaceMigration,
				CanDestroy:                  org.Permissions.CanDestroy,
				CanManageRunTasks:           org.Permissions.CanManageRunTasks,
				CanTraverse:                 org.Permissions.CanTraverse,
				CanUpdate:                   org.Permissions.CanUpdate,
			}
		}

		return v.renderer.RenderJSON(output)
	}

	// Terminal mode: render key fields in order
	properties := []PropertyPair{
		{Key: "Name", Value: org.Name},
		{Key: "ID", Value: org.ExternalID},
		{Key: "Email", Value: org.Email},
		{Key: "Created At", Value: org.CreatedAt},
		{Key: "Collaborator Auth Policy", Value: org.CollaboratorAuthPolicy},
		{Key: "Cost Estimation Enabled", Value: org.CostEstimationEnabled},
		{Key: "Owners Team SAML Role ID", Value: org.OwnersTeamSAMLRoleID},
		{Key: "SAML Enabled", Value: org.SAMLEnabled},
		{Key: "Session Remember Minutes", Value: org.SessionRemember},
		{Key: "Session Timeout Minutes", Value: org.SessionTimeout},
		{Key: "Two Factor Conformant", Value: org.TwoFactorConformant},
		{Key: "Trial Expires At", Value: org.TrialExpiresAt},
		{Key: "Default Execution Mode", Value: org.DefaultExecutionMode},
		{Key: "Is Unified", Value: org.IsUnified},
	}

	err := v.renderer.RenderProperties(properties)
	if err != nil {
		return err
	}

	// Add permissions section if available
	if org.Permissions != nil {
		permissions := []PropertyPair{
			{Key: "Can Create Team", Value: org.Permissions.CanCreateTeam},
			{Key: "Can Create Workspace", Value: org.Permissions.CanCreateWorkspace},
			{Key: "Can Create Workspace Migration", Value: org.Permissions.CanCreateWorkspaceMigration},
			{Key: "Can Destroy", Value: org.Permissions.CanDestroy},
			{Key: "Can Manage Run Tasks", Value: org.Permissions.CanManageRunTasks},
			{Key: "Can Traverse", Value: org.Permissions.CanTraverse},
			{Key: "Can Update", Value: org.Permissions.CanUpdate},
		}
		err = v.renderer.RenderTags("Permissions", permissions)
		if err != nil {
			return err
		}
	}

	return nil
}
