// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

// VariableDeleteView handles rendering for variable delete command
type VariableDeleteView struct {
	*BaseView
}

func NewVariableDeleteView() *VariableDeleteView {
	return &VariableDeleteView{
		BaseView: NewBaseView(),
	}
}

// variableDeleteOutput is a JSON-safe representation of a delete operation
type variableDeleteOutput struct {
	Status string `json:"status"`
	Key    string `json:"key"`
}

// Render renders a successful delete operation
func (v *VariableDeleteView) Render(key string) error {
	if v.IsJSON() {
		output := variableDeleteOutput{
			Status: "Success",
			Key:    key,
		}
		return v.renderer.RenderJSON(output)
	}

	// Terminal mode: render status
	properties := []PropertyPair{
		{Key: "Status", Value: "Success"},
		{Key: "Key", Value: key},
	}

	return v.renderer.RenderProperties(properties)
}
