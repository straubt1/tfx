// Copyright © 2021 Tom Straub <github.com/straubt1>

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Workspace Runs",
		Long:  "Work with Runs of a TFx Workspace.",
	}

	runListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Runs",
		Long:  "List Runs of a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList()
		},
		PreRun: bindPFlags,
	}

	runCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create Run",
		Long:  "Create Run for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate()
		},
		PreRun: bindPFlags,
	}

	runShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show Run",
		Long:  "Show Run details for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runShow()
		},
		PreRun: bindPFlags,
	}
)

func init() {
	// All `tfx run` commands
	runCmd.PersistentFlags().StringP("workspaceName", "w", "", "Workspace name")
	runCmd.MarkPersistentFlagRequired("workspaceName")

	// `tfx run create`
	runCreateCmd.Flags().StringP("directory", "d", "./", "Directory of Terraform (defaults to current directory)")
	runCreateCmd.Flags().StringP("message", "m", "", "Run Message (optional)")

	// `tfx run show`
	runShowCmd.Flags().StringP("runId", "i", "", "Run Id (i.e. run-*)")
	runShowCmd.MarkFlagRequired("runId")

	rootCmd.AddCommand(runCmd)
	runCmd.AddCommand(runListCmd)
	runCmd.AddCommand(runCreateCmd)
	runCmd.AddCommand(runShowCmd)
}

func runList() error {
	// Validate flags
	orgName := *viperString("tfeOrganization")
	wsName := *viperString("workspaceName")
	client, ctx := getClientContext()

	// Read workspace
	fmt.Print("Reading Workspace ID for Name: ", color.GreenString(wsName), " ...")
	w, err := client.Workspaces.Read(ctx, orgName, wsName)
	if err != nil {
		logError(err, "failed to read workspace id")
	}
	fmt.Println(" Found:", color.BlueString(w.ID))

	// Get all config versions and show the current config
	run, err := client.Runs.List(ctx, w.ID, &tfe.RunListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
		Include: []tfe.RunIncludeOpt{"workspace"}, // To get TF Version
	})
	if err != nil {
		logError(err, "failed to list runs")
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Id", "Status", "Terraform Version", "Created"})
	for _, i := range run.Items {
		t.AppendRow(table.Row{i.ID, i.Status, i.Workspace.TerraformVersion, i.CreatedAt.String()})
	}
	t.SetStyle(table.StyleRounded)
	t.Render()

	return nil
}

func runCreate() error {
	// Validate flags
	orgName := *viperString("tfeOrganization")
	wsName := *viperString("workspaceName")
	message := *viperString("message")
	client, ctx := getClientContext()

	// Read workspace
	fmt.Print("Reading Workspace ", color.GreenString(wsName), " for ID...")
	w, err := client.Workspaces.Read(ctx, orgName, wsName)
	if err != nil {
		logError(err, "failed to read workspace id")
	}
	fmt.Println(" Found:", color.BlueString(w.ID))

	// Create Config Version
	fmt.Print("Creating Run ...")
	run, err := client.Runs.Create(ctx, tfe.RunCreateOptions{
		Workspace: w,
		IsDestroy: tfe.Bool(false),
		Message:   tfe.String(message),
		// ConfigurationVersion: {}
		// without a CV the run will kick off automatically with the last know CV
		// this is likely not ideal
	})
	if err != nil {
		logError(err, "failed to create a run")
	}
	fmt.Println(" ID:", run.ID)

	return nil
}

func runShow() error {
	// Validate flags
	runID := *viperString("runId")
	client, ctx := getClientContext()

	// Read Run
	fmt.Print("Reading Run for ID ", color.GreenString(runID), "...")
	run, err := client.Runs.ReadWithOptions(ctx, runID, &tfe.RunReadOptions{
		Include: []tfe.RunIncludeOpt{"workspace"}, // To get TF Version
	})
	if err != nil {
		logError(err, "failed to read run")
	}
	fmt.Println(" run Found")
	fmt.Println(color.BlueString("ID:          "), run.ID)
	fmt.Println(color.BlueString("Status:      "), run.Status)
	fmt.Println(color.BlueString("Message:     "), run.Message)
	fmt.Println(color.BlueString("TF Version:  "), run.Workspace.TerraformVersion)
	fmt.Println(color.BlueString("Created:     "), run.CreatedAt.String())

	return nil
}
