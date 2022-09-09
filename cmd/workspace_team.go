// Copyright © 2021 Tom Straub <github.com/straubt1>

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package cmd

import (
	"math"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// `tfx workspace team` commands
	workspaceTeamCmd = &cobra.Command{
		Use:   "team",
		Short: "Team Commands",
		Long:  "Commands to work with Workspace Teams.",
	}

	// `tfx workspace team list` command
	workspaceTeamListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Teams",
		Long:  "List Teams in a Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			m := *viperInt("max-items")
			if *viperBool("all") {
				m = math.MaxInt
			}

			return workspaceTeamList(
				getTfxClientContext(),
				*viperString("workspace-name"),
				m)
		},
	}
)

func init() {
	// `tfx variable list` command
	workspaceTeamListCmd.Flags().StringP("workspace-name", "w", "", "Name of the Workspace")
	workspaceTeamListCmd.Flags().IntP("max-items", "m", 100, "Max number of results (optional)")
	workspaceTeamListCmd.Flags().BoolP("all", "a", false, "Retrieve all results regardless of maxItems flag (optional)")
	workspaceTeamListCmd.MarkFlagRequired("workspace-name")

	workspaceCmd.AddCommand(workspaceTeamCmd)
	workspaceTeamCmd.AddCommand(workspaceTeamListCmd)
}

func workspaceTeamListAll(c TfxClientContext, workspaceId string, maxItems int) ([]*tfe.TeamAccess, error) {
	pageSize := 100
	if maxItems < 100 {
		pageSize = maxItems // Only get what we need in one page
	}

	allItems := []*tfe.TeamAccess{}
	opts := tfe.TeamAccessListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 1, PageSize: pageSize},
		WorkspaceID: workspaceId,
	}
	for {
		items, err := c.Client.TeamAccess.List(c.Context, &opts)
		if err != nil {
			return nil, err
		}

		allItems = append(allItems, items.Items...)
		if len(allItems) >= maxItems {
			break // Hit the max, break. For maxItems > 100 it is possible to return more than max in this approach
		}

		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage
	}

	return allItems, nil
}

func workspaceTeamList(c TfxClientContext, workspaceName string, maxItems int) error {
	o.AddMessageUserProvided("List Variables for Workspace:", workspaceName)
	workspaceId, err := getWorkspaceId(c, workspaceName)
	if err != nil {
		return errors.Wrap(err, "unable to read workspace id")
	}

	items, err := workspaceTeamListAll(c, workspaceId, maxItems)
	if err != nil {
		return errors.Wrap(err, "failed to list teams")
	}

	o.AddTableHeader("Name", "Team Id", "Team Access Id", "Access Type", "Runs", "Workspace Locking", "Sentinel Mocks", "Run Tasks", "Variables", "State Versions")
	for _, i := range items {
		t, err := c.Client.Teams.Read(c.Context, i.Team.ID)
		if err != nil {
			return errors.Wrap(err, "failed to find team name")
		}
		o.AddTableRows(t.Name, i.Team.ID, i.ID, i.Access, i.Runs, i.WorkspaceLocking, i.SentinelMocks, i.RunTasks, i.Variables, i.StateVersions)
	}

	return nil
}
