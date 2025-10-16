// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// VariableListFlags holds all flags for the variable list command
type VariableListFlags struct {
	WorkspaceName string
}

// VariableShowFlags holds all flags for the variable show command
type VariableShowFlags struct {
	WorkspaceName string
	Key           string
}

// VariableCreateFlags holds all flags for the variable create command
type VariableCreateFlags struct {
	WorkspaceName string
	Key           string
	Value         string
	ValueFile     string
	Description   string
	Env           bool
	HCL           bool
	Sensitive     bool
}

// VariableUpdateFlags holds all flags for the variable update command
type VariableUpdateFlags struct {
	WorkspaceName string
	Key           string
	Value         string
	ValueFile     string
	Description   string
	Env           bool
	HCL           bool
	Sensitive     bool
}

// VariableDeleteFlags holds all flags for the variable delete command
type VariableDeleteFlags struct {
	WorkspaceName string
	Key           string
}

// ParseVariableListFlags creates a VariableListFlags from the current command context
func ParseVariableListFlags(cmd *cobra.Command) (*VariableListFlags, error) {
	return &VariableListFlags{
		WorkspaceName: viper.GetString("workspace-name"),
	}, nil
}

// ParseVariableShowFlags creates a VariableShowFlags from the current command context
func ParseVariableShowFlags(cmd *cobra.Command) (*VariableShowFlags, error) {
	return &VariableShowFlags{
		WorkspaceName: viper.GetString("workspace-name"),
		Key:           viper.GetString("key"),
	}, nil
}

// ParseVariableCreateFlags creates a VariableCreateFlags from the current command context
func ParseVariableCreateFlags(cmd *cobra.Command) (*VariableCreateFlags, error) {
	return &VariableCreateFlags{
		WorkspaceName: viper.GetString("workspace-name"),
		Key:           viper.GetString("key"),
		Value:         viper.GetString("value"),
		ValueFile:     viper.GetString("value-file"),
		Description:   viper.GetString("description"),
		Env:           viper.GetBool("env"),
		HCL:           viper.GetBool("hcl"),
		Sensitive:     viper.GetBool("sensitive"),
	}, nil
}

// ParseVariableUpdateFlags creates a VariableUpdateFlags from the current command context
func ParseVariableUpdateFlags(cmd *cobra.Command) (*VariableUpdateFlags, error) {
	return &VariableUpdateFlags{
		WorkspaceName: viper.GetString("workspace-name"),
		Key:           viper.GetString("key"),
		Value:         viper.GetString("value"),
		ValueFile:     viper.GetString("value-file"),
		Description:   viper.GetString("description"),
		Env:           viper.GetBool("env"),
		HCL:           viper.GetBool("hcl"),
		Sensitive:     viper.GetBool("sensitive"),
	}, nil
}

// ParseVariableDeleteFlags creates a VariableDeleteFlags from the current command context
func ParseVariableDeleteFlags(cmd *cobra.Command) (*VariableDeleteFlags, error) {
	return &VariableDeleteFlags{
		WorkspaceName: viper.GetString("workspace-name"),
		Key:           viper.GetString("key"),
	}, nil
}
