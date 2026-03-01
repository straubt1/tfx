// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// `tfx release` commands
	releaseCmd = &cobra.Command{
		Use:   "release",
		Short: "Release commands",
		Long:  "Work with releases needed for airgap Terraform Enterprise installations.",
	}
)

func init() {
	rootCmd.AddCommand(releaseCmd)
}
