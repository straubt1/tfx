// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

// ── Phase 7c CV file message types ───────────────────────────────────────────

// cvFilesLoadedMsg carries the flat sorted file list for a config version.
type cvFilesLoadedMsg struct{ files []cvFile }

// cvFileContentLoadedMsg carries the lines and base name of a file read from disk.
type cvFileContentLoadedMsg struct {
	lines []string
	name  string
}

// cvFileErrMsg carries an error from a CV archive download or file read.
type cvFileErrMsg struct{ err error }

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
// Kept as a simple path helper (does not use the cacheDir() helper to avoid
// creating directories on every call).
func svJsonCachePath(svID string) string {
	p, err := stateJSONPath(svID)
	if err != nil {
		// Fallback using home dir directly.
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".tfx", "cache", "state", svID+".json")
	}
	return p
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

// ── CV archive commands ───────────────────────────────────────────────────────

// walkCVExtractDir walks the extracted directory and returns a flat sorted list
// of cvFile entries (filepath.WalkDir visits in lexicographic order).
func walkCVExtractDir(dir string) ([]cvFile, error) {
	var files []cvFile
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip unreadable entries
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil || rel == "." {
			return nil // skip root itself
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		files = append(files, cvFile{
			relPath: rel,
			size:    info.Size(),
			isDir:   d.IsDir(),
		})
		return nil
	})
	return files, err
}

// loadCVFiles downloads (or loads from cache) the config version archive and
// returns its file list. Set force=true to bypass the on-disk cache.
func loadCVFiles(c *client.TfxClient, cvID string, force bool) tea.Cmd {
	return func() tea.Msg {
		archivePath, err := cvArchivePath(cvID)
		if err != nil {
			return cvFileErrMsg{err: err}
		}
		extractDir, err := cvExtractDir(cvID)
		if err != nil {
			return cvFileErrMsg{err: err}
		}

		// Cache hit: extract dir exists and has content.
		if !force {
			if files, err := walkCVExtractDir(extractDir); err == nil && len(files) > 0 {
				return cvFilesLoadedMsg{files: files}
			}
		}

		// Force re-download: remove existing cache entries.
		if force {
			_ = os.RemoveAll(extractDir)
			_ = os.Remove(archivePath)
			if err := os.MkdirAll(extractDir, 0700); err != nil {
				return cvFileErrMsg{err: err}
			}
		}

		// Download the archive.
		b, err := data.DownloadConfigurationVersion(c, cvID)
		if err != nil {
			return cvFileErrMsg{err: err}
		}

		// Write archive to disk.
		if err := os.WriteFile(archivePath, b, 0600); err != nil {
			return cvFileErrMsg{err: err}
		}

		// Extract.
		if err := extractTarGz(archivePath, extractDir); err != nil {
			return cvFileErrMsg{err: err}
		}

		// Walk and return file list.
		files, err := walkCVExtractDir(extractDir)
		if err != nil {
			return cvFileErrMsg{err: err}
		}
		return cvFilesLoadedMsg{files: files}
	}
}

// loadCVFileContent reads a single file from the already-extracted CV directory.
func loadCVFileContent(cvID string, f cvFile) tea.Cmd {
	return func() tea.Msg {
		extractDir, err := cvExtractDir(cvID)
		if err != nil {
			return cvFileErrMsg{err: err}
		}
		b, err := os.ReadFile(filepath.Join(extractDir, f.relPath))
		if err != nil {
			return cvFileErrMsg{err: err}
		}
		lines := strings.Split(string(b), "\n")
		return cvFileContentLoadedMsg{lines: lines, name: filepath.Base(f.relPath)}
	}
}

// ── Instance info / health check messages ────────────────────────────────────

// healthCheckLoadedMsg carries the health check response as a flat string map.
type healthCheckLoadedMsg map[string]string

// healthCheckErrMsg carries an error from the health check fetch.
type healthCheckErrMsg struct{ err error }

// loadHealthCheck fetches /_health_check?full=1 from the configured host.
// The endpoint returns a JSON object whose values may be strings (e.g. "UP")
// or booleans; all values are coerced to strings for display.
func loadHealthCheck(c *client.TfxClient) tea.Cmd {
	return func() tea.Msg {
		url := fmt.Sprintf("https://%s/_health_check?full=1", c.Hostname)
		req, err := http.NewRequestWithContext(c.Context, http.MethodGet, url, nil)
		if err != nil {
			return healthCheckErrMsg{err}
		}
		req.Header.Set("Authorization", "Bearer "+c.Token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return healthCheckErrMsg{err}
		}
		defer resp.Body.Close()

		// Decode into interface{} to handle any JSON shape (string, bool, etc.).
		var raw map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
			return healthCheckErrMsg{err}
		}
		result := make(map[string]string, len(raw))
		for k, v := range raw {
			switch val := v.(type) {
			case string:
				result[k] = val
			default:
				b, _ := json.Marshal(val)
				result[k] = string(b)
			}
		}
		return healthCheckLoadedMsg(result)
	}
}
