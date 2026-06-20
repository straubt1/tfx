// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package data

import (
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/output"
)

// FetchVariableSetVariables fetches all variables in a variable set using pagination.
func FetchVariableSetVariables(c *client.TfxClient, variableSetID string) ([]*tfe.VariableSetVariable, error) {
	output.Get().Logger().Debug("Fetching variable set variables", "variableSetID", variableSetID)

	return client.FetchAll(c.Context, func(pageNumber int) ([]*tfe.VariableSetVariable, *client.Pagination, error) {
		output.Get().Logger().Trace("Fetching variable set variables page", "variableSetID", variableSetID, "page", pageNumber)

		opts := &tfe.VariableSetVariableListOptions{
			ListOptions: tfe.ListOptions{PageNumber: pageNumber, PageSize: 100},
		}

		result, err := c.Client.VariableSetVariables.List(c.Context, variableSetID, opts)
		if err != nil {
			output.Get().Logger().Error("Failed to fetch variable set variables page", "variableSetID", variableSetID, "page", pageNumber, "error", err)
			return nil, nil, err
		}

		output.Get().Logger().Trace("Variable set variables page fetched", "variableSetID", variableSetID, "page", pageNumber, "count", len(result.Items))
		return result.Items, client.NewPaginationFromTFE(result.Pagination), nil
	})
}

// FetchVariableSetVariable fetches a single variable by key in the specified variable set.
func FetchVariableSetVariable(c *client.TfxClient, variableSetID string, key string) (*tfe.VariableSetVariable, error) {
	output.Get().Logger().Debug("Fetching variable set variable by key", "variableSetID", variableSetID, "key", key)

	variableID, err := GetVariableSetVariableID(c, variableSetID, key)
	if err != nil {
		output.Get().Logger().Error("Failed to get variable set variable ID", "variableSetID", variableSetID, "key", key, "error", err)
		return nil, err
	}

	variable, err := c.Client.VariableSetVariables.Read(c.Context, variableSetID, variableID)
	if err != nil {
		output.Get().Logger().Error("Failed to fetch variable set variable", "variableSetID", variableSetID, "variableID", variableID, "error", err)
		return nil, err
	}

	output.Get().Logger().Debug("Variable set variable fetched successfully", "variableSetID", variableSetID, "key", key, "variableID", variableID)
	return variable, nil
}

// CreateVariableSetVariable creates a new variable in a variable set.
func CreateVariableSetVariable(c *client.TfxClient, variableSetID string, opts tfe.VariableSetVariableCreateOptions) (*tfe.VariableSetVariable, error) {
	output.Get().Logger().Debug("Creating variable set variable", "variableSetID", variableSetID, "key", *opts.Key)

	variable, err := c.Client.VariableSetVariables.Create(c.Context, variableSetID, &opts)
	if err != nil {
		output.Get().Logger().Error("Failed to create variable set variable", "variableSetID", variableSetID, "key", *opts.Key, "error", err)
		return nil, err
	}

	output.Get().Logger().Debug("Variable set variable created successfully", "variableSetID", variableSetID, "key", *opts.Key, "variableID", variable.ID)
	return variable, nil
}

// UpdateVariableSetVariable updates an existing variable in a variable set.
func UpdateVariableSetVariable(c *client.TfxClient, variableSetID, variableID string, opts tfe.VariableSetVariableUpdateOptions) (*tfe.VariableSetVariable, error) {
	output.Get().Logger().Debug("Updating variable set variable", "variableSetID", variableSetID, "variableID", variableID)

	variable, err := c.Client.VariableSetVariables.Update(c.Context, variableSetID, variableID, &opts)
	if err != nil {
		output.Get().Logger().Error("Failed to update variable set variable", "variableSetID", variableSetID, "variableID", variableID, "error", err)
		return nil, err
	}

	output.Get().Logger().Debug("Variable set variable updated successfully", "variableSetID", variableSetID, "variableID", variableID)
	return variable, nil
}

// DeleteVariableSetVariable deletes a variable from a variable set.
func DeleteVariableSetVariable(c *client.TfxClient, variableSetID, variableID string) error {
	output.Get().Logger().Debug("Deleting variable set variable", "variableSetID", variableSetID, "variableID", variableID)

	err := c.Client.VariableSetVariables.Delete(c.Context, variableSetID, variableID)
	if err != nil {
		output.Get().Logger().Error("Failed to delete variable set variable", "variableSetID", variableSetID, "variableID", variableID, "error", err)
		return err
	}

	output.Get().Logger().Debug("Variable set variable deleted successfully", "variableSetID", variableSetID, "variableID", variableID)
	return nil
}

// GetVariableSetVariableID retrieves the variable ID from a variable set by variable key.
func GetVariableSetVariableID(c *client.TfxClient, variableSetID string, key string) (string, error) {
	output.Get().Logger().Debug("Getting variable set variable ID by key", "variableSetID", variableSetID, "key", key)

	variables, err := FetchVariableSetVariables(c, variableSetID)
	if err != nil {
		output.Get().Logger().Error("Failed to fetch variable set variables", "variableSetID", variableSetID, "error", err)
		return "", errors.Wrap(err, "failed to fetch variable set variables")
	}

	for _, v := range variables {
		if v.Key == key {
			output.Get().Logger().Debug("Variable set variable ID found", "variableSetID", variableSetID, "key", key, "variableID", v.ID)
			return v.ID, nil
		}
	}

	output.Get().Logger().Warn("Variable set variable key not found", "variableSetID", variableSetID, "key", key)
	return "", errors.Errorf("variable with key %q not found in variable set %s", key, variableSetID)
}
