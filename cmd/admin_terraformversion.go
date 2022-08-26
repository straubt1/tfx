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
	"io"
	"net/http"
	"strconv"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	//  `tfx admin terraform-version ` commands
	tfvCmd = &cobra.Command{
		Use:     "terraform-version",
		Aliases: []string{"tfv"},
		Short:   "Terraform Version Commands",
		Long:    "Work with Terraform Versions of a TFE Installation.",
	}

	// `tfx admin terraform-version list`
	tfvListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Terraform Versions",
		Long:  "List Terraform Versions of a TFE Installation.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tfvList(
				getTfxClientContext(),
				*viperString("search"))
		},
	}

	tfvCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create Terraform Version",
		Long:  "Create Terraform Version for a TFE Installation.",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := viperSemanticVersionString("version")
			if err != nil {
				return errors.New("failed to parse semantic version")
			}
			_, err = viperShaString("sha")
			if err != nil {
				return errors.New("failed to parse semantic version")
			}

			return tfvCreate(
				getTfxClientContext(),
				*viperString("version"),
				*viperString("url"),
				*viperString("sha"),
				*viperBool("official"),
				!*viperBool("disable"),
				*viperBool("beta"))
		},
	}

	tfvCreateOfficialCmd = &cobra.Command{
		Use:   "official",
		Short: "Create Terraform Version Official",
		Long:  "Create Terraform Version Official for a TFE Installation.",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := viperSemanticVersionString("version")
			if err != nil {
				return errors.New("failed to parse semantic version")
			}

			return tfvCreateOfficial(
				getTfxClientContext(),
				*viperString("version"),
				!*viperBool("disable"),
				*viperBool("beta"))
		},
	}

	tfvShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show Terraform Version",
		Long:  "Show Terraform Version details for a TFE Installation.",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := viperSemanticVersionString("version")
			if err != nil {
				return errors.New("failed to parse semantic version")
			}

			return tfvShow(
				getTfxClientContext(),
				*viperString("version"))
		},
	}

	tfvDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete Terraform Version",
		Long:  "Delete Terraform Version for a TFE Installation.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tfvDelete(
				getTfxClientContext(),
				*viperString("version"))
		},
	}

	tfvDisableCmd = &cobra.Command{
		Use:   "disable",
		Short: "Disable Terraform Version",
		Long:  "Disable Terraform Version for a TFE Installation.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tfvDisable(
				getTfxClientContext(),
				viperStringSlice("versions"))
		},
	}

	tfvDisableAllCmd = &cobra.Command{
		Use:   "all",
		Short: "Disable All Terraform Versions",
		Long:  "Disable All Terraform Versions for a TFE Installation.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tfvDisableAll(
				getTfxClientContext())
		},
	}

	tfvEnableCmd = &cobra.Command{
		Use:   "enable",
		Short: "Enable Terraform Version",
		Long:  "Enable Terraform Version for a TFE Installation.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tfvEnable(
				getTfxClientContext(),
				viperStringSlice("versions"))
		},
	}

	tfvEnableAllCmd = &cobra.Command{
		Use:   "all",
		Short: "Disable All Terraform Versions",
		Long:  "Disable All Terraform Versions for a TFE Installation.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tfvEnableAll(
				getTfxClientContext())
		},
	}
)

