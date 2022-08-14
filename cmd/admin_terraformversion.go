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
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	tfvCmd = &cobra.Command{
		Use:   "terraform-version",
		Aliases: []string{"tfv"},
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

	tfvCreateOfficialCmd = &cobra.Command{
		Use:   "official",
		Short: "Create Terraform Version Official",
		Long:  "Create Terraform Version Official for a TFx install.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tfvCreateOfficial()
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

	tfvDisableAllCmd = &cobra.Command{
		Use:   "all",
		Short: "Disable All Terraform Versions",
		Long:  "Disable All Terraform Versions for a TFx install.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tfvDisableAll()
		},
		PreRun: bindPFlags,
	}

	tfvEnableCmd = &cobra.Command{
		Use:   "enable",
		Short: "Enable Terraform Version",
		Long:  "Enable Terraform Version for a TFx install.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tfvEnable()
		},
		PreRun: bindPFlags,
	}

	tfvEnableAllCmd = &cobra.Command{
		Use:   "all",
		Short: "Disable All Terraform Versions",
		Long:  "Disable All Terraform Versions for a TFx install.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tfvEnableAll()
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

	// `tfx tfv create official`
	tfvCreateOfficialCmd.Flags().StringP("version", "v", "", "Version of Terraform (i.e. 0.15.0)")

	// `tfx tfv show`
	tfvShowCmd.Flags().StringP("versionId", "i", "", "Terraform Version Id (i.e. tool-*)")
	tfvShowCmd.Flags().StringP("version", "v", "", "Terraform Version (i.e. 0.15.0)")

	// `tfx tfv delete`
	tfvDeleteCmd.Flags().StringP("versionId", "i", "", "Terraform Version Id (i.e. tool-*)")
	tfvDeleteCmd.Flags().StringP("version", "v", "", "Terraform Version (i.e. 0.15.0)")

	// `tfx tfv disable`
	tfvDisableCmd.Flags().StringSliceP("versions", "v", []string{}, "Versions to disable, can be comma separated (i.e. 0.11.0,0.11.1)")
	tfvDisableCmd.MarkFlagRequired("versions")

	// `tfx tfv enable`
	tfvEnableCmd.Flags().StringSliceP("versions", "v", []string{}, "Versions to enable, can be comma separated (i.e. 0.11.0,0.11.1)")
	tfvEnableCmd.MarkFlagRequired("versions")

	adminCmd.AddCommand(tfvCmd)
	tfvCmd.AddCommand(tfvListCmd)
	tfvCmd.AddCommand(tfvCreateCmd)
	tfvCreateCmd.AddCommand(tfvCreateOfficialCmd)
	tfvCmd.AddCommand(tfvShowCmd)
	tfvCmd.AddCommand(tfvDeleteCmd)
	tfvCmd.AddCommand(tfvDisableCmd)
	tfvDisableCmd.AddCommand(tfvDisableAllCmd)
	tfvCmd.AddCommand(tfvEnableCmd)
	tfvEnableCmd.AddCommand(tfvEnableAllCmd)
}

