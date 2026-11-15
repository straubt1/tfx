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
	// `tfx workspace assessments` command
	workspaceAssessmentsCmd = &cobra.Command{
		Use:     "assessments",
		Aliases: []string{"health", "drift-detection"},
		Short:   "Health Assessment Commands",
		Long:    "Enable or disable Health Assessments (drift detection and continuous validation) on Workspaces in a TFx Organization.",
	}

	// `tfx workspace assessments enable` command
	workspaceAssessmentsEnableCmd = &cobra.Command{
		Use:   "enable",
		Short: "Enable Health Assessments",
		Long:  "Enable Health Assessments on one or more Workspaces by name.",
		Example: `
tfx workspace assessments enable --name my-workspace

tfx workspace assessments enable --name ws-a,ws-b,ws-c`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseWorkspaceAssessmentsFlags(cmd)
			if err != nil {
				return err
			}
			return workspaceAssessmentsSet(cmdConfig, true)
		},
	}

	// `tfx workspace assessments enable all` command
	workspaceAssessmentsEnableAllCmd = &cobra.Command{
		Use:   "all",
		Short: "Enable Health Assessments on All Workspaces",
		Long:  "Enable Health Assessments on all Workspaces in a TFx Organization. Optional filter flags select a subset.",
		Example: `
tfx workspace assessments enable all

tfx workspace assessments enable all --project-id prj-abc123

tfx workspace assessments enable all --wildcard-name "*-prod" --exclude-tags legacy`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseWorkspaceAssessmentsAllFlags(cmd)
			if err != nil {
				return err
			}
			return workspaceAssessmentsSetAll(cmdConfig, true)
		},
	}

	// `tfx workspace assessments disable` command
	workspaceAssessmentsDisableCmd = &cobra.Command{
		Use:   "disable",
		Short: "Disable Health Assessments",
		Long:  "Disable Health Assessments on one or more Workspaces by name.",
		Example: `
tfx workspace assessments disable --name my-workspace

tfx workspace assessments disable --name ws-a,ws-b,ws-c`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseWorkspaceAssessmentsFlags(cmd)
			if err != nil {
				return err
			}
			return workspaceAssessmentsSet(cmdConfig, false)
		},
	}

	// `tfx workspace assessments disable all` command
	workspaceAssessmentsDisableAllCmd = &cobra.Command{
		Use:   "all",
		Short: "Disable Health Assessments on All Workspaces",
		Long:  "Disable Health Assessments on all Workspaces in a TFx Organization. Optional filter flags select a subset.",
		Example: `
tfx workspace assessments disable all

tfx workspace assessments disable all --tags agent-pool:onprem`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseWorkspaceAssessmentsAllFlags(cmd)
			if err != nil {
				return err
			}
			return workspaceAssessmentsSetAll(cmdConfig, false)
		},
	}
)

func init() {
	for _, c := range []*cobra.Command{workspaceAssessmentsEnableCmd, workspaceAssessmentsDisableCmd} {
		c.Flags().StringSliceP("name", "n", []string{}, "Workspace name(s), can be comma separated or repeated")
		c.MarkFlagRequired("name")
	}

	for _, c := range []*cobra.Command{workspaceAssessmentsEnableAllCmd, workspaceAssessmentsDisableAllCmd} {
		c.Flags().StringP("search", "s", "", "Search string anywhere in the Workspace Name (optional).")
		c.Flags().StringP("wildcard-name", "w", "", "Wildcard search string for Workspace Name, Examples: *-prod or prod-* (optional).")
		c.Flags().String("project-id", "", "Filter on Workspaces in this Project (optional).")
		c.Flags().String("tags", "", "Filter on Workspaces with this tag (optional).")
		c.Flags().String("exclude-tags", "", "Filter out Workspaces with this tag (optional).")
	}

	workspaceCmd.AddCommand(workspaceAssessmentsCmd)
	workspaceAssessmentsCmd.AddCommand(workspaceAssessmentsEnableCmd)
	workspaceAssessmentsEnableCmd.AddCommand(workspaceAssessmentsEnableAllCmd)
	workspaceAssessmentsCmd.AddCommand(workspaceAssessmentsDisableCmd)
	workspaceAssessmentsDisableCmd.AddCommand(workspaceAssessmentsDisableAllCmd)
}

func assessmentsAction(enabled bool) string {
	if enabled {
		return "Enabling"
	}
	return "Disabling"
}

func workspaceAssessmentsSet(cmdConfig *flags.WorkspaceAssessmentsFlags, enabled bool) error {
	v := view.NewWorkspaceAssessmentsView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("%s health assessments on %d workspace(s) in organization '%s'", assessmentsAction(enabled), len(cmdConfig.Names), c.OrganizationName)

	results := make([]view.WorkspaceAssessmentsResult, 0, len(cmdConfig.Names))
	for _, name := range cmdConfig.Names {
		status, err := data.SetWorkspaceAssessmentsByName(c, c.OrganizationName, name, enabled)
		if err != nil {
			status = err.Error()
		}
		results = append(results, view.WorkspaceAssessmentsResult{Name: name, Status: status})
	}

	return v.Render(results)
}

func workspaceAssessmentsSetAll(cmdConfig *flags.WorkspaceAssessmentsAllFlags, enabled bool) error {
	v := view.NewWorkspaceAssessmentsView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("%s health assessments on workspaces in organization '%s'", assessmentsAction(enabled), c.OrganizationName)
	printActiveFilters(v.BaseView, cmdConfig.ListFlags())

	workspaces, err := data.FetchWorkspaces(c, c.OrganizationName, cmdConfig.ListFlags())
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to list workspaces"))
	}

	results := make([]view.WorkspaceAssessmentsResult, 0, len(workspaces))
	for _, ws := range workspaces {
		status, err := data.SetWorkspaceAssessments(c, ws, enabled)
		if err != nil {
			status = err.Error()
		}
		results = append(results, view.WorkspaceAssessmentsResult{Name: ws.Name, Status: status})
	}

	return v.Render(results)
}
