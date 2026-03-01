// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

//go:build integration
// +build integration

package integration

import (
	"os"
	"testing"
)

func TestVariableSetList(t *testing.T) {
	hostname, token, org := setupTest(t)

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "list all variable sets",
			args: []string{"variable-set", "list"},
		},
		{
			name: "list variable sets with search",
			args: []string{"variable-set", "list", "--search", "test"},
		},
		{
			name: "list variable sets via alias",
			args: []string{"varset", "list"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := executeCommand(t, tt.args, hostname, token, org)
			if err != nil {
				t.Errorf("Command failed: %v", err)
			}
		})
	}
}

func TestVariableSetCRUD(t *testing.T) {
	hostname, token, org := setupTest(t)

	// Create
	var createdID string
	t.Run("create variable set", func(t *testing.T) {
		err := executeCommand(t,
			[]string{"variable-set", "create", "--name", "tfx-integration-test-varset", "--description", "Created by TFx integration test"},
			hostname, token, org,
		)
		if err != nil {
			t.Fatalf("create variable set failed: %v", err)
		}
	})

	// List to get ID
	t.Run("list variable sets after create", func(t *testing.T) {
		err := executeCommand(t,
			[]string{"variable-set", "list", "--search", "tfx-integration-test-varset"},
			hostname, token, org,
		)
		if err != nil {
			t.Errorf("list after create failed: %v", err)
		}
	})

	// Show (only if TEST_VARSET_ID is provided)
	testVarSetID := os.Getenv("TEST_VARSET_ID")
	if testVarSetID != "" {
		createdID = testVarSetID
		t.Run("show variable set", func(t *testing.T) {
			err := executeCommand(t,
				[]string{"variable-set", "show", "--id", createdID},
				hostname, token, org,
			)
			if err != nil {
				t.Errorf("show variable set failed: %v", err)
			}
		})

		// Delete
		t.Run("delete variable set", func(t *testing.T) {
			err := executeCommand(t,
				[]string{"variable-set", "delete", "--id", createdID},
				hostname, token, org,
			)
			if err != nil {
				t.Errorf("delete variable set failed: %v", err)
			}
		})
	} else {
		t.Log("TEST_VARSET_ID not set — skipping show and delete tests. Set TEST_VARSET_ID to the ID of a variable set to test these.")
	}
}

func TestVariableSetShow(t *testing.T) {
	hostname, token, org := setupTest(t)

	testVarSetID := os.Getenv("TEST_VARSET_ID")
	if testVarSetID == "" {
		t.Skip("TEST_VARSET_ID not set — skipping variable-set show test")
	}

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "show variable set",
			args: []string{"variable-set", "show", "--id", testVarSetID},
		},
		{
			name: "show variable set json",
			args: []string{"variable-set", "show", "--id", testVarSetID, "--json"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := executeCommand(t, tt.args, hostname, token, org)
			if err != nil {
				t.Errorf("Command failed: %v", err)
			}
		})
	}
}
