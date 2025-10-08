package data

import (
	tfe "github.com/hashicorp/go-tfe"
	"github.com/straubt1/tfx/client"
)

// FetchOrganizations fetches all organizations using pagination
func FetchOrganizations(c *client.TfxClient) ([]*tfe.Organization, error) {
	return client.FetchAll(c.Context, func(pageNumber int) ([]*tfe.Organization, *client.Pagination, error) {
		opts := &tfe.OrganizationListOptions{
			ListOptions: tfe.ListOptions{PageNumber: pageNumber, PageSize: 100},
		}

		result, err := c.Client.Organizations.List(c.Context, opts)
		if err != nil {
			return nil, nil, err
		}

		return result.Items, client.NewPaginationFromTFE(result.Pagination), nil
	})
}