func tfvList() error {
	client, ctx := getClientContext()

	// Read all versions through pagination
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

func tfvCreateOfficial() error {
	version, err := viperSemanticVersionString("version")
	if err != nil {
		logError(err, "failed to parse semantic version")
	}

	fmt.Print("Looking for official Terraform Version: ", color.GreenString(version), " ...")
	url := fmt.Sprintf(
		"https://releases.hashicorp.com/terraform/%s/terraform_%s_linux_amd64.zip",
		version,
		version,
	)
	// read checksum file
	urlSha := fmt.Sprintf(
		"https://releases.hashicorp.com/terraform/%s/terraform_%s_SHA256SUMS",
		version,
		version,
	)
	clientChecksum := &http.Client{}
	req, err := http.NewRequest("GET", urlSha, nil)
	if err != nil {
		logError(err, "failed to create checksum request")
	}

	resp, err := clientChecksum.Do(req)
	if err != nil || resp.StatusCode != 200 {
		logError(err, "failed to request checksum")
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logError(err, "failed to read checksum")
	}
	// split by new line
	var sha string
	lines := strings.Split(string(b), "\n")
	for _, l := range lines {
		// looks for linux version
		if strings.Contains(l, "linux_amd64") {
			// only grab the checksum
			innerLines := strings.Split(l, " ")
			sha = innerLines[0]
			break
		}
	}
	fmt.Println("Found")

	client, ctx := getClientContext()

	// Create Terraform Version
	fmt.Print("Creating Terraform Version ...")
	tfv, err := client.Admin.TerraformVersions.Create(ctx, tfe.AdminTerraformVersionCreateOptions{
		Version:  tfe.String(version),
		URL:      tfe.String(url),
		Sha:      tfe.String(sha),
		Official: tfe.Bool(true),
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
	vID := *viperString("versionId")
	v := *viperString("version")
	if vID == "" && v == "" {
		logError(errors.New(""), "version or version id must be supplied")
	} else if vID != "" && v != "" {
		logError(errors.New(""), "only one can be supplied [version or version id]")
	}
	client, ctx := getClientContext()

	var tfv *tfe.AdminTerraformVersion
	var err error
	if vID != "" {
		// Read Terraform Version
		fmt.Print("Reading Terraform Version with ID ", color.GreenString(vID), "...")
		tfv, err = client.Admin.TerraformVersions.Read(ctx, vID)
		if err != nil {
			logError(err, "failed to find version id")
		}
		fmt.Println(" Found")
	} else {
		fmt.Print("Reading Terraform Version  ", color.GreenString(v), "...")
		tfv, err = getTerraformVersion(ctx, client, v)
		if err != nil {
			logError(err, "failed to find version")
		}
		fmt.Println(" Found")
	}
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
	vID := *viperString("versionId")
	v := *viperString("version")
	if vID == "" && v == "" {
		logError(errors.New(""), "version or version id must be supplied")
	} else if vID != "" && v != "" {
		logError(errors.New(""), "only one can be supplied [version or version id]")
	}
	client, ctx := getClientContext()

	var tfv *tfe.AdminTerraformVersion
	var err error
	if vID != "" {
		// Read Terraform Version
		fmt.Print("Reading Terraform Version with ID ", color.GreenString(vID), "...")
		tfv, err = client.Admin.TerraformVersions.Read(ctx, vID)
		if err != nil {
			logError(err, "failed to find version id")
		}
		fmt.Println(" Found")
	} else {
		fmt.Print("Reading Terraform Version  ", color.GreenString(v), "...")
		tfv, err = getTerraformVersion(ctx, client, v)
		if err != nil {
			logError(err, "failed to find version")
		}
		fmt.Println(" Found")
	}

	// TODO: Need to verify an update wont bring these back
	if tfv.Official {
		fmt.Println("Forcing Terraform Version to be unofficial")
		tfv, err = client.Admin.TerraformVersions.Update(ctx, tfv.ID, tfe.AdminTerraformVersionUpdateOptions{
			Official: tfe.Bool(false),
		})
		if err != nil {
			logError(err, "failed to set version to official = false")
		}
	}

	// Delete Terraform Version
	fmt.Print("Deleting Terraform Version ID ", color.GreenString(tfv.ID), "...")
	err = client.Admin.TerraformVersions.Delete(ctx, tfv.ID)
	if err != nil {
		logError(err, "failed to delete version")
	}
	fmt.Println(" Deleted")

	return nil
}

func tfvDisable() error {
	versions := viperStringSlice("versions")

	setTfvEnabledFlag(versions, false)

	return nil
}

func tfvDisableAll() error {
	fmt.Println("Disabling All Terraform Versions that are not in use")
	setTfvEnabledFlagAll(false)

	return nil
}

func tfvEnable() error {
	versions := viperStringSlice("versions")

	setTfvEnabledFlag(versions, true)

	return nil
}

func tfvEnableAll() error {
	fmt.Println("Enabling All Terraform Versions")
	setTfvEnabledFlagAll(true)

	return nil
}

func setTfvEnabledFlag(versions []string, enabled bool) error {
	client, ctx := getClientContext()
	allTFV, err := getAllTerraformVersions(ctx, client)
	if err != nil {
		return errors.New("failed to read all terraform versions")
	}

	// loop on passed in versions
	for _, v := range versions {
		var foundVersion *tfe.AdminTerraformVersion
		for _, s := range allTFV {
			if s.Version == v {
				foundVersion = s
			}
		}
		// not found, skip version
		if foundVersion == nil {
			logWarning(err, "failed to find terraform version: "+v)
			continue
		}

		if foundVersion.Enabled == enabled { //already set, skip
			fmt.Println("Terraform Version", color.BlueString(foundVersion.Version), "is already", color.GreenString(strconv.FormatBool(enabled)))
			continue
		}
		if !enabled && foundVersion.Usage > 0 { //can not disable a version with usage
			fmt.Println(color.RedString("Terraform Version in use"), color.BlueString(foundVersion.Version))
			continue
		}

		_, err = client.Admin.TerraformVersions.Update(ctx, foundVersion.ID, tfe.AdminTerraformVersionUpdateOptions{
			Enabled: tfe.Bool(enabled),
		})
		if err == nil {
			fmt.Println("Terraform Versions", color.BlueString(foundVersion.Version), "is now set to", color.GreenString(strconv.FormatBool(enabled)))
		} else {
			fmt.Println(color.RedString("Unable to update Terraform Version "), color.BlueString(foundVersion.Version))
		}
	}

	return nil
}

func setTfvEnabledFlagAll(enabled bool) error {
	client, ctx := getClientContext()
	allTFV, err := getAllTerraformVersions(ctx, client)
	if err != nil {
		return errors.New("failed to read all terraform versions")
	}

	// loop on all versions
	for _, v := range allTFV {
		if v.Enabled == enabled { //already set, skip
			fmt.Println("Terraform Version", color.BlueString(v.Version), "is already", color.GreenString(strconv.FormatBool(enabled)))
			continue
		}
		if !enabled && v.Usage > 0 { //can not disable a version with usage
			fmt.Println(color.RedString("Terraform Version in use"), color.BlueString(v.Version))
			continue
		}

		_, err = client.Admin.TerraformVersions.Update(ctx, v.ID, tfe.AdminTerraformVersionUpdateOptions{
			Enabled: tfe.Bool(enabled),
		})
		if err == nil {
			fmt.Println("Terraform Versions", color.BlueString(v.Version), "is now set to", color.GreenString(strconv.FormatBool(enabled)))
		} else {
			fmt.Println(color.RedString("Unable to update Terraform Version "), color.BlueString(v.Version))
		}
	}

	return nil
}
