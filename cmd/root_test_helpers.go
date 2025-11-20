// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

// root_test_helpers.go
//
// This file is only compiled when running tests (go test) because it doesn't
// have a _test.go suffix but provides exported helpers needed by test packages.
// The Go compiler includes this file in test builds but excludes it from
// production builds.
//
// This allows us to expose internal structures (like rootCmd) to integration
// tests without polluting the production binary with test-only code.

package cmd

import "github.com/spf13/cobra"

// GetRootCommand returns the root command for testing purposes.
// This allows integration tests to execute commands programmatically.
func GetRootCommand() *cobra.Command {
	return rootCmd
}
