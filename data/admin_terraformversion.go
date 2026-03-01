// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/output"
)

// FetchTerraformVersions fetches all Terraform versions with optional filtering
func FetchTerraformVersions(c *client.TfxClient, filter string, search string) ([]*tfe.AdminTerraformVersion, error) {
	output.Get().Logger().Debug("Fetching Terraform versions", "filter", filter, "search", search)

	return client.FetchAll(c.Context, func(pageNumber int) ([]*tfe.AdminTerraformVersion, *client.Pagination, error) {
		output.Get().Logger().Trace("Fetching Terraform versions page", "page", pageNumber)

		opts := &tfe.AdminTerraformVersionsListOptions{
			ListOptions: tfe.ListOptions{PageNumber: pageNumber, PageSize: 100},
			Filter:      filter,
			Search:      search,
		}

		result, err := c.Client.Admin.TerraformVersions.List(c.Context, opts)
		if err != nil {
			output.Get().Logger().Error("Failed to fetch Terraform versions page", "page", pageNumber, "error", err)
			return nil, nil, err
		}

		output.Get().Logger().Trace("Terraform versions page fetched", "page", pageNumber, "count", len(result.Items))
		return result.Items, client.NewPaginationFromTFE(result.Pagination), nil
	})
}

// FetchTerraformVersion fetches a single Terraform version by version string
// Uses filter which returns exact match
func FetchTerraformVersion(c *client.TfxClient, version string) (*tfe.AdminTerraformVersion, error) {
	output.Get().Logger().Debug("Fetching Terraform version", "version", version)

	items, err := FetchTerraformVersions(c, version, "")
	if err != nil {
		output.Get().Logger().Error("Failed to fetch Terraform version", "version", version, "error", err)
		return nil, errors.Wrap(err, "failed to fetch terraform version")
	}

	if len(items) == 0 {
		return nil, errors.New("terraform version not found")
	} else if len(items) > 1 {
		// unlikely to ever hit this, but just in case
		return nil, errors.New("too many terraform versions found")
	}

	output.Get().Logger().Debug("Terraform version fetched successfully", "version", version)
	return items[0], nil
}

// CreateTerraformVersion creates a new Terraform version
func CreateTerraformVersion(c *client.TfxClient, version string, url string, sha string, official bool, enabled bool, beta bool) (*tfe.AdminTerraformVersion, error) {
	output.Get().Logger().Debug("Creating Terraform version", "version", version, "official", official, "enabled", enabled, "beta", beta)

	opts := tfe.AdminTerraformVersionCreateOptions{
		Version:  tfe.String(version),
		URL:      tfe.String(url),
		Sha:      tfe.String(sha),
		Official: tfe.Bool(official),
		Enabled:  tfe.Bool(enabled),
		Beta:     tfe.Bool(beta),
	}

	tfv, err := c.Client.Admin.TerraformVersions.Create(c.Context, opts)
	if err != nil {
		output.Get().Logger().Error("Failed to create Terraform version", "version", version, "error", err)
		return nil, errors.Wrap(err, "unable to create terraform version")
	}

	output.Get().Logger().Debug("Terraform version created successfully", "version", version, "id", tfv.ID)
	return tfv, nil
}

