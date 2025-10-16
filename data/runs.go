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
