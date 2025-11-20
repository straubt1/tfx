//go:build ignore

// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/fatih/color"
	"github.com/hashicorp/go-slug"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"github.com/straubt1/tfx/client"
)

var (
	planCmd = &cobra.Command{
		Use:   "plan",
		Short: "Workspace Plan",
		Long:  "Work with Plans of a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPlan()
		},
		PreRun: bindPFlags,
	}

	planExportCmd = &cobra.Command{
		Use:   "export",
		Short: "Export plan",
		Long:  "Export plan details for a Plan.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPlanExport()
		},
		PreRun: bindPFlags,
	}
)

func init() {
	// All `tfx plan` commands
	planCmd.Flags().StringP("workspace-name", "w", "", "Workspace name")
	planCmd.Flags().StringP("directory", "d", "./", "Directory of Terraform (optional, defaults to current directory)")
	planCmd.Flags().StringP("configuration-id", "i", "", "Configuration Version Id (optional, i.e. cv-*)")
	planCmd.Flags().StringP("message", "m", "", "Run Message (optional)")
	planCmd.Flags().Bool("speculative", false, "Perform a Speculative Plan (optional)")
	planCmd.Flags().Bool("destroy", false, "Perform a Destroy Plan (optional)")
	planCmd.Flags().StringSlice("env", []string{}, "Environment variables to write to the Workspace. Can be supplied multiple times. (optional, i.e. '--env='AWS_REGION=us-east1')")
	planCmd.MarkFlagRequired("workspace-name")

	planExportCmd.Flags().StringP("plan-id", "i", "", "Plan Id (i.e. plan-*)")
	planExportCmd.Flags().StringP("directory", "d", "", "Directory to download export to (optional, defaults to a temp directory)")
	planExportCmd.MarkFlagRequired("plan-id")
	planExportCmd.MarkFlagRequired("directory")

	rootCmd.AddCommand(planCmd)
	planCmd.AddCommand(planExportCmd)
}

func runPlan() error {
	// Validate flags
	wsName := *viperString("workspace-name")
	dir := *viperString("directory")
	configID := *viperString("configuration-id")
	message := *viperString("message")
	isSpeculative := *viperBool("speculative")
	isDestroy := *viperBool("destroy")
	envs, err := viperStringSliceMap("env")
	if err != nil {
		return err
	}
	// message := "Plan created by TFx"
	c, err := client.NewFromViper()
	if err != nil {
		return err
	}

	fmt.Println("Remote Terraform Plan, speculative plan:", color.GreenString(strconv.FormatBool(isSpeculative)),
		" destroy plan:", color.GreenString(strconv.FormatBool(isDestroy)))
	fmt.Println()

	// Read workspace
	fmt.Print("Reading Workspace ", color.GreenString(wsName), " for ID...")
	w, err := c.Client.Workspaces.Read(c.Context, c.OrganizationName, wsName)
	if err != nil {
		return err
	}
	fmt.Println(" Found:", color.BlueString(w.ID))

	// Update environment variables
	err = createOrUpdateEnvVariables(c.Context, c.Client, w.ID, envs)
	if err != nil {
		return err
	}

	// Create config version
	cv, err := createOrReadConfigurationVersion(c.Context, c.Client, w.ID, configID, dir, isSpeculative)
	if err != nil {
		return err
	}

	// Create run
	var r *tfe.Run
	r, err = c.Client.Runs.Create(c.Context, tfe.RunCreateOptions{
		IsDestroy:            tfe.Bool(isDestroy),
		Message:              tfe.String(message),
		ConfigurationVersion: cv,
		Workspace:            w,
		// TargetAddrs:          []string{"random_pet.two"},
	})
	if err != nil {
		return err
	}
	fmt.Println("Workspace Run Created, Run Id:", color.BlueString(r.ID), "Config Version:", color.BlueString(r.ConfigurationVersion.ID))
	fmt.Println("Navigate:", "https://"+c.Hostname+"/app/"+c.OrganizationName+"/workspaces/"+wsName+"/runs/"+r.ID)
	fmt.Println()

	// Get Run Logs
	err = getRunLogs(c.Context, c.Client, r.Plan.ID)
	if err != nil {
		return err
	}

	// Get Cost Estimation Logs (if any)
	err = getCostEstimationLogs(c.Context, c.Client, r)
	if err != nil {
		return err
	}

	// Get Policy Logs (if any)
	err = getPolicyLogs(c.Context, c.Client, r)
	if err != nil {
		return err
	}

	fmt.Println("Run Complete:", r.ID)

	return nil
}

func runPlanExport() error {
	planID := *viperString("plan-id")
	directory := *viperString("directory")
	c, err := client.NewFromViper()
	if err != nil {
		return err
	}

	fmt.Print("Reading Plan Export for Plan ID ", color.GreenString(planID), " ...")
	plan, err := c.Client.Plans.Read(c.Context, planID)
	if err != nil {
		return err
	}
	fmt.Println(" Found")

	var planExportID string
	if plan.Exports == nil {
		fmt.Print("Creating Plan Export ...")
		planExport, err := c.Client.PlanExports.Create(c.Context, tfe.PlanExportCreateOptions{
			Plan:     plan,
			DataType: tfe.PlanExportType(tfe.PlanExportSentinelMockBundleV0),
		})
		if err != nil {
			return err
		}
		planExportID = planExport.ID
	} else {
		fmt.Print("Found existing Plan Export ...")
		planExportID = plan.Exports[0].ID // Just grab the first one?
	}
	fmt.Println("ID ", color.BlueString(planExportID))
	buff, err := c.Client.PlanExports.Download(c.Context, planExportID)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(buff)

	// Determine a directory to unpack the slug contents into.
	if directory != "" {
		directory, err = filepath.Abs(directory)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("Directory not supplied, creating a temp directory")
		dst, err := ioutil.TempDir("", "slug")
		if err != nil {
			return err
		}
		directory = dst
	}

	if err := slug.Unpack(reader, directory); err != nil {
		return err
	}

	fmt.Println("Downloaded: ", color.BlueString(directory))

	return nil
}
