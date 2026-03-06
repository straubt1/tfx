// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/cmd/flags"
	"github.com/straubt1/tfx/data"
)

// ── Phase 7 detail message types ──────────────────────────────────────────────

// runDetailLoadedMsg carries a fully-fetched run (with Plan, Apply, CV + ingress includes).
type runDetailLoadedMsg *tfe.Run

// svJsonLoadedMsg carries the lines of a downloaded (and pretty-printed) state JSON.
type svJsonLoadedMsg struct{ lines []string }

// svJsonErrMsg carries an error from the state JSON download.
type svJsonErrMsg struct{ err error }

// ── Spinner ───────────────────────────────────────────────────────────────────

// spinnerTickMsg advances the animated loading spinner by one frame.
type spinnerTickMsg struct{}

// tickSpinner returns a command that sleeps briefly then fires spinnerTickMsg.
// Chain it from Update() while m.loading == true to animate the spinner.
func tickSpinner() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(80 * time.Millisecond)
		return spinnerTickMsg{}
	}
}

// ── Data messages ─────────────────────────────────────────────────────────────

// orgsLoadedMsg carries the fetched organization list.
type orgsLoadedMsg []*tfe.Organization

// projectsLoadedMsg carries the fetched project list.
type projectsLoadedMsg []*tfe.Project

// workspacesLoadedMsg carries the fetched workspace list.
type workspacesLoadedMsg []*tfe.Workspace

// runsLoadedMsg carries the fetched run list.
type runsLoadedMsg []*tfe.Run

// variablesLoadedMsg carries the fetched variable list.
type variablesLoadedMsg []*tfe.Variable

// configVersionsLoadedMsg carries the fetched configuration version list.
type configVersionsLoadedMsg []*tfe.ConfigurationVersion

// stateVersionsLoadedMsg carries the fetched state version list.
type stateVersionsLoadedMsg []*tfe.StateVersion

// fetchErrMsg wraps any error returned from an async fetch.
type fetchErrMsg struct{ err error }

// ── Commands ──────────────────────────────────────────────────────────────────

func loadOrganizations(c *client.TfxClient) tea.Cmd {
	return func() tea.Msg {
		orgs, err := data.FetchOrganizations(c, "")
		if err != nil {
			return fetchErrMsg{err}
		}
		return orgsLoadedMsg(orgs)
	}
}

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

func loadVariables(c *client.TfxClient, workspaceID string) tea.Cmd {
	return func() tea.Msg {
		vars, err := data.FetchVariables(c, workspaceID)
		if err != nil {
			return fetchErrMsg{err}
		}
		return variablesLoadedMsg(vars)
	}
}

func loadConfigVersions(c *client.TfxClient, orgName, wsName string) tea.Cmd {
	return func() tea.Msg {
		cvs, err := data.FetchConfigurationVersions(c, orgName, wsName, 50)
		if err != nil {
			return fetchErrMsg{err}
		}
		return configVersionsLoadedMsg(cvs)
	}
}

func loadStateVersions(c *client.TfxClient, orgName, wsName string) tea.Cmd {
	return func() tea.Msg {
		svs, err := data.FetchStateVersions(c, orgName, wsName, 50)
		if err != nil {
			return fetchErrMsg{err}
		}
		return stateVersionsLoadedMsg(svs)
	}
}

// loadRunDetail fetches a run with full includes (Plan, Apply, ConfigurationVersion + ingress).
// The result silently updates selectedRun without changing the current view or loading state.
func loadRunDetail(c *client.TfxClient, runID string) tea.Cmd {
	return func() tea.Msg {
		run, err := c.Client.Runs.ReadWithOptions(c.Context, runID, &tfe.RunReadOptions{
			Include: []tfe.RunIncludeOpt{
				tfe.RunPlan,
				tfe.RunApply,
				tfe.RunConfigVer,
				tfe.RunConfigVerIngress,
			},
		})
		if err != nil {
			// Swallow the error silently — partial data from the list is still shown.
			return nil
		}
		return runDetailLoadedMsg(run)
	}
}

// svJsonCachePath returns the on-disk cache path for a state version's JSON.
// Path: ~/.tfx/cache/state/<svID>.json
func svJsonCachePath(svID string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".tfx", "cache", "state", svID+".json")
}

// loadStateVersionJson downloads (or loads from cache) the state JSON for svID.
// If force is true the cached file is ignored and a fresh download is performed.
func loadStateVersionJson(c *client.TfxClient, svID string, force bool) tea.Cmd {
	return func() tea.Msg {
		cacheFile := svJsonCachePath(svID)

		// Try cache first (unless force re-download requested).
		if !force {
			if b, err := os.ReadFile(cacheFile); err == nil {
				lines := strings.Split(string(b), "\n")
				return svJsonLoadedMsg{lines: lines}
			}
		}

		// Download via the data layer.
		b, err := data.DownloadStateVersion(c, svID)
		if err != nil {
			return svJsonErrMsg{err: err}
		}

		// Pretty-print the JSON for readability.
		var pretty bytes.Buffer
		if jerr := json.Indent(&pretty, b, "", "  "); jerr == nil {
			b = pretty.Bytes()
		}

		// Write to cache (best-effort; ignore errors).
		if dir := filepath.Dir(cacheFile); dir != "" {
			if merr := os.MkdirAll(dir, 0755); merr == nil {
				_ = os.WriteFile(cacheFile, b, 0644)
			}
		}

		lines := strings.Split(string(b), "\n")
		return svJsonLoadedMsg{lines: lines}
	}
}
