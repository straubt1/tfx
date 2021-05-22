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
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

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

	tfvDisableCmd = &cobra.Command{
		Use:   "disable",
		Short: "Disable Terraform Version",
		Long:  "Disable Terraform Version for a TFx install.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tfvDisable()
		},
		PreRun: bindPFlags,
	}
)

func init() {
	// `tfx tfv list`
	// tfvListCmd.Flags().StringP("enabled", "e", "", "Filter on enabled Terraform Versions, if set must be ['true', 'false']")
	// tfvListCmd.Flags().StringP("official", "o", "", "Filter on official Terraform Versions, if set must be ['true', 'false']")
	// tfvListCmd.Flags().StringP("beta", "b", "", "Filter on beta Terraform Versions, if set must be ['true', 'false']")

	// `tfx tfv create`
	tfvCreateCmd.Flags().StringP("version", "v", "", "Version of Terraform (i.e. 0.15.0)")
	tfvCreateCmd.Flags().StringP("url", "u", "", "Url of a hosted file containing Terraform (i.e. https://terraform.io...)")
	tfvCreateCmd.Flags().StringP("sha", "s", "", "Sha checksum of the file at the url, must be 64 characters long")
	tfvCreateCmd.MarkFlagRequired("version")
	tfvCreateCmd.MarkFlagRequired("url")
	tfvCreateCmd.MarkFlagRequired("sha")

	// `tfx tfv show`
	tfvShowCmd.Flags().StringP("versionId", "i", "", "Terraform Version Id (i.e. tool-*)")
	tfvShowCmd.MarkFlagRequired("versionId")

	// `tfx tfv delete`
	tfvDeleteCmd.Flags().StringP("versionId", "i", "", "Terraform Version Id (i.e. tool-*)")
	tfvDeleteCmd.MarkFlagRequired("versionId")

	// `tfx tfv disable`
	tfvDisableCmd.Flags().BoolP("all", "a", false, "Disable All")
	tfvDisableCmd.Flags().StringSlice("versions", []string{}, "Versions to disable, comma seperated (i.e. 0.11.0,0.11.1)")

	rootCmd.AddCommand(tfvCmd)
	tfvCmd.AddCommand(tfvListCmd)
	tfvCmd.AddCommand(tfvCreateCmd)
	tfvCmd.AddCommand(tfvShowCmd)
	tfvCmd.AddCommand(tfvDeleteCmd)
	tfvCmd.AddCommand(tfvDisableCmd)
}

func tfvList() error {
	// Validate flags
	// orgName := *viperString("tfeOrganization")
	// wsName := *viperString("workspaceName")
	// filterEnabled := *viperString("enabled")
	client, ctx := getClientContext()

	// Read all versions through pagination
	// var err error
	tfvItems, err := getAllTerraformVersions(ctx, client)
	if err != nil {
		logError(err, "failed to read all terraform versions")
	}
	var tfvItemsFiltered []*tfe.AdminTerraformVersion

	// TODO://implement filter
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
	fmt.Println("Terraform Versions Found: ", color.BlueString(strconv.Itoa(len(tfvItemsFiltered))))
	return nil
}

func tfvCreate() error {
	url := *viperString("url")
	sha := *viperString("sha")
	if len(sha) != 64 {
		logError(errors.New(""), "sha must be 64 characters long")
	}
	// Attempt to prevent a non semantic version from being created
	version, err := viperSemanticVersionString("version")
	if err != nil {
		logError(err, "failed to parse semantic version")
	}
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
	fmt.Println(" Found")
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

	// TODO: force official provider to false, then delete.
	// Need to verify an update wont bring these back
	// client.Admin.TerraformVersions.Update(ctx, vId, tfe.AdminTerraformVersionUpdateOptions{
	// 	Official: tfe.Bool(false),
	// })
	// Delete Terraform Version
	fmt.Print("Deleting Terraform Version ID ", color.GreenString(vId), "...")
	err := client.Admin.TerraformVersions.Delete(ctx, vId)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(" Deleted")

	return nil
}

func tfvDisable() error {
	// Validate flags
	all := *viperBool("all")
	versions := viperStringSlice("versions")
	if len(versions) == 0 && !all {
		logError(errors.New(""), "No Versions provided")
	}

	if all {
		fmt.Println("Disabling All Terraform Versions that are not in use")
	}
	client, ctx := getClientContext()

	allTFV, err := getAllTerraformVersions(ctx, client)
	if err != nil {
		logError(err, "failed to read all terraform versions")
	}

	for _, s := range allTFV {
		// if not all, then see if version is passed in
		if !all {
			found := false
			for _, v := range versions {
				if s.Version == v {
					found = true
				}
			}
			// not found, skip version
			if !found {
				continue
			}
		}

		currentTFV, err := client.Admin.TerraformVersions.Read(ctx, s.ID)
		if err != nil { //this should never happen
			logWarning(err, "failed to read terraform version: "+s.Version)
		}
		if !currentTFV.Enabled { //already disabled, skip
			fmt.Println("Terraform Version already disabled: ", color.BlueString(s.Version))
			continue
		}
		if currentTFV.Usage > 0 { //can not disable a version with usage
			fmt.Println(color.RedString("Terraform Version in use: "), color.BlueString(s.Version))
			continue
		}

		_, err = client.Admin.TerraformVersions.Update(ctx, s.ID, tfe.AdminTerraformVersionUpdateOptions{
			Enabled: tfe.Bool(false),
		})
		if err == nil {
			fmt.Println("Terraform Versions disabled: ", color.BlueString(s.Version))
		} else {
			fmt.Println(color.RedString("Unable to update Terraform Version: "), color.BlueString(s.Version))
		}
	}
	_ = all

	return nil
}
