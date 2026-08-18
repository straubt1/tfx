// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// userChangedCommandFlags records local flags that were set on the command line
// before postInitCommands copies Viper values onto Cobra flags (which marks
// them Changed even when the user did not pass them).
var userChangedCommandFlags map[string]bool

func commandFlagKey(cmd *cobra.Command, name string) string {
	return cmd.CommandPath() + "\x00" + name
}

// CaptureCommandFlagChanges snapshots which local flags were explicitly passed
// on each command. Call this before postInitCommands.
func CaptureCommandFlagChanges(root *cobra.Command) {
	userChangedCommandFlags = make(map[string]bool)
	captureCommandFlagChanges(root)
}

func captureCommandFlagChanges(cmd *cobra.Command) {
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			userChangedCommandFlags[commandFlagKey(cmd, f.Name)] = true
		}
	})
	for _, child := range cmd.Commands() {
		captureCommandFlagChanges(child)
	}
}

// ResetCapturedFlagChanges clears the snapshot. Tests only.
func ResetCapturedFlagChanges() {
	userChangedCommandFlags = nil
}

// UserPassed reports whether name was set on the command line for cmd.
// After CaptureCommandFlagChanges, this ignores later Flag.Set calls (Viper leak).
// Before a capture, it falls back to cmd.Flags().Changed.
func UserPassed(cmd *cobra.Command, name string) bool {
	if cmd == nil {
		return false
	}
	if userChangedCommandFlags != nil {
		return userChangedCommandFlags[commandFlagKey(cmd, name)]
	}
	return cmd.Flags().Changed(name)
}
