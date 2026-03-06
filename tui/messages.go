// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	tea "charm.land/bubbletea/v2"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/cmd/flags"
	"github.com/straubt1/tfx/data"
)

// projectsLoadedMsg carries the fetched project list.
type projectsLoadedMsg []*tfe.Project

// workspacesLoadedMsg carries the fetched workspace list.
type workspacesLoadedMsg []*tfe.Workspace

// runsLoadedMsg carries the fetched run list.
type runsLoadedMsg []*tfe.Run

// fetchErrMsg wraps any error returned from an async fetch.
type fetchErrMsg struct{ err error }

func loadProjects(c *client.TfxClient, org string) tea.Cmd {
	return func() tea.Msg {
		projects, err := data.FetchProjects(c, org, "")
		if err != nil {
			return fetchErrMsg{err}
		}
		return projectsLoadedMsg(projects)
	}
}

func loadWorkspaces(c *client.TfxClient, org, projectID string) tea.Cmd {
	return func() tea.Msg {
		opts := &flags.WorkspaceListFlags{ProjectID: projectID}
		workspaces, err := data.FetchWorkspaces(c, org, opts)
		if err != nil {
			return fetchErrMsg{err}
		}
		return workspacesLoadedMsg(workspaces)
	}
}

func loadRuns(c *client.TfxClient, workspaceID string) tea.Cmd {
	return func() tea.Msg {
		runs, err := data.FetchRunsForWorkspace(c, workspaceID, 50)
		if err != nil {
			return fetchErrMsg{err}
		}
		return runsLoadedMsg(runs)
	}
}
