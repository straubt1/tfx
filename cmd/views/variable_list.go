// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

// VariableListView handles rendering for variable list command
type VariableListView struct {
	*BaseView
}

func NewVariableListView() *VariableListView {
	return &VariableListView{
		BaseView: NewBaseView(),
	}
}

// variableListOutput is a JSON-safe representation of a variable for list views
type variableListOutput struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	Value       string `json:"value"`
	Sensitive   bool   `json:"sensitive"`
	HCL         bool   `json:"hcl"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

// Render renders variables for a workspace
func (v *VariableListView) Render(workspaceName string, variables []*tfe.Variable) error {
	if v.IsJSON() {
		// Convert to JSON-safe representation
		output := make([]variableListOutput, len(variables))
		for i, variable := range variables {
			output[i] = variableListOutput{
				ID:          variable.ID,
				Key:         variable.Key,
				Value:       variable.Value,
				Sensitive:   variable.Sensitive,
				HCL:         variable.HCL,
				Category:    string(variable.Category),
				Description: variable.Description,
			}
		}
		return v.renderer.RenderJSON(output)
	}

	// Terminal mode: render as table
	headers := []string{"Id", "Key", "Value", "Sensitive", "HCL", "Category", "Description"}
	rows := make([][]interface{}, len(variables))

	for i, variable := range variables {
		// Limit value to 20 characters for display
		value := variable.Value
		if len(value) > 20 {
			value = value[:20] + "..."
		}

		// Limit description to 20 characters for display
		description := variable.Description
		if len(description) > 20 {
			description = description[:20] + "..."
		}

		rows[i] = []interface{}{
			variable.ID,
			variable.Key,
			value,
			variable.Sensitive,
			variable.HCL,
			variable.Category,
			description,
		}
	}

	return v.renderer.RenderTable(headers, rows)
}
