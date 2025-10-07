package client

import tfe "github.com/hashicorp/go-tfe"

// FetchOrganizations fetches all organizations using pagination
func (c *TfxClient) FetchOrganizations() ([]*tfe.Organization, error) {
	return FetchAll(c.Context, func(pageNumber int) ([]*tfe.Organization, *Pagination, error) {
		opts := &tfe.OrganizationListOptions{
			ListOptions: tfe.ListOptions{PageNumber: pageNumber, PageSize: 100},
		}

		result, err := c.Client.Organizations.List(c.Context, opts)
		if err != nil {
			return nil, nil, err
		}

		return result.Items, NewPaginationFromTFE(result.Pagination), nil
	})
}
