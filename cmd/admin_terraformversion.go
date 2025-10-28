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
	// `tfx admin terraform-version` commands
	tfvCmd = &cobra.Command{
		Use:     "terraform-version",
		Aliases: []string{"tfv"},
		Short:   "Terraform Version Commands",
		Long:    "Work with Terraform Versions in a TFE Installation",
		Example: `
List all Terraform versions:
tfx admin terraform-version list

Search for a specific version:
tfx admin terraform-version list --search 1.5

Create a custom Terraform version:
tfx admin terraform-version create --version 1.5.7 --url https://... --sha abc123...

Create an official Terraform version:
tfx admin terraform-version create official --version 1.5.7

Show a Terraform version:
tfx admin terraform-version show --version 1.5.7

Delete a Terraform version:
tfx admin terraform-version delete --version 1.5.7

Enable specific versions:
tfx admin terraform-version enable --versions 1.5.0,1.5.1

Disable specific versions:
tfx admin terraform-version disable --versions 1.4.0,1.4.1`,
	}

	// `tfx admin terraform-version list` command
	tfvListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Terraform Versions",
		Long:  "List Terraform Versions in a TFE Installation.",
		Example: `
tfx admin terraform-version list

tfx admin terraform-version list --search 1.5`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseAdminTerraformVersionListFlags(cmd)
			if err != nil {
				return err
			}
			return tfvList(cmdConfig)
		},
	}

	// `tfx admin terraform-version create` command
	tfvCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create Terraform Version",
		Long:  "Create a custom Terraform Version for a TFE Installation.",
		Example: `
tfx admin terraform-version create --version 1.5.7 --url https://example.com/terraform.zip --sha abc123...`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseAdminTerraformVersionCreateFlags(cmd)
			if err != nil {
				return err
			}
			return tfvCreate(cmdConfig)
		},
	}

	// `tfx admin terraform-version create official` command
	tfvCreateOfficialCmd = &cobra.Command{
		Use:   "official",
		Short: "Create Official Terraform Version",
		Long:  "Create a Terraform Version from official HashiCorp releases.",
		Example: `
tfx admin terraform-version create official --version 1.5.7

tfx admin terraform-version create official --version 1.6.0 --beta`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseAdminTerraformVersionCreateOfficialFlags(cmd)
			if err != nil {
				return err
			}
			return tfvCreateOfficial(cmdConfig)
		},
	}

	// `tfx admin terraform-version show` command
	tfvShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show Terraform Version",
		Long:  "Show Terraform Version details for a TFE Installation.",
		Example: `
tfx admin terraform-version show --version 1.5.7`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseAdminTerraformVersionShowFlags(cmd)
			if err != nil {
				return err
			}
			return tfvShow(cmdConfig)
		},
	}

	// `tfx admin terraform-version delete` command
	tfvDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete Terraform Version",
		Long:  "Delete a Terraform Version from a TFE Installation.",
		Example: `
tfx admin terraform-version delete --version 1.5.7`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseAdminTerraformVersionDeleteFlags(cmd)
			if err != nil {
				return err
			}
			return tfvDelete(cmdConfig)
		},
	}

	// `tfx admin terraform-version disable` command
	tfvDisableCmd = &cobra.Command{
		Use:   "disable",
		Short: "Disable Terraform Versions",
		Long:  "Disable one or more Terraform Versions in a TFE Installation.",
		Example: `
tfx admin terraform-version disable --versions 1.4.0,1.4.1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseAdminTerraformVersionEnableDisableFlags(cmd)
			if err != nil {
				return err
			}
			return tfvDisable(cmdConfig)
		},
	}

	// `tfx admin terraform-version disable all` command
	tfvDisableAllCmd = &cobra.Command{
		Use:   "all",
		Short: "Disable All Terraform Versions",
		Long:  "Disable All Terraform Versions in a TFE Installation.",
		Example: `
tfx admin terraform-version disable all`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseAdminTerraformVersionEnableDisableFlags(cmd)
			if err != nil {
				return err
			}
			return tfvDisableAll(cmdConfig)
		},
	}

	// `tfx admin terraform-version enable` command
	tfvEnableCmd = &cobra.Command{
		Use:   "enable",
		Short: "Enable Terraform Versions",
		Long:  "Enable one or more Terraform Versions in a TFE Installation.",
		Example: `
tfx admin terraform-version enable --versions 1.5.0,1.5.1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseAdminTerraformVersionEnableDisableFlags(cmd)
			if err != nil {
				return err
			}
			return tfvEnable(cmdConfig)
		},
	}

	// `tfx admin terraform-version enable all` command
	tfvEnableAllCmd = &cobra.Command{
		Use:   "all",
		Short: "Enable All Terraform Versions",
		Long:  "Enable All Terraform Versions in a TFE Installation.",
		Example: `
tfx admin terraform-version enable all`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseAdminTerraformVersionEnableDisableFlags(cmd)
			if err != nil {
				return err
			}
			return tfvEnableAll(cmdConfig)
		},
	}
)

