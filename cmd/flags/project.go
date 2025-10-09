package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ProjectListFlags holds all flags for the project list command
type ProjectListFlags struct {
	Search string
	All    bool
}

// ProjectShowFlags holds all flags for the project show command
type ProjectShowFlags struct {
	ID   string
	Name string
}

// ParseProjectListFlags creates a ProjectListFlags from the current command context
func ParseProjectListFlags(cmd *cobra.Command) (*ProjectListFlags, error) {
	return &ProjectListFlags{
		Search: viper.GetString("search"),
		All:    viper.GetBool("all"),
	}, nil
}

// ParseProjectShowFlags creates a ProjectShowFlags from the current command context
func ParseProjectShowFlags(cmd *cobra.Command) (*ProjectShowFlags, error) {
	return &ProjectShowFlags{
		ID:   viper.GetString("id"),
		Name: viper.GetString("name"),
	}, nil
}
