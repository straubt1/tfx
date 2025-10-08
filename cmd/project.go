/*
Copyright © 2021 Tom Straub <github.com/straubt1>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/straubt1/tfx/client"
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
			c, err := client.NewFromViper()
			if err != nil {
				return err
			}

			if *viperBool("all") {
				return projectListAll(
					c,
					*viperString("search"))
			} else {
				return projectList(
					c,
					*viperString("search"))
			}
		},
	}

	// `tfx project show` command
	projectShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show project details",
		Long:  "Show Project in a TFx Organization.",
		Example: `tfx project show --id prj-abc123
tfx project show --name myprojectname`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.NewFromViper()
			if err != nil {
				return err
			}

			projectID := *viperString("id")
			projectName := *viperString("name")

			// Validate that exactly one of id or name is provided
			if projectID == "" && projectName == "" {
				return errors.New("either --id or --name must be provided")
			}
			if projectID != "" && projectName != "" {
				return errors.New("only one of --id or --name can be provided")
			}

			return projectShow(c, projectID, projectName)
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

	rootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectShowCmd)

}

func projectListAll(c *client.TfxClient, searchString string) error {
	o.AddMessageUserProvided("List Projects for all available Organizations", "")

	projects, err := c.FetchProjectsAcrossOrgs(searchString)
	if err != nil {
		return err
	}

	o.AddTableHeader("Organization", "Name", "Id", "Description")
	for _, p := range projects {
		o.AddTableRows(p.Organization.Name, p.Name, p.ID, p.Description)
	}

	return nil
}

func projectList(c *client.TfxClient, searchString string) error {
	o.AddMessageUserProvided("List Projects for Organization:", c.OrganizationName)

	projects, err := c.FetchProjects(c.OrganizationName, searchString)
	if err != nil {
		return errors.Wrap(err, "failed to list projects")
	}

	o.AddTableHeader("Name", "Id", "Description")
	for _, p := range projects {
		o.AddTableRows(p.Name, p.ID, p.Description)
	}

	return nil
}

func projectShow(c *client.TfxClient, projectID string, projectName string) error {
	var p *tfe.Project
	var err error

	readOptions := &tfe.ProjectReadOptions{
		Include: []tfe.ProjectIncludeOpt{
			tfe.ProjectEffectiveTagBindings,
		},
	}

	if projectID != "" {
		o.AddMessageUserProvided("Project ID:", projectID)
		p, err = c.FetchProject(projectID, readOptions)
	} else {
		o.AddMessageUserProvided("Project Name:", projectName)
		p, err = c.FetchProjectByName(c.OrganizationName, projectName, readOptions)
	}

	if err != nil {
		logError(err, "failed to read project")
	}

	o.AddDeferredMessageRead("Name", p.Name)
	o.AddDeferredMessageRead("ID", p.ID)
	o.AddDeferredMessageRead("Description", p.Description)
	o.AddDeferredMessageRead("DefaultExecutionMode", p.DefaultExecutionMode)

	var duration string
	if p.AutoDestroyActivityDuration.IsSpecified() {
		if duration, err = p.AutoDestroyActivityDuration.Get(); err == nil {
		}
	}
	o.AddDeferredMessageRead("Auto Destroy Activity Duration", duration)

	tags := make(map[string]interface{})
	for _, i := range p.EffectiveTagBindings {
		tags[i.Key] = i.Value
	}
	o.AddDeferredMapMessageRead("Tags", tags)

	return nil
}
