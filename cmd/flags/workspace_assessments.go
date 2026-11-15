// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// WorkspaceAssessmentsFlags holds flags for `tfx workspace assessments enable|disable`
type WorkspaceAssessmentsFlags struct {
	Names []string
}

// WorkspaceAssessmentsAllFlags holds flags for `tfx workspace assessments enable|disable all`
type WorkspaceAssessmentsAllFlags struct {
	Search       string
	WildcardName string
	ProjectID    string
	Tags         string
	ExcludeTags  string
}

// ParseWorkspaceAssessmentsFlags creates a WorkspaceAssessmentsFlags from the current command context
func ParseWorkspaceAssessmentsFlags(cmd *cobra.Command) (*WorkspaceAssessmentsFlags, error) {
	return &WorkspaceAssessmentsFlags{
		Names: viper.GetStringSlice("name"),
	}, nil
}

// ParseWorkspaceAssessmentsAllFlags creates a WorkspaceAssessmentsAllFlags from the current command context
func ParseWorkspaceAssessmentsAllFlags(cmd *cobra.Command) (*WorkspaceAssessmentsAllFlags, error) {
	return &WorkspaceAssessmentsAllFlags{
		Search:       viper.GetString("search"),
		WildcardName: viper.GetString("wildcard-name"),
		ProjectID:    viper.GetString("project-id"),
		Tags:         viper.GetString("tags"),
		ExcludeTags:  viper.GetString("exclude-tags"),
	}, nil
}

// ListFlags converts the bulk filter flags into workspace list options
func (f *WorkspaceAssessmentsAllFlags) ListFlags() *WorkspaceListFlags {
	return &WorkspaceListFlags{
		Search:       f.Search,
		WildcardName: f.WildcardName,
		ProjectID:    f.ProjectID,
		Tags:         f.Tags,
		ExcludeTags:  f.ExcludeTags,
	}
}
