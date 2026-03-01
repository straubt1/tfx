// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// `tfx registry` commands
	registryCmd = &cobra.Command{
		Use:   "registry",
		Short: "Private Registry Commands",
		Long:  "Commands to work with the Private Registry of a TFx Organization.",
	}
)

func init() {
	rootCmd.AddCommand(registryCmd)
}
