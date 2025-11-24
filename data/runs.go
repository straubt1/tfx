// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/output"
)

// FetchRun retrieves a run by ID with optional includes
func FetchRun(c *client.TfxClient, runID string) (*tfe.Run, error) {
	output.Get().Logger().Debug("Fetching run by ID", "runID", runID)

	run, err := c.Client.Runs.ReadWithOptions(c.Context, runID, &tfe.RunReadOptions{
		Include: []tfe.RunIncludeOpt{},
	})
	if err != nil {
		output.Get().Logger().Error("Failed to fetch run", "runID", runID, "error", err)
		return nil, err
	}

	output.Get().Logger().Debug("Run fetched successfully", "runID", runID, "status", run.Status)
	return run, nil
}

// FetchRunsForWorkspace lists runs for a workspace limited by maxItems
func FetchRunsForWorkspace(c *client.TfxClient, workspaceID string, maxItems int) ([]*tfe.Run, error) {
	output.Get().Logger().Debug("Fetching runs for workspace", "workspaceID", workspaceID, "maxItems", maxItems)

	// Determine page size: fetch only what we need if <=100
	pageSize := 100
	if maxItems > 0 && maxItems < 100 {
		pageSize = maxItems
	}

	var all []*tfe.Run
	opts := &tfe.RunListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 1, PageSize: pageSize},
		Operation:   "plan_only,plan_and_apply,refresh_only,destroy,empty_apply",
		Include:     []tfe.RunIncludeOpt{},
	}
	for {
		res, err := c.Client.Runs.List(c.Context, workspaceID, opts)
		if err != nil {
			output.Get().Logger().Error("Failed to list runs", "workspaceID", workspaceID, "page", opts.PageNumber, "error", err)
			return nil, err
		}
		all = append(all, res.Items...)

		// Stop if reached maxItems
		if maxItems > 0 && len(all) >= maxItems {
			break
		}

		if res.CurrentPage >= res.TotalPages {
			break
		}
		opts.PageNumber = res.NextPage
	}

	// Trim to exact maxItems if we slightly over-fetched
	if maxItems > 0 && len(all) > maxItems {
		all = all[:maxItems]
	}

	output.Get().Logger().Debug("Runs fetched", "workspaceID", workspaceID, "count", len(all))
	return all, nil
}

// CreateRun creates a run for a workspace, optionally using a specific configuration version
func CreateRun(c *client.TfxClient, orgName, workspaceName, message, configurationVersionID string) (*tfe.Run, error) {
	output.Get().Logger().Debug("Creating run", "organization", orgName, "workspaceName", workspaceName, "cvID", configurationVersionID)

	w, err := c.Client.Workspaces.Read(c.Context, orgName, workspaceName)
	if err != nil {
		output.Get().Logger().Error("Failed to read workspace", "organization", orgName, "workspaceName", workspaceName, "error", err)
		return nil, err
	}

	var cv *tfe.ConfigurationVersion
	if configurationVersionID != "" {
		cv, err = c.Client.ConfigurationVersions.Read(c.Context, configurationVersionID)
		if err != nil {
			output.Get().Logger().Error("Failed to read configuration version", "cvID", configurationVersionID, "error", err)
			return nil, err
		}
	}

	run, err := c.Client.Runs.Create(c.Context, tfe.RunCreateOptions{
		Workspace:            w,
		IsDestroy:            tfe.Bool(false),
		Message:              tfe.String(message),
		ConfigurationVersion: cv, // may be nil
	})
	if err != nil {
		output.Get().Logger().Error("Failed to create run", "organization", orgName, "workspaceName", workspaceName, "error", err)
		return nil, err
	}

	output.Get().Logger().Debug("Run created", "runID", run.ID)
	return run, nil
}

// DiscardRun discards a run by ID
func DiscardRun(c *client.TfxClient, runID string) error {
	output.Get().Logger().Debug("Discarding run", "runID", runID)
	if err := c.Client.Runs.Discard(c.Context, runID, tfe.RunDiscardOptions{Comment: tfe.String("Discarded by tfx")}); err != nil {
		output.Get().Logger().Error("Failed to discard run", "runID", runID, "error", err)
		return err
	}
	return nil
}

// CancelRun cancels a run by ID
func CancelRun(c *client.TfxClient, runID string) error {
	output.Get().Logger().Debug("Canceling run", "runID", runID)
	if err := c.Client.Runs.Cancel(c.Context, runID, tfe.RunCancelOptions{Comment: tfe.String("Canceled via TFx")}); err != nil {
		output.Get().Logger().Error("Failed to cancel run", "runID", runID, "error", err)
		return err
	}
	return nil
}

// GetLatestRunID returns the most recent run ID for a workspace
func GetLatestRunID(c *client.TfxClient, workspaceID string) (string, error) {
	output.Get().Logger().Debug("Getting latest run ID", "workspaceID", workspaceID)

	res, err := c.Client.Runs.List(c.Context, workspaceID, &tfe.RunListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 1, PageSize: 1},
	})
	if err != nil {
		output.Get().Logger().Error("Failed to list runs for latest", "workspaceID", workspaceID, "error", err)
		return "", err
	}
	if res == nil || len(res.Items) != 1 {
		output.Get().Logger().Error("Latest run not found", "workspaceID", workspaceID)
		return "", tfe.ErrResourceNotFound
	}
	return res.Items[0].ID, nil
}

// FetchPlan retrieves a plan by ID
func FetchPlan(c *client.TfxClient, planID string) (*tfe.Plan, error) {
	output.Get().Logger().Debug("Fetching plan by ID", "planID", planID)

	plan, err := c.Client.Plans.Read(c.Context, planID)
	if err != nil {
		output.Get().Logger().Error("Failed to fetch plan", "planID", planID, "error", err)
		return nil, err
	}

	output.Get().Logger().Debug("Plan fetched successfully", "planID", planID, "status", plan.Status)
	return plan, nil
}

// FetchPlanLogs retrieves logs for a plan by ID
func FetchPlanLogs(c *client.TfxClient, planID string) ([]string, error) {
	output.Get().Logger().Debug("Fetching plan logs", "planID", planID)

	logsReader, err := c.Client.Plans.Logs(c.Context, planID)
	if err != nil {
		output.Get().Logger().Error("Failed to fetch plan logs", "planID", planID, "error", err)
		return nil, err
	}

	// Read all logs into memory and split by lines
	logBytes := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := logsReader.Read(buf)
		if n > 0 {
			logBytes = append(logBytes, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	// Split logs by newline
	logContent := string(logBytes)
	var lines []string
	for _, line := range strings.Split(logContent, "\n") {
		if line != "" {
			lines = append(lines, line)
		}
	}

	output.Get().Logger().Debug("Plan logs fetched successfully", "planID", planID, "lineCount", len(lines))
	return lines, nil
}

// FetchPlanJSONOutput retrieves the JSON output for a plan by ID
func FetchPlanJSONOutput(c *client.TfxClient, planID string) ([]byte, error) {
	output.Get().Logger().Debug("Fetching plan JSON output", "planID", planID)

	jsonOutput, err := c.Client.Plans.ReadJSONOutput(c.Context, planID)
	if err != nil {
		output.Get().Logger().Error("Failed to fetch plan JSON output", "planID", planID, "error", err)
		return nil, err
	}

	output.Get().Logger().Debug("Plan JSON output fetched successfully", "planID", planID, "size", len(jsonOutput))
	return jsonOutput, nil
}
