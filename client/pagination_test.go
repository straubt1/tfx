package client

import (
	"context"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
)

func TestFetchAll(t *testing.T) {
	ctx := context.Background()

	t.Run("single page", func(t *testing.T) {
		callCount := 0
		fetcher := func(pageNumber int) ([]string, *Pagination, error) {
			callCount++
			if pageNumber != 1 {
				t.Errorf("expected pageNumber 1, got %d", pageNumber)
			}
			return []string{"item1", "item2"}, &Pagination{
				CurrentPage: 1,
				NextPage:    0,
				TotalPages:  1,
			}, nil
		}

		results, err := FetchAll(ctx, fetcher)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if callCount != 1 {
			t.Errorf("expected 1 call, got %d", callCount)
		}

		if len(results) != 2 {
			t.Errorf("expected 2 items, got %d", len(results))
		}
	})

	t.Run("multiple pages", func(t *testing.T) {
		callCount := 0
		fetcher := func(pageNumber int) ([]string, *Pagination, error) {
			callCount++
			switch pageNumber {
			case 1:
				return []string{"item1", "item2"}, &Pagination{
					CurrentPage: 1,
					NextPage:    2,
					TotalPages:  3,
				}, nil
			case 2:
				return []string{"item3", "item4"}, &Pagination{
					CurrentPage: 2,
					NextPage:    3,
					TotalPages:  3,
				}, nil
			case 3:
				return []string{"item5"}, &Pagination{
					CurrentPage: 3,
					NextPage:    0,
					TotalPages:  3,
				}, nil
			default:
				t.Errorf("unexpected page number: %d", pageNumber)
				return nil, nil, nil
			}
		}

		results, err := FetchAll(ctx, fetcher)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if callCount != 3 {
			t.Errorf("expected 3 calls, got %d", callCount)
		}

		if len(results) != 5 {
			t.Errorf("expected 5 items, got %d", len(results))
		}
	})

	t.Run("handles errors", func(t *testing.T) {
		fetcher := func(pageNumber int) ([]string, *Pagination, error) {
			return nil, nil, context.DeadlineExceeded
		}

		results, err := FetchAll(ctx, fetcher)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if results != nil {
			t.Errorf("expected nil results on error, got %v", results)
		}
	})
}

func TestNewPaginationFromTFE(t *testing.T) {
	tfePagination := &tfe.Pagination{
		CurrentPage: 2,
		NextPage:    3,
		TotalPages:  5,
	}

	pagination := NewPaginationFromTFE(tfePagination)

	if pagination.CurrentPage != 2 {
		t.Errorf("expected CurrentPage 2, got %d", pagination.CurrentPage)
	}
	if pagination.NextPage != 3 {
		t.Errorf("expected NextPage 3, got %d", pagination.NextPage)
	}
	if pagination.TotalPages != 5 {
		t.Errorf("expected TotalPages 5, got %d", pagination.TotalPages)
	}
}
