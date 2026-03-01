// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

//go:build integration
// +build integration

// Copyright © 2025 Tom Straub <github.com/straubt1>

package integration

import (
	"os"
	"testing"
)

func TestWorkspaceList(t *testing.T) {
	hostname, token, org := setupTest(t)

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "list all workspaces",
			args: []string{"workspace", "list"},
		},
		{
			name: "list workspaces with search",
			args: []string{"workspace", "list", "--search", "test"},
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

func TestWorkspaceShow(t *testing.T) {
	hostname, token, org := setupTest(t)

	// Requires a workspace name to test against
	testWorkspace := os.Getenv("TEST_WORKSPACE_NAME")
	if testWorkspace == "" {
		t.Skip("TEST_WORKSPACE_NAME not set - skipping workspace show test")
	}

	err := executeCommand(t,
		[]string{"workspace", "show", "--name", testWorkspace},
		hostname, token, org,
	)

	if err != nil {
		t.Errorf("workspace show command failed: %v", err)
	}
}