// CreateOfficialTerraformVersion creates a Terraform version from the official HashiCorp releases
func CreateOfficialTerraformVersion(c *client.TfxClient, version string, enabled bool, beta bool) (*tfe.AdminTerraformVersion, error) {
	output.Get().Logger().Debug("Creating official Terraform version", "version", version)

	// Build URLs for official releases
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

	// Fetch the SHA checksum from HashiCorp releases
	output.Get().Logger().Debug("Fetching SHA checksum from HashiCorp releases", "url", urlSha)

	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", urlSha, nil)
	if err != nil {
		output.Get().Logger().Error("Failed to create HTTP request", "error", err)
		return nil, errors.Wrap(err, "failed to find official terraform version")
	}

	resp, err := httpClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		output.Get().Logger().Error("Failed to fetch SHA checksum", "statusCode", resp.StatusCode, "error", err)
		return nil, errors.New("failed to find official terraform version")
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		output.Get().Logger().Error("Failed to read SHA checksum response", "error", err)
		return nil, errors.Wrap(err, "failed to read checksum")
	}

	// Parse the checksum file to find the linux_amd64 version
	var sha string
	lines := strings.Split(string(b), "\n")
	for _, l := range lines {
		if strings.Contains(l, "linux_amd64") {
			// Extract the checksum (first part before space)
			parts := strings.Split(l, " ")
			if len(parts) > 0 {
				sha = parts[0]
				break
			}
		}
	}

	if sha == "" {
		return nil, errors.New("failed to find SHA checksum for linux_amd64")
	}

	output.Get().Logger().Debug("SHA checksum found", "sha", sha)

	// Create the version
	return CreateTerraformVersion(c, version, url, sha, true, enabled, beta)
}

// DeleteTerraformVersion deletes a Terraform version
func DeleteTerraformVersion(c *client.TfxClient, version string) error {
	output.Get().Logger().Debug("Deleting Terraform version", "version", version)

	tfv, err := FetchTerraformVersion(c, version)
	if err != nil {
		return errors.Wrap(err, "failed to find terraform version")
	}

	// If the version is official, it must be set to unofficial before deletion.
	// Note: The TFE platform may re-sync official versions on upgrade; deleted official
	// versions could reappear. This is a known platform limitation.
	if tfv.Official {
		output.Get().Logger().Debug("Setting Terraform version to unofficial before deletion", "version", version)

		tfv, err = c.Client.Admin.TerraformVersions.Update(c.Context, tfv.ID, tfe.AdminTerraformVersionUpdateOptions{
			Official: tfe.Bool(false),
		})
		if err != nil {
			output.Get().Logger().Error("Failed to set version to unofficial", "version", version, "error", err)
			return errors.Wrap(err, "failed to set version to official to false")
		}
	}

	// Delete the version
	err = c.Client.Admin.TerraformVersions.Delete(c.Context, tfv.ID)
	if err != nil {
		output.Get().Logger().Error("Failed to delete Terraform version", "version", version, "error", err)
		return errors.Wrap(err, "failed to delete version")
	}

	output.Get().Logger().Debug("Terraform version deleted successfully", "version", version)
	return nil
}

// UpdateTerraformVersions enables or disables multiple Terraform versions
func UpdateTerraformVersions(c *client.TfxClient, versions []string, enabled bool) (map[string]string, error) {
	output.Get().Logger().Debug("Updating Terraform versions", "count", len(versions), "enabled", enabled)

	results := make(map[string]string)

	opts := tfe.AdminTerraformVersionUpdateOptions{
		Enabled: tfe.Bool(enabled),
	}

	for _, v := range versions {
		output.Get().Logger().Trace("Processing version", "version", v)

		tfv, err := FetchTerraformVersion(c, v)
		if err != nil {
			output.Get().Logger().Warn("Failed to find Terraform version", "version", v, "error", err)
			results[v] = "failed to find terraform version"
			continue
		}

		// Cannot disable a version with usage
		if !enabled && tfv.Usage > 0 {
			output.Get().Logger().Warn("Cannot disable Terraform version with usage", "version", v, "usage", tfv.Usage)
			results[v] = "unable to disable a terraform version in use"
			continue
		}

		tfv, err = c.Client.Admin.TerraformVersions.Update(c.Context, tfv.ID, opts)
		if err != nil {
			output.Get().Logger().Error("Failed to update Terraform version", "version", v, "error", err)
			return results, errors.Wrap(err, "failed to update terraform version")
		}

		results[v] = fmt.Sprintf("%t", tfv.Enabled)
	}

	output.Get().Logger().Debug("Terraform versions updated successfully", "results", len(results))
	return results, nil
}
