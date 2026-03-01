// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/cmd/flags"
	view "github.com/straubt1/tfx/cmd/views"
	"github.com/straubt1/tfx/data"
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
			cmdConfig, err := flags.ParseTeamListFlags(cmd)
			if err != nil {
				return err
			}
			return workspaceTeamList(cmdConfig)
		},
	}
)

func init() {
	// `tfx workspace team list` command flags
	workspaceTeamListCmd.Flags().StringP("name", "n", "", "Name of the Workspace")
	workspaceTeamListCmd.MarkFlagRequired("name")

	workspaceCmd.AddCommand(workspaceTeamCmd)
	workspaceTeamCmd.AddCommand(workspaceTeamListCmd)
}

func workspaceTeamList(cmdConfig *flags.TeamListFlags) error {
	// Create view for rendering
	v := view.NewTeamListView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Listing teams for workspace '%s'", cmdConfig.WorkspaceName)

	workspaceID, err := data.GetWorkspaceID(c, c.OrganizationName, cmdConfig.WorkspaceName)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "unable to read workspace id"))
	}

	// Fetch all items (no max limit)
	teamAccess, err := data.FetchWorkspaceTeamAccess(c, workspaceID, 0)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to list teams"))
	}

	// Resolve team names
	namesIface, err := data.GetTeamAccessNames(c, teamAccess)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to find team name"))
	}

	// Convert []interface{} to []string
	names := make([]string, len(namesIface))
	for i, n := range namesIface {
		if s, ok := n.(string); ok {
			names[i] = s
		} else {
			names[i] = ""
		}
	}

	return v.Render(teamAccess, names)
}
