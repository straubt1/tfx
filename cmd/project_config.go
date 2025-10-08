package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/straubt1/tfx/client"
)

// ProjectListConfig holds all flags for the project list command
type ProjectListConfig struct {
	Client *client.TfxClient
	Search string
	All    bool
}

// ProjectShowConfig holds all flags for the project show command
type ProjectShowConfig struct {
	Client *client.TfxClient
	ID     string
	Name   string
}

// NewProjectListConfig creates a ProjectListConfig from the current command context
func NewProjectListConfig(cmd *cobra.Command) (*ProjectListConfig, error) {
	c, err := client.NewFromViper()
	if err != nil {
		return nil, err
	}
	return &ProjectListConfig{
		Client: c,
		Search: viper.GetString("search"),
		All:    viper.GetBool("all"),
	}, nil
}

// NewProjectShowConfig creates a ProjectShowConfig from the current command context
func NewProjectShowConfig(cmd *cobra.Command) (*ProjectShowConfig, error) {
	c, err := client.NewFromViper()
	if err != nil {
		return nil, err
	}
	return &ProjectShowConfig{
		Client: c,
		ID:     viper.GetString("id"),
		Name:   viper.GetString("name"),
	}, nil
}
