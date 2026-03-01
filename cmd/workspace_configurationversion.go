// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	"bytes"

	"github.com/hashicorp/go-slug"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/cmd/flags"
	view "github.com/straubt1/tfx/cmd/views"
	"github.com/straubt1/tfx/data"
	pkgfile "github.com/straubt1/tfx/pkg/file"
)

var (
	// `tfx workspace configuration-version` commands
	cvCmd = &cobra.Command{
		Use:     "configuration-version",
		Aliases: []string{"cv"},
		Short:   "Configuration Version Commands",
		Long:    "Work with Configuration Versions of a TFx Workspace.",
	}

	// `tfx workspace configuration-version list` command
	cvListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Configuration Versions",
		Long:  "List Configuration Versions of a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseCVListFlags(cmd)
			if err != nil {
				return err
			}
			return cvList(cmdConfig)
		},
	}

	// `tfx workspace configuration-version create` command
	cvCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create Configuration Version",
		Long:  "Create Configuration Version for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseCVCreateFlags(cmd)
			if err != nil {
				return err
			}
			return cvCreate(cmdConfig)
		},
	}

	// `tfx workspace configuration-version show` command
	cvShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show Configuration Version",
		Long:  "Show Configuration Version details for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseCVShowFlags(cmd)
			if err != nil {
				return err
			}
			return cvShow(cmdConfig)
		},
	}

	// `tfx workspace configuration-version download` command
	cvDownloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download the Configuration Version",
		Long:  "Download the Configuration Version code for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseCVDownloadFlags(cmd)
			if err != nil {
				return err
			}
			return cvDownload(cmdConfig)
		},
	}
)

func init() {
	// `tfx workspace configuration-version` commands

	// `tfx workspace configuration-version list` command
	cvListCmd.Flags().StringP("name", "n", "", "Workspace name")
	cvListCmd.Flags().IntP("max-items", "m", 10, "Max number of results (optional)")
	cvListCmd.MarkFlagRequired("name")

	// `tfx cv create`
	cvCreateCmd.Flags().StringP("name", "n", "", "Workspace name")
	cvCreateCmd.Flags().StringP("directory", "d", "./", "Directory of Terraform (optional, defaults to current directory)")
	cvCreateCmd.Flags().BoolP("speculative", "s", false, "Perform a Speculative Plan (optional, defaults to false)")
	cvCreateCmd.MarkFlagRequired("name")

	// `tfx cv show`
	cvShowCmd.Flags().StringP("id", "i", "", "Configuration Version Id (i.e. cv-*)")
	cvShowCmd.MarkFlagRequired("id")

	// `tfx cv download`
	cvDownloadCmd.Flags().StringP("id", "i", "", "Configuration Version Id (i.e. cv-*)")
	cvDownloadCmd.Flags().StringP("directory", "d", "", "Directory to download Configuration Version to (optional, defaults to a temp directory)")
	cvDownloadCmd.MarkFlagRequired("id")

	workspaceCmd.AddCommand(cvCmd)
	cvCmd.AddCommand(cvListCmd)
	cvCmd.AddCommand(cvCreateCmd)
	cvCmd.AddCommand(cvShowCmd)
	cvCmd.AddCommand(cvDownloadCmd)
}

func cvList(cmdConfig *flags.CVListFlags) error {
	v := view.NewConfigVersionListView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Listing configuration versions for workspace '%s'", cmdConfig.WorkspaceName)

	items, err := data.FetchConfigurationVersions(c, c.OrganizationName, cmdConfig.WorkspaceName, cmdConfig.MaxItems)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to list configuration versions"))
	}

	return v.Render(items)
}

func cvCreate(cmdConfig *flags.CVCreateFlags) error {
	v := view.NewConfigVersionCreateView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	if !pkgfile.IsDirectory(cmdConfig.Directory) {
		return v.RenderError(errors.New("directory file does not exist"))
	}

	v.PrintCommandHeader("Creating configuration version for workspace '%s'", cmdConfig.WorkspaceName)

	cv, err := data.CreateConfigurationVersion(c, c.OrganizationName, cmdConfig.WorkspaceName, cmdConfig.Directory, cmdConfig.Speculative)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to create configuration version"))
	}

	return v.Render(cv)
}

func cvShow(cmdConfig *flags.CVShowFlags) error {
	v := view.NewConfigVersionShowView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Showing configuration version '%s'", cmdConfig.ID)

	cv, err := data.FetchConfigurationVersion(c, cmdConfig.ID)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to read configuration version from provided id"))
	}

	return v.Render(cv)
}

func cvDownload(cmdConfig *flags.CVDownloadFlags) error {
	v := view.NewConfigVersionDownloadView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Downloading configuration version '%s'", cmdConfig.ID)

	directory, err := pkgfile.GetDirectory(cmdConfig.Directory, cmdConfig.ID)
	if err != nil {
		return v.RenderError(err)
	}

	buff, err := data.DownloadConfigurationVersion(c, cmdConfig.ID)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to download configuration version"))
	}

	// Unpack slug to directory
	reader := bytes.NewReader(buff)
	if err := slug.Unpack(reader, directory); err != nil {
		return v.RenderError(errors.Wrap(err, "failed to unpack configuration version slug"))
	}

	return v.Render(directory)
}
