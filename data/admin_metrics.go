// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	"time"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/cmd/flags"
	view "github.com/straubt1/tfx/cmd/views"
	"github.com/straubt1/tfx/output"
)

// FetchWorkspaceMetrics fetches metrics for all workspaces in an organization since a given time
func FetchWorkspaceMetrics(c *client.TfxClient, orgName string, runSinceTime time.Time) ([]view.MetricsWorkspace, error) {
	output.Get().Logger().Debug("Fetching workspace metrics", "organization", orgName, "since", runSinceTime)

	// Fetch all workspaces for the organization
	wsFlags := &flags.WorkspaceListFlags{
		Search:    "",
		ProjectID: "",
		All:       false,
	}
	workspaces, err := FetchWorkspaces(c, orgName, wsFlags)
	if err != nil {
		output.Get().Logger().Error("Failed to fetch workspaces", "organization", orgName, "error", err)
		return nil, errors.Wrap(err, "failed to list workspaces")
	}

	output.Get().Logger().Debug("Workspaces fetched", "organization", orgName, "count", len(workspaces))

	var result []view.MetricsWorkspace

	for _, ws := range workspaces {
		output.Get().Logger().Trace("Processing workspace metrics", "workspace", ws.Name, "id", ws.ID)

		wsResult := view.MetricsWorkspace{
			Name: ws.Name,
			ID:   ws.ID,
		}

		// Fetch runs for the workspace
		// Note: Using PageSize 100 as in original implementation
		// TODO: Consider using client.FetchAll for full pagination
		runs, err := c.Client.Runs.List(c.Context, ws.ID, &tfe.RunListOptions{
			ListOptions: tfe.ListOptions{
				PageSize: 100,
			},
			Include: []tfe.RunIncludeOpt{},
		})
		if err != nil {
			output.Get().Logger().Error("Failed to fetch runs", "workspace", ws.Name, "error", err)
			return nil, errors.Wrapf(err, "failed to list runs for workspace %s", ws.Name)
		}

		// Process each run
		for _, r := range runs.Items {
			// Skip runs outside the time frame
			if runSinceTime.After(r.CreatedAt) {
				continue
			}

			wsResult.RunCount++

			// Count run statuses
			switch r.Status {
			case "errored":
				wsResult.RunErroredCount++
			case "canceled", "force_canceled":
				wsResult.RunCancelledCount++
			case "discarded":
				wsResult.RunDiscardedCount++
			}

			// Process policy checks
			wsResult.PolicyCheckCount += len(r.PolicyChecks)
			for _, p := range r.PolicyChecks {
				pFull, err := c.Client.PolicyChecks.Read(c.Context, p.ID)
				if err != nil {
					output.Get().Logger().Warn("Failed to read policy check", "policyCheckID", p.ID, "error", err)
					continue
				}
				if pFull != nil && pFull.Result != nil {
					wsResult.PoliciesPassCount += pFull.Result.Passed
					wsResult.PoliciesFailCount += pFull.Result.TotalFailed
				}
			}
		}

		output.Get().Logger().Trace("Workspace metrics processed", "workspace", ws.Name, "runs", wsResult.RunCount)
		result = append(result, wsResult)
	}

	output.Get().Logger().Debug("Workspace metrics fetched successfully", "organization", orgName, "workspaceCount", len(result))
	return result, nil
}
