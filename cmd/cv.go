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
	"os"

	"github.com/fatih/color"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	cvCmd = &cobra.Command{
		Use:   "cv",
		Short: "Configuration Versions",
		Long:  "Work with Configration Versions of a TFx Workspace.",
	}

	cvListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Configuration Versions",
		Long:  "List Configuration Versions of a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cvList()
		},
		PreRun: bindPFlags,
	}

	cvCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create Configuration Version",
		Long:  "Create Configuration Version for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cvCreate()
		},
		PreRun: bindPFlags,
	}

	cvShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show Configuration Version",
		Long:  "Show Configuration Version details for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cvShow()
		},
		PreRun: bindPFlags,
	}
)

func init() {
	// All `tfx cv` commands
	cvCmd.PersistentFlags().StringP("workspaceName", "w", "", "Workspace name")
	cvCmd.MarkPersistentFlagRequired("workspaceName")

	// `tfx cv create`
	cvCreateCmd.Flags().StringP("directory", "d", "./", "Directory of Terraform (optional, defaults to current directory)")
	cvCreateCmd.Flags().Bool("speculative", false, "Perform a Speculative Plan (optional, defaults to false)")

	// `tfx cv show`
	cvShowCmd.Flags().StringP("configurationId", "i", "", "Configuration Version Id (i.e. cv-*)")
	cvShowCmd.MarkPersistentFlagRequired("configurationId")

	rootCmd.AddCommand(cvCmd)
	cvCmd.AddCommand(cvListCmd)
	cvCmd.AddCommand(cvCreateCmd)
	cvCmd.AddCommand(cvShowCmd)
}

func cvList() error {
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

	// Get all config versions
	cv, err := client.ConfigurationVersions.List(ctx, w.ID, tfe.ConfigurationVersionListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 10,
		},
	})
	if err != nil {
		logError(err, "failed to list configuration versions")
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Id", "Speculative", "Status"})
	for _, i := range cv.Items {
		t.AppendRow(table.Row{i.ID, i.Speculative, i.Status})
	}
	t.SetStyle(table.StyleRounded)
	t.Render()

	return nil
}

func cvCreate() error {
	// Validate flags
	orgName := *viperString("tfeOrganization")
	wsName := *viperString("workspaceName")
	dir := *viperString("directory")
	isSpeculative := *viperBool("speculative")
	client, ctx := getClientContext()

	// Read workspace
	fmt.Print("Reading Workspace ", color.GreenString(wsName), " for ID...")
	w, err := client.Workspaces.Read(ctx, orgName, wsName)
	if err != nil {
		logError(err, "failed to read workspace id")
	}
	fmt.Println(" Found:", w.ID)

	// Create Config Version
	fmt.Print("Creating Configuration Version ...")
	cv, err := client.ConfigurationVersions.Create(ctx, w.ID, tfe.ConfigurationVersionCreateOptions{
		AutoQueueRuns: tfe.Bool(false),
		Speculative:   tfe.Bool(isSpeculative),
	})
	if err != nil {
		logError(err, "failed to create configuration version")
	}
	fmt.Println(" ID:", cv.ID)

	// Upload to Config Version
	fmt.Print("Uploading code from directory ", color.GreenString(dir), " ...")
	err = client.ConfigurationVersions.Upload(ctx, cv.UploadURL, dir)
	if err != nil {
		logError(err, "failed to upload code")
	}
	fmt.Println(" Done")

	return nil
}

func cvShow() error {
	// Validate flags
	configId := *viperString("configurationId")
	client, ctx := getClientContext()

	// Read Config Version
	fmt.Print("Reading Configuration for ID ", color.GreenString(configId), " ...")
	cv, err := client.ConfigurationVersions.Read(ctx, configId)
	if err != nil {
		logError(err, "failed to read configuration version")
	}
	fmt.Println(" CV Found")
	fmt.Println(color.BlueString("ID:          "), cv.ID)
	fmt.Println(color.BlueString("Status:      "), cv.Status)
	fmt.Println(color.BlueString("Speculative: "), cv.Speculative)
	fmt.Println(color.BlueString("Source:      "), cv.Source)
	fmt.Println(color.BlueString("Error:       "), cv.ErrorMessage)

	return nil
}
