package data

import (
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/logger"
)

// FetchProjects fetches all projects for a given organization using pagination
func FetchProjects(c *client.TfxClient, orgName string, searchString string) ([]*tfe.Project, error) {
	logger.Debug("Fetching projects", "organization", orgName, "searchString", searchString)

	return client.FetchAll(c.Context, func(pageNumber int) ([]*tfe.Project, *client.Pagination, error) {
		logger.Trace("Fetching projects page", "organization", orgName, "page", pageNumber)

		opts := &tfe.ProjectListOptions{
			ListOptions: tfe.ListOptions{PageNumber: pageNumber, PageSize: 100},
			Query:       searchString,
			Include: []tfe.ProjectIncludeOpt{
				tfe.ProjectEffectiveTagBindings,
			},
		}

		result, err := c.Client.Projects.List(c.Context, orgName, opts)
		if err != nil {
			logger.Error("Failed to fetch projects page", "organization", orgName, "page", pageNumber, "error", err)
			return nil, nil, err
		}

		logger.Trace("Projects page fetched", "organization", orgName, "page", pageNumber, "count", len(result.Items))
		return result.Items, client.NewPaginationFromTFE(result.Pagination), nil
	})
}

// FetchProjectsAcrossOrgs fetches projects across all organizations
func FetchProjectsAcrossOrgs(c *client.TfxClient, searchString string) ([]*tfe.Project, error) {
	logger.Info("Fetching projects across all organizations", "searchString", searchString)

	orgs, err := FetchOrganizations(c, "")
	if err != nil {
		logger.Error("Failed to fetch organizations", "error", err)
		return nil, errors.Wrap(err, "failed to list organizations")
	}

	logger.Debug("Organizations fetched", "count", len(orgs))

	var allProjects []*tfe.Project
	for _, org := range orgs {
		logger.Debug("Fetching projects for organization", "organization", org.Name)

		projects, err := FetchProjects(c, org.Name, searchString)
		if err != nil {
			logger.Error("Failed to fetch projects", "organization", org.Name, "error", err)
			return nil, errors.Wrapf(err, "failed to list projects for organization %s", org.Name)
		}

		logger.Debug("Projects fetched for organization", "organization", org.Name, "count", len(projects))
		allProjects = append(allProjects, projects...)
	}

	logger.Info("All projects fetched successfully", "totalProjects", len(allProjects), "organizations", len(orgs))
	return allProjects, nil
}

// FetchProject fetches a single project by ID
func FetchProject(c *client.TfxClient, projectID string) (*tfe.Project, error) {
	logger.Debug("Fetching project by ID", "projectID", projectID)

	options := &tfe.ProjectReadOptions{
		Include: []tfe.ProjectIncludeOpt{
			tfe.ProjectEffectiveTagBindings,
		},
	}

	project, err := c.Client.Projects.ReadWithOptions(c.Context, projectID, *options)
	if err != nil {
		logger.Error("Failed to fetch project", "projectID", projectID, "error", err)
		return nil, err
	}

	logger.Debug("Project fetched successfully", "projectID", projectID, "name", project.Name)
	return project, nil
}

// FetchProjectByName fetches a single project by name in the specified organization
func FetchProjectByName(c *client.TfxClient, orgName string, projectName string) (*tfe.Project, error) {
	logger.Debug("Fetching project by name", "organization", orgName, "projectName", projectName)

	projects, err := FetchProjects(c, orgName, projectName)
	if err != nil {
		logger.Error("Failed to fetch projects", "organization", orgName, "projectName", projectName, "error", err)
		return nil, errors.Wrap(err, "failed to fetch projects")
	}

	logger.Trace("Projects search completed", "organization", orgName, "resultsCount", len(projects))

	// Find exact match in case there are multiple results from the search
	for _, p := range projects {
		if p.Name == projectName {
			logger.Debug("Project found by name", "organization", orgName, "projectName", projectName, "projectID", p.ID)
			// Fetch full project details with options
			return FetchProject(c, p.ID)
		}
	}

	logger.Warn("Project not found by name", "organization", orgName, "projectName", projectName)
	return nil, errors.Errorf("project with name '%s' not found in organization '%s'", projectName, orgName)
}
