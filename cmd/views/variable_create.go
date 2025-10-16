// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

// VariableCreateView handles rendering for variable create command
type VariableCreateView struct {
	*BaseView
}

func NewVariableCreateView() *VariableCreateView {
	return &VariableCreateView{
		BaseView: NewBaseView(),
	}
}

// Render renders a created variable's details
func (v *VariableCreateView) Render(variable *tfe.Variable) error {
	if v.IsJSON() {
		// JSON mode: convert to JSON-safe structure
		output := variableShowOutput{
			ID:          variable.ID,
			Key:         variable.Key,
			Value:       variable.Value,
			Sensitive:   variable.Sensitive,
			HCL:         variable.HCL,
			Category:    string(variable.Category),
			Description: variable.Description,
		}
		return v.renderer.RenderJSON(output)
	}

	// Terminal mode: render key fields in order
	properties := []PropertyPair{
		{Key: "ID", Value: variable.ID},
		{Key: "Key", Value: variable.Key},
		{Key: "Value", Value: variable.Value},
		{Key: "Sensitive", Value: variable.Sensitive},
		{Key: "HCL", Value: variable.HCL},
		{Key: "Category", Value: variable.Category},
		{Key: "Description", Value: variable.Description},
	}

	return v.renderer.RenderProperties(properties)
}
