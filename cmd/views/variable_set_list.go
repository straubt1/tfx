// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

type VariableSetListView struct{ *BaseView }

func NewVariableSetListView() *VariableSetListView {
	return &VariableSetListView{NewBaseView()}
}

type variableSetParentOutput struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// parseVariableSetParent safely extracts parent info from a VariableSet.
// Returns an empty struct on any error since the Parent relation is BETA and may be absent.
func parseVariableSetParent(vs *tfe.VariableSet) (out variableSetParentOutput) {
	defer func() {
		if r := recover(); r != nil {
			out = variableSetParentOutput{}
		}
	}()
	if vs.Parent == nil {
		return
	}
	if vs.Parent.Organization != nil {
		// Organization uses Name as its primary/ID field in the jsonapi schema.
		return variableSetParentOutput{Type: "organization", ID: vs.Parent.Organization.Name}
	}
	if vs.Parent.Project != nil {
		return variableSetParentOutput{Type: "project", ID: vs.Parent.Project.ID}
	}
	return
}

type variableSetListOutput struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Global      bool                    `json:"global"`
	Priority    bool                    `json:"priority"`
	Parent      variableSetParentOutput `json:"parent"`
}

func (v *VariableSetListView) Render(items []*tfe.VariableSet) error {
	if v.IsJSON() {
		out := make([]variableSetListOutput, len(items))
		for i, vs := range items {
			out[i] = variableSetListOutput{ID: vs.ID, Name: vs.Name, Description: vs.Description, Global: vs.Global, Priority: vs.Priority, Parent: parseVariableSetParent(vs)}
		}
		return v.Output().RenderJSON(out)
	}
	headers := []string{"Name", "ID", "Global", "Priority", "Parent"}
	rows := make([][]interface{}, len(items))
	for i, vs := range items {
		parent := parseVariableSetParent(vs)
		parentDisplay := ""
		if parent.Type != "" {
			parentDisplay = parent.Type + ":" + parent.ID
		}
		rows[i] = []interface{}{vs.Name, vs.ID, vs.Global, vs.Priority, parentDisplay}
	}
	return v.Output().RenderTable(headers, rows)
}
