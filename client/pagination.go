package client

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

// Pagination represents the pagination information from TFE API responses
type Pagination struct {
	CurrentPage int
	NextPage    int
	TotalPages  int
}

// FetchAll is a generic pagination helper that fetches all pages of results from the TFE API.
// It accepts a fetcher function that takes a page number and returns items, pagination info, and an error.
//
// Example usage:
//
//	projects, err := FetchAll(ctx, func(pageNumber int) ([]*tfe.Project, *Pagination, error) {
//	    opts := &tfe.ProjectListOptions{
//	        ListOptions: tfe.ListOptions{PageNumber: pageNumber, PageSize: 100},
//	    }
//	    result, err := client.Projects.List(ctx, "org-name", opts)
//	    if err != nil {
//	        return nil, nil, err
//	    }
//	    return result.Items, &Pagination{
//	        CurrentPage: result.CurrentPage,
//	        NextPage:    result.NextPage,
//	        TotalPages:  result.TotalPages,
//	    }, nil
//	})
func FetchAll[T any](ctx context.Context, fetcher func(pageNumber int) ([]T, *Pagination, error)) ([]T, error) {
	var allItems []T
	pageNumber := 1

	for {
		items, pagination, err := fetcher(pageNumber)
		if err != nil {
			return nil, err
		}

		allItems = append(allItems, items...)

		if pagination.CurrentPage >= pagination.TotalPages {
			break
		}
		pageNumber = pagination.NextPage
	}

	return allItems, nil
}

// NewPaginationFromTFE converts TFE pagination to our Pagination type
func NewPaginationFromTFE(p *tfe.Pagination) *Pagination {
	return &Pagination{
		CurrentPage: p.CurrentPage,
		NextPage:    p.NextPage,
		TotalPages:  p.TotalPages,
	}
}
