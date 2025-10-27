// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

// TeamListView handles rendering for team list command
type TeamListView struct {
	*BaseView
}

func NewTeamListView() *TeamListView {
	return &TeamListView{BaseView: NewBaseView()}
}

// teamAccessOutput is a JSON-safe representation for team access list
type teamAccessOutput struct {
	TeamName         string `json:"teamName"`
	TeamID           string `json:"teamId"`
	TeamAccessID     string `json:"teamAccessId"`
	AccessType       string `json:"accessType"`
	Runs             string `json:"runs"`
	WorkspaceLocking bool   `json:"workspaceLocking"`
	SentinelMocks    string `json:"sentinelMocks"`
	RunTasks         bool   `json:"runTasks"`
	Variables        string `json:"variables"`
	StateVersions    string `json:"stateVersions"`
}

// Render renders team access entries and names
func (v *TeamListView) Render(access []*tfe.TeamAccess, teamNames []string) error {
	if v.IsJSON() {
		output := make([]teamAccessOutput, len(access))
		for i, a := range access {
			output[i] = teamAccessOutput{
				TeamName:         safeIndex(teamNames, i),
				TeamID:           a.Team.ID,
				TeamAccessID:     a.ID,
				AccessType:       string(a.Access),
				Runs:             string(a.Runs),
				WorkspaceLocking: a.WorkspaceLocking,
				SentinelMocks:    string(a.SentinelMocks),
				RunTasks:         a.RunTasks,
				Variables:        string(a.Variables),
				StateVersions:    string(a.StateVersions),
			}
		}
		return v.Output().RenderJSON(output)
	}

	headers := []string{"Name", "Team Id", "Team Access Id", "Access Type", "Runs", "Workspace Locking", "Sentinel Mocks", "Run Tasks", "Variables", "State Versions"}
	rows := make([][]interface{}, len(access))
	for i, a := range access {
		rows[i] = []interface{}{
			safeIndex(teamNames, i),
			a.Team.ID,
			a.ID,
			a.Access,
			a.Runs,
			a.WorkspaceLocking,
			a.SentinelMocks,
			a.RunTasks,
			a.Variables,
			a.StateVersions,
		}
	}
	return v.Output().RenderTable(headers, rows)
}

// safeIndex returns the element at index i if within bounds, otherwise empty string
func safeIndex(arr []string, i int) string {
	if i >= 0 && i < len(arr) {
		return arr[i]
	}
	return ""
}
