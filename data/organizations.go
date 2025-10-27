package data

import (
	tfe "github.com/hashicorp/go-tfe"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/output"
)

// OrganizationListOptions holds options for listing organizations
type OrganizationListOptions struct {
	Search string
}

// FetchOrganizationsWithOptions fetches all organizations using pagination with options
func FetchOrganizationsWithOptions(c *client.TfxClient, options *OrganizationListOptions) ([]*tfe.Organization, error) {
	searchString := ""
	if options != nil {
		searchString = options.Search
	}

	output.Get().Logger().Debug("Fetching organizations", "searchString", searchString)

	orgs, err := client.FetchAll(c.Context, func(pageNumber int) ([]*tfe.Organization, *client.Pagination, error) {
		output.Get().Logger().Trace("Fetching organizations page", "page", pageNumber)

		opts := &tfe.OrganizationListOptions{
			ListOptions: tfe.ListOptions{PageNumber: pageNumber, PageSize: 100},
			Query:       searchString,
		}

		result, err := c.Client.Organizations.List(c.Context, opts)
		if err != nil {
			output.Get().Logger().Error("Failed to fetch organizations page", "page", pageNumber, "error", err)
			return nil, nil, err
		}

		output.Get().Logger().Trace("Organizations page fetched", "page", pageNumber, "count", len(result.Items))
		return result.Items, client.NewPaginationFromTFE(result.Pagination), nil
	})

	if err != nil {
		return nil, err
	}

	// Log all organization names
	orgNames := make([]string, len(orgs))
	for i, org := range orgs {
		orgNames[i] = org.Name
	}
	output.Get().Logger().Info("Found organizations", "count", len(orgs), "names", orgNames)

	return orgs, nil
}

// FetchOrganizations fetches all organizations using pagination with optional search
// Deprecated: Use FetchOrganizationsWithOptions instead
func FetchOrganizations(c *client.TfxClient, searchString string) ([]*tfe.Organization, error) {
	options := &OrganizationListOptions{
		Search: searchString,
	}
	return FetchOrganizationsWithOptions(c, options)
}

// FetchOrganization fetches a single organization by name
func FetchOrganization(c *client.TfxClient, orgName string, options *tfe.OrganizationReadOptions) (*tfe.Organization, error) {
	output.Get().Logger().Debug("Fetching organization by name", "organization", orgName)

	org, err := c.Client.Organizations.ReadWithOptions(c.Context, orgName, *options)
	if err != nil {
		output.Get().Logger().Error("Failed to fetch organization", "organization", orgName, "error", err)
		return nil, err
	}

	output.Get().Logger().Debug("Organization fetched successfully", "organization", orgName, "id", org.ExternalID)
	return org, nil
}
