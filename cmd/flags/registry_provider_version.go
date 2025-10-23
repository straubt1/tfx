// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RegistryProviderVersionListFlags holds flags for provider version list
type RegistryProviderVersionListFlags struct {
	Name string
}

// RegistryProviderVersionCreateFlags holds flags for provider version create
type RegistryProviderVersionCreateFlags struct {
	Name       string
	Version    string
	KeyID      string
	Shasums    string
	ShasumsSig string
}

// RegistryProviderVersionShowFlags holds flags for provider version show
type RegistryProviderVersionShowFlags struct {
	Name    string
	Version string
}

// RegistryProviderVersionDeleteFlags holds flags for provider version delete
type RegistryProviderVersionDeleteFlags struct {
	Name    string
	Version string
}

func ParseRegistryProviderVersionListFlags(cmd *cobra.Command) (*RegistryProviderVersionListFlags, error) {
	return &RegistryProviderVersionListFlags{
		Name: viper.GetString("name"),
	}, nil
}

func ParseRegistryProviderVersionCreateFlags(cmd *cobra.Command) (*RegistryProviderVersionCreateFlags, error) {
	return &RegistryProviderVersionCreateFlags{
		Name:       viper.GetString("name"),
		Version:    viper.GetString("version"),
		KeyID:      viper.GetString("key-id"),
		Shasums:    viper.GetString("shasums"),
		ShasumsSig: viper.GetString("shasums-sig"),
	}, nil
}

func ParseRegistryProviderVersionShowFlags(cmd *cobra.Command) (*RegistryProviderVersionShowFlags, error) {
	return &RegistryProviderVersionShowFlags{
		Name:    viper.GetString("name"),
		Version: viper.GetString("version"),
	}, nil
}

func ParseRegistryProviderVersionDeleteFlags(cmd *cobra.Command) (*RegistryProviderVersionDeleteFlags, error) {
	return &RegistryProviderVersionDeleteFlags{
		Name:    viper.GetString("name"),
		Version: viper.GetString("version"),
	}, nil
}
