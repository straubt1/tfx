// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/straubt1/tfx/output"
	"github.com/straubt1/tfx/tui"
)

var loginCmd = &cobra.Command{
	Use:   "login [hostname]",
	Short: "Authenticate to HCP Terraform or Terraform Enterprise",
	Long: `Authenticate to HCP Terraform or Terraform Enterprise.

Opens your browser to the API token creation page, prompts you to paste the
token, then saves it (along with the selected organization) to ~/.tfx.hcl.

If no hostname is provided the default is app.terraform.io.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runLogin,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func runLogin(cmd *cobra.Command, args []string) error {
	output.Get().DisableSpinner()

	hostname := "app.terraform.io"
	if len(args) == 1 {
		hostname = strings.TrimSpace(args[0])
	}

	return tui.RunLogin(hostname)
}
