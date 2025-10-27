//go:build integration
// +build integration

package integration

import (
	"os"
	"testing"
)

func TestProjectList(t *testing.T) {
	hostname, token, org := setupTest(t)

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "list all projects",
			args: []string{"project", "list"},
		},
		{
			name: "list projects with search",
			args: []string{"project", "list", "--search", "test"},
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

func TestProjectShow(t *testing.T) {
	hostname, token, org := setupTest(t)

	// Test with project name if provided
	testProjectName := os.Getenv("TEST_PROJECT_NAME")
	if testProjectName != "" {
		t.Run("show project by name", func(t *testing.T) {
			err := executeCommand(t,
				[]string{"project", "show", "--name", testProjectName},
				hostname, token, org,
			)

			if err != nil {
				t.Errorf("project show by name command failed: %v", err)
			}
		})
	}

	// Test with project ID if provided
	testProjectID := os.Getenv("TEST_PROJECT_ID")
	if testProjectID != "" {
		t.Run("show project by id", func(t *testing.T) {
			err := executeCommand(t,
				[]string{"project", "show", "--id", testProjectID},
				hostname, token, org,
			)

			if err != nil {
				t.Errorf("project show by id command failed: %v", err)
			}
		})
	}

	// Skip if neither is provided
	if testProjectName == "" && testProjectID == "" {
		t.Skip("TEST_PROJECT_NAME or TEST_PROJECT_ID not set - skipping project show test")
	}
}
