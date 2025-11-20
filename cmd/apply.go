//go:build ignore

// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"github.com/straubt1/tfx/client"
)

var (
	applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "Applies a Workspace Run.",
		Long:  `Creates a Run Apply based on an existing Run Plan and displays its Apply logs.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runApply()
		},
		PreRun: bindPFlags,
	}
)

func init() {
	// All `tfx apply` commands
	applyCmd.PersistentFlags().StringP("run-id", "i", "", "Run Id to apply")

	applyCmd.MarkPersistentFlagRequired("run-id")

	rootCmd.AddCommand(applyCmd)
}

func runApply() error {
	// var err error

	// Validate flags
	runID := *viperString("run-id")

	c, err := client.NewFromViper()
	if err != nil {
		return err
	}

	// Verify run can be applied
	var r *tfe.Run
	r, err = c.Client.Runs.Read(c.Context, runID)
	if err != nil {
		return err
	}

	if !runCanBeApplied(string(r.Status)) {
		return errors.New("run id " + r.ID + " can not be applied. status: " + string(r.Status))
	}

	// Create Apply
	err = c.Client.Runs.Apply(c.Context, r.ID, tfe.RunApplyOptions{
		Comment: tfe.String("TFx did the apply"),
	})
	if err != nil {
		return err
	}

	// Retrieve workspace name using the ID for the URL
	workspace, err := c.Client.Workspaces.ReadByID(c.Context, r.Workspace.ID)
	if err != nil {
		return err
	}
	wsName := workspace.Name

	fmt.Println("Workspace Apply Created, Apply Id:", color.BlueString(r.Apply.ID))
	fmt.Println("Navigate:", "https://"+c.Hostname+"/app/"+c.OrganizationName+"/workspaces/"+wsName+"/runs/"+r.ID)
	fmt.Println()

	err = getApplyLogs(c.Context, c.Client, r.Apply.ID)
	if err != nil {
		return err
	}

	fmt.Println("Apply Complete:", r.Apply.ID)
	return nil
}
