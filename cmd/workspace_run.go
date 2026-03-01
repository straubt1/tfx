// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/cmd/flags"
	view "github.com/straubt1/tfx/cmd/views"
	"github.com/straubt1/tfx/data"
)

var (
	// `tfx workspace run` commands
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Workspace Runs",
		Long:  "Work with Runs of a TFx Workspace.",
	}

	// `tfx workspace run list` command
	runListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Runs",
		Long:  "List Runs of a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseRunListFlags(cmd)
			if err != nil {
				return err
			}
			return runList(cmdConfig)
		},
	}

	// `tfx workspace run create` command
	runCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create Run",
		Long:  "Create Run for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseRunCreateFlags(cmd)
			if err != nil {
				return err
			}
			return runCreate(cmdConfig)
		},
	}

	// `tfx workspace run show` command
	runShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show Run",
		Long:  "Show Run details for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseRunShowFlags(cmd)
			if err != nil {
				return err
			}
			return runShow(cmdConfig)
		},
	}

	// `tfx workspace run discard` command
	runDiscardCmd = &cobra.Command{
		Use:   "discard",
		Short: "Discard Run",
		Long:  "Discard Run for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseRunDiscardFlags(cmd)
			if err != nil {
				return err
			}
			return runDiscard(cmdConfig)
		},
	}

	// `tfx workspace run cancel` command
	runCancelCmd = &cobra.Command{
		Use:   "cancel",
		Short: "Cancel Latest Run",
		Long:  "Cancel Latest Run for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseRunCancelFlags(cmd)
			if err != nil {
				return err
			}
			return runCancel(cmdConfig)
		},
	}
)

func init() {
	// `tfx workspace run` commands

	// `tfx workspace run list` command
	runListCmd.Flags().StringP("name", "n", "", "Workspace name")
	runListCmd.Flags().IntP("max-items", "m", 10, "Max number of results (optional)")
	runListCmd.MarkFlagRequired("name")

	// `tfx workspace run create` command
	runCreateCmd.Flags().StringP("name", "n", "", "Workspace name")
	// runCreateCmd.Flags().StringP("directory", "d", "./", "Directory of Terraform (defaults to current directory)")
	runCreateCmd.Flags().StringP("message", "m", "", "Run Message (optional)")
	runCreateCmd.Flags().StringP("configuration-version-id", "i", "", "Configuration Version (optional)")
	runCreateCmd.MarkFlagRequired("name")

	// `tfx workspace run show` command
	runShowCmd.Flags().StringP("id", "i", "", "Run Id (i.e. run-*)")
	runShowCmd.MarkFlagRequired("id")

	// `tfx workspace run discard` command
	runDiscardCmd.Flags().StringP("id", "i", "", "Run Id (i.e. run-*)")
	runDiscardCmd.MarkFlagRequired("id")

	// `tfx workspace run cancel` command
	runCancelCmd.Flags().StringP("name", "n", "", "Workspace name")
	runCancelCmd.MarkFlagRequired("name")

	workspaceCmd.AddCommand(runCmd)
	runCmd.AddCommand(runListCmd)
	runCmd.AddCommand(runCreateCmd)
	runCmd.AddCommand(runShowCmd)
	runCmd.AddCommand(runDiscardCmd)
	runCmd.AddCommand(runCancelCmd)
}

func runList(cmdConfig *flags.RunListFlags) error {
	v := view.NewRunListView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Listing runs for workspace '%s'", cmdConfig.WorkspaceName)

	workspaceID, err := data.GetWorkspaceID(c, c.OrganizationName, cmdConfig.WorkspaceName)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "unable to read workspace id"))
	}

	runs, err := data.FetchRunsForWorkspace(c, workspaceID, cmdConfig.MaxItems)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to list runs"))
	}

	return v.Render(runs)
}

func runCreate(cmdConfig *flags.RunCreateFlags) error {
	v := view.NewRunCreateView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Creating run for workspace '%s'", cmdConfig.WorkspaceName)

	run, err := data.CreateRun(c, c.OrganizationName, cmdConfig.WorkspaceName, cmdConfig.Message, cmdConfig.ConfigurationVersionID)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to create run"))
	}

	// Build link to run
	link := fmt.Sprintf("https://%s/app/%s/workspaces/%s/runs/%s", c.Hostname, c.OrganizationName, cmdConfig.WorkspaceName, run.ID)

	return v.Render(run, link)
}

func runShow(cmdConfig *flags.RunShowFlags) error {
	v := view.NewRunShowView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Showing run '%s'", cmdConfig.ID)

	run, err := data.FetchRun(c, cmdConfig.ID)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to read run from id"))
	}

	return v.Render(run)
}

func runDiscard(cmdConfig *flags.RunDiscardFlags) error {
	v := view.NewRunDiscardView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Discarding run '%s'", cmdConfig.ID)

	if err := data.DiscardRun(c, cmdConfig.ID); err != nil {
		return v.RenderError(errors.Wrap(err, "failed to discard run"))
	}

	return v.Render(cmdConfig.ID)
}

func runCancel(cmdConfig *flags.RunCancelFlags) error {
	v := view.NewRunCancelView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Canceling latest run for workspace '%s'", cmdConfig.WorkspaceName)

	workspaceID, err := data.GetWorkspaceID(c, c.OrganizationName, cmdConfig.WorkspaceName)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "unable to read workspace id"))
	}

	runID, err := data.GetLatestRunID(c, workspaceID)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to get latest run id"))
	}

	if err := data.CancelRun(c, runID); err != nil {
		return v.RenderError(errors.Wrap(err, "failed to cancel run"))
	}

	return v.Render(runID)
}
