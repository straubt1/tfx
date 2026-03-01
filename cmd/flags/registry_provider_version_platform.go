// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RegistryProviderVersionPlatformListFlags holds flags for platform list
type RegistryProviderVersionPlatformListFlags struct {
	Name    string
	Version string
}

// RegistryProviderVersionPlatformCreateFlags holds flags for platform create
type RegistryProviderVersionPlatformCreateFlags struct {
	Name     string
	Version  string
	OS       string
	Arch     string
	Filename string
}

// RegistryProviderVersionPlatformShowFlags holds flags for platform show
type RegistryProviderVersionPlatformShowFlags struct {
	Name    string
	Version string
	OS      string
	Arch    string
}

// RegistryProviderVersionPlatformDeleteFlags holds flags for platform delete
type RegistryProviderVersionPlatformDeleteFlags struct {
	Name    string
	Version string
	OS      string
	Arch    string
}

func ParseRegistryProviderVersionPlatformListFlags(cmd *cobra.Command) (*RegistryProviderVersionPlatformListFlags, error) {
	return &RegistryProviderVersionPlatformListFlags{
		Name:    viper.GetString("name"),
		Version: viper.GetString("version"),
	}, nil
}

func ParseRegistryProviderVersionPlatformCreateFlags(cmd *cobra.Command) (*RegistryProviderVersionPlatformCreateFlags, error) {
	return &RegistryProviderVersionPlatformCreateFlags{
		Name:     viper.GetString("name"),
		Version:  viper.GetString("version"),
		OS:       viper.GetString("os"),
		Arch:     viper.GetString("arch"),
		Filename: viper.GetString("filename"),
	}, nil
}

func ParseRegistryProviderVersionPlatformShowFlags(cmd *cobra.Command) (*RegistryProviderVersionPlatformShowFlags, error) {
	return &RegistryProviderVersionPlatformShowFlags{
		Name:    viper.GetString("name"),
		Version: viper.GetString("version"),
		OS:      viper.GetString("os"),
		Arch:    viper.GetString("arch"),
	}, nil
}

func ParseRegistryProviderVersionPlatformDeleteFlags(cmd *cobra.Command) (*RegistryProviderVersionPlatformDeleteFlags, error) {
	return &RegistryProviderVersionPlatformDeleteFlags{
		Name:    viper.GetString("name"),
		Version: viper.GetString("version"),
		OS:      viper.GetString("os"),
		Arch:    viper.GetString("arch"),
	}, nil
}
