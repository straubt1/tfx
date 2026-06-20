// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

// VariableSetVariableListView handles rendering for variable set variable list command
type VariableSetVariableListView struct {
	*BaseView
}

func NewVariableSetVariableListView() *VariableSetVariableListView {
	return &VariableSetVariableListView{
		BaseView: NewBaseView(),
	}
}

// variableSetVariableListOutput is a JSON-safe representation of a variable set variable for list views
type variableSetVariableListOutput struct {
	ID        string `json:"id"`
	Key       string `json:"key"`
	Category  string `json:"category"`
	Sensitive bool   `json:"sensitive"`
	HCL       bool   `json:"hcl"`
}

// Render renders variables for a variable set
func (v *VariableSetVariableListView) Render(variables []*tfe.VariableSetVariable) error {
	if v.IsJSON() {
		output := make([]variableSetVariableListOutput, len(variables))
		for i, variable := range variables {
			output[i] = variableSetVariableListOutput{
				ID:        variable.ID,
				Key:       variable.Key,
				Category:  string(variable.Category),
				Sensitive: variable.Sensitive,
				HCL:       variable.HCL,
			}
		}
		return v.Output().RenderJSON(output)
	}

	headers := []string{"Key", "ID", "Category", "Sensitive", "HCL"}
	rows := make([][]interface{}, len(variables))
	for i, variable := range variables {
		rows[i] = []interface{}{
			variable.Key,
			variable.ID,
			variable.Category,
			variable.Sensitive,
			variable.HCL,
		}
	}

	return v.Output().RenderTable(headers, rows)
}
