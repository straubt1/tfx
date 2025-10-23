// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/logger"
)

// FetchConfigurationVersions lists configuration versions for a workspace name with max-items
func FetchConfigurationVersions(c *client.TfxClient, orgName, workspaceName string, maxItems int) ([]*tfe.ConfigurationVersion, error) {
	logger.Debug("Fetching configuration versions", "organization", orgName, "workspace", workspaceName, "maxItems", maxItems)

	// Resolve workspace ID
	workspaceID, err := GetWorkspaceID(c, orgName, workspaceName)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read workspace id")
	}

	if maxItems == 0 {
		maxItems = 100
	}

	pageSize := 100
	if maxItems > 0 && maxItems < 100 {
		pageSize = maxItems
	}

	var all []*tfe.ConfigurationVersion
	opts := &tfe.ConfigurationVersionListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 1, PageSize: pageSize},
		Include:     []tfe.ConfigVerIncludeOpt{"ingress_attributes"},
	}

	for {
		res, err := c.Client.ConfigurationVersions.List(c.Context, workspaceID, opts)
		if err != nil {
			logger.Error("Failed to list configuration versions", "workspaceID", workspaceID, "page", opts.PageNumber, "error", err)
			return nil, err
		}

		all = append(all, res.Items...)
		if maxItems > 0 && len(all) >= maxItems {
			break
		}

		if res.CurrentPage >= res.TotalPages {
			break
		}
		opts.PageNumber = res.NextPage
	}

	if maxItems > 0 && len(all) > maxItems {
		all = all[:maxItems]
	}

	logger.Debug("Configuration versions fetched", "count", len(all))
	return all, nil
}

// CreateConfigurationVersion creates a new configuration version in a workspace and uploads code from a directory
func CreateConfigurationVersion(c *client.TfxClient, orgName, workspaceName, directory string, speculative bool) (*tfe.ConfigurationVersion, error) {
	logger.Debug("Creating configuration version", "organization", orgName, "workspace", workspaceName, "speculative", speculative)

	workspaceID, err := GetWorkspaceID(c, orgName, workspaceName)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read workspace id")
	}

	cv, err := c.Client.ConfigurationVersions.Create(c.Context, workspaceID, tfe.ConfigurationVersionCreateOptions{
		AutoQueueRuns: tfe.Bool(false),
		Speculative:   tfe.Bool(speculative),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create configuration version")
	}

	if err := c.Client.ConfigurationVersions.Upload(c.Context, cv.UploadURL, directory); err != nil {
		return nil, errors.Wrap(err, "failed to upload code to the configuration version")
	}

	logger.Info("Configuration version created", "configurationVersionID", cv.ID)
	return cv, nil
}

// FetchConfigurationVersion reads a configuration version with includes
func FetchConfigurationVersion(c *client.TfxClient, configurationID string) (*tfe.ConfigurationVersion, error) {
	logger.Debug("Fetching configuration version", "configurationVersionID", configurationID)

	cv, err := c.Client.ConfigurationVersions.ReadWithOptions(c.Context, configurationID, &tfe.ConfigurationVersionReadOptions{
		Include: []tfe.ConfigVerIncludeOpt{"ingress_attributes"},
	})
	if err != nil {
		logger.Error("Failed to read configuration version", "configurationVersionID", configurationID, "error", err)
		return nil, err
	}
	return cv, nil
}

// DownloadConfigurationVersion downloads the configuration version slug bytes
func DownloadConfigurationVersion(c *client.TfxClient, configurationID string) ([]byte, error) {
	logger.Debug("Downloading configuration version", "configurationVersionID", configurationID)

	data, err := c.Client.ConfigurationVersions.Download(c.Context, configurationID)
	if err != nil {
		logger.Error("Failed to download configuration version", "configurationVersionID", configurationID, "error", err)
		return nil, err
	}
	return data, nil
}
