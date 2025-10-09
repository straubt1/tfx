package data

import (
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/straubt1/tfx/client"
)

// FetchProjects fetches all projects for a given organization using pagination
func FetchProjects(c *client.TfxClient, orgName string, searchString string) ([]*tfe.Project, error) {
	return client.FetchAll(c.Context, func(pageNumber int) ([]*tfe.Project, *client.Pagination, error) {
		opts := &tfe.ProjectListOptions{
			ListOptions: tfe.ListOptions{PageNumber: pageNumber, PageSize: 100},
			Query:       searchString,
			Include: []tfe.ProjectIncludeOpt{
				tfe.ProjectEffectiveTagBindings,
			},
		}

		result, err := c.Client.Projects.List(c.Context, orgName, opts)
		if err != nil {
			return nil, nil, err
		}

		return result.Items, client.NewPaginationFromTFE(result.Pagination), nil
	})
}

// FetchProjectsAcrossOrgs fetches projects across all organizations
func FetchProjectsAcrossOrgs(c *client.TfxClient, searchString string) ([]*tfe.Project, error) {
	orgs, err := FetchOrganizations(c, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to list organizations")
	}

	var allProjects []*tfe.Project
	for _, org := range orgs {
		projects, err := FetchProjects(c, org.Name, searchString)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to list projects for organization %s", org.Name)
		}
		allProjects = append(allProjects, projects...)
	}

	return allProjects, nil
}

// FetchProject fetches a single project by ID
func FetchProject(c *client.TfxClient, projectID string) (*tfe.Project, error) {
	options := &tfe.ProjectReadOptions{
		Include: []tfe.ProjectIncludeOpt{
			tfe.ProjectEffectiveTagBindings,
		},
	}
	return c.Client.Projects.ReadWithOptions(c.Context, projectID, *options)
}

// FetchProjectByName fetches a single project by name in the specified organization
func FetchProjectByName(c *client.TfxClient, orgName string, projectName string) (*tfe.Project, error) {
	projects, err := FetchProjects(c, orgName, projectName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch projects")
	}

	// Find exact match in case there are multiple results from the search
	for _, p := range projects {
		if p.Name == projectName {
			// Fetch full project details with options
			return FetchProject(c, p.ID)
		}
	}

	return nil, errors.Errorf("project with name '%s' not found in organization '%s'", projectName, orgName)
}
