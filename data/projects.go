// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/output"
)

// FetchProjects fetches all projects for a given organization using pagination
func FetchProjects(c *client.TfxClient, orgName string, searchString string) ([]*tfe.Project, error) {
	output.Get().Logger().Debug("Fetching projects", "organization", orgName, "searchString", searchString)

	return client.FetchAll(c.Context, func(pageNumber int) ([]*tfe.Project, *client.Pagination, error) {
		output.Get().Logger().Trace("Fetching projects page", "organization", orgName, "page", pageNumber)

		opts := &tfe.ProjectListOptions{
			ListOptions: tfe.ListOptions{PageNumber: pageNumber, PageSize: 100},
			Query:       searchString,
			Include: []tfe.ProjectIncludeOpt{
				tfe.ProjectEffectiveTagBindings,
			},
		}

		result, err := c.Client.Projects.List(c.Context, orgName, opts)
		if err != nil {
			output.Get().Logger().Error("Failed to fetch projects page", "organization", orgName, "page", pageNumber, "error", err)
			return nil, nil, err
		}

		output.Get().Logger().Trace("Projects page fetched", "organization", orgName, "page", pageNumber, "count", len(result.Items))
		return result.Items, client.NewPaginationFromTFE(result.Pagination), nil
	})
}

// FetchProjectsWithOrgScope fetches projects for either a single organization or across all organizations
// If allOrgs is true, it will fetch across all organizations. Otherwise it will fetch for the specified orgName.
func FetchProjectsWithOrgScope(c *client.TfxClient, orgName string, options *ProjectListOptions) ([]*tfe.Project, error) {
	if !options.All {
		// Fetch projects for a single organization
		output.Get().Logger().Info("Fetching projects for organization", "organization", orgName, "searchString", options.Search)
		return FetchProjects(c, orgName, options.Search)
	}

	// Fetch projects across all organizations
	output.Get().Logger().Info("Fetching projects across all organizations", "searchString", options.Search)

	orgs, err := FetchOrganizations(c, "")
	if err != nil {
		output.Get().Logger().Error("Failed to fetch organizations", "error", err)
		return nil, errors.Wrap(err, "failed to list organizations")
	}

	output.Get().Logger().Debug("Organizations fetched", "count", len(orgs))

	var allProjects []*tfe.Project
	for _, org := range orgs {
		output.Get().Logger().Debug("Fetching projects for organization", "organization", org.Name)

		projects, err := FetchProjects(c, org.Name, options.Search)
		if err != nil {
			output.Get().Logger().Error("Failed to fetch projects", "organization", org.Name, "error", err)
			return nil, errors.Wrapf(err, "failed to list projects for organization %s", org.Name)
		}

		output.Get().Logger().Debug("Projects fetched for organization", "organization", org.Name, "count", len(projects))
		allProjects = append(allProjects, projects...)
	}

	output.Get().Logger().Info("All projects fetched successfully", "totalProjects", len(allProjects), "organizations", len(orgs))
	return allProjects, nil
}

// ProjectListOptions holds options for listing projects
type ProjectListOptions struct {
	Search string
	All    bool
}

// FetchProjectsAcrossOrgs fetches projects across all organizations
// Deprecated: Use FetchProjectsWithOrgScope instead
func FetchProjectsAcrossOrgs(c *client.TfxClient, searchString string) ([]*tfe.Project, error) {
	options := &ProjectListOptions{
		Search: searchString,
		All:    true,
	}
	return FetchProjectsWithOrgScope(c, "", options)
}

// FetchProject fetches a single project by ID
func FetchProject(c *client.TfxClient, projectID string) (*tfe.Project, error) {
	output.Get().Logger().Debug("Fetching project by ID", "projectID", projectID)

	options := &tfe.ProjectReadOptions{
		Include: []tfe.ProjectIncludeOpt{
			tfe.ProjectEffectiveTagBindings,
		},
	}

	project, err := c.Client.Projects.ReadWithOptions(c.Context, projectID, *options)
	if err != nil {
		output.Get().Logger().Error("Failed to fetch project", "projectID", projectID, "error", err)
		return nil, err
	}

	output.Get().Logger().Debug("Project fetched successfully", "projectID", projectID, "name", project.Name)
	return project, nil
}

// FetchProjectByName fetches a single project by name in the specified organization
func FetchProjectByName(c *client.TfxClient, orgName string, projectName string) (*tfe.Project, error) {
	output.Get().Logger().Debug("Fetching project by name", "organization", orgName, "projectName", projectName)

	projects, err := FetchProjects(c, orgName, projectName)
	if err != nil {
		output.Get().Logger().Error("Failed to fetch projects", "organization", orgName, "projectName", projectName, "error", err)
		return nil, errors.Wrap(err, "failed to fetch projects")
	}

	output.Get().Logger().Trace("Projects search completed", "organization", orgName, "resultsCount", len(projects))

	// Find exact match in case there are multiple results from the search
	for _, p := range projects {
		if p.Name == projectName {
			output.Get().Logger().Debug("Project found by name", "organization", orgName, "projectName", projectName, "projectID", p.ID)
			// Fetch full project details with options
			return FetchProject(c, p.ID)
		}
	}

	output.Get().Logger().Warn("Project not found by name", "organization", orgName, "projectName", projectName)
	return nil, errors.Errorf("project with name '%s' not found in organization '%s'", projectName, orgName)
}
