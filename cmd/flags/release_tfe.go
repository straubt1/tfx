// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ReleaseTfeListFlags holds all flags for the release tfe list command
type ReleaseTfeListFlags struct {
	TfeLicensePath string
	RegistryURL    string
	MaxItems       int
	All            bool
	StableOnly     bool
}

// ParseReleaseTfeListFlags creates a ReleaseTfeListFlags from the current command context
func ParseReleaseTfeListFlags(cmd *cobra.Command) (*ReleaseTfeListFlags, error) {
	return &ReleaseTfeListFlags{
		TfeLicensePath: viper.GetString("tfe-license-path"),
		RegistryURL:    viper.GetString("registry-url"),
		MaxItems:       viper.GetInt("max-items"),
		All:            viper.GetBool("all"),
		StableOnly:     viper.GetBool("stable-only"),
	}, nil
}

// ReleaseTfeShowFlags holds all flags for the release tfe show command
type ReleaseTfeShowFlags struct {
	Tag            string
	TfeLicensePath string
	RegistryURL    string
}

// ParseReleaseTfeShowFlags creates a ReleaseTfeShowFlags from the current command context
func ParseReleaseTfeShowFlags(cmd *cobra.Command) (*ReleaseTfeShowFlags, error) {
	return &ReleaseTfeShowFlags{
		Tag:            viper.GetString("tag"),
		TfeLicensePath: viper.GetString("tfe-license-path"),
		RegistryURL:    viper.GetString("registry-url"),
	}, nil
}

// ReleaseTfeDownloadFlags holds all flags for the release tfe download command
type ReleaseTfeDownloadFlags struct {
	Tag            string
	TfeLicensePath string
	RegistryURL    string
	Output         string
}

// ParseReleaseTfeDownloadFlags creates a ReleaseTfeDownloadFlags from the current command context
func ParseReleaseTfeDownloadFlags(cmd *cobra.Command) (*ReleaseTfeDownloadFlags, error) {
	return &ReleaseTfeDownloadFlags{
		Tag:            viper.GetString("tag"),
		TfeLicensePath: viper.GetString("tfe-license-path"),
		RegistryURL:    viper.GetString("registry-url"),
		Output:         viper.GetString("output"),
	}, nil
}
