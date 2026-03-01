// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type CVListFlags struct {
	WorkspaceName string
	MaxItems      int
}

type CVCreateFlags struct {
	WorkspaceName string
	Directory     string
	Speculative   bool
}

type CVShowFlags struct {
	ID string
}

type CVDownloadFlags struct {
	ID        string
	Directory string
}

func ParseCVListFlags(cmd *cobra.Command) (*CVListFlags, error) {
	return &CVListFlags{
		WorkspaceName: viper.GetString("name"),
		MaxItems:      viper.GetInt("max-items"),
	}, nil
}

func ParseCVCreateFlags(cmd *cobra.Command) (*CVCreateFlags, error) {
	return &CVCreateFlags{
		WorkspaceName: viper.GetString("name"),
		Directory:     viper.GetString("directory"),
		Speculative:   viper.GetBool("speculative"),
	}, nil
}

func ParseCVShowFlags(cmd *cobra.Command) (*CVShowFlags, error) {
	return &CVShowFlags{ID: viper.GetString("id")}, nil
}

func ParseCVDownloadFlags(cmd *cobra.Command) (*CVDownloadFlags, error) {
	return &CVDownloadFlags{
		ID:        viper.GetString("id"),
		Directory: viper.GetString("directory"),
	}, nil
}
