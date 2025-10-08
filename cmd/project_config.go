package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ProjectListConfig holds all flags for the project list command
type ProjectListConfig struct {
	Search string
	All    bool
}

// ProjectShowConfig holds all flags for the project show command
type ProjectShowConfig struct {
	ID   string
	Name string
}

// NewProjectListConfig creates a ProjectListConfig from the current command context
func NewProjectListConfig(cmd *cobra.Command) (*ProjectListConfig, error) {
	return &ProjectListConfig{
		Search: viper.GetString("search"),
		All:    viper.GetBool("all"),
	}, nil
}

// NewProjectShowConfig creates a ProjectShowConfig from the current command context
func NewProjectShowConfig(cmd *cobra.Command) (*ProjectShowConfig, error) {
	return &ProjectShowConfig{
		ID:   viper.GetString("id"),
		Name: viper.GetString("name"),
	}, nil
}
