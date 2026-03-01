// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package data

import (
	tfe "github.com/hashicorp/go-tfe"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/output"
)

// ListVariableSets lists all variable sets for an organization.
func ListVariableSets(c *client.TfxClient, orgName string, search string) ([]*tfe.VariableSet, error) {
	output.Get().Logger().Debug("Listing variable sets", "org", orgName, "search", search)

	return client.FetchAll(c.Context, func(pageNumber int) ([]*tfe.VariableSet, *client.Pagination, error) {
		opts := &tfe.VariableSetListOptions{
			ListOptions: tfe.ListOptions{PageNumber: pageNumber, PageSize: 100},
			Query:       search,
		}
		res, err := c.Client.VariableSets.List(c.Context, orgName, opts)
		if err != nil {
			return nil, nil, err
		}
		return res.Items, client.NewPaginationFromTFE(res.Pagination), nil
	})
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
func CreateVariableSet(c *client.TfxClient, orgName, name, description string, global, priority bool) (*tfe.VariableSet, error) {
	output.Get().Logger().Debug("Creating variable set", "org", orgName, "name", name)

	return c.Client.VariableSets.Create(c.Context, orgName, &tfe.VariableSetCreateOptions{
		Name:        &name,
		Description: &description,
		Global:      &global,
		Priority:    &priority,
	})
}

// DeleteVariableSet deletes a variable set by ID.
func DeleteVariableSet(c *client.TfxClient, variableSetID string) error {
	output.Get().Logger().Debug("Deleting variable set", "id", variableSetID)
	return c.Client.VariableSets.Delete(c.Context, variableSetID)
}
