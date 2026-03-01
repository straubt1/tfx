// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AdminGPGListFlags holds all flags for the admin gpg list command
type AdminGPGListFlags struct {
	Namespace    string
	RegistryName string
}

// AdminGPGCreateFlags holds all flags for the admin gpg create command
type AdminGPGCreateFlags struct {
	Namespace    string
	PublicKey    string
	RegistryName string
}

// AdminGPGShowFlags holds all flags for the admin gpg show command
type AdminGPGShowFlags struct {
	Namespace    string
	ID           string
	RegistryName string
}

// AdminGPGDeleteFlags holds all flags for the admin gpg delete command
type AdminGPGDeleteFlags struct {
	Namespace    string
	ID           string
	RegistryName string
}

// ParseAdminGPGListFlags creates AdminGPGListFlags from the current command context
func ParseAdminGPGListFlags(cmd *cobra.Command) (*AdminGPGListFlags, error) {
	return &AdminGPGListFlags{
		Namespace:    viper.GetString("namespace"),
		RegistryName: viper.GetString("registry-name"),
	}, nil
}

// ParseAdminGPGCreateFlags creates AdminGPGCreateFlags from the current command context
func ParseAdminGPGCreateFlags(cmd *cobra.Command) (*AdminGPGCreateFlags, error) {
	return &AdminGPGCreateFlags{
		Namespace:    viper.GetString("namespace"),
		PublicKey:    viper.GetString("public-key"),
		RegistryName: viper.GetString("registry-name"),
	}, nil
}

// ParseAdminGPGShowFlags creates AdminGPGShowFlags from the current command context
func ParseAdminGPGShowFlags(cmd *cobra.Command) (*AdminGPGShowFlags, error) {
	return &AdminGPGShowFlags{
		Namespace:    viper.GetString("namespace"),
		ID:           viper.GetString("id"),
		RegistryName: viper.GetString("registry-name"),
	}, nil
}

// ParseAdminGPGDeleteFlags creates AdminGPGDeleteFlags from the current command context
func ParseAdminGPGDeleteFlags(cmd *cobra.Command) (*AdminGPGDeleteFlags, error) {
	return &AdminGPGDeleteFlags{
		Namespace:    viper.GetString("namespace"),
		ID:           viper.GetString("id"),
		RegistryName: viper.GetString("registry-name"),
	}, nil
}
