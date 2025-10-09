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
	AutoDestroyActivityDuration string            `json:"autoDestroyActivityDuration,omitempty"`
	IsUnified                   bool              `json:"isUnified"`
	DefaultAgentPool            *string           `json:"defaultAgentPool,omitempty"`
	Tags                        map[string]string `json:"tags,omitempty"`
}

// RenderAll renders all projects across organizations
func (v *ProjectListView) RenderAll(projects []*tfe.Project) error {
	if v.IsJSON() {
		// Convert to JSON-safe representation
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
		return v.renderer.RenderJSON(output)
	}

	// Terminal mode: render as table
	headers := []string{"Organization", "Name", "ID", "Description"}
	rows := make([][]interface{}, len(projects))

	for i, p := range projects {
		orgName := ""
		if p.Organization != nil {
			orgName = p.Organization.Name
		}
		rows[i] = []interface{}{orgName, p.Name, p.ID, p.Description}
	}

	return v.renderer.RenderTable(headers, rows)
}

// Render renders projects for a single organization
func (v *ProjectListView) Render(orgName string, projects []*tfe.Project) error {
	if v.IsJSON() {
		// Convert to JSON-safe representation
		output := make([]projectListOutput, len(projects))
		for i, p := range projects {
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
		return v.renderer.RenderJSON(output)
	}

	// Terminal mode: render as table
	headers := []string{"Name", "ID", "Description"}
	rows := make([][]interface{}, len(projects))

	for i, p := range projects {
		rows[i] = []interface{}{p.Name, p.ID, p.Description}
	}

	return v.renderer.RenderTable(headers, rows)
}
