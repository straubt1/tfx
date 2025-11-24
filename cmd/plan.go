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
	// `tfx plan` commands
	planCmd = &cobra.Command{
		Use:   "plan",
		Short: "Plans",
		Long:  "Work with Plans of a TFx Workspace.",
	}

	// `tfx plan show` command
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

	// `tfx plan logs` command
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

	// `tfx plan jsonoutput` command
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
)

func init() {
	// `tfx plan show` command
	planShowCmd.Flags().StringP("id", "i", "", "Plan Id (i.e. plan-*)")
	planShowCmd.MarkFlagRequired("id")

	// `tfx plan logs` command
	planLogsCmd.Flags().StringP("id", "i", "", "Plan Id (i.e. plan-*)")
	planLogsCmd.MarkFlagRequired("id")

	// `tfx plan jsonoutput` command
	planJSONOutputCmd.Flags().StringP("id", "i", "", "Plan Id (i.e. plan-*)")
	planJSONOutputCmd.MarkFlagRequired("id")

	rootCmd.AddCommand(planCmd)
	planCmd.AddCommand(planShowCmd)
	planCmd.AddCommand(planLogsCmd)
	planCmd.AddCommand(planJSONOutputCmd)
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

	planExport, err := c.Client.PlanExports.Create(c.Context, tfe.PlanExportCreateOptions{
		Plan:     plan,
		DataType: tfe.PlanExportType(tfe.PlanExportSentinelMockBundleV0),
	})

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
