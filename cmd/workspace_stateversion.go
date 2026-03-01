// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/cmd/flags"
	view "github.com/straubt1/tfx/cmd/views"
	"github.com/straubt1/tfx/data"
	pkgfile "github.com/straubt1/tfx/pkg/file"
)

var (
	// `tfx workspace state` commands
	stateCmd = &cobra.Command{
		Use:     "state-version",
		Aliases: []string{"sv"},
		Short:   "State Version Commands",
		Long:    "Work with State Versions of a TFx Workspace.",
	}

	// `tfx workspace state list` command
	stateListCmd = &cobra.Command{
		Use:   "list",
		Short: "List State Versions",
		Long:  "List State Versions of a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseStateListFlags(cmd)
			if err != nil {
				return err
			}
			return stateList(cmdConfig)
		},
	}

	// `tfx workspace state create` command
	stateCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create State Version",
		Long:  "Create State Version for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseStateCreateFlags(cmd)
			if err != nil {
				return err
			}
			return stateCreate(cmdConfig)
		},
	}

	// `tfx workspace state show` command
	stateShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show State Version",
		Long:  "Show State Version details for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseStateShowFlags(cmd)
			if err != nil {
				return err
			}
			return stateShow(cmdConfig)
		},
	}

	// `tfx workspace state download` command
	stateDownloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download State Version",
		Long:  "Download State Version for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseStateDownloadFlags(cmd)
			if err != nil {
				return err
			}
			return stateDownload(cmdConfig)
		},
	}
)

func init() {
	// `tfx workspace state list` command
	stateListCmd.Flags().StringP("name", "n", "", "Workspace name")
	stateListCmd.Flags().IntP("max-items", "m", 10, "Max number of results (optional)")
	stateListCmd.MarkFlagRequired("name")

	// `tfx workspace state create` command
	stateCreateCmd.Flags().StringP("name", "n", "", "Workspace name")
	stateCreateCmd.Flags().StringP("filename", "f", "", "Filename of the state file to create")
	stateCreateCmd.MarkFlagRequired("name")
	stateCreateCmd.MarkFlagRequired("filename")

	// `tfx workspace state show` command
	stateShowCmd.Flags().StringP("state-id", "i", "", "State Version Id (i.e. sv-*)")
	stateShowCmd.MarkFlagRequired("state-id")

	// `tfx workspace state download` command
	stateDownloadCmd.Flags().StringP("state-id", "i", "", "State Version Id (i.e. sv-*)")
	stateDownloadCmd.Flags().StringP("directory", "d", "", "Directory of download state version (optional, defaults to a temp directory)")
	stateDownloadCmd.Flags().StringP("filename", "f", "", "Filename to save State Version as (optional)")
	stateDownloadCmd.MarkFlagRequired("state-id")

	workspaceCmd.AddCommand(stateCmd)
	stateCmd.AddCommand(stateListCmd)
	stateCmd.AddCommand(stateDownloadCmd)
	stateCmd.AddCommand(stateCreateCmd)
	stateCmd.AddCommand(stateShowCmd)
}

func stateList(cmdConfig *flags.StateListFlags) error {
	v := view.NewStateVersionListView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Listing state versions for workspace '%s'", cmdConfig.WorkspaceName)

	items, err := data.FetchStateVersions(c, c.OrganizationName, cmdConfig.WorkspaceName, cmdConfig.MaxItems)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to list state versions"))
	}

	return v.Render(items)
}

func stateCreate(cmdConfig *flags.StateCreateFlags) error {
	v := view.NewStateVersionCreateView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Creating state version for workspace '%s'", cmdConfig.WorkspaceName)

	sv, err := data.CreateStateVersionFromFile(c, c.OrganizationName, cmdConfig.WorkspaceName, cmdConfig.Filename)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to create state version"))
	}

	return v.Render(sv)
}

func stateShow(cmdConfig *flags.StateShowFlags) error {
	v := view.NewStateVersionShowView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Showing state version '%s'", cmdConfig.StateID)

	sv, err := data.FetchStateVersion(c, cmdConfig.StateID)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to read state version"))
	}

	return v.Render(sv)
}

func stateDownload(cmdConfig *flags.StateDownloadFlags) error {
	v := view.NewStateVersionDownloadView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Downloading state version '%s'", cmdConfig.StateID)

	// Determine filename to save
	directory, err := pkgfile.GetDirectory(cmdConfig.Directory, cmdConfig.StateID)
	if err != nil {
		return v.RenderError(err)
	}
	filename := cmdConfig.Filename
	if filename == "" {
		filename = directory + "/" + cmdConfig.StateID + ".state"
	}

	buff, err := data.DownloadStateVersion(c, cmdConfig.StateID)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to download state version"))
	}

	if err := os.WriteFile(filename, buff, 0644); err != nil {
		return v.RenderError(errors.Wrap(err, "failed to save state version"))
	}

	return v.Render(filename)
}
