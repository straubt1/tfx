// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package data

import (
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/output"
)

// VariableSetScope selects which variable sets to list or resolve.
type VariableSetScope struct {
	All              bool
	OrganizationName string
	ProjectName      string
	WorkspaceName    string
	Search           string
}

// VariableSetCreateParams holds options for creating a variable set.
type VariableSetCreateParams struct {
	Name, Description string
	Global, Priority  bool
	ProjectName       string // non-empty = project-owned (Parent.Project)
	WorkspaceName     string // non-empty = ApplyToWorkspaces after create
}

func resolveVariableSetOrg(scope VariableSetScope, defaultOrg string) string {
	if scope.OrganizationName != "" {
		return scope.OrganizationName
	}
	return defaultOrg
}

// ListVariableSets lists all variable sets for an organization.
func ListVariableSets(c *client.TfxClient, orgName string, search string) ([]*tfe.VariableSet, error) {
	output.Get().Logger().Debug("Listing variable sets", "org", orgName, "search", search)

	return client.FetchAll(c.Context, func(pageNumber int) ([]*tfe.VariableSet, *client.Pagination, error) {
		output.Get().Logger().Trace("Listing variable sets page", "org", orgName, "page", pageNumber)

		opts := &tfe.VariableSetListOptions{
			ListOptions: tfe.ListOptions{PageNumber: pageNumber, PageSize: 100},
			Query:       search,
		}
		res, err := c.Client.VariableSets.List(c.Context, orgName, opts)
		if err != nil {
			output.Get().Logger().Error("Failed to list variable sets page", "org", orgName, "page", pageNumber, "error", err)
			return nil, nil, err
		}
		return res.Items, client.NewPaginationFromTFE(res.Pagination), nil
	})
}

// ListVariableSetsWithScope lists variable sets according to the given scope.
func ListVariableSetsWithScope(c *client.TfxClient, defaultOrg string, scope VariableSetScope) ([]*tfe.VariableSet, error) {
	if scope.All {
		output.Get().Logger().Debug("Listing variable sets across all organizations", "search", scope.Search)

		orgs, err := FetchOrganizations(c, "")
		if err != nil {
			return nil, errors.Wrap(err, "failed to list organizations")
		}

		var all []*tfe.VariableSet
		for _, org := range orgs {
			output.Get().Logger().Trace("Listing variable sets for organization", "org", org.Name)

			sets, err := ListVariableSets(c, org.Name, scope.Search)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to list variable sets for organization %s", org.Name)
			}
			all = append(all, sets...)
		}
		return all, nil
	}

	orgName := resolveVariableSetOrg(scope, defaultOrg)

	if scope.WorkspaceName != "" {
		return ListVariableSetsForWorkspace(c, orgName, scope.WorkspaceName, scope.Search)
	}
	if scope.ProjectName != "" {
		return ListVariableSetsForProject(c, orgName, scope.ProjectName, scope.Search)
	}

	return ListVariableSets(c, orgName, scope.Search)
}

// ListVariableSetsForWorkspace lists variable sets applied to a workspace.
func ListVariableSetsForWorkspace(c *client.TfxClient, orgName, workspaceName, search string) ([]*tfe.VariableSet, error) {
	output.Get().Logger().Debug("Listing variable sets for workspace", "org", orgName, "workspace", workspaceName, "search", search)

	workspaceID, err := GetWorkspaceID(c, orgName, workspaceName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to resolve workspace %q", workspaceName)
	}

	return client.FetchAll(c.Context, func(pageNumber int) ([]*tfe.VariableSet, *client.Pagination, error) {
		output.Get().Logger().Trace("Listing variable sets for workspace page", "workspaceID", workspaceID, "page", pageNumber)

		opts := &tfe.VariableSetListOptions{
			ListOptions: tfe.ListOptions{PageNumber: pageNumber, PageSize: 100},
			Query:       search,
		}
		res, err := c.Client.VariableSets.ListForWorkspace(c.Context, workspaceID, opts)
		if err != nil {
			output.Get().Logger().Error("Failed to list variable sets for workspace page", "workspaceID", workspaceID, "page", pageNumber, "error", err)
			return nil, nil, err
		}
		return res.Items, client.NewPaginationFromTFE(res.Pagination), nil
	})
}

