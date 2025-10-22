// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

type StateVersionShowView struct{ *BaseView }

func NewStateVersionShowView() *StateVersionShowView { return &StateVersionShowView{NewBaseView()} }

func (v *StateVersionShowView) Render(sv *tfe.StateVersion) error {
	if v.IsJSON() {
		outputs := make(map[string]interface{}, len(sv.Outputs))
		for _, o := range sv.Outputs {
			outputs[o.Name] = o.Value
		}
		var runID interface{}
		if sv.Run != nil {
			runID = sv.Run.ID
		} else {
			runID = nil
		}
		return v.renderer.RenderJSON(map[string]interface{}{
			"id":               sv.ID,
			"created":          FormatDateTime(sv.CreatedAt),
			"terraformVersion": sv.TerraformVersion,
			"serial":           sv.Serial,
			"stateVersion":     sv.StateVersion,
			"runId":            runID,
			"outputs":          outputs,
		})
	}
	var runID string
	if sv.Run != nil {
		runID = sv.Run.ID
	}
	props := []PropertyPair{
		{Key: "ID", Value: sv.ID},
		{Key: "Created", Value: FormatDateTime(sv.CreatedAt)},
		{Key: "Terraform Version", Value: sv.TerraformVersion},
		{Key: "Serial", Value: sv.Serial},
		{Key: "State Version", Value: sv.StateVersion},
		{Key: "Run Id", Value: runID},
	}
	if err := v.renderer.RenderProperties(props); err != nil {
		return err
	}
	// Render outputs (if any)
	tags := make([]PropertyPair, 0, len(sv.Outputs))
	for _, o := range sv.Outputs {
		tags = append(tags, PropertyPair{Key: o.Name, Value: o.Value})
	}
	return v.renderer.RenderTags("Outputs", tags)
}
