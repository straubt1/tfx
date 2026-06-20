// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

//go:build integration
// +build integration

package integration

import (
	"os"
	"testing"
)

const (
	integrationVarsetName = "tfx-integration-test-varset"
	integrationVarKey    = "tfx-integration-test-key"
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
		{
			name: "list variable sets across all organizations",
			args: []string{"varset", "list", "--all"},
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

	testProjectName := os.Getenv("TEST_PROJECT_NAME")
	if testProjectName != "" {
		t.Run("list variable sets for project", func(t *testing.T) {
			err := executeCommand(t,
				[]string{"varset", "list", "--project-name", testProjectName},
				hostname, token, org,
			)
			if err != nil {
				t.Errorf("Command failed: %v", err)
			}
		})
	} else {
		t.Run("list variable sets for project", func(t *testing.T) {
			t.Skip("TEST_PROJECT_NAME not set — skipping project-scoped list test")
		})
	}

	testWorkspaceName := os.Getenv("TEST_WORKSPACE_NAME")
	if testWorkspaceName != "" {
		t.Run("list variable sets for workspace", func(t *testing.T) {
			err := executeCommand(t,
				[]string{"varset", "list", "--workspace-name", testWorkspaceName},
				hostname, token, org,
			)
			if err != nil {
				t.Errorf("Command failed: %v", err)
			}
		})
	} else {
		t.Run("list variable sets for workspace", func(t *testing.T) {
			t.Skip("TEST_WORKSPACE_NAME not set — skipping workspace-scoped list test")
		})
	}
}

func TestVariableSetCRUD(t *testing.T) {
	hostname, token, org := setupTest(t)

	t.Run("create variable set", func(t *testing.T) {
		err := executeCommand(t,
			[]string{"variable-set", "create", "--name", integrationVarsetName, "--description", "Created by TFx integration test"},
			hostname, token, org,
		)
		if err != nil {
			t.Fatalf("create variable set failed: %v", err)
		}
	})

	t.Run("list variable sets after create", func(t *testing.T) {
		err := executeCommand(t,
			[]string{"variable-set", "list", "--search", integrationVarsetName},
			hostname, token, org,
		)
		if err != nil {
			t.Errorf("list after create failed: %v", err)
		}
	})

	t.Run("show variable set by name", func(t *testing.T) {
		err := executeCommand(t,
			[]string{"varset", "show", "--name", integrationVarsetName},
			hostname, token, org,
		)
		if err != nil {
			t.Errorf("show variable set by name failed: %v", err)
		}
	})

	t.Run("delete variable set by name", func(t *testing.T) {
		err := executeCommand(t,
			[]string{"varset", "delete", "--name", integrationVarsetName},
			hostname, token, org,
		)
		if err != nil {
			t.Errorf("delete variable set by name failed: %v", err)
		}
	})
}

func TestVariableSetVariableCRUD(t *testing.T) {
	hostname, token, org := setupTest(t)

	t.Run("create variable set", func(t *testing.T) {
		err := executeCommand(t,
			[]string{"varset", "create", "--name", integrationVarsetName, "--description", "Created by TFx integration test"},
			hostname, token, org,
		)
		if err != nil {
			t.Fatalf("create variable set failed: %v", err)
		}
	})

	t.Run("create variable in variable set", func(t *testing.T) {
		err := executeCommand(t,
			[]string{
				"varset", "variable", "create",
				"--varset-name", integrationVarsetName,
				"--key", integrationVarKey,
				"--value", "integration-test",
				"--description", "Created by TFx integration test",
			},
			hostname, token, org,
		)
		if err != nil {
			t.Fatalf("create variable failed: %v", err)
		}
	})

	t.Run("list variables in variable set", func(t *testing.T) {
		err := executeCommand(t,
			[]string{"varset", "variable", "list", "--varset-name", integrationVarsetName},
			hostname, token, org,
		)
		if err != nil {
			t.Errorf("list variables failed: %v", err)
		}
	})

	t.Run("delete variable from variable set", func(t *testing.T) {
		err := executeCommand(t,
			[]string{
				"varset", "variable", "delete",
				"--varset-name", integrationVarsetName,
				"--key", integrationVarKey,
			},
			hostname, token, org,
		)
		if err != nil {
			t.Errorf("delete variable failed: %v", err)
		}
	})

	t.Run("delete variable set", func(t *testing.T) {
		err := executeCommand(t,
			[]string{"varset", "delete", "--name", integrationVarsetName},
			hostname, token, org,
		)
		if err != nil {
			t.Errorf("delete variable set failed: %v", err)
		}
	})
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