func init() {
	// `tfx admin terraform-version list` flags
	tfvListCmd.Flags().StringP("search", "s", "", "Search string for partial version string (optional).")
	// tfvListCmd.Flags().StringP("enabled", "e", "", "Filter on enabled Terraform Versions, if set must be ['true', 'false']")
	// tfvListCmd.Flags().StringP("official", "o", "", "Filter on official Terraform Versions, if set must be ['true', 'false']")
	// tfvListCmd.Flags().StringP("beta", "b", "", "Filter on beta Terraform Versions, if set must be ['true', 'false']")

	// `tfx tfv show`
	tfvShowCmd.Flags().StringP("version", "v", "", "Terraform Version (i.e. 0.15.0)")
	tfvShowCmd.MarkFlagRequired("version")

	// `tfx tfv create`
	tfvCreateCmd.Flags().StringP("version", "v", "", "Version of Terraform (i.e. 0.15.0)")
	tfvCreateCmd.Flags().StringP("url", "u", "", "Url of a hosted file containing Terraform (i.e. https://terraform.io...)")
	tfvCreateCmd.Flags().StringP("sha", "s", "", "Sha checksum of the file at the url, must be 64 characters long")
	tfvCreateCmd.Flags().BoolP("official", "", false, "Terraform Version is official (optional)")
	tfvCreateCmd.Flags().BoolP("disable", "", false, "Created Terraform Version will be disabled (optional)")
	tfvCreateCmd.Flags().BoolP("beta", "", false, "Terraform Version is beta (optional)")
	tfvCreateCmd.MarkFlagRequired("version")
	tfvCreateCmd.MarkFlagRequired("url")
	tfvCreateCmd.MarkFlagRequired("sha")

	// `tfx tfv create official`
	tfvCreateOfficialCmd.Flags().StringP("version", "v", "", "Version of Terraform (i.e. 0.15.0)")
	tfvCreateOfficialCmd.Flags().BoolP("disable", "", false, "Created Terraform Version will be disabled (optional)")
	tfvCreateOfficialCmd.Flags().BoolP("beta", "", false, "Terraform Version is beta (optional)")
	tfvCreateOfficialCmd.MarkFlagRequired("version")

	// `tfx tfv delete`
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

func adminTFVListAll(c TfxClientContext, filter string, search string) ([]*tfe.AdminTerraformVersion, error) {
	allItems := []*tfe.AdminTerraformVersion{}
	opts := tfe.AdminTerraformVersionsListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 1, PageSize: 100},
		Filter:      filter,
		Search:      search,
	}
	for {
		items, err := c.Client.Admin.TerraformVersions.List(c.Context, &opts)
		if err != nil {
			return nil, err
		}

		allItems = append(allItems, items.Items...)
		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage
	}

	return allItems, nil
}

func adminTFVGetVersion(c TfxClientContext, version string) (*tfe.AdminTerraformVersion, error) {
	// Use the list all function, filter will return based on an exact version match
	items, err := adminTFVListAll(c, version, "")
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, errors.New("terraform version not found")
	} else if len(items) > 1 {
		// unlikely to ever hit this, but just in case
		return nil, errors.New("too many terraform versions found")
	}

	return items[0], nil
}

func tfvList(c TfxClientContext, search string) error {
	o.AddMessageUserProvided("List Terraform Versions for TFE", "")
	items, err := adminTFVListAll(c, "", search)
	if err != nil {
		return errors.Wrap(err, "failed to list terraform versions")
	}

	o.AddTableHeader("Version", "Id", "Enabled", "Official", "Usage", "Deprecated")
	for _, i := range items {
		o.AddTableRows(i.Version, i.ID, i.Enabled, i.Official, i.Usage, i.Deprecated)
	}

	return nil
}

func tfvCreate(c TfxClientContext, version string, url string,
	sha string, isOfficial bool, isEnabled bool, isBeta bool) error {
	o.AddMessageUserProvided("Create Terraform Version:", version)
	tfv, err := c.Client.Admin.TerraformVersions.Create(c.Context, tfe.AdminTerraformVersionCreateOptions{
		Version:  tfe.String(version),
		URL:      tfe.String(url),
		Sha:      tfe.String(sha),
		Official: tfe.Bool(isOfficial),
		Enabled:  tfe.Bool(isEnabled),
		Beta:     tfe.Bool(isBeta),
	})
	if err != nil {
		return errors.Wrap(err, "unable to create terraform version")
	}

	o.AddDeferredMessageRead("Version", tfv.Version)
	o.AddDeferredMessageRead("ID", tfv.ID)
	o.AddDeferredMessageRead("URL", tfv.URL)
	o.AddDeferredMessageRead("Sha", tfv.Sha)
	o.AddDeferredMessageRead("Enabled", tfv.Enabled)
	o.AddDeferredMessageRead("Beta", tfv.Beta)

	return nil
}

