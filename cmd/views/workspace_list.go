// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	"time"

	tfe "github.com/hashicorp/go-tfe"
)

// FormatDateTime formats a time.Time to a consistent string format
func FormatDateTime(t time.Time) string {
	return t.Format("Mon Jan _2 15:04 2006")
}

// WorkspaceListView handles rendering for workspace list command
type WorkspaceListView struct {
	*BaseView
}

func NewWorkspaceListView() *WorkspaceListView {
	return &WorkspaceListView{
		BaseView: NewBaseView(),
	}
}

// workspaceListOutput is a JSON-safe representation of a workspace for list views
type workspaceListOutput struct {
	Organization      string `json:"organization,omitempty"`
	Name              string `json:"name"`
	ID                string `json:"id"`
	ResourceCount     int    `json:"resourceCount"`
	CurrentRunCreated string `json:"currentRunCreated,omitempty"`
	CurrentRunStatus  string `json:"currentRunStatus,omitempty"`
	Repository        string `json:"repository,omitempty"`
	Locked            bool   `json:"locked"`
}

// RenderAll renders all workspaces across organizations
func (v *WorkspaceListView) RenderAll(workspaces []*tfe.Workspace) error {
	if v.IsJSON() {
		// Convert to JSON-safe representation
		output := make([]workspaceListOutput, len(workspaces))
		for i, w := range workspaces {
			orgName := ""
			if w.Organization != nil {
				orgName = w.Organization.Name
			}

			currentRunCreated := ""
			currentRunStatus := ""
			if w.CurrentRun != nil {
				currentRunCreated = FormatDateTime(w.CurrentRun.CreatedAt)
				currentRunStatus = string(w.CurrentRun.Status)
			}

			repository := ""
			if w.VCSRepo != nil {
				repository = w.VCSRepo.DisplayIdentifier
			}

			output[i] = workspaceListOutput{
				Organization:      orgName,
				Name:              w.Name,
				ID:                w.ID,
				ResourceCount:     w.ResourceCount,
				CurrentRunCreated: currentRunCreated,
				CurrentRunStatus:  currentRunStatus,
				Repository:        repository,
				Locked:            w.Locked,
			}
		}
		return v.renderer.RenderJSON(output)
	}

	// Terminal mode: render as table
	headers := []string{"Organization", "Name", "Id", "Status"}
	rows := make([][]interface{}, len(workspaces))

	for i, w := range workspaces {
		orgName := ""
		if w.Organization != nil {
			orgName = w.Organization.Name
		}

		currentRunStatus := ""
		if w.CurrentRun != nil {
			currentRunStatus = string(w.CurrentRun.Status)
		}

		rows[i] = []interface{}{orgName, w.Name, w.ID, currentRunStatus}
	}

	return v.renderer.RenderTable(headers, rows)
}

// Render renders workspaces for a single organization
func (v *WorkspaceListView) Render(orgName string, workspaces []*tfe.Workspace) error {
	if v.IsJSON() {
		// Convert to JSON-safe representation
		output := make([]workspaceListOutput, len(workspaces))
		for i, w := range workspaces {
			currentRunCreated := ""
			currentRunStatus := ""
			if w.CurrentRun != nil {
				currentRunCreated = FormatDateTime(w.CurrentRun.CreatedAt)
				currentRunStatus = string(w.CurrentRun.Status)
			}

			repository := ""
			if w.VCSRepo != nil {
				repository = w.VCSRepo.DisplayIdentifier
			}

			output[i] = workspaceListOutput{
				Name:              w.Name,
				ID:                w.ID,
				ResourceCount:     w.ResourceCount,
				CurrentRunCreated: currentRunCreated,
				CurrentRunStatus:  currentRunStatus,
				Repository:        repository,
				Locked:            w.Locked,
			}
		}
		return v.renderer.RenderJSON(output)
	}

	// Terminal mode: render as table
	headers := []string{"Name", "Id", "Status"}
	rows := make([][]interface{}, len(workspaces))

	for i, w := range workspaces {
		currentRunStatus := ""
		if w.CurrentRun != nil {
			currentRunStatus = string(w.CurrentRun.Status)
		}

		rows[i] = []interface{}{w.Name, w.ID, currentRunStatus}
	}

	return v.renderer.RenderTable(headers, rows)
}