func init() {
	// `tfx admin terraform-version list` flags
	tfvListCmd.Flags().StringP("search", "s", "", "Search string for partial version string (optional)")

	// `tfx admin terraform-version show` flags
	tfvShowCmd.Flags().StringP("version", "v", "", "Terraform Version (e.g., 1.5.0)")
	tfvShowCmd.MarkFlagRequired("version")

	// `tfx admin terraform-version create` flags
	tfvCreateCmd.Flags().StringP("version", "v", "", "Version of Terraform (e.g., 1.5.0)")
	tfvCreateCmd.Flags().StringP("url", "u", "", "URL of a hosted file containing Terraform (e.g., https://terraform.io...)")
	tfvCreateCmd.Flags().StringP("sha", "s", "", "SHA-256 checksum of the file at the url (must be 64 characters long)")
	tfvCreateCmd.Flags().BoolP("official", "", false, "Terraform Version is official (optional)")
	tfvCreateCmd.Flags().BoolP("disable", "", false, "Created Terraform Version will be disabled (optional)")
	tfvCreateCmd.Flags().BoolP("beta", "", false, "Terraform Version is beta (optional)")
	tfvCreateCmd.MarkFlagRequired("version")
	tfvCreateCmd.MarkFlagRequired("url")
	tfvCreateCmd.MarkFlagRequired("sha")

	// `tfx admin terraform-version create official` flags
	tfvCreateOfficialCmd.Flags().StringP("version", "v", "", "Version of Terraform (e.g., 1.5.0)")
	tfvCreateOfficialCmd.Flags().BoolP("disable", "", false, "Created Terraform Version will be disabled (optional)")
	tfvCreateOfficialCmd.Flags().BoolP("beta", "", false, "Terraform Version is beta (optional)")
	tfvCreateOfficialCmd.MarkFlagRequired("version")

	// `tfx admin terraform-version delete` flags
	tfvDeleteCmd.Flags().StringP("version", "v", "", "Terraform Version (e.g., 1.5.0)")
	tfvDeleteCmd.MarkFlagRequired("version")

	// `tfx admin terraform-version disable` flags
	tfvDisableCmd.Flags().StringSliceP("versions", "v", []string{}, "Versions to disable, can be comma separated (e.g., 1.4.0,1.4.1)")
	tfvDisableCmd.MarkFlagRequired("versions")

	// `tfx admin terraform-version enable` flags
	tfvEnableCmd.Flags().StringSliceP("versions", "v", []string{}, "Versions to enable, can be comma separated (e.g., 1.5.0,1.5.1)")
	tfvEnableCmd.MarkFlagRequired("versions")

	adminCmd.AddCommand(tfvCmd)
	tfvCmd.AddCommand(tfvListCmd)
	tfvCmd.AddCommand(tfvCreateCmd)
	tfvCreateCmd.AddCommand(tfvCreateOfficialCmd)
	tfvCmd.AddCommand(tfvShowCmd)
	tfvCmd.AddCommand(tfvDeleteCmd)
	tfvCmd.AddCommand(tfvDisableCmd)
	tfvDisableCmd.AddCommand(tfvDisableAllCmd)
	tfvCmd.AddCommand(tfvEnableCmd)
	tfvEnableCmd.AddCommand(tfvEnableAllCmd)
}

func tfvList(cmdConfig *flags.AdminTerraformVersionListFlags) error {
	// Create view for rendering
	v := view.NewAdminTerraformVersionListView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	// Print command header
	if cmdConfig.Search != "" {
		v.PrintCommandHeader("Listing Terraform versions matching '%s'", cmdConfig.Search)
	} else {
		v.PrintCommandHeader("Listing all Terraform versions")
	}

	// Fetch Terraform versions
	versions, err := data.FetchTerraformVersions(c, "", cmdConfig.Search)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to list terraform versions"))
	}

	return v.Render(versions)
}

