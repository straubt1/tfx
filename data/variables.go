// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	tfe "github.com/hashicorp/go-tfe"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/logger"
)

// FetchVariables fetches all variables for a given workspace using pagination
func FetchVariables(c *client.TfxClient, workspaceID string) ([]*tfe.Variable, error) {
	logger.Debug("Fetching variables", "workspaceID", workspaceID)

	return client.FetchAll(c.Context, func(pageNumber int) ([]*tfe.Variable, *client.Pagination, error) {
		logger.Trace("Fetching variables page", "workspaceID", workspaceID, "page", pageNumber)

		opts := &tfe.VariableListOptions{
			ListOptions: tfe.ListOptions{PageNumber: pageNumber, PageSize: 100},
		}

		result, err := c.Client.Variables.List(c.Context, workspaceID, opts)
		if err != nil {
			logger.Error("Failed to fetch variables page", "workspaceID", workspaceID, "page", pageNumber, "error", err)
			return nil, nil, err
		}

		logger.Trace("Variables page fetched", "workspaceID", workspaceID, "page", pageNumber, "count", len(result.Items))
		return result.Items, client.NewPaginationFromTFE(result.Pagination), nil
	})
}

// FetchVariable fetches a single variable by key in the specified workspace
func FetchVariable(c *client.TfxClient, workspaceID string, key string) (*tfe.Variable, error) {
	logger.Debug("Fetching variable by key", "workspaceID", workspaceID, "key", key)

	variableID, err := GetVariableID(c, workspaceID, key)
	if err != nil {
		logger.Error("Failed to get variable ID", "workspaceID", workspaceID, "key", key, "error", err)
		return nil, err
	}

	variable, err := c.Client.Variables.Read(c.Context, workspaceID, variableID)
	if err != nil {
		logger.Error("Failed to fetch variable", "workspaceID", workspaceID, "variableID", variableID, "error", err)
		return nil, err
	}

	logger.Debug("Variable fetched successfully", "workspaceID", workspaceID, "key", key, "variableID", variableID)
	return variable, nil
}

// CreateVariable creates a new variable in a workspace
func CreateVariable(c *client.TfxClient, workspaceID string, opts tfe.VariableCreateOptions) (*tfe.Variable, error) {
	logger.Debug("Creating variable", "workspaceID", workspaceID, "key", *opts.Key)

	variable, err := c.Client.Variables.Create(c.Context, workspaceID, opts)
	if err != nil {
		logger.Error("Failed to create variable", "workspaceID", workspaceID, "key", *opts.Key, "error", err)
		return nil, err
	}

	logger.Debug("Variable created successfully", "workspaceID", workspaceID, "key", *opts.Key, "variableID", variable.ID)
	return variable, nil
}

// UpdateVariable updates an existing variable in a workspace
func UpdateVariable(c *client.TfxClient, workspaceID string, variableID string, opts tfe.VariableUpdateOptions) (*tfe.Variable, error) {
	logger.Debug("Updating variable", "workspaceID", workspaceID, "variableID", variableID)

	variable, err := c.Client.Variables.Update(c.Context, workspaceID, variableID, opts)
	if err != nil {
		logger.Error("Failed to update variable", "workspaceID", workspaceID, "variableID", variableID, "error", err)
		return nil, err
	}

	logger.Debug("Variable updated successfully", "workspaceID", workspaceID, "variableID", variableID)
	return variable, nil
}

// DeleteVariable deletes a variable from a workspace
func DeleteVariable(c *client.TfxClient, workspaceID string, variableID string) error {
	logger.Debug("Deleting variable", "workspaceID", workspaceID, "variableID", variableID)

	err := c.Client.Variables.Delete(c.Context, workspaceID, variableID)
	if err != nil {
		logger.Error("Failed to delete variable", "workspaceID", workspaceID, "variableID", variableID, "error", err)
		return err
	}

	logger.Debug("Variable deleted successfully", "workspaceID", workspaceID, "variableID", variableID)
	return nil
}

// GetVariableID retrieves the variable ID from a workspace by variable key
func GetVariableID(c *client.TfxClient, workspaceID string, key string) (string, error) {
	logger.Debug("Getting variable ID by key", "workspaceID", workspaceID, "key", key)

	// variables, err := FetchVariables(c, workspaceID)
	variable, err := c.Client.Variables.Read(c.Context, workspaceID, key)
	if err != nil {
		logger.Error("Variable key not found", "workspaceID", workspaceID, "key", key)
		return "", err
	}

	return variable.ID, nil
}

// GetWorkspaceID retrieves the workspace ID from organization and workspace name
func GetWorkspaceID(c *client.TfxClient, orgName string, workspaceName string) (string, error) {
	logger.Debug("Getting workspace ID", "organization", orgName, "workspaceName", workspaceName)

	workspace, err := c.Client.Workspaces.Read(c.Context, orgName, workspaceName)
	if err != nil {
		logger.Error("Failed to fetch workspace", "organization", orgName, "workspaceName", workspaceName, "error", err)
		return "", err
	}

	logger.Debug("Workspace ID found", "organization", orgName, "workspaceName", workspaceName, "workspaceID", workspace.ID)
	return workspace.ID, nil
}
