// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	tfe "github.com/hashicorp/go-tfe"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/logger"
)

// FetchRun retrieves a run by ID with optional includes
func FetchRun(c *client.TfxClient, runID string) (*tfe.Run, error) {
	logger.Debug("Fetching run by ID", "runID", runID)

	run, err := c.Client.Runs.ReadWithOptions(c.Context, runID, &tfe.RunReadOptions{
		Include: []tfe.RunIncludeOpt{},
	})
	if err != nil {
		logger.Error("Failed to fetch run", "runID", runID, "error", err)
		return nil, err
	}

	logger.Debug("Run fetched successfully", "runID", runID, "status", run.Status)
	return run, nil
}

// FetchRunsForWorkspace lists runs for a workspace limited by maxItems
func FetchRunsForWorkspace(c *client.TfxClient, workspaceID string, maxItems int) ([]*tfe.Run, error) {
	logger.Debug("Fetching runs for workspace", "workspaceID", workspaceID, "maxItems", maxItems)

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
			logger.Error("Failed to list runs", "workspaceID", workspaceID, "page", opts.PageNumber, "error", err)
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

	logger.Debug("Runs fetched", "workspaceID", workspaceID, "count", len(all))
	return all, nil
}

// CreateRun creates a run for a workspace, optionally using a specific configuration version
func CreateRun(c *client.TfxClient, orgName, workspaceName, message, configurationVersionID string) (*tfe.Run, error) {
	logger.Debug("Creating run", "organization", orgName, "workspaceName", workspaceName, "cvID", configurationVersionID)

	w, err := c.Client.Workspaces.Read(c.Context, orgName, workspaceName)
	if err != nil {
		logger.Error("Failed to read workspace", "organization", orgName, "workspaceName", workspaceName, "error", err)
		return nil, err
	}

	var cv *tfe.ConfigurationVersion
	if configurationVersionID != "" {
		cv, err = c.Client.ConfigurationVersions.Read(c.Context, configurationVersionID)
		if err != nil {
			logger.Error("Failed to read configuration version", "cvID", configurationVersionID, "error", err)
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
		logger.Error("Failed to create run", "organization", orgName, "workspaceName", workspaceName, "error", err)
		return nil, err
	}

	logger.Debug("Run created", "runID", run.ID)
	return run, nil
}

// DiscardRun discards a run by ID
func DiscardRun(c *client.TfxClient, runID string) error {
	logger.Debug("Discarding run", "runID", runID)
	if err := c.Client.Runs.Discard(c.Context, runID, tfe.RunDiscardOptions{Comment: tfe.String("Discarded by tfx")}); err != nil {
		logger.Error("Failed to discard run", "runID", runID, "error", err)
		return err
	}
	return nil
}

// CancelRun cancels a run by ID
func CancelRun(c *client.TfxClient, runID string) error {
	logger.Debug("Canceling run", "runID", runID)
	if err := c.Client.Runs.Cancel(c.Context, runID, tfe.RunCancelOptions{Comment: tfe.String("Canceled via TFx")}); err != nil {
		logger.Error("Failed to cancel run", "runID", runID, "error", err)
		return err
	}
	return nil
}

// GetLatestRunID returns the most recent run ID for a workspace
func GetLatestRunID(c *client.TfxClient, workspaceID string) (string, error) {
	logger.Debug("Getting latest run ID", "workspaceID", workspaceID)

	res, err := c.Client.Runs.List(c.Context, workspaceID, &tfe.RunListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 1, PageSize: 1},
	})
	if err != nil {
		logger.Error("Failed to list runs for latest", "workspaceID", workspaceID, "error", err)
		return "", err
	}
	if res == nil || len(res.Items) != 1 {
		logger.Error("Latest run not found", "workspaceID", workspaceID)
		return "", tfe.ErrResourceNotFound
	}
	return res.Items[0].ID, nil
}
