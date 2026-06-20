// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	"fmt"

	tfe "github.com/hashicorp/go-tfe"
)

// VariableSetVariableShowView handles rendering for variable set variable show command
type VariableSetVariableShowView struct {
	*BaseView
}

func NewVariableSetVariableShowView() *VariableSetVariableShowView {
	return &VariableSetVariableShowView{
		BaseView: NewBaseView(),
	}
}

// variableSetVariableShowOutput is a JSON-safe representation of a variable set variable
type variableSetVariableShowOutput struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	Value       string `json:"value"`
	Sensitive   bool   `json:"sensitive"`
	HCL         bool   `json:"hcl"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

// Render renders a single variable set variable's details
func (v *VariableSetVariableShowView) Render(variable *tfe.VariableSetVariable) error {
	if variable == nil {
		return v.RenderError(fmt.Errorf("variable not found"))
	}

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
