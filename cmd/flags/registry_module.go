// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RegistryModuleListFlags holds flags for module list
type RegistryModuleListFlags struct {
	MaxItems int
	All      bool
}

// RegistryModuleCreateFlags holds flags for module create
type RegistryModuleCreateFlags struct {
	Name     string
	Provider string
}

// RegistryModuleShowFlags holds flags for module show
type RegistryModuleShowFlags struct {
	Name     string
	Provider string
}

// RegistryModuleDeleteFlags holds flags for module delete
type RegistryModuleDeleteFlags struct {
	Name     string
	Provider string
}

func ParseRegistryModuleListFlags(cmd *cobra.Command) (*RegistryModuleListFlags, error) {
	return &RegistryModuleListFlags{
		MaxItems: viper.GetInt("max-items"),
		All:      viper.GetBool("all"),
	}, nil
}

func ParseRegistryModuleCreateFlags(cmd *cobra.Command) (*RegistryModuleCreateFlags, error) {
	return &RegistryModuleCreateFlags{
		Name:     viper.GetString("name"),
		Provider: viper.GetString("provider"),
	}, nil
}

func ParseRegistryModuleShowFlags(cmd *cobra.Command) (*RegistryModuleShowFlags, error) {
	return &RegistryModuleShowFlags{
		Name:     viper.GetString("name"),
		Provider: viper.GetString("provider"),
	}, nil
}

func ParseRegistryModuleDeleteFlags(cmd *cobra.Command) (*RegistryModuleDeleteFlags, error) {
	return &RegistryModuleDeleteFlags{
		Name:     viper.GetString("name"),
		Provider: viper.GetString("provider"),
	}, nil
}
