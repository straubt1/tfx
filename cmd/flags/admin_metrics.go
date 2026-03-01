// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AdminMetricsWorkspaceFlags holds all flags for the admin metrics workspace command
type AdminMetricsWorkspaceFlags struct {
	Since string
}

// ParseAdminMetricsWorkspaceFlags creates AdminMetricsWorkspaceFlags from the current command context
func ParseAdminMetricsWorkspaceFlags(cmd *cobra.Command) (*AdminMetricsWorkspaceFlags, error) {
	return &AdminMetricsWorkspaceFlags{
		Since: viper.GetString("since"),
	}, nil
}
