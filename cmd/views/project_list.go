// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
	"github.com/logrusorgru/aurora"
)

// ProjectListView handles rendering for project list command
type ProjectListView struct {
	*BaseView
}

func NewProjectListView(isJSON bool) *ProjectListView {
	return &ProjectListView{
		BaseView: NewBaseView(isJSON),
	}
}

// RenderAll renders all projects across organizations
func (v *ProjectListView) RenderAll(projects []*tfe.Project) error {
	// Print context message (only in terminal mode)
	v.renderer.Message("List Projects for all available Organizations")

	if v.IsJSON() {
		return v.renderer.RenderJSON(projects)
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
	// Print context message with highlighted org name
	v.renderer.Message("List Projects for Organization: %s", aurora.Green(orgName))

	for _, p := range projects {
		p.EffectiveTagBindings = make([]*tfe.EffectiveTagBinding, 0)
	}
	if v.IsJSON() {
		return v.renderer.RenderJSON(projects)
	}

	// Terminal mode: render as table
	headers := []string{"Name", "ID", "Description"}
	rows := make([][]interface{}, len(projects))

	for i, p := range projects {
		rows[i] = []interface{}{p.Name, p.ID, p.Description}
	}

	return v.renderer.RenderTable(headers, rows)
}
