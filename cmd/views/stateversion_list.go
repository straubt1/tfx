// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

type StateVersionListView struct{ *BaseView }

func NewStateVersionListView() *StateVersionListView { return &StateVersionListView{NewBaseView()} }

type stateVersionRow struct {
	ID               string `json:"id"`
	TerraformVersion string `json:"terraformVersion"`
	Serial           int64  `json:"serial"`
	RunID            string `json:"runId"`
	Created          string `json:"created"`
}

func (v *StateVersionListView) Render(items []*tfe.StateVersion) error {
	if v.IsJSON() {
		out := make([]stateVersionRow, len(items))
		for i, sv := range items {
			runID := ""
			if sv.Run != nil {
				runID = sv.Run.ID
			}
			out[i] = stateVersionRow{
				ID:               sv.ID,
				TerraformVersion: sv.TerraformVersion,
				Serial:           sv.Serial,
				RunID:            runID,
				Created:          FormatDateTime(sv.CreatedAt),
			}
		}
		return v.Output().RenderJSON(out)
	}
	headers := []string{"Id", "Terraform Version", "Serial", "Run Id", "Created"}
	rows := make([][]interface{}, len(items))
	for i, sv := range items {
		runID := ""
		if sv.Run != nil {
			runID = sv.Run.ID
		}
		rows[i] = []interface{}{sv.ID, sv.TerraformVersion, sv.Serial, runID, FormatDateTime(sv.CreatedAt)}
	}
	return v.Output().RenderTable(headers, rows)
}
