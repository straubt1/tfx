// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package view

type VariableSetDeleteView struct{ *BaseView }

func NewVariableSetDeleteView() *VariableSetDeleteView {
	return &VariableSetDeleteView{NewBaseView()}
}

func (v *VariableSetDeleteView) Render(id string) error {
	if v.IsJSON() {
		return v.Output().RenderJSON(map[string]interface{}{"status": "Success", "id": id})
	}
	return v.Output().RenderProperties([]PropertyPair{{Key: "Status", Value: "Success"}, {Key: "ID", Value: id}})
}
