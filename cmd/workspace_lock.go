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
	// `tfx workspace lock` command
	workspaceLockCmd = &cobra.Command{
		Use:   "lock",
		Short: "Lock a Workspace",
		Long:  "Lock a Workspace in a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseWorkspaceLockFlags(cmd)
			if err != nil {
				return err
			}
			return workspaceLock(cmdConfig)
		},
	}

	workspaceLockAllCmd = &cobra.Command{
		Use:   "all",
		Short: "Lock All Workspaces",
		Long:  "Lock All Workspaces in a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseWorkspaceLockAllFlags(cmd)
			if err != nil {
				return err
			}
			return workspaceLockAll(cmdConfig)
		},
	}

	// `tfx workspace unlock` command
	workspaceUnlockCmd = &cobra.Command{
		Use:   "unlock",
		Short: "Unlock a Workspace",
		Long:  "Unlock a Workspace in a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseWorkspaceUnlockFlags(cmd)
			if err != nil {
				return err
			}
			return workspaceUnlock(cmdConfig)
		},
	}

	// `tfx workspace unlock all` command
	workspaceUnlockAllCmd = &cobra.Command{
		Use:   "all",
		Short: "Unlock All Workspaces",
		Long:  "Unlock All Workspaces in a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseWorkspaceUnlockAllFlags(cmd)
			if err != nil {
				return err
			}
			return workspaceUnlockAll(cmdConfig)
		},
	}
)

func init() {
	// `tfx workspace lock`
	workspaceLockCmd.Flags().StringP("name", "n", "", "Workspace name")
	workspaceLockCmd.MarkFlagRequired("name")

	// `tfx workspace lock all`
	workspaceLockAllCmd.Flags().StringP("search", "s", "", "Search string for Workspace Name (optional).")

	// `tfx workspace unlock`
	workspaceUnlockCmd.Flags().StringP("name", "n", "", "Workspace name")
	workspaceUnlockCmd.MarkFlagRequired("name")

	// `tfx workspace unlock all`
	workspaceUnlockAllCmd.Flags().StringP("search", "s", "", "Search string for Workspace Name (optional).")

	workspaceCmd.AddCommand(workspaceLockCmd)
	workspaceLockCmd.AddCommand(workspaceLockAllCmd)
	workspaceCmd.AddCommand(workspaceUnlockCmd)
	workspaceUnlockCmd.AddCommand(workspaceUnlockAllCmd)
}

func workspaceLock(cmdConfig *flags.WorkspaceLockFlags) error {
	v := view.NewWorkspaceLockView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Locking workspace '%s'", cmdConfig.Name)

	status, err := data.SetWorkspaceLock(c, c.OrganizationName, cmdConfig.Name, true)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "unable to lock workspace"))
	}

	return v.RenderSingle(cmdConfig.Name, status)
}

func workspaceLockAll(cmdConfig *flags.WorkspaceLockAllFlags) error {
	v := view.NewWorkspaceLockView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Locking workspaces in organization '%s'", c.OrganizationName)
	if cmdConfig.Search != "" {
		v.PrintCommandFilter("search: %s", cmdConfig.Search)
	}

	opts := &flags.WorkspaceListFlags{Search: cmdConfig.Search}
	workspaces, err := data.FetchWorkspaces(c, c.OrganizationName, opts)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to list workspaces"))
	}

	results := make([]view.WorkspaceLockResult, 0, len(workspaces))
	for _, ws := range workspaces {
		status, err := data.SetWorkspaceLock(c, c.OrganizationName, ws.Name, true)
		if err != nil {
			results = append(results, view.WorkspaceLockResult{Name: ws.Name, Status: err.Error()})
		} else {
			results = append(results, view.WorkspaceLockResult{Name: ws.Name, Status: status})
		}
	}

	return v.RenderBulk(results)
}

func workspaceUnlock(cmdConfig *flags.WorkspaceUnlockFlags) error {
	v := view.NewWorkspaceLockView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Unlocking workspace '%s'", cmdConfig.Name)

	status, err := data.SetWorkspaceLock(c, c.OrganizationName, cmdConfig.Name, false)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "unable to unlock workspace"))
	}

	return v.RenderSingle(cmdConfig.Name, status)
}

func workspaceUnlockAll(cmdConfig *flags.WorkspaceUnlockAllFlags) error {
	v := view.NewWorkspaceLockView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Unlocking workspaces in organization '%s'", c.OrganizationName)
	if cmdConfig.Search != "" {
		v.PrintCommandFilter("search: %s", cmdConfig.Search)
	}

	opts := &flags.WorkspaceListFlags{Search: cmdConfig.Search}
	workspaces, err := data.FetchWorkspaces(c, c.OrganizationName, opts)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to list workspaces"))
	}

	results := make([]view.WorkspaceLockResult, 0, len(workspaces))
	for _, ws := range workspaces {
		status, err := data.SetWorkspaceLock(c, c.OrganizationName, ws.Name, false)
		if err != nil {
			results = append(results, view.WorkspaceLockResult{Name: ws.Name, Status: err.Error()})
		} else {
			results = append(results, view.WorkspaceLockResult{Name: ws.Name, Status: status})
		}
	}

	return v.RenderBulk(results)
}
