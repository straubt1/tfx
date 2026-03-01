// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

type VariableSetShowView struct{ *BaseView }

func NewVariableSetShowView() *VariableSetShowView {
	return &VariableSetShowView{NewBaseView()}
}

type variableSetRefOutput struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type variableSetVarOutput struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type variableSetShowOutput struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Global      bool                    `json:"global"`
	Priority    bool                    `json:"priority"`
	Parent      variableSetParentOutput `json:"parent"`
	Workspaces  []variableSetRefOutput  `json:"workspaces"`
	Projects    []variableSetRefOutput  `json:"projects"`
	Variables   []variableSetVarOutput  `json:"variables"`
}

func (v *VariableSetShowView) Render(vs *tfe.VariableSet) error {
	if v.IsJSON() {
		out := variableSetShowOutput{
			ID:          vs.ID,
			Name:        vs.Name,
			Description: vs.Description,
			Global:      vs.Global,
			Priority:    vs.Priority,
			Parent:      parseVariableSetParent(vs),
			Workspaces:  make([]variableSetRefOutput, len(vs.Workspaces)),
			Projects:    make([]variableSetRefOutput, len(vs.Projects)),
			Variables:   make([]variableSetVarOutput, len(vs.Variables)),
		}
		for i, ws := range vs.Workspaces {
			out.Workspaces[i] = variableSetRefOutput{ID: ws.ID, Name: ws.Name}
		}
		for i, p := range vs.Projects {
			out.Projects[i] = variableSetRefOutput{ID: p.ID, Name: p.Name}
		}
		for i, vv := range vs.Variables {
			out.Variables[i] = variableSetVarOutput{ID: vv.ID, Key: vv.Key, Value: vv.Value}
		}
		return v.Output().RenderJSON(out)
	}

	parent := parseVariableSetParent(vs)
	parentDisplay := ""
	if parent.Type != "" {
		parentDisplay = parent.Type + ":" + parent.ID
	}
	props := []PropertyPair{
		{Key: "ID", Value: vs.ID},
		{Key: "Name", Value: vs.Name},
		{Key: "Description", Value: vs.Description},
		{Key: "Global", Value: vs.Global},
		{Key: "Priority", Value: vs.Priority},
		{Key: "Parent", Value: parentDisplay},
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
