//go:build integration
// +build integration

package data

import (
	"os"
	"testing"

	"github.com/straubt1/tfx/client"
)

// getIntegrationTestConfig retrieves test configuration from environment variables
// and skips the test if any required variables are missing
func getIntegrationTestConfig(t *testing.T) (hostname, token, organization string) {
	hostname = os.Getenv("TFE_HOSTNAME")
	token = os.Getenv("TFE_TOKEN")
	organization = os.Getenv("TFE_ORGANIZATION")

	if hostname == "" || token == "" || organization == "" {
		t.Skip("Skipping integration test: TFE_HOSTNAME, TFE_TOKEN, and TFE_ORGANIZATION must be set")
	}

	return hostname, token, organization
}

func TestFetchOrganizations_Integration(t *testing.T) {
	hostname, token, org := getIntegrationTestConfig(t)

	c, err := client.New(hostname, token, org)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	t.Run("fetch all organizations", func(t *testing.T) {
		orgs, err := FetchOrganizations(c, "")
		if err != nil {
			t.Fatalf("FetchOrganizations() error = %v", err)
		}

		if len(orgs) == 0 {
			t.Error("Expected at least one organization, got none")
		}

		t.Logf("Successfully fetched %d organizations", len(orgs))

		// Verify the test organization is in the list
		found := false
		for _, o := range orgs {
			if o.Name == org {
				found = true
				t.Logf("Found test organization: %s (email: %s)", o.Name, o.Email)
				break
			}
		}

		if !found {
			t.Errorf("Test organization %s not found in results", org)
		}
	})

	t.Run("verify organization details", func(t *testing.T) {
		orgs, err := FetchOrganizations(c, "")
		if err != nil {
			t.Fatalf("FetchOrganizations() error = %v", err)
		}

		// Verify each organization has required fields
		for _, o := range orgs {
			if o.Name == "" {
				t.Error("Organization has empty name")
			}
			if o.Email == "" {
				t.Logf("Warning: Organization %s has no email", o.Name)
			}
			t.Logf("Organization: %s (External ID: %s)", o.Name, o.ExternalID)
		}
	})
}
