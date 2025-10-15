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

// organizationCmd represents the organization command
var (
	// `tfx organization` commands
	organizationCmd = &cobra.Command{
		Use:     "organization",
		Aliases: []string{"org"},
		Short:   "Organization Commands",
		Long:    "Work with TFx Organizations",
		Example: `
List all Organizations:
tfx organization list

Show an Organization by name:
tfx organization show --name "my-org"`,
	}

	// `tfx organization list` command
	organizationListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Organizations",
		Long:  "List all Organizations available to the authenticated user.",
		Example: `
tfx organization list
tfx organization list --search "my-org"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseOrganizationListFlags(cmd)
			if err != nil {
				return err
			}
			return organizationList(cmdConfig)
		},
	}

	// `tfx organization show` command
	organizationShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show organization details",
		Long:  "Show Organization details.",
		Example: `
tfx organization show --name myorganization
tfx organization show -n myorganization`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseOrganizationShowFlags(cmd)
			if err != nil {
				return err
			}
			return organizationShow(cmdConfig)
		},
	}
)

func init() {
	// `tfx organization list`
	organizationListCmd.Flags().StringP("search", "s", "", "Search string for Organization Name (optional).")

	// `tfx organization show`
	organizationShowCmd.Flags().StringP("name", "n", "", "Name of the organization.")
	organizationShowCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(organizationCmd)
	organizationCmd.AddCommand(organizationListCmd)
	organizationCmd.AddCommand(organizationShowCmd)
}

func organizationList(cmdConfig *flags.OrganizationListFlags) error {
	// Create view for rendering
	v := view.NewOrganizationListView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	// Print command header before API call
	if cmdConfig.Search != "" {
		v.PrintCommandHeader("Listing organizations matching '%s'", cmdConfig.Search)
	} else {
		v.PrintCommandHeader("Listing all organizations")
	}

	orgs, err := data.FetchOrganizations(c, cmdConfig.Search)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to list organizations"))
	}

	return v.Render(orgs)
}

func organizationShow(cmdConfig *flags.OrganizationShowFlags) error {
	// Create view for rendering
	v := view.NewOrganizationShowView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Showing organization '%s'", cmdConfig.Name)

	readOptions := &tfe.OrganizationReadOptions{
		Include: []tfe.OrganizationIncludeOpt{},
	}

	org, err := data.FetchOrganization(c, cmdConfig.Name, readOptions)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to read organization"))
	}

	return v.Render(org)
}
