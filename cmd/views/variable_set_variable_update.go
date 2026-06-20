// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

// VariableSetVariableUpdateView handles rendering for variable set variable update command
type VariableSetVariableUpdateView struct {
	*BaseView
}

func NewVariableSetVariableUpdateView() *VariableSetVariableUpdateView {
	return &VariableSetVariableUpdateView{
		BaseView: NewBaseView(),
	}
}

// Render renders an updated variable set variable's details
func (v *VariableSetVariableUpdateView) Render(variable *tfe.VariableSetVariable) error {
	if v.IsJSON() {
		output := variableSetVariableShowOutput{
			ID:          variable.ID,
			Key:         variable.Key,
			Value:       variable.Value,
			Sensitive:   variable.Sensitive,
			HCL:         variable.HCL,
			Category:    string(variable.Category),
			Description: variable.Description,
		}
		return v.Output().RenderJSON(output)
	}

	properties := []PropertyPair{
		{Key: "ID", Value: variable.ID},
		{Key: "Key", Value: variable.Key},
		{Key: "Value", Value: variable.Value},
		{Key: "Sensitive", Value: variable.Sensitive},
		{Key: "HCL", Value: variable.HCL},
		{Key: "Category", Value: variable.Category},
		{Key: "Description", Value: variable.Description},
	}

	return v.Output().RenderProperties(properties)
}
