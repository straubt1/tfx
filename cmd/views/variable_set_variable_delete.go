// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

// VariableSetVariableDeleteView handles rendering for variable set variable delete command
type VariableSetVariableDeleteView struct {
	*BaseView
}

func NewVariableSetVariableDeleteView() *VariableSetVariableDeleteView {
	return &VariableSetVariableDeleteView{
		BaseView: NewBaseView(),
	}
}

// variableSetVariableDeleteOutput is a JSON-safe representation of a delete operation
type variableSetVariableDeleteOutput struct {
	Status string `json:"status"`
	Key    string `json:"key"`
}

// Render renders a successful delete operation
func (v *VariableSetVariableDeleteView) Render(key string) error {
	if v.IsJSON() {
		output := variableSetVariableDeleteOutput{
			Status: "Success",
			Key:    key,
		}
		return v.Output().RenderJSON(output)
	}

	properties := []PropertyPair{
		{Key: "Status", Value: "Success"},
		{Key: "Key", Value: key},
	}

	return v.Output().RenderProperties(properties)
}
