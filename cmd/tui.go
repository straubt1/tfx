// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/straubt1/tfx/tui"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the interactive TUI",
	Long:  "Launch an interactive terminal UI for browsing TFE/HCP Terraform resources.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.Run()
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
