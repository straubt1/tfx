/*
Copyright © 2021 Tom Straub <tstraub@hashicorp.com>

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
	"log"
	"os"

	"github.com/fatih/color"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	tfvCmd = &cobra.Command{
		Use:   "tfv",
		Short: "Terraform Versions",
		Long:  "Work with Terraform Versions of a TFx install.",
	}

	tfvListCmd = &cobra.Command{
		Use: "list",
		// Aliases: []string{"ls"},
		Short: "List Terraform Versions",
		Long:  "List Terraform Versions of a TFx install.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tfvList()
		},
		PreRun: bindPFlags,
	}

	tfvCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create Terraform Version",
		Long:  "Create Terraform Version for a TFx install.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tfvCreate()
		},
		PreRun: bindPFlags,
	}

	tfvShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show Terraform Version",
		Long:  "Show Terraform Version details for a TFx install.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tfvShow()
		},
		PreRun: bindPFlags,
	}

	tfvDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete Terraform Version",
		Long:  "Delete Terraform Version for a TFx install.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tfvDelete()
		},
		PreRun: bindPFlags,
	}
)

func init() {
	// All `tfx tfv` commands
	// tfvCmd.PersistentFlags().StringP("workspaceName", "w", "", "Workspace name")

	// `tfx tfv list`
	tfvListCmd.Flags().BoolP("enabled", "e", false, "Is enabled")

	// `tfx tfv create`
	tfvCreateCmd.Flags().String("version", "", "Version of Terraform (i.e. 0.15.0)")
	tfvCreateCmd.Flags().String("url", "", "Url of a hosted file containing Terraform (i.e. https://terraform.io...)")
	tfvCreateCmd.Flags().String("sha", "", "Sha checksum of the file at the url, must be 64 characters long")

	// `tfx tfv show`
	tfvShowCmd.Flags().StringP("versionId", "i", "", "Terraform Version Id (i.e. tool-*)")

	// `tfx tfv delete`
	tfvDeleteCmd.Flags().StringP("versionId", "i", "", "Terraform Version Id (i.e. tool-*)")

	rootCmd.AddCommand(tfvCmd)
	tfvCmd.AddCommand(tfvListCmd)
	tfvCmd.AddCommand(tfvCreateCmd)
	tfvCmd.AddCommand(tfvShowCmd)
	tfvCmd.AddCommand(tfvDeleteCmd)
}

func tfvList() error {
	// Validate flags
	// orgName := *viperString("tfeOrganization")
	// wsName := *viperString("workspaceName")
	// filterEnabled := viperBool("enabled")
	client, ctx := getClientContext()

	// Read all versions through pagination
	var err error
	var tfv *tfe.AdminTerraformVersionsList
	var tfvItems []*tfe.AdminTerraformVersion
	var tfvItemsFiltered []*tfe.AdminTerraformVersion
	pageNumber := 1
	for {
		tfv, err = client.Admin.TerraformVersions.List(ctx, tfe.AdminTerraformVersionsListOptions{
			ListOptions: tfe.ListOptions{
				PageSize:   100,
				PageNumber: pageNumber,
			},
		})
		if err != nil {
			log.Fatal(err)
		}

		tfvItems = append(tfvItems, tfv.Items...)
		if tfv.NextPage == 0 {
			break
		}
		pageNumber++
	}

	// filter
	for i, s := range tfvItems {
		// if s.Enabled {
		if s != nil {
			tfvItemsFiltered = append(tfvItemsFiltered, tfvItems[i])
		}
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Version", "Id", "Enabled", "Official", "Usage"})
	for _, i := range tfvItemsFiltered {
		t.AppendRow(table.Row{i.Version, i.ID, i.Enabled, i.Official, i.Usage})
	}
	t.SetStyle(table.StyleRounded)
	t.Render()
	// fmt.Println("Terraform Versions Found: ", color.BlueString(string(len(tfvItemsFiltered))))
	return nil
}

func tfvCreate() error {
	// Validate flags
	version := *viperString("version")
	url := *viperString("url")
	sha := *viperString("sha")
	client, ctx := getClientContext()

	// Create Terraform Version
	fmt.Print("Creating Terraform Version ...")
	tfv, err := client.Admin.TerraformVersions.Create(ctx, tfe.AdminTerraformVersionCreateOptions{
		Version:  tfe.String(version),
		URL:      tfe.String(url),
		Sha:      tfe.String(sha),
		Official: tfe.Bool(false),
		Enabled:  tfe.Bool(true),
		Beta:     tfe.Bool(false),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(" ID:", tfv.ID)

	return nil
}

func tfvShow() error {
	// Validate flags
	vId := *viperString("versionId")

	client, ctx := getClientContext()

	// Read Terraform Version
	fmt.Print("Reading Terraform Version with ID ", color.GreenString(vId), "...")
	tfv, err := client.Admin.TerraformVersions.Read(ctx, vId)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(" ID Found")
	fmt.Println(color.BlueString("Version:   "), tfv.Version)
	fmt.Println(color.BlueString("ID:        "), tfv.ID)
	fmt.Println(color.BlueString("URL:       "), tfv.URL)
	fmt.Println(color.BlueString("Sha:       "), tfv.Sha)
	fmt.Println(color.BlueString("Enabled:   "), tfv.Enabled)
	fmt.Println(color.BlueString("Official:  "), tfv.Official)
	fmt.Println(color.BlueString("Beta:      "), tfv.Beta)

	return nil
}

func tfvDelete() error {
	// Validate flags
	vId := *viperString("versionId")
	client, ctx := getClientContext()

	// Delete Terraform Version
	fmt.Print("Deleting Terraform Version ID ", color.GreenString(vId), "...")
	err := client.Admin.TerraformVersions.Delete(ctx, vId)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(" tfv Deleted")

	return nil
}
