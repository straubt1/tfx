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
		tapePath, _ := cmd.Flags().GetString("tape")
		return tui.Run(tapePath)
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
	tuiCmd.Flags().String("tape", "", "Record TUI input to a .tape file for VHS (e.g. debug/demo.tape)")
	tuiCmd.Flags().MarkHidden("tape")
}
