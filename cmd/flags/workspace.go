package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// WorkspaceListFlags holds all flags for the workspace list command
type WorkspaceListFlags struct {
	Search     string
	Repository string
	RunStatus  string
	ProjectID  string
	All        bool
}

// WorkspaceShowFlags holds all flags for the workspace show command
type WorkspaceShowFlags struct {
	Name string
}

// ParseWorkspaceListFlags creates a WorkspaceListFlags from the current command context
func ParseWorkspaceListFlags(cmd *cobra.Command) (*WorkspaceListFlags, error) {
	return &WorkspaceListFlags{
		Search:     viper.GetString("search"),
		Repository: viper.GetString("repository"),
		RunStatus:  viper.GetString("run-status"),
		ProjectID:  viper.GetString("project-id"),
		All:        viper.GetBool("all"),
	}, nil
}

// ParseWorkspaceShowFlags creates a WorkspaceShowFlags from the current command context
func ParseWorkspaceShowFlags(cmd *cobra.Command) (*WorkspaceShowFlags, error) {
	return &WorkspaceShowFlags{
		Name: viper.GetString("name"),
	}, nil
}
