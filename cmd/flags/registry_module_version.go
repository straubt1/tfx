// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RegistryModuleVersionListFlags holds flags for module version list
type RegistryModuleVersionListFlags struct {
	Name     string
	Provider string
}

// RegistryModuleVersionCreateFlags holds flags for module version create
type RegistryModuleVersionCreateFlags struct {
	Name      string
	Provider  string
	Version   string
	Directory string
}

// RegistryModuleVersionDeleteFlags holds flags for module version delete
type RegistryModuleVersionDeleteFlags struct {
	Name     string
	Provider string
	Version  string
}

// RegistryModuleVersionDownloadFlags holds flags for module version download
type RegistryModuleVersionDownloadFlags struct {
	Name      string
	Provider  string
	Version   string
	Directory string
}

func ParseRegistryModuleVersionListFlags(cmd *cobra.Command) (*RegistryModuleVersionListFlags, error) {
	return &RegistryModuleVersionListFlags{
		Name:     viper.GetString("name"),
		Provider: viper.GetString("provider"),
	}, nil
}

func ParseRegistryModuleVersionCreateFlags(cmd *cobra.Command) (*RegistryModuleVersionCreateFlags, error) {
	return &RegistryModuleVersionCreateFlags{
		Name:      viper.GetString("name"),
		Provider:  viper.GetString("provider"),
		Version:   viper.GetString("version"),
		Directory: viper.GetString("directory"),
	}, nil
}

func ParseRegistryModuleVersionDeleteFlags(cmd *cobra.Command) (*RegistryModuleVersionDeleteFlags, error) {
	return &RegistryModuleVersionDeleteFlags{
		Name:     viper.GetString("name"),
		Provider: viper.GetString("provider"),
		Version:  viper.GetString("version"),
	}, nil
}

func ParseRegistryModuleVersionDownloadFlags(cmd *cobra.Command) (*RegistryModuleVersionDownloadFlags, error) {
	return &RegistryModuleVersionDownloadFlags{
		Name:      viper.GetString("name"),
		Provider:  viper.GetString("provider"),
		Version:   viper.GetString("version"),
		Directory: viper.GetString("directory"),
	}, nil
}
