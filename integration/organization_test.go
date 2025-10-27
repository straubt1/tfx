//go:build integration
// +build integration

package integration

import (
	"testing"
)

func TestOrganizationList(t *testing.T) {
	hostname, token, org := setupTest(t)

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "list all organizations",
			args: []string{"organization", "list"},
		},
		{
			name: "list organizations with search filter",
			args: []string{"organization", "list", "--search", org},
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

func TestOrganizationShow(t *testing.T) {
	hostname, token, org := setupTest(t)

	err := executeCommand(t,
		[]string{"organization", "show", "--name", org},
		hostname, token, org,
	)

	if err != nil {
		t.Errorf("organization show command failed: %v", err)
	}
}
