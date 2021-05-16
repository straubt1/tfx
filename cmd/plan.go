/*
Copyright © 2021 Tom Straub <github.com/straubt1>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
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
)

func init() {
	// All `tfx plan` commands
	planCmd.PersistentFlags().StringP("workspaceName", "w", "", "Workspace name")
	planCmd.Flags().StringP("directory", "d", "./", "Directory of Terraform (optional, defaults to current directory)")
	planCmd.Flags().StringP("configurationId", "i", "", "Configuration Version Id (optional, i.e. cv-*)")
	planCmd.Flags().Bool("speculative", false, "Perform a Speculative Plan (optional)")
	planCmd.Flags().Bool("destroy", false, "Perform a Destroy Plan (optional)")
	planCmd.Flags().StringSlice("env", []string{}, "Environment variables to write to the Workspace. Can be suplied multiple times. (optional, i.e. '--env='AWS_REGION=us-east1')")

	planCmd.MarkPersistentFlagRequired("workspaceName")

	rootCmd.AddCommand(planCmd)
}

func runPlan() error {
	// Validate flags
	hostname := *viperString("tfeHostname")
	orgName := *viperString("tfeOrganization")
	wsName := *viperString("workspaceName")
	dir := *viperString("directory")
	configId := *viperString("configurationId")
	isSpeculative := *viperBool("speculative")
	isDestroy := *viperBool("destroy")
	envs, err := viperStringSliceMap("env")
	if err != nil {
		logError(err, "failed to parse provided environment variables")
	}
	message := "Plan created by TFx"
	client, ctx := getClientContext()

	fmt.Println("Remote Terraform Plan, speculative plan:", color.GreenString(strconv.FormatBool(isSpeculative)),
		" destroy plan:", color.GreenString(strconv.FormatBool(isDestroy)))
	fmt.Println()

	// Read workspace
	fmt.Print("Reading Workspace ", color.GreenString(wsName), " for ID...")
	w, err := client.Workspaces.Read(ctx, orgName, wsName)
	if err != nil {
		logError(err, "failed to read workspace id")
	}
	fmt.Println(" Found:", color.BlueString(w.ID))

	// Update environment variables
	err = createOrUpdateEnvVariables(ctx, client, w.ID, envs)
	if err != nil {
		logError(err, "failed to write environment variables to workspace")
	}

	// Create config version
	cv, err := createOrReadConfigurationVersion(ctx, client, w.ID, configId, dir, isSpeculative)
	if err != nil {
		logError(err, "failed to create configuration version")
	}

	// Create run
	var r *tfe.Run
	r, err = client.Runs.Create(ctx, tfe.RunCreateOptions{
		IsDestroy:            tfe.Bool(isDestroy),
		Message:              tfe.String(message),
		ConfigurationVersion: cv,
		Workspace:            w,
		// TargetAddrs:          []string{"random_pet.two"},
	})
	if err != nil {
		logError(err, "failed to create run")
	}
	fmt.Println("Workspace Run Created, Run Id:", color.BlueString(r.ID), "Config Version:", color.BlueString(r.ConfigurationVersion.ID))
	fmt.Println("Navigate:", "https://"+hostname+"/app/"+orgName+"/workspaces/"+wsName+"/runs/"+r.ID)
	fmt.Println()

	// Get Run Logs
	err = getRunLogs(ctx, client, r.Plan.ID)
	if err != nil {
		logError(err, "failed to read run logs")
	}

	// Get Cost Estimation Logs (if any)
	err = getCostEstimationLogs(ctx, client, r)
	if err != nil {
		logError(err, "failed to read cost estimation logs")
	}

	// Get Policy Logs (if any)
	err = getPolicyLogs(ctx, client, r)
	if err != nil {
		logError(err, "failed to read policy logs")
	}

	fmt.Println("Run Complete:", r.ID)

	return nil
}
