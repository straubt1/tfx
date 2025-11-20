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
	hostname := *viperString("tfeHostname")
	orgName := *viperString("tfeOrganization")
	runID := *viperString("run-id")

	client, ctx := getClientContext()

	// Verify run can be applied
	var r *tfe.Run
	r, err := client.Runs.Read(ctx, runID)
	if err != nil {
		logError(err, "failed to read run")
	}

	if !runCanBeApplied(string(r.Status)) {
		logError(errors.New("run id "+r.ID+" can not be applied. status: "+string(r.Status)),
			"Unable to apply run")
	}

	// Create Apply
	err = client.Runs.Apply(ctx, r.ID, tfe.RunApplyOptions{
		Comment: tfe.String("TFx did the apply"),
	})
	if err != nil {
		logError(err, "failed to create apply")
	}

	// Retrieve workspace name using the ID for the URL
	workspace, err := client.Workspaces.ReadByID(ctx, r.Workspace.ID)
	if err != nil {
		logError(err, "failed to read workspace")
	}
	wsName := workspace.Name

	fmt.Println("Workspace Apply Created, Apply Id:", color.BlueString(r.Apply.ID))
	fmt.Println("Navigate:", "https://"+hostname+"/app/"+orgName+"/workspaces/"+wsName+"/runs/"+r.ID)
	fmt.Println()

	getApplyLogs(ctx, client, r.Apply.ID)
	if err != nil {
		logError(err, "failed to read apply logs")
	}

	fmt.Println("Apply Complete:", r.Apply.ID)
	return nil
}
