// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type WorkspaceLockFlags struct {
	Name string
}

type WorkspaceLockAllFlags struct {
	Search string
}

type WorkspaceUnlockFlags struct {
	Name string
}

type WorkspaceUnlockAllFlags struct {
	Search string
}

func ParseWorkspaceLockFlags(cmd *cobra.Command) (*WorkspaceLockFlags, error) {
	return &WorkspaceLockFlags{Name: viper.GetString("name")}, nil
}

func ParseWorkspaceLockAllFlags(cmd *cobra.Command) (*WorkspaceLockAllFlags, error) {
	return &WorkspaceLockAllFlags{Search: viper.GetString("search")}, nil
}

func ParseWorkspaceUnlockFlags(cmd *cobra.Command) (*WorkspaceUnlockFlags, error) {
	return &WorkspaceUnlockFlags{Name: viper.GetString("name")}, nil
}

func ParseWorkspaceUnlockAllFlags(cmd *cobra.Command) (*WorkspaceUnlockAllFlags, error) {
	return &WorkspaceUnlockAllFlags{Search: viper.GetString("search")}, nil
}
