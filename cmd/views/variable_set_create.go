// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

type VariableSetCreateView struct{ *BaseView }

func NewVariableSetCreateView() *VariableSetCreateView {
	return &VariableSetCreateView{NewBaseView()}
}

func (v *VariableSetCreateView) Render(vs *tfe.VariableSet) error {
	if v.IsJSON() {
		return v.Output().RenderJSON(variableSetListOutput{ID: vs.ID, Name: vs.Name, Description: vs.Description, Global: vs.Global, Priority: vs.Priority})
	}
	props := []PropertyPair{
		{Key: "ID", Value: vs.ID},
		{Key: "Name", Value: vs.Name},
		{Key: "Description", Value: vs.Description},
		{Key: "Global", Value: vs.Global},
		{Key: "Priority", Value: vs.Priority},
	}
	return v.Output().RenderProperties(props)
}
