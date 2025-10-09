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

// projectCmd represents the project command
var (
	// `tfx project` commands
	projectCmd = &cobra.Command{
		Use:     "project",
		Aliases: []string{"prj"},
		Short:   "Project Commands",
		Long:    "Work with TFx Projects",
		Example: `
List all Projects in all Organizations:
tfx project list --all

List all Projects in all Organizations with a search string:
tfx project list --all --search "my-project"

List all projects specified in tfeOrganization:
tfx project list

List projects specified in tfeOrganization with a search string:
tfx project list --search "my-project"`,
	}

	// `tfx project list` command
	projectListCmd = &cobra.Command{
		Use:     "list",
		Short:   "List Projects",
		Long:    "List Projects in a TFx Organization.",
		Example: `tfx project list --search "my-project"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseProjectListFlags(cmd)
			if err != nil {
				return err
			}

			if cmdConfig.All {
				return projectListAll(cmdConfig)
			} else {
				return projectList(cmdConfig)
			}
		},
	}

	// `tfx project show` command
	projectShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show project details",
		Long:  "Show Project in a TFx Organization.",
		Example: `
tfx project show --id prj-abc123
tfx project show --name myprojectname`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseProjectShowFlags(cmd)
			if err != nil {
				return err
			}
			return projectShow(cmdConfig)
		},
	}
)

func init() {
	// `tfx project list`
	projectListCmd.Flags().StringP("search", "s", "", "Search string for Project Name (optional).")
	projectListCmd.Flags().BoolP("all", "a", false, "List All Organizations Projects (optional).")

	// `tfx project show`
	projectShowCmd.Flags().StringP("id", "i", "", "ID of the project.")
	projectShowCmd.Flags().StringP("name", "n", "", "Name of the project.")
	projectShowCmd.MarkFlagsMutuallyExclusive("id", "name")
	projectShowCmd.MarkFlagsOneRequired("id", "name")

	rootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectShowCmd)
}

func projectListAll(cmdConfig *flags.ProjectListFlags) error {
	c, err := client.NewFromViper()
	if err != nil {
		return err
	}

	o.AddMessageUserProvided("List Projects for all available Organizations", "")
	projects, err := data.FetchProjectsAcrossOrgs(c, cmdConfig.Search)
	if err != nil {
		return err
	}

	o.AddTableHeader("Organization", "Name", "Id", "Description")
	for _, p := range projects {
		o.AddTableRows(p.Organization.Name, p.Name, p.ID, p.Description)
	}

	return nil
}

func projectList(cmdConfig *flags.ProjectListFlags) error {
	c, err := client.NewFromViper()
	if err != nil {
		return err
	}

	o.AddMessageUserProvided("List Projects for Organization:", c.OrganizationName)
	projects, err := data.FetchProjects(c, c.OrganizationName, cmdConfig.Search)
	if err != nil {
		return errors.Wrap(err, "failed to list projects")
	}

	o.AddTableHeader("Name", "Id", "Description")
	for _, p := range projects {
		o.AddTableRows(p.Name, p.ID, p.Description)
	}

	return nil
}

func projectShow(cmdConfig *flags.ProjectShowFlags) error {
	// Create view for rendering
	v := view.NewProjectShowView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	var p *tfe.Project
	if cmdConfig.ID != "" {
		v.PrintCommandHeader("Showing project '%s' in organization '%s'", cmdConfig.ID, c.OrganizationName)
		p, err = data.FetchProject(c, cmdConfig.ID)
	} else {
		v.PrintCommandHeader("Showing project '%s' in organization '%s'", cmdConfig.Name, c.OrganizationName)
		p, err = data.FetchProjectByName(c, c.OrganizationName, cmdConfig.Name)
	}

	if err != nil {
		return v.RenderError(err)
	}

	return v.Render(c.OrganizationName, p)
}
