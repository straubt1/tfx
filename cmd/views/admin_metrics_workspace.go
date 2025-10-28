// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	"time"
)

// AdminMetricsWorkspaceView handles rendering for admin metrics workspace command
type AdminMetricsWorkspaceView struct {
	*BaseView
}

func NewAdminMetricsWorkspaceView() *AdminMetricsWorkspaceView {
	return &AdminMetricsWorkspaceView{
		BaseView: NewBaseView(),
	}
}

// MetricsWorkspace represents metrics for a single workspace
type MetricsWorkspace struct {
	Name              string `json:"name"`
	ID                string `json:"id"`
	RunCount          int    `json:"runCount"`
	RunErroredCount   int    `json:"runErroredCount"`
	RunDiscardedCount int    `json:"runDiscardedCount"`
	RunCancelledCount int    `json:"runCancelledCount"`
	PolicyCheckCount  int    `json:"policyCheckCount"`
	PoliciesPassCount int    `json:"policiesPassCount"`
	PoliciesFailCount int    `json:"policiesFailCount"`
}

// MetricsWorkspaceResult holds the complete metrics result
type MetricsWorkspaceResult struct {
	Workspaces []MetricsWorkspace `json:"workspaces"`
	Since      time.Time          `json:"since"`
	QueryTime  string             `json:"queryTime"`
}

// Render renders workspace metrics
func (v *AdminMetricsWorkspaceView) Render(result *MetricsWorkspaceResult) error {
	if v.IsJSON() {
		return v.Output().RenderJSON(result)
	}

	// Terminal mode: render as table
	headers := []string{"Name", "Total Runs", "Errored Runs", "Discarded Runs", "Cancelled Runs"}
	var rows [][]interface{}

	for _, ws := range result.Workspaces {
		// Skip workspaces with no runs
		if ws.RunCount == 0 {
			continue
		}
		rows = append(rows, []interface{}{
			ws.Name,
			ws.RunCount,
			ws.RunErroredCount,
			ws.RunDiscardedCount,
			ws.RunCancelledCount,
		})
	}

	if err := v.Output().RenderTable(headers, rows); err != nil {
		return err
	}

	// Print query time
	v.Output().MessageCommandHeader("\nMetrics Query Time: %s", result.QueryTime)

	return nil
}
