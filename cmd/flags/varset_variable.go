// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// VariableSetVariableListFlags holds flags for the variable-set variable list command.
type VariableSetVariableListFlags struct {
	VarsetName string
	VarsetID   string
	VariableSetScopeFlags
}

// VariableSetVariableShowFlags holds flags for the variable-set variable show command.
type VariableSetVariableShowFlags struct {
	VarsetName string
	VarsetID   string
	Key        string
	VariableSetScopeFlags
}

// VariableSetVariableCreateFlags holds flags for the variable-set variable create command.
type VariableSetVariableCreateFlags struct {
	VarsetName  string
	VarsetID    string
	Key         string
	Value       string
	ValueFile   string
	Description string
	Env         bool
	HCL         bool
	Sensitive   bool
	VariableSetScopeFlags
}

// VariableSetVariableUpdateFlags holds flags for the variable-set variable update command.
type VariableSetVariableUpdateFlags struct {
	VarsetName  string
	VarsetID    string
	Key         string
	Value       string
	ValueFile   string
	Description string
	Env         bool
	HCL         bool
	Sensitive   bool
	VariableSetScopeFlags
}

// VariableSetVariableDeleteFlags holds flags for the variable-set variable delete command.
type VariableSetVariableDeleteFlags struct {
	VarsetName string
	VarsetID   string
	Key        string
	VariableSetScopeFlags
}

func ParseVariableSetVariableListFlags(cmd *cobra.Command) (*VariableSetVariableListFlags, error) {
	return &VariableSetVariableListFlags{
		VarsetName:            viper.GetString("varset-name"),
		VarsetID:              viper.GetString("varset-id"),
		VariableSetScopeFlags: parseVariableSetScopeFlagsWithoutSearch(),
	}, nil
}

func ParseVariableSetVariableShowFlags(cmd *cobra.Command) (*VariableSetVariableShowFlags, error) {
	return &VariableSetVariableShowFlags{
		VarsetName:            viper.GetString("varset-name"),
		VarsetID:              viper.GetString("varset-id"),
		Key:                   viper.GetString("key"),
		VariableSetScopeFlags: parseVariableSetScopeFlagsWithoutSearch(),
	}, nil
}

func ParseVariableSetVariableCreateFlags(cmd *cobra.Command) (*VariableSetVariableCreateFlags, error) {
	return &VariableSetVariableCreateFlags{
		VarsetName:            viper.GetString("varset-name"),
		VarsetID:              viper.GetString("varset-id"),
		Key:                   viper.GetString("key"),
		Value:                 viper.GetString("value"),
		ValueFile:             viper.GetString("value-file"),
		Description:           viper.GetString("description"),
		Env:                   viper.GetBool("env"),
		HCL:                   viper.GetBool("hcl"),
		Sensitive:             viper.GetBool("sensitive"),
		VariableSetScopeFlags: parseVariableSetScopeFlagsWithoutSearch(),
	}, nil
}

func ParseVariableSetVariableUpdateFlags(cmd *cobra.Command) (*VariableSetVariableUpdateFlags, error) {
	return &VariableSetVariableUpdateFlags{
		VarsetName:            viper.GetString("varset-name"),
		VarsetID:              viper.GetString("varset-id"),
		Key:                   viper.GetString("key"),
		Value:                 viper.GetString("value"),
		ValueFile:             viper.GetString("value-file"),
		Description:           viper.GetString("description"),
		Env:                   viper.GetBool("env"),
		HCL:                   viper.GetBool("hcl"),
		Sensitive:             viper.GetBool("sensitive"),
		VariableSetScopeFlags: parseVariableSetScopeFlagsWithoutSearch(),
	}, nil
}

func ParseVariableSetVariableDeleteFlags(cmd *cobra.Command) (*VariableSetVariableDeleteFlags, error) {
	return &VariableSetVariableDeleteFlags{
		VarsetName:            viper.GetString("varset-name"),
		VarsetID:              viper.GetString("varset-id"),
		Key:                   viper.GetString("key"),
		VariableSetScopeFlags: parseVariableSetScopeFlagsWithoutSearch(),
	}, nil
}
