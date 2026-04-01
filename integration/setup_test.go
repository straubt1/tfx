// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

//go:build integration
// +build integration

// Copyright © 2025 Tom Straub <github.com/straubt1>

package integration

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/straubt1/tfx/cmd"
)

// setupTest retrieves test configuration from environment variables
// and skips the test if any required variables are missing
func setupTest(t *testing.T) (hostname, token, organization string) {
	hostname = os.Getenv("TFE_HOSTNAME")
	token = os.Getenv("TFE_TOKEN")
	organization = os.Getenv("TFE_ORGANIZATION")

	if hostname == "" || token == "" || organization == "" {
		t.Skip("Skipping integration test: TFE_HOSTNAME, TFE_TOKEN, and TFE_ORGANIZATION must be set")
	}

	return hostname, token, organization
}

// executeCommand runs a tfx command with the given arguments and credentials
// Returns error if command failed
func executeCommand(t *testing.T, args []string, hostname, token, organization string) error {
	// Get the root command
	// Note: This requires exposing cmd.GetRootCommand() or using cmd.Execute directly
	rootCmd := getRootCommand()

	// Build full args with credentials
	fullArgs := append([]string{
		"--hostname", hostname,
		"--token", token,
		"--organization", organization,
	}, args...)

	rootCmd.SetArgs(fullArgs)

	// Capture output for debugging
	var stdout, stderr bytes.Buffer
	rootCmd.SetOut(&stdout)
	rootCmd.SetErr(&stderr)

	// Execute the command
	err := rootCmd.Execute()

	// Log output for debugging (visible with -v flag)
	if stdout.Len() > 0 {
		t.Logf("STDOUT:\n%s", stdout.String())
	}
	if stderr.Len() > 0 {
		t.Logf("STDERR:\n%s", stderr.String())
	}

	return err
}

// getRootCommand returns the root command for testing
// This is a wrapper to handle the cmd package structure
func getRootCommand() *cobra.Command {
	// We'll need to call cmd.GetRootCommand() once that's exposed
	// For now, we'll use a workaround by accessing the exported Execute function
	return cmd.GetRootCommand()
}
