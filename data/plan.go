// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/output"
)

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
