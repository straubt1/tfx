// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	"time"

	"github.com/araddon/dateparse"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/cmd/flags"
	view "github.com/straubt1/tfx/cmd/views"
	"github.com/straubt1/tfx/data"
)

// metricsCmd represents the metrics command
var (
	// `tfx admin metrics` command
	metricsCmd = &cobra.Command{
		Use:    "metrics",
		Short:  "Read metrics about TFx Usage",
		Long:   "Read details about how TFx is being used. This command can take a while to execute.",
		Hidden: true, // hide until this is better defined
	}

	// `tfx admin metrics workspace` command
	metricsWorkspaceCmd = &cobra.Command{
		Use:   "workspace",
		Short: "Read metrics about TFx Workspace Usage",
		Long:  "Read details about how TFx Workspaces are being used. This command can take a while to execute.",
		Example: `
tfx admin metrics workspace

tfx admin metrics workspace --since "01/31/2021 10:30"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseAdminMetricsWorkspaceFlags(cmd)
			if err != nil {
				return err
			}
			return metricsWorkspace(cmdConfig)
		},
	}
)

func init() {
	// `tfx admin metrics workspace` flags
	metricsWorkspaceCmd.Flags().StringP("since", "s", "", "Start time when querying runs in the format MM/DD/YYYY hh:mm:ss. Examples: ['01/31/2021 10:30', '02/28/2021 10:30 AM', '03/20/2021'] (optional).")

	adminCmd.AddCommand(metricsCmd)
	metricsCmd.AddCommand(metricsWorkspaceCmd)
}

func metricsWorkspace(cmdConfig *flags.AdminMetricsWorkspaceFlags) error {
	// Create view for rendering
	v := view.NewAdminMetricsWorkspaceView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	// Parse the since time
	since, err := parseTime(cmdConfig.Since)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to parse since time"))
	}

	// Print command header before API call
	if cmdConfig.Since != "" {
		v.PrintCommandHeader("Getting metrics for all workspaces in organization '%s' since '%s'", c.OrganizationName, cmdConfig.Since)
	} else {
		v.PrintCommandHeader("Getting metrics for all workspaces in organization '%s'", c.OrganizationName)
	}
	v.PrintCommandHeader("This can take some time to complete...")

	// Start timer
	start := time.Now()

	// Fetch workspace metrics
	workspaceMetrics, err := data.FetchWorkspaceMetrics(c, c.OrganizationName, since)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to get workspace metrics"))
	}

	// Calculate query time
	elapsed := time.Since(start)

	// Build result
	result := &view.MetricsWorkspaceResult{
		Workspaces: workspaceMetrics,
		Since:      since,
		QueryTime:  elapsed.String(),
	}

	return v.Render(result)
}

// parseTime parses a time string, returning zero time if empty
func parseTime(s string) (time.Time, error) {
	zeroTime := time.Time{}
	// nothing passed, include all time
	if s == "" {
		return zeroTime, nil
	}
	// "3/1/2014 10:25 PM"
	t, err := dateparse.ParseLocal(s)
	if err != nil {
		return zeroTime, err
	}

	return t, nil
}
