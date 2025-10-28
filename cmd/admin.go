// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// `tfx admin` commands
	adminCmd = &cobra.Command{
		Use:   "admin",
		Short: "Admin Commands",
		Long:  "Work with TFx Admin Operations",
		Example: `
List metrics for all workspaces in an organization:
tfx admin metrics workspace

List metrics for all workspaces since a specific date:
tfx admin metrics workspace --since "01/31/2021"`,
	}
)

func init() {
	rootCmd.AddCommand(adminCmd)
}
