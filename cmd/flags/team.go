// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TeamListFlags holds all flags for the team list command
type TeamListFlags struct {
	WorkspaceName string
}

// ParseTeamListFlags creates a TeamListFlags from the current command context
func ParseTeamListFlags(cmd *cobra.Command) (*TeamListFlags, error) {
	return &TeamListFlags{
		WorkspaceName: viper.GetString("name"),
	}, nil
}
