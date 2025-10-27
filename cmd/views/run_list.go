// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

// RunListView handles rendering for run list command
type RunListView struct{ *BaseView }

func NewRunListView() *RunListView { return &RunListView{NewBaseView()} }

type runListOutput struct {
	ID                   string `json:"id"`
	ConfigurationVersion string `json:"configurationVersion"`
	Status               string `json:"status"`
	PlanOnly             bool   `json:"planOnly"`
	TerraformVersion     string `json:"terraformVersion"`
	Created              string `json:"created"`
	Message              string `json:"message"`
}

func (v *RunListView) Render(items []*tfe.Run) error {
	if v.IsJSON() {
		out := make([]runListOutput, len(items))
		for i, r := range items {
			cv := ""
			if r.ConfigurationVersion != nil {
				cv = r.ConfigurationVersion.ID
			}
			out[i] = runListOutput{
				ID:                   r.ID,
				ConfigurationVersion: cv,
				Status:               string(r.Status),
				PlanOnly:             r.PlanOnly,
				TerraformVersion:     r.TerraformVersion,
				Created:              FormatDateTime(r.CreatedAt),
				Message:              r.Message,
			}
		}
		return v.Output().RenderJSON(out)
	}

	headers := []string{"Id", "Configuration Version", "Created"}
	rows := make([][]interface{}, len(items))
	for i, r := range items {
		cv := ""
		if r.ConfigurationVersion != nil {
			cv = r.ConfigurationVersion.ID
		}
		rows[i] = []interface{}{r.ID, cv, FormatDateTime(r.CreatedAt)}
	}
	return v.Output().RenderTable(headers, rows)
}
