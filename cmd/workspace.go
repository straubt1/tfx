// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/cmd/flags"
	view "github.com/straubt1/tfx/cmd/views"
	"github.com/straubt1/tfx/data"
)

// workspaceCmd represents the workspace command
var (
	// `tfx workspace` commands
	workspaceCmd = &cobra.Command{
		Use:     "workspace",
		Aliases: []string{"ws"},
		Short:   "Workspace Commands",
		Long:    "Work with TFx Workspaces",
	}

	// `tfx workspace list` command
	workspaceListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Workspaces",
		Long:  "List Workspaces in a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseWorkspaceListFlags(cmd)
			if err != nil {
				return err
			}

			if cmdConfig.All {
				return workspaceListAll(cmdConfig)
			} else {
				return workspaceList(cmdConfig)
			}
		},
	}

	// `tfx workspace show` command
	workspaceShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show Workspace",
		Long:  "Show Workspace in a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseWorkspaceShowFlags(cmd)
			if err != nil {
				return err
			}
			return workspaceShow(cmdConfig)
		},
	}
)

func init() {
	// `tfx workspace list`
	workspaceListCmd.Flags().StringP("search", "s", "", "Search string anywhere in the Workspace Name (optional).")
	workspaceListCmd.Flags().StringP("wildcard-name", "w", "", "Wildcard search string for Workspace Name, Examples: *-prod or prod-* (optional).")
	workspaceListCmd.Flags().StringP("project-id", "p", "", "Filter on Workspaces in this Project (optional).")
	workspaceListCmd.Flags().String("run-status", "", "Filter on Workspaces with this current run status (optional).")
	workspaceListCmd.Flags().String("tags", "", "Filter on Workspaces with this tag (optional).")
	workspaceListCmd.Flags().String("exclude-tags", "", "Filter out Workspaces with this tag (optional).")
	workspaceListCmd.Flags().BoolP("all", "a", false, "List All Organizations Workspaces (optional).")

	// remove?
	workspaceListCmd.Flags().StringP("repository", "r", "", "Filter on Repository Identifier (i.e. username/repo_name) (optional).")

	// `tfx workspace show`
	workspaceShowCmd.Flags().StringP("name", "n", "", "Name of the workspace.")
	workspaceShowCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(workspaceCmd)
	workspaceCmd.AddCommand(workspaceListCmd)
	workspaceCmd.AddCommand(workspaceShowCmd)
}

func workspaceList(cmdConfig *flags.WorkspaceListFlags) error {
	// Create view for rendering
	v := view.NewWorkspaceListView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	// Print command header before API call
	if cmdConfig.ProjectID != "" {
		v.PrintCommandHeader("Listing workspaces in organization '%s' and project ID '%s'", c.OrganizationName, cmdConfig.ProjectID)
	} else {
		v.PrintCommandHeader("Listing workspaces in organization '%s'", c.OrganizationName)
	}

	workspaces, err := data.FetchWorkspaces(c, c.OrganizationName, cmdConfig)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to list workspaces"))
	}

	// Apply client-side filters if needed (repository filter not supported by API)
	if cmdConfig.Repository != "" {
		workspaces = data.FilterWorkspaces(workspaces, "", cmdConfig.Repository)
	}

	return v.Render(c.OrganizationName, workspaces)
}

func workspaceListAll(cmdConfig *flags.WorkspaceListFlags) error {
	// Create view for rendering
	v := view.NewWorkspaceListView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	// Print command header before API call
	if cmdConfig.ProjectID != "" {
		v.PrintCommandHeader("Listing workspaces across all organizations and project ID '%s'", cmdConfig.ProjectID)
	} else {
		v.PrintCommandHeader("Listing workspaces across all organizations")
	}

	workspaces, err := data.FetchWorkspacesAcrossOrgs(c, cmdConfig)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to list workspaces"))
	}

	// Apply client-side filters if needed (repository filter not supported by API)
	if cmdConfig.Repository != "" {
		workspaces = data.FilterWorkspaces(workspaces, "", cmdConfig.Repository)
	}

	return v.RenderAll(workspaces)
}

func workspaceShow(cmdConfig *flags.WorkspaceShowFlags) error {
	// Create view for rendering
	v := view.NewWorkspaceShowView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Showing workspace '%s' in organization '%s'", cmdConfig.Name, c.OrganizationName)

	workspace, err := data.FetchWorkspace(c, c.OrganizationName, cmdConfig.Name)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to read workspace"))
	}

	// Fetch current run if present
	var currentRun *tfe.Run
	if workspace.CurrentRun != nil {
		currentRun, err = c.Client.Runs.ReadWithOptions(c.Context, workspace.CurrentRun.ID, &tfe.RunReadOptions{
			Include: []tfe.RunIncludeOpt{},
		})
		if err != nil {
			return v.RenderError(errors.Wrap(err, "failed to read workspace current run"))
		}
	}

	// Fetch remote state consumers if not global
	var remoteStateConsumers []*tfe.Workspace
	if !workspace.GlobalRemoteState {
		remoteStateConsumers, err = data.FetchWorkspaceRemoteStateConsumers(c, workspace.ID)
		if err != nil {
			return v.RenderError(errors.Wrap(err, "failed to list remote state consumers"))
		}
	}

	// Fetch team access
	teamAccess, err := data.FetchWorkspaceTeamAccess(c, workspace.ID, 0)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to list teams"))
	}

	// Get team names from team access
	var teamNames []interface{}
	if len(teamAccess) > 0 {
		teamNames, err = data.GetTeamAccessNames(c, teamAccess)
		if err != nil {
			return v.RenderError(errors.Wrap(err, "failed to find team name"))
		}
	}

	return v.Render(c.OrganizationName, workspace, currentRun, teamNames, remoteStateConsumers)
}

// Legacy helper functions for backward compatibility with other workspace commands
// These are thin wrappers around the data layer functions
// TODO: Delete these in future major version

func workspaceListAllForOrganization(c TfxClientContext, orgName string, searchString string, projectID string) ([]*tfe.Workspace, error) {
	tfxClient := &client.TfxClient{
		Client:           c.Client,
		Context:          c.Context,
		OrganizationName: c.OrganizationName,
	}
	// TODO: do we need the same filters as without -a
	options := &flags.WorkspaceListFlags{
		Search:    searchString,
		ProjectID: projectID,
	}
	return data.FetchWorkspaces(tfxClient, orgName, options)
}

func organizationListAll(c TfxClientContext) ([]*tfe.Organization, error) {
	tfxClient := &client.TfxClient{
		Client:           c.Client,
		Context:          c.Context,
		OrganizationName: c.OrganizationName,
	}
	return data.FetchOrganizations(tfxClient, "")
}
