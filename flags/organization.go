package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// OrganizationShowFlags holds all flags for the organization show command
type OrganizationShowFlags struct {
	Name string
}

// ParseOrganizationShowFlags creates a OrganizationShowFlags from the current command context
func ParseOrganizationShowFlags(cmd *cobra.Command) (*OrganizationShowFlags, error) {
	return &OrganizationShowFlags{
		Name: viper.GetString("name"),
	}, nil
}