// ListVariableSetsForProject lists variable sets applied to a project.
func ListVariableSetsForProject(c *client.TfxClient, orgName, projectName, search string) ([]*tfe.VariableSet, error) {
	output.Get().Logger().Debug("Listing variable sets for project", "org", orgName, "project", projectName, "search", search)

	project, err := FetchProjectByName(c, orgName, projectName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to resolve project %q", projectName)
	}

	return client.FetchAll(c.Context, func(pageNumber int) ([]*tfe.VariableSet, *client.Pagination, error) {
		output.Get().Logger().Trace("Listing variable sets for project page", "projectID", project.ID, "page", pageNumber)

		opts := &tfe.VariableSetListOptions{
			ListOptions: tfe.ListOptions{PageNumber: pageNumber, PageSize: 100},
			Query:       search,
		}
		res, err := c.Client.VariableSets.ListForProject(c.Context, project.ID, opts)
		if err != nil {
			output.Get().Logger().Error("Failed to list variable sets for project page", "projectID", project.ID, "page", pageNumber, "error", err)
			return nil, nil, err
		}
		return res.Items, client.NewPaginationFromTFE(res.Pagination), nil
	})
}

// GetVariableSetByName finds a variable set by exact name within the given scope.
func GetVariableSetByName(c *client.TfxClient, orgName string, scope VariableSetScope, name string) (*tfe.VariableSet, error) {
	output.Get().Logger().Debug("Getting variable set by name", "org", orgName, "name", name, "scope", scope)

	sets, err := ListVariableSetsWithScope(c, orgName, scope)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list variable sets")
	}

	for _, vs := range sets {
		if vs.Name == name {
			output.Get().Logger().Debug("Variable set found by name", "name", name, "id", vs.ID)
			return vs, nil
		}
	}

	return nil, errors.Errorf("variable set with name %q not found", name)
}

// ResolveVariableSet reads a variable set by ID or resolves it by name within the given scope.
func ResolveVariableSet(c *client.TfxClient, defaultOrg string, scope VariableSetScope, name, id string) (*tfe.VariableSet, error) {
	if id != "" {
		return ReadVariableSet(c, id)
	}
	if name == "" {
		return nil, errors.New("variable set name or id is required")
	}
	return GetVariableSetByName(c, defaultOrg, scope, name)
}

// ReadVariableSet reads a single variable set by ID, including its workspaces, projects, and variables.
func ReadVariableSet(c *client.TfxClient, variableSetID string) (*tfe.VariableSet, error) {
	output.Get().Logger().Debug("Reading variable set", "id", variableSetID)

	include := []tfe.VariableSetIncludeOpt{
		tfe.VariableSetWorkspaces,
		tfe.VariableSetProjects,
		tfe.VariableSetVars,
	}
	return c.Client.VariableSets.Read(c.Context, variableSetID, &tfe.VariableSetReadOptions{
		Include: &include,
	})
}

// CreateVariableSet creates a new variable set in the given organization.
func CreateVariableSet(c *client.TfxClient, orgName string, params VariableSetCreateParams) (*tfe.VariableSet, error) {
	output.Get().Logger().Debug("Creating variable set", "org", orgName, "name", params.Name)

	opts := &tfe.VariableSetCreateOptions{
		Name:        &params.Name,
		Description: &params.Description,
		Global:      &params.Global,
		Priority:    &params.Priority,
	}

	if params.ProjectName != "" {
		project, err := FetchProjectByName(c, orgName, params.ProjectName)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to resolve project %q", params.ProjectName)
		}
		opts.Parent = &tfe.Parent{Project: project}
	}

	vs, err := c.Client.VariableSets.Create(c.Context, orgName, opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create variable set")
	}

	if params.WorkspaceName != "" {
		workspaceID, err := GetWorkspaceID(c, orgName, params.WorkspaceName)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to resolve workspace %q", params.WorkspaceName)
		}

		err = c.Client.VariableSets.ApplyToWorkspaces(c.Context, vs.ID, &tfe.VariableSetApplyToWorkspacesOptions{
			Workspaces: []*tfe.Workspace{{ID: workspaceID}},
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to apply variable set to workspace")
		}
	}

	output.Get().Logger().Debug("Variable set created successfully", "org", orgName, "name", params.Name, "id", vs.ID)
	return vs, nil
}

// DeleteVariableSet deletes a variable set by ID.
func DeleteVariableSet(c *client.TfxClient, variableSetID string) error {
	output.Get().Logger().Debug("Deleting variable set", "id", variableSetID)
	return c.Client.VariableSets.Delete(c.Context, variableSetID)
}