func tfvCreateOfficial(c TfxClientContext, version string, isEnabled bool, isBeta bool) error {
	url := fmt.Sprintf(
		"https://releases.hashicorp.com/terraform/%s/terraform_%s_linux_amd64.zip",
		version,
		version,
	)
	urlSha := fmt.Sprintf(
		"https://releases.hashicorp.com/terraform/%s/terraform_%s_SHA256SUMS",
		version,
		version,
	)

	o.AddMessageUserProvided("Searching for official Terraform Version:", version)
	clientChecksum := &http.Client{}
	req, err := http.NewRequest("GET", urlSha, nil)
	if err != nil {
		return errors.Wrap(err, "failed find official terraform version")
	}
	resp, err := clientChecksum.Do(req)
	if err != nil || resp.StatusCode != 200 {
		// if this fails, assume the version does not exist
		return errors.New("failed find official terraform version")
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read checksum")
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
	o.AddMessageUserProvided("Terraform Version SHASUM:", sha)
	err = tfvCreate(c, version, url, sha, true, isEnabled, isBeta)
	if err != nil {
		return errors.Wrap(err, "unable to create terraform version")
	}

	return nil
}

func tfvShow(c TfxClientContext, version string) error {
	o.AddMessageUserProvided("Show Terraform Version:", version)
	tfv, err := adminTFVGetVersion(c, version)
	if err != nil {
		return errors.Wrap(err, "failed to find terraform version")
	}

	o.AddDeferredMessageRead("Version", tfv.Version)
	o.AddDeferredMessageRead("ID", tfv.ID)
	o.AddDeferredMessageRead("URL", tfv.URL)
	o.AddDeferredMessageRead("Sha", tfv.Sha)
	o.AddDeferredMessageRead("Enabled", tfv.Enabled)
	o.AddDeferredMessageRead("Beta", tfv.Beta)

	return nil
}

func tfvDelete(c TfxClientContext, version string) error {
	o.AddMessageUserProvided("Delete Terraform Version:", version)
	tfv, err := adminTFVGetVersion(c, version)
	if err != nil {
		return errors.Wrap(err, "failed to find terraform version")
	}

	// TODO: Need to verify an update wont bring these back
	if tfv.Official {
		o.AddMessageUserProvided("Forcing Terraform Version to be unofficial", "")
		tfv, err = c.Client.Admin.TerraformVersions.Update(c.Context, tfv.ID, tfe.AdminTerraformVersionUpdateOptions{
			Official: tfe.Bool(false),
		})
		if err != nil {
			return errors.Wrap(err, "failed to set version to official to false")
		}
	}

	// Delete Terraform Version
	err = c.Client.Admin.TerraformVersions.Delete(c.Context, tfv.ID)
	if err != nil {
		return errors.Wrap(err, "failed to delete version")
	}

	o.AddMessageUserProvided("Variable Deleted:", version)
	o.AddDeferredMessageRead("Status", "Success")

	return nil
}

func tfvDisable(c TfxClientContext, versions []string) error {
	o.AddMessageUserProvided("Disable Terraform Versions:", versions)
	return tfvUpdateVersions(c, versions, false)
}

func tfvDisableAll(c TfxClientContext) error {
	o.AddMessageUserProvided("Disable All Terraform Versions", "")
	items, err := adminTFVListAll(c, "", "")
	if err != nil {
		return errors.Wrap(err, "failed to list terraform versions")
	}
	versions := []string{}
	for _, v := range items {
		versions = append(versions, v.Version)
	}

	return tfvUpdateVersions(c, versions, false)
}

func tfvEnable(c TfxClientContext, versions []string) error {
	o.AddMessageUserProvided("Enable Terraform Versions:", versions)
	return tfvUpdateVersions(c, versions, true)
}

func tfvEnableAll(c TfxClientContext) error {
	o.AddMessageUserProvided("Enable All Terraform Versions", "")
	items, err := adminTFVListAll(c, "", "")
	if err != nil {
		return errors.Wrap(err, "failed to list terraform versions")
	}
	versions := []string{}
	for _, v := range items {
		versions = append(versions, v.Version)
	}

	return tfvUpdateVersions(c, versions, true)
}

func tfvUpdateVersions(c TfxClientContext, versions []string, enabled bool) error {
	opts := tfe.AdminTerraformVersionUpdateOptions{
		Enabled: tfe.Bool(enabled),
	}

	for _, v := range versions {
		tfv, err := adminTFVGetVersion(c, v)
		if err != nil {
			o.AddDeferredMessageRead(v, "failed to find terraform version")
			continue
		}
		if !enabled && tfv.Usage > 0 { //can not disable a version with usage
			o.AddDeferredMessageRead(v, "unable to disable a terraform version in use")
			continue
		}

		tfv, err = c.Client.Admin.TerraformVersions.Update(c.Context, tfv.ID, opts)
		if err == nil {
			o.AddDeferredMessageRead(v, strconv.FormatBool(tfv.Enabled))
		} else {
			return errors.Wrap(err, "failed to update terraform version")
		}
	}

	return nil
}
