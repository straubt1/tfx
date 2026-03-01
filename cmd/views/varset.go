// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

type VariableSetListView struct{ *BaseView }
type VariableSetShowView struct{ *BaseView }
type VariableSetCreateView struct{ *BaseView }
type VariableSetDeleteView struct{ *BaseView }

func NewVariableSetListView() *VariableSetListView {
	return &VariableSetListView{NewBaseView()}
}
func NewVariableSetShowView() *VariableSetShowView {
	return &VariableSetShowView{NewBaseView()}
}
func NewVariableSetCreateView() *VariableSetCreateView {
	return &VariableSetCreateView{NewBaseView()}
}
func NewVariableSetDeleteView() *VariableSetDeleteView {
	return &VariableSetDeleteView{NewBaseView()}
}

type variableSetListOutput struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Global      bool   `json:"global"`
	Priority    bool   `json:"priority"`
}

func (v *VariableSetListView) Render(items []*tfe.VariableSet) error {
	if v.IsJSON() {
		out := make([]variableSetListOutput, len(items))
		for i, vs := range items {
			out[i] = variableSetListOutput{vs.ID, vs.Name, vs.Description, vs.Global, vs.Priority}
		}
		return v.Output().RenderJSON(out)
	}
	headers := []string{"Name", "ID", "Global", "Priority", "Description"}
	rows := make([][]interface{}, len(items))
	for i, vs := range items {
		rows[i] = []interface{}{vs.Name, vs.ID, vs.Global, vs.Priority, vs.Description}
	}
	return v.Output().RenderTable(headers, rows)
}

func (v *VariableSetShowView) Render(vs *tfe.VariableSet) error {
	if v.IsJSON() {
		return v.Output().RenderJSON(vs)
	}

	props := []PropertyPair{
		{Key: "ID", Value: vs.ID},
		{Key: "Name", Value: vs.Name},
		{Key: "Description", Value: vs.Description},
		{Key: "Global", Value: vs.Global},
		{Key: "Priority", Value: vs.Priority},
	}
	if err := v.Output().RenderProperties(props); err != nil {
		return err
	}

	// Workspaces
	workspaceNames := make([]PropertyPair, len(vs.Workspaces))
	for i, ws := range vs.Workspaces {
		workspaceNames[i] = PropertyPair{Key: ws.Name, Value: ws.ID}
	}
	if err := v.Output().RenderTags("Workspaces", workspaceNames); err != nil {
		return err
	}

	// Projects
	projectNames := make([]PropertyPair, len(vs.Projects))
	for i, p := range vs.Projects {
		projectNames[i] = PropertyPair{Key: p.Name, Value: p.ID}
	}
	if err := v.Output().RenderTags("Projects", projectNames); err != nil {
		return err
	}

	// Variables
	varPairs := make([]PropertyPair, len(vs.Variables))
	for i, vv := range vs.Variables {
		varPairs[i] = PropertyPair{Key: vv.Key, Value: vv.ID}
	}
	return v.Output().RenderTags("Variables", varPairs)
}

func (v *VariableSetCreateView) Render(vs *tfe.VariableSet) error {
	if v.IsJSON() {
		return v.Output().RenderJSON(vs)
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

func (v *VariableSetDeleteView) Render(id string) error {
	if v.IsJSON() {
		return v.Output().RenderJSON(map[string]interface{}{"status": "Success", "id": id})
	}
	return v.Output().RenderProperties([]PropertyPair{{Key: "Status", Value: "Success"}, {Key: "ID", Value: id}})
}
