// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// VariableSetListFlags holds flags for the variable-set list command.
type VariableSetListFlags struct {
	Search string
}

// VariableSetShowFlags holds flags for the variable-set show command.
type VariableSetShowFlags struct {
	ID string
}

// VariableSetCreateFlags holds flags for the variable-set create command.
type VariableSetCreateFlags struct {
	Name        string
	Description string
	Global      bool
	Priority    bool
}

// VariableSetDeleteFlags holds flags for the variable-set delete command.
type VariableSetDeleteFlags struct {
	ID string
}

func ParseVariableSetListFlags(cmd *cobra.Command) (*VariableSetListFlags, error) {
	return &VariableSetListFlags{
		Search: viper.GetString("search"),
	}, nil
}

func ParseVariableSetShowFlags(cmd *cobra.Command) (*VariableSetShowFlags, error) {
	return &VariableSetShowFlags{
		ID: viper.GetString("id"),
	}, nil
}

func ParseVariableSetCreateFlags(cmd *cobra.Command) (*VariableSetCreateFlags, error) {
	return &VariableSetCreateFlags{
		Name:        viper.GetString("name"),
		Description: viper.GetString("description"),
		Global:      viper.GetBool("global"),
		Priority:    viper.GetBool("priority"),
	}, nil
}

func ParseVariableSetDeleteFlags(cmd *cobra.Command) (*VariableSetDeleteFlags, error) {
	return &VariableSetDeleteFlags{
		ID: viper.GetString("id"),
	}, nil
}
