package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// OrganizationListFlags holds all flags for the organization list command
type OrganizationListFlags struct {
	Search string
}

// OrganizationShowFlags holds all flags for the organization show command
type OrganizationShowFlags struct {
	Name string
}

// ParseOrganizationListFlags creates a OrganizationListFlags from the current command context
func ParseOrganizationListFlags(cmd *cobra.Command) (*OrganizationListFlags, error) {
	return &OrganizationListFlags{
		Search: viper.GetString("search"),
	}, nil
}

// ParseOrganizationShowFlags creates a OrganizationShowFlags from the current command context
func ParseOrganizationShowFlags(cmd *cobra.Command) (*OrganizationShowFlags, error) {
	return &OrganizationShowFlags{
		Name: viper.GetString("name"),
	}, nil
}
