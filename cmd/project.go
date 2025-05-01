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

			if *viperBool("all") {
				return projectListAll(
					getTfxClientContext(),
					*viperString("search"))
			} else {
				return projectList(
					getTfxClientContext(),
					*viperString("search"))
			}
		},
	}
)

func init() {
	// `tfx project list`
	projectListCmd.Flags().StringP("search", "s", "", "Search string for Project Name (optional).")
	projectListCmd.Flags().BoolP("all", "a", false, "List All Organizations Projects (optional).")

	rootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectListCmd)

}

func projectListAll(c TfxClientContext, searchString string) error {
	o.AddMessageUserProvided("List Projects for all available Organizations", "")
	orgs, err := organizationPrjListAll(c, searchString)
	if err != nil {
		logError(err, "failed to list organizations")
	}

	var allProjectList []*tfe.Project
	for _, v := range orgs {
		projectList, err := projectListAllForOrganization(c, v.Name, searchString)
		if err != nil {
			logError(err, "failed to list projects for organization")
		}

		allProjectList = append(allProjectList, projectList...)
	}

	o.AddTableHeader("Organization", "Name", "Id", "Description")
	for _, i := range allProjectList {
		o.AddTableRows(i.Organization.Name, i.Name, i.ID, i.Description)
	}

	return nil
}

func projectListAllForOrganization(c TfxClientContext, orgName string, searchString string) ([]*tfe.Project, error) {
	allItems := []*tfe.Project{}
	opts := tfe.ProjectListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 1, PageSize: 100},
		Query:       searchString,
	}
	for {
		items, err := c.Client.Projects.List(c.Context, orgName, &opts)
		if err != nil {
			return nil, err
		}

		allItems = append(allItems, items.Items...)
		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage
	}

	return allItems, nil
}

func organizationPrjListAll(c TfxClientContext, searchString string) ([]*tfe.Organization, error) {
	allItems := []*tfe.Organization{}
	opts := tfe.OrganizationListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100}}
	for {
		items, err := c.Client.Organizations.List(c.Context, &opts)
		if err != nil {
			return nil, err
		}

		allItems = append(allItems, items.Items...)
		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage
	}

	return allItems, nil
}

func projectList(c TfxClientContext, searchString string) error {
	o.AddMessageUserProvided("List Projects for Organization:", c.OrganizationName)
	items, err := projectListAllForOrganization(c, c.OrganizationName, searchString)
	if err != nil {
		return errors.Wrap(err, "failed to list projects")
	}

	o.AddTableHeader("Name", "Id", "Description")
	for _, i := range items {
		o.AddTableRows(i.Name, i.ID, i.Description)
	}

	return nil
}
