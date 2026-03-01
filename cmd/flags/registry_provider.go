// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RegistryProviderListFlags holds flags for provider list
type RegistryProviderListFlags struct {
	MaxItems int
	All      bool
}

// RegistryProviderCreateFlags holds flags for provider create
type RegistryProviderCreateFlags struct {
	Name string
}

// RegistryProviderShowFlags holds flags for provider show
type RegistryProviderShowFlags struct {
	Name string
}

// RegistryProviderDeleteFlags holds flags for provider delete
type RegistryProviderDeleteFlags struct {
	Name string
}

func ParseRegistryProviderListFlags(cmd *cobra.Command) (*RegistryProviderListFlags, error) {
	return &RegistryProviderListFlags{
		MaxItems: viper.GetInt("max-items"),
		All:      viper.GetBool("all"),
	}, nil
}

func ParseRegistryProviderCreateFlags(cmd *cobra.Command) (*RegistryProviderCreateFlags, error) {
	return &RegistryProviderCreateFlags{
		Name: viper.GetString("name"),
	}, nil
}

func ParseRegistryProviderShowFlags(cmd *cobra.Command) (*RegistryProviderShowFlags, error) {
	return &RegistryProviderShowFlags{
		Name: viper.GetString("name"),
	}, nil
}

func ParseRegistryProviderDeleteFlags(cmd *cobra.Command) (*RegistryProviderDeleteFlags, error) {
	return &RegistryProviderDeleteFlags{
		Name: viper.GetString("name"),
	}, nil
}
