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

var (
	// `tfx workspace plan` commands
	planCmd = &cobra.Command{
		Use:   "plan",
		Short: "Plans",
		Long:  "Work with Plans of a TFx Workspace.",
	}

	// `tfx workspace plan show` command
	planShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show Plan",
		Long:  "Show Plan details for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParsePlanShowFlags(cmd)
			if err != nil {
				return err
			}
			return planShow(cmdConfig)
		},
	}

	// `tfx workspace plan logs` command
	planLogsCmd = &cobra.Command{
		Use:   "logs",
		Short: "Show Plan Logs",
		Long:  "Show Plan logs for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParsePlanLogsFlags(cmd)
			if err != nil {
				return err
			}
			return planLogs(cmdConfig)
		},
	}

	// `tfx workspace plan jsonoutput` command
	planJSONOutputCmd = &cobra.Command{
		Use:   "jsonoutput",
		Short: "Show Plan JSON Output",
		Long:  "Show Plan JSON output for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParsePlanJSONOutputFlags(cmd)
			if err != nil {
				return err
			}
			return planJSONOutput(cmdConfig)
		},
	}

	// `tfx workspace plan create` command
	planCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create Plan",
		Long:  "Create a new Plan for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParsePlanCreateFlags(cmd)
			if err != nil {
				return err
			}
			return planCreate(cmdConfig)
		},
	}

)

func init() {
	// `tfx workspace plan show` command
	planShowCmd.Flags().StringP("id", "i", "", "Plan Id (i.e. plan-*)")
	planShowCmd.MarkFlagRequired("id")

	// `tfx workspace plan logs` command
	planLogsCmd.Flags().StringP("id", "i", "", "Plan Id (i.e. plan-*)")
	planLogsCmd.MarkFlagRequired("id")

	// `tfx workspace plan jsonoutput` command
	planJSONOutputCmd.Flags().StringP("id", "i", "", "Plan Id (i.e. plan-*)")
	planJSONOutputCmd.MarkFlagRequired("id")

	// `tfx workspace plan create` command
	planCreateCmd.Flags().StringP("workspace-name", "w", "", "Workspace name")
	planCreateCmd.Flags().StringP("directory", "d", "./", "Directory of Terraform (optional, defaults to current directory)")
	planCreateCmd.Flags().StringP("configuration-id", "i", "", "Configuration Version Id (optional, i.e. cv-*)")
	planCreateCmd.Flags().StringP("message", "m", "", "Run Message (optional)")
	planCreateCmd.Flags().Bool("speculative", false, "Perform a Speculative Plan (optional)")
	planCreateCmd.Flags().Bool("destroy", false, "Perform a Destroy Plan (optional)")
	planCreateCmd.Flags().StringSlice("env", []string{}, "Environment variables to write to the Workspace. Can be supplied multiple times. (optional, i.e. '--env='AWS_REGION=us-east1')")
	planCreateCmd.MarkFlagRequired("workspace-name")

	workspaceCmd.AddCommand(planCmd)
	planCmd.AddCommand(planShowCmd)
	planCmd.AddCommand(planLogsCmd)
	planCmd.AddCommand(planJSONOutputCmd)
	planCmd.AddCommand(planCreateCmd)
}

func planShow(cmdConfig *flags.PlanShowFlags) error {
	v := view.NewPlanShowView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Showing plan '%s'", cmdConfig.ID)

	plan, err := data.FetchPlan(c, cmdConfig.ID)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to read plan from id"))
	}

	return v.Render(plan)
}

func planLogs(cmdConfig *flags.PlanLogsFlags) error {
	v := view.NewPlanLogsView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Showing logs for plan '%s'", cmdConfig.ID)

	logs, err := data.FetchPlanLogs(c, cmdConfig.ID)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to read plan logs"))
	}

	return v.Render(logs)
}

func planJSONOutput(cmdConfig *flags.PlanJSONOutputFlags) error {
	v := view.NewPlanJSONOutputView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Showing JSON output for plan '%s'", cmdConfig.ID)

	jsonOutput, err := data.FetchPlanJSONOutput(c, cmdConfig.ID)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to read plan JSON output"))
	}

	return v.Render(jsonOutput)
}

func planCreate(cmdConfig *flags.PlanCreateFlags) error {
	// Create a view for output
	v := view.NewPlanCreateView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	// Read workspace by name
	workspace, err := c.Client.Workspaces.Read(c.Context, c.OrganizationName, cmdConfig.WorkspaceName)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to read workspace"))
	}

	v.PrintCommandHeader("Creating plan for workspace '%s' (%s)", cmdConfig.WorkspaceName, workspace.ID)

	// Create a new run with plan
	runOptions := tfe.RunCreateOptions{
		IsDestroy: tfe.Bool(cmdConfig.Destroy),
		Workspace: workspace,
	}

	// Add optional message if provided
	if cmdConfig.Message != "" {
		runOptions.Message = tfe.String(cmdConfig.Message)
	}

	// Create the run
	v.Output().Message("Creating new Plan...")
	run, err := c.Client.Runs.Create(c.Context, runOptions)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to create run"))
	}

	// Fetch and display the plan details
	plan, err := data.FetchPlan(c, run.Plan.ID)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to fetch plan details"))
	}

	return v.Render(plan, &view.PlanCreateRenderOptions{
		RunID:        run.ID,
		PlanID:       run.Plan.ID,
		Hostname:     c.Hostname,
		Organization: c.OrganizationName,
		Workspace:    cmdConfig.WorkspaceName,
	})
}
