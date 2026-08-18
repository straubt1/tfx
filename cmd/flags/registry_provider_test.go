// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package flags

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestParseRegistryProviderDownloadFlagsDirectoryDefault(t *testing.T) {
	t.Cleanup(ResetCapturedFlagChanges)
	viper.Reset()

	cmd := &cobra.Command{Use: "download"}
	cmd.Flags().StringP("directory", "d", "", "")
	CaptureCommandFlagChanges(cmd)

	got, err := ParseRegistryProviderDownloadFlags(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if got.Directory != DefaultPublicProviderDirectory {
		t.Fatalf("directory = %q, want %q", got.Directory, DefaultPublicProviderDirectory)
	}
}

func TestParseRegistryProviderDownloadFlagsDirectoryUserPassed(t *testing.T) {
	t.Cleanup(ResetCapturedFlagChanges)
	viper.Reset()

	cmd := &cobra.Command{Use: "download"}
	cmd.Flags().StringP("directory", "d", "", "")
	if err := cmd.Flags().Set("directory", "./custom-providers"); err != nil {
		t.Fatal(err)
	}
	CaptureCommandFlagChanges(cmd)

	got, err := ParseRegistryProviderDownloadFlags(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if got.Directory != "./custom-providers" {
		t.Fatalf("directory = %q, want ./custom-providers", got.Directory)
	}
}

func TestParseRegistryProviderDownloadFlagsIgnoresLeakedDirectory(t *testing.T) {
	t.Cleanup(ResetCapturedFlagChanges)
	viper.Reset()
	viper.Set("directory", "./leaked")

	cmd := &cobra.Command{Use: "download"}
	cmd.Flags().StringP("directory", "d", "", "")
	CaptureCommandFlagChanges(cmd)
	if err := cmd.Flags().Set("directory", "./leaked"); err != nil {
		t.Fatal(err)
	}

	got, err := ParseRegistryProviderDownloadFlags(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if got.Directory != DefaultPublicProviderDirectory {
		t.Fatalf("directory = %q, want default %q", got.Directory, DefaultPublicProviderDirectory)
	}
}
