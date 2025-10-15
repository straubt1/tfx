package data

import (
	tfe "github.com/hashicorp/go-tfe"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/logger"
)

// FetchOrganizations fetches all organizations using pagination with optional search
func FetchOrganizations(c *client.TfxClient, searchString string) ([]*tfe.Organization, error) {
	logger.Debug("Fetching organizations", "searchString", searchString)

	return client.FetchAll(c.Context, func(pageNumber int) ([]*tfe.Organization, *client.Pagination, error) {
		logger.Trace("Fetching organizations page", "page", pageNumber)

		opts := &tfe.OrganizationListOptions{
			ListOptions: tfe.ListOptions{PageNumber: pageNumber, PageSize: 100},
			Query:       searchString,
		}

		result, err := c.Client.Organizations.List(c.Context, opts)
		if err != nil {
			logger.Error("Failed to fetch organizations page", "page", pageNumber, "error", err)
			return nil, nil, err
		}

		logger.Trace("Organizations page fetched", "page", pageNumber, "count", len(result.Items))
		return result.Items, client.NewPaginationFromTFE(result.Pagination), nil
	})
}

// FetchOrganization fetches a single organization by name
func FetchOrganization(c *client.TfxClient, orgName string, options *tfe.OrganizationReadOptions) (*tfe.Organization, error) {
	logger.Debug("Fetching organization by name", "organization", orgName)

	org, err := c.Client.Organizations.ReadWithOptions(c.Context, orgName, *options)
	if err != nil {
		logger.Error("Failed to fetch organization", "organization", orgName, "error", err)
		return nil, err
	}

	logger.Debug("Organization fetched successfully", "organization", orgName, "id", org.ExternalID)
	return org, nil
}
