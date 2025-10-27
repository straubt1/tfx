// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

// ProjectListView handles rendering for project list command
type ProjectListView struct {
	*BaseView
}

func NewProjectListView() *ProjectListView {
	return &ProjectListView{
		BaseView: NewBaseView(),
	}
}

// projectListOutput is a JSON-safe representation of a project for list views
type projectListOutput struct {
	Organization                string            `json:"organization,omitempty"`
	Name                        string            `json:"name"`
	ID                          string            `json:"id"`
	Description                 string            `json:"description"`
	DefaultExecutionMode        string            `json:"defaultExecutionMode"`
	AutoDestroyActivityDuration string            `json:"autoDestroyActivityDuration"`
	IsUnified                   bool              `json:"isUnified"`
	DefaultAgentPool            *string           `json:"defaultAgentPool,omitempty"`
	Tags                        map[string]string `json:"tags,omitempty"`
}

// Render renders projects for a single organization or across all organizations
// If includeOrgColumn is true, the organization column will be included in the terminal output
func (v *ProjectListView) Render(projects []*tfe.Project, includeOrgColumn bool) error {
	if v.IsJSON() {
		// Convert to JSON-safe representation (always includes organization)
		output := make([]projectListOutput, len(projects))
		for i, p := range projects {
			orgName := ""
			if p.Organization != nil {
				orgName = p.Organization.Name
			}

			// Get auto-destroy duration if specified
			duration, _ := GetIfSpecified(p.AutoDestroyActivityDuration)

			// Get default agent pool if specified
			var agentPoolName *string
			if p.DefaultAgentPool != nil {
				agentPoolName = &p.DefaultAgentPool.Name
			}

			// Extract tags
			var tags map[string]string
			if len(p.EffectiveTagBindings) > 0 {
				tags = make(map[string]string)
				for _, tag := range p.EffectiveTagBindings {
					tags[tag.Key] = tag.Value
				}
			}

			output[i] = projectListOutput{
				Organization:                orgName,
				Name:                        p.Name,
				ID:                          p.ID,
				Description:                 p.Description,
				DefaultExecutionMode:        p.DefaultExecutionMode,
				AutoDestroyActivityDuration: duration,
				IsUnified:                   p.IsUnified,
				DefaultAgentPool:            agentPoolName,
				Tags:                        tags,
			}
		}
		return v.Output().RenderJSON(output)
	}

	// Terminal mode: render as table
	var headers []string
	if includeOrgColumn {
		headers = []string{"Organization", "Name", "ID"}
	} else {
		headers = []string{"Name", "ID"}
	}

	rows := make([][]interface{}, len(projects))
	for i, p := range projects {
		if includeOrgColumn {
			orgName := ""
			if p.Organization != nil {
				orgName = p.Organization.Name
			}
			rows[i] = []interface{}{orgName, p.Name, p.ID}
		} else {
			rows[i] = []interface{}{p.Name, p.ID}
		}
	}

	return v.Output().RenderTable(headers, rows)
}
