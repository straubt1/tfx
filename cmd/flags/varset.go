// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// VariableSetScopeFlags holds shared scope flags for variable set commands.
type VariableSetScopeFlags struct {
	All              bool
	OrganizationName string
	ProjectName      string
	WorkspaceName    string
	Search           string
}

// ResolveVarsetOrg returns the organization name from flags or the client default.
func ResolveVarsetOrg(flagOrg, clientOrg string) string {
	if flagOrg != "" {
		return flagOrg
	}
	return clientOrg
}

func parseVariableSetScopeFlags() VariableSetScopeFlags {
	return VariableSetScopeFlags{
		All:              viper.GetBool("all"),
		OrganizationName: viper.GetString("organization-name"),
		ProjectName:      viper.GetString("project-name"),
		WorkspaceName:    viper.GetString("workspace-name"),
		Search:           viper.GetString("search"),
	}
}

// VariableSetListFlags holds flags for the variable-set list command.
type VariableSetListFlags struct {
	VariableSetScopeFlags
}

// VariableSetShowFlags holds flags for the variable-set show command.
type VariableSetShowFlags struct {
	ID   string
	Name string
	VariableSetScopeFlags
}

// VariableSetCreateFlags holds flags for the variable-set create command.
type VariableSetCreateFlags struct {
	Name             string
	Description      string
	Global           bool
	Priority         bool
	OrganizationName string
	ProjectName      string
	WorkspaceName    string
}

// VariableSetDeleteFlags holds flags for the variable-set delete command.
type VariableSetDeleteFlags struct {
	ID   string
	Name string
	VariableSetScopeFlags
}

func ParseVariableSetListFlags(cmd *cobra.Command) (*VariableSetListFlags, error) {
	return &VariableSetListFlags{
		VariableSetScopeFlags: parseVariableSetScopeFlags(),
	}, nil
}

func ParseVariableSetShowFlags(cmd *cobra.Command) (*VariableSetShowFlags, error) {
	return &VariableSetShowFlags{
		ID:                    viper.GetString("id"),
		Name:                  viper.GetString("name"),
		VariableSetScopeFlags: parseVariableSetScopeFlagsWithoutSearch(),
	}, nil
}

func ParseVariableSetCreateFlags(cmd *cobra.Command) (*VariableSetCreateFlags, error) {
	return &VariableSetCreateFlags{
		Name:             viper.GetString("name"),
		Description:      viper.GetString("description"),
		Global:           viper.GetBool("global"),
		Priority:         viper.GetBool("priority"),
		OrganizationName: viper.GetString("organization-name"),
		ProjectName:      viper.GetString("project-name"),
		WorkspaceName:    viper.GetString("workspace-name"),
	}, nil
}

func ParseVariableSetDeleteFlags(cmd *cobra.Command) (*VariableSetDeleteFlags, error) {
	return &VariableSetDeleteFlags{
		ID:                    viper.GetString("id"),
		Name:                  viper.GetString("name"),
		VariableSetScopeFlags: parseVariableSetScopeFlagsWithoutSearch(),
	}, nil
}

func parseVariableSetScopeFlagsWithoutSearch() VariableSetScopeFlags {
	f := parseVariableSetScopeFlags()
	f.Search = ""
	return f
}
