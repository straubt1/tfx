// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RunListFlags holds flags for list runs
type RunListFlags struct {
	WorkspaceName string
	MaxItems      int
}

// RunCreateFlags holds flags for create run
type RunCreateFlags struct {
	WorkspaceName          string
	Message                string
	ConfigurationVersionID string
}

// RunShowFlags holds flags for show run
type RunShowFlags struct {
	ID string
}

// RunDiscardFlags holds flags for discard run
type RunDiscardFlags struct {
	ID string
}

// RunCancelFlags holds flags for cancel run
type RunCancelFlags struct {
	WorkspaceName string
}

func ParseRunListFlags(cmd *cobra.Command) (*RunListFlags, error) {
	return &RunListFlags{
		WorkspaceName: viper.GetString("workspace-name"),
		MaxItems:      viper.GetInt("max-items"),
	}, nil
}

func ParseRunCreateFlags(cmd *cobra.Command) (*RunCreateFlags, error) {
	return &RunCreateFlags{
		WorkspaceName:          viper.GetString("workspace-name"),
		Message:                viper.GetString("message"),
		ConfigurationVersionID: viper.GetString("configuration-version-id"),
	}, nil
}

func ParseRunShowFlags(cmd *cobra.Command) (*RunShowFlags, error) {
	return &RunShowFlags{ID: viper.GetString("id")}, nil
}

func ParseRunDiscardFlags(cmd *cobra.Command) (*RunDiscardFlags, error) {
	return &RunDiscardFlags{ID: viper.GetString("id")}, nil
}

func ParseRunCancelFlags(cmd *cobra.Command) (*RunCancelFlags, error) {
	return &RunCancelFlags{WorkspaceName: viper.GetString("workspace-name")}, nil
}
