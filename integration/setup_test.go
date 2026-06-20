// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

//go:build integration
// +build integration

package integration

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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

const localIntegrationProfile = "local"

// integrationSkipCleanup reports whether integration tests should leave created resources
// in place. Set TFX_INTEGRATION_NO_CLEANUP=1 (or "true"/"yes") to skip delete steps
// and t.Cleanup removal.
func integrationSkipCleanup() bool {
	switch os.Getenv("TFX_INTEGRATION_NO_CLEANUP") {
	case "1", "true", "yes":
		return true
	default:
		return false
	}
}

// setupProfileIntegrationTest skips unless TFX_INTEGRATION_PROFILE matches profileName.
// Use with a ~/.tfx.hcl profile (e.g. "local" for local.tfe.rocks).
func setupProfileIntegrationTest(t *testing.T, profileName string) {
	t.Helper()
	if os.Getenv("TFX_INTEGRATION_PROFILE") != profileName {
		t.Skipf("Skipping profile integration test: set TFX_INTEGRATION_PROFILE=%q (requires profile %q in .tfx.hcl)", profileName, profileName)
	}
}

// executeCommand runs a tfx command with the given arguments and credentials.
func executeCommand(t *testing.T, args []string, hostname, token, organization string) error {
	t.Helper()
	prefix := []string{
		"--hostname", hostname,
		"--token", token,
		"--organization", organization,
	}
	return runRootCommand(t, append(prefix, args...))
}

// executeCommandWithProfile runs a tfx command using a named config profile.
func executeCommandWithProfile(t *testing.T, profileName string, args []string) error {
	t.Helper()
	prefix := []string{"--profile", profileName}
	if configFile := os.Getenv("TFX_CONFIG_FILE"); configFile != "" {
		prefix = append(prefix, "--config-file", configFile)
	}
	return runRootCommand(t, append(prefix, args...))
}

func runRootCommand(t *testing.T, args []string) error {
	t.Helper()
	viper.Reset()

	rootCmd := getRootCommand()
	resetCommandFlags(rootCmd)
	rootCmd.SetArgs(args)

	var stdout, stderr bytes.Buffer
	rootCmd.SetOut(&stdout)
	rootCmd.SetErr(&stderr)

	err := rootCmd.Execute()

	if stdout.Len() > 0 {
		t.Logf("STDOUT:\n%s", stdout.String())
	}
	if stderr.Len() > 0 {
		t.Logf("STDERR:\n%s", stderr.String())
	}

	return err
}

func resetCommandFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		_ = f.Value.Set(f.DefValue)
		f.Changed = false
	})
	for _, sub := range cmd.Commands() {
		resetCommandFlags(sub)
	}
}

// getRootCommand returns the root command for testing
// This is a wrapper to handle the cmd package structure
func getRootCommand() *cobra.Command {
	// We'll need to call cmd.GetRootCommand() once that's exposed
	// For now, we'll use a workaround by accessing the exported Execute function
	return cmd.GetRootCommand()
}
