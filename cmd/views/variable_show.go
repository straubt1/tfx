// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	"fmt"

	tfe "github.com/hashicorp/go-tfe"
)

// VariableShowView handles rendering for variable show command
type VariableShowView struct {
	*BaseView
}

func NewVariableShowView() *VariableShowView {
	return &VariableShowView{
		BaseView: NewBaseView(),
	}
}

// variableShowOutput is a JSON-safe representation of a variable
type variableShowOutput struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	Value       string `json:"value"`
	Sensitive   bool   `json:"sensitive"`
	HCL         bool   `json:"hcl"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

// Render renders a single variable's details
func (v *VariableShowView) Render(variable *tfe.Variable) error {
	if variable == nil {
		return v.RenderError(fmt.Errorf("variable not found"))
	}

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
