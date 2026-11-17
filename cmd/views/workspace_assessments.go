// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

// WorkspaceAssessmentsResult is the outcome of enabling/disabling health assessments on one workspace
type WorkspaceAssessmentsResult struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type WorkspaceAssessmentsView struct{ *BaseView }

func NewWorkspaceAssessmentsView() *WorkspaceAssessmentsView {
	return &WorkspaceAssessmentsView{NewBaseView()}
}

// Render prints the per-workspace results as a table, or a JSON array in JSON mode
func (v *WorkspaceAssessmentsView) Render(results []WorkspaceAssessmentsResult) error {
	if v.IsJSON() {
		return v.Output().RenderJSON(results)
	}
	headers := []string{"Workspace", "Status"}
	rows := make([][]interface{}, len(results))
	for i, r := range results {
		rows[i] = []interface{}{r.Name, r.Status}
	}
	return v.Output().RenderTable(headers, rows)
}
