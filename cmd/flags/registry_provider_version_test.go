// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package flags

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func versionCreateCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "create"}
	cmd.Flags().StringP("directory", "d", "", "")
	cmd.Flags().Int("concurrency", 4, "")
	return cmd
}

func TestParseRegistryProviderVersionCreateFlagsDirectoryEmptyByDefault(t *testing.T) {
	t.Cleanup(ResetCapturedFlagChanges)
	viper.Reset()

	cmd := versionCreateCmd()
	CaptureCommandFlagChanges(cmd)

	got, err := ParseRegistryProviderVersionCreateFlags(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if got.Directory != "" {
		t.Fatalf("directory = %q, want empty", got.Directory)
	}
	if got.Concurrency != 4 {
		t.Fatalf("concurrency = %d, want 4", got.Concurrency)
	}
}

func TestParseRegistryProviderVersionCreateFlagsDirectoryUserPassed(t *testing.T) {
	t.Cleanup(ResetCapturedFlagChanges)
	viper.Reset()

	cmd := versionCreateCmd()
	if err := cmd.Flags().Set("directory", "./providers/hashicorp/azurerm/5.0.0"); err != nil {
		t.Fatal(err)
	}
	CaptureCommandFlagChanges(cmd)

	got, err := ParseRegistryProviderVersionCreateFlags(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if got.Directory != "./providers/hashicorp/azurerm/5.0.0" {
		t.Fatalf("directory = %q", got.Directory)
	}
}

func TestParseRegistryProviderVersionCreateFlagsIgnoresLeakedDirectory(t *testing.T) {
	t.Cleanup(ResetCapturedFlagChanges)
	viper.Reset()
	viper.Set("directory", "./providers")

	cmd := versionCreateCmd()
	CaptureCommandFlagChanges(cmd)
	if err := cmd.Flags().Set("directory", "./providers"); err != nil {
		t.Fatal(err)
	}

	got, err := ParseRegistryProviderVersionCreateFlags(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if got.Directory != "" {
		t.Fatalf("leaked directory = %q, want empty", got.Directory)
	}
}
