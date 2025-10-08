package client

import (
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
)

// FetchProjects fetches all projects for a given organization using pagination
func (c *TfxClient) FetchProjects(orgName string, searchString string) ([]*tfe.Project, error) {
	return FetchAll(c.Context, func(pageNumber int) ([]*tfe.Project, *Pagination, error) {
		opts := &tfe.ProjectListOptions{
			ListOptions: tfe.ListOptions{PageNumber: pageNumber, PageSize: 100},
			Query:       searchString,
		}

		result, err := c.Client.Projects.List(c.Context, orgName, opts)
		if err != nil {
			return nil, nil, err
		}

		return result.Items, NewPaginationFromTFE(result.Pagination), nil
	})
}

// FetchProjectsAcrossOrgs fetches projects across all organizations
func (c *TfxClient) FetchProjectsAcrossOrgs(searchString string) ([]*tfe.Project, error) {
	orgs, err := c.FetchOrganizations()
	if err != nil {
		return nil, errors.Wrap(err, "failed to list organizations")
	}

	var allProjects []*tfe.Project
	for _, org := range orgs {
		projects, err := c.FetchProjects(org.Name, searchString)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to list projects for organization %s", org.Name)
		}
		allProjects = append(allProjects, projects...)
	}

	return allProjects, nil
}

// FetchProject fetches a single project by ID
func (c *TfxClient) FetchProject(projectID string, options *tfe.ProjectReadOptions) (*tfe.Project, error) {
	return c.Client.Projects.ReadWithOptions(c.Context, projectID, *options)
}

// FetchProjectByName fetches a single project by name in the specified organization
func (c *TfxClient) FetchProjectByName(orgName string, projectName string, options *tfe.ProjectReadOptions) (*tfe.Project, error) {
	projects, err := c.FetchProjects(orgName, projectName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch projects")
	}

	// Find exact match in case there are multiple results from the search
	for _, p := range projects {
		if p.Name == projectName {
			// Fetch full project details with options
			return c.FetchProject(p.ID, options)
		}
	}

	return nil, errors.Errorf("project with name '%s' not found in organization '%s'", projectName, orgName)
}