func tfvShow(cmdConfig *flags.AdminTerraformVersionShowFlags) error {
	// Create view for rendering
	v := view.NewAdminTerraformVersionShowView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	// Print command header
	v.PrintCommandHeader("Showing Terraform version '%s'", cmdConfig.Version)

	// Fetch Terraform version
	tfv, err := data.FetchTerraformVersion(c, cmdConfig.Version)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to find terraform version"))
	}

	return v.Render(tfv)
}

func tfvCreate(cmdConfig *flags.AdminTerraformVersionCreateFlags) error {
	// Create view for rendering
	v := view.NewAdminTerraformVersionCreateView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	// Print command header
	v.PrintCommandHeader("Creating Terraform version '%s'", cmdConfig.Version)

	// Create Terraform version
	tfv, err := data.CreateTerraformVersion(c, cmdConfig.Version, cmdConfig.URL, cmdConfig.SHA, cmdConfig.Official, cmdConfig.Enabled, cmdConfig.Beta)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "unable to create terraform version"))
	}

	return v.Render(tfv)
}

func tfvCreateOfficial(cmdConfig *flags.AdminTerraformVersionCreateOfficialFlags) error {
	// Create view for rendering
	v := view.NewAdminTerraformVersionCreateView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	// Print command header
	v.PrintCommandHeader("Creating official Terraform version '%s'", cmdConfig.Version)
	v.PrintCommandHeader("Searching for official Terraform version in HashiCorp releases...")

	// Create official Terraform version
	tfv, err := data.CreateOfficialTerraformVersion(c, cmdConfig.Version, cmdConfig.Enabled, cmdConfig.Beta)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "unable to create terraform version"))
	}

	return v.Render(tfv)
}

func tfvDelete(cmdConfig *flags.AdminTerraformVersionDeleteFlags) error {
	// Create view for rendering
	v := view.NewAdminTerraformVersionDeleteView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	// Print command header
	v.PrintCommandHeader("Deleting Terraform version '%s'", cmdConfig.Version)

	// Delete Terraform version
	err = data.DeleteTerraformVersion(c, cmdConfig.Version)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to delete version"))
	}

	return v.Render()
}

func tfvDisable(cmdConfig *flags.AdminTerraformVersionEnableDisableFlags) error {
	// Create view for rendering
	v := view.NewAdminTerraformVersionUpdateView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	// Print command header
	v.PrintCommandHeader("Disabling Terraform versions: %v", cmdConfig.Versions)

	// Disable versions
	results, err := data.UpdateTerraformVersions(c, cmdConfig.Versions, false)
	if err != nil {
		return v.RenderError(err)
	}

	return v.Render(results)
}

func tfvDisableAll(cmdConfig *flags.AdminTerraformVersionEnableDisableFlags) error {
	// Create view for rendering
	v := view.NewAdminTerraformVersionUpdateView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	// Print command header
	v.PrintCommandHeader("Disabling all Terraform versions")

	// Fetch all versions
	items, err := data.FetchTerraformVersions(c, "", "")
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to list terraform versions"))
	}

	// Extract version strings
	versions := make([]string, len(items))
	for i, v := range items {
		versions[i] = v.Version
	}

	// Disable all versions
	results, err := data.UpdateTerraformVersions(c, versions, false)
	if err != nil {
		return v.RenderError(err)
	}

	return v.Render(results)
}

func tfvEnable(cmdConfig *flags.AdminTerraformVersionEnableDisableFlags) error {
	// Create view for rendering
	v := view.NewAdminTerraformVersionUpdateView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	// Print command header
	v.PrintCommandHeader("Enabling Terraform versions: %v", cmdConfig.Versions)

	// Enable versions
	results, err := data.UpdateTerraformVersions(c, cmdConfig.Versions, true)
	if err != nil {
		return v.RenderError(err)
	}

	return v.Render(results)
}

func tfvEnableAll(cmdConfig *flags.AdminTerraformVersionEnableDisableFlags) error {
	// Create view for rendering
	v := view.NewAdminTerraformVersionUpdateView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	// Print command header
	v.PrintCommandHeader("Enabling all Terraform versions")

	// Fetch all versions
	items, err := data.FetchTerraformVersions(c, "", "")
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to list terraform versions"))
	}

	// Extract version strings
	versions := make([]string, len(items))
	for i, v := range items {
		versions[i] = v.Version
	}

	// Enable all versions
	results, err := data.UpdateTerraformVersions(c, versions, true)
	if err != nil {
		return v.RenderError(err)
	}

	return v.Render(results)
}
