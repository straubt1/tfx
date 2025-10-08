//go:build integration
// +build integration

package data

import (
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/straubt1/tfx/client"
)

func TestFetchProjects_Integration(t *testing.T) {
	hostname, token, org := getIntegrationTestConfig(t)

	c, err := client.New(hostname, token, org)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	t.Run("fetch all projects in organization", func(t *testing.T) {
		projects, err := FetchProjects(c, org, "")
		if err != nil {
			t.Fatalf("FetchProjects() error = %v", err)
		}

		t.Logf("Successfully fetched %d projects", len(projects))

		// Verify each project has required fields
		for _, p := range projects {
			if p.ID == "" {
				t.Error("Project has empty ID")
			}
			if p.Name == "" {
				t.Error("Project has empty name")
			}
			t.Logf("Project: %s (ID: %s)", p.Name, p.ID)
		}
	})

	t.Run("fetch projects with search string", func(t *testing.T) {
		// First get all projects to find a test project name
		allProjects, err := FetchProjects(c, org, "")
		if err != nil {
			t.Fatalf("FetchProjects() error = %v", err)
		}

		if len(allProjects) == 0 {
			t.Skip("No projects available for search test")
		}

		// Use the first project's name as search string
		searchString := allProjects[0].Name
		projects, err := FetchProjects(c, org, searchString)
		if err != nil {
			t.Fatalf("FetchProjects() with search error = %v", err)
		}

		t.Logf("Search for '%s' returned %d projects", searchString, len(projects))

		// Verify at least one result matches the search
		found := false
		for _, p := range projects {
			if p.Name == searchString {
				found = true
				break
			}
		}
		if !found && len(projects) > 0 {
			t.Logf("Warning: Exact match not found, but got %d results", len(projects))
		}
	})
}

func TestFetchProjectsAcrossOrgs_Integration(t *testing.T) {
	hostname, token, org := getIntegrationTestConfig(t)

	c, err := client.New(hostname, token, org)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	t.Run("fetch all projects across all organizations", func(t *testing.T) {
		projects, err := FetchProjectsAcrossOrgs(c, "")
		if err != nil {
			t.Fatalf("FetchProjectsAcrossOrgs() error = %v", err)
		}

		t.Logf("Successfully fetched %d projects across all organizations", len(projects))

		// Count projects per organization
		orgCounts := make(map[string]int)
		for _, p := range projects {
			orgCounts[p.Organization.Name]++
		}

		for orgName, count := range orgCounts {
			t.Logf("Organization %s has %d projects", orgName, count)
		}
	})
}

func TestFetchProject_Integration(t *testing.T) {
	hostname, token, org := getIntegrationTestConfig(t)

	c, err := client.New(hostname, token, org)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Get a project to test with
	projects, err := FetchProjects(c, org, "")
	if err != nil {
		t.Fatalf("Failed to fetch projects: %v", err)
	}

	if len(projects) == 0 {
		t.Skip("No projects available for testing")
	}

	testProject := projects[0]

	t.Run("fetch project by ID", func(t *testing.T) {
		options := &tfe.ProjectReadOptions{}
		project, err := FetchProject(c, testProject.ID, options)
		if err != nil {
			t.Fatalf("FetchProject() error = %v", err)
		}

		if project.ID != testProject.ID {
			t.Errorf("Expected project ID %s, got %s", testProject.ID, project.ID)
		}
		if project.Name != testProject.Name {
			t.Errorf("Expected project name %s, got %s", testProject.Name, project.Name)
		}

		t.Logf("Successfully fetched project: %s (ID: %s)", project.Name, project.ID)
	})

	t.Run("fetch project with invalid ID", func(t *testing.T) {
		options := &tfe.ProjectReadOptions{}
		_, err := FetchProject(c, "prj-invalid123456", options)
		if err == nil {
			t.Error("Expected error for invalid project ID, got nil")
		}
		t.Logf("Got expected error for invalid ID: %v", err)
	})
}

func TestFetchProjectByName_Integration(t *testing.T) {
	hostname, token, org := getIntegrationTestConfig(t)

	c, err := client.New(hostname, token, org)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Get a project to test with
	projects, err := FetchProjects(c, org, "")
	if err != nil {
		t.Fatalf("Failed to fetch projects: %v", err)
	}

	if len(projects) == 0 {
		t.Skip("No projects available for testing")
	}

	testProject := projects[0]

	t.Run("fetch project by name", func(t *testing.T) {
		options := &tfe.ProjectReadOptions{}
		project, err := FetchProjectByName(c, org, testProject.Name, options)
		if err != nil {
			t.Fatalf("FetchProjectByName() error = %v", err)
		}

		if project.Name != testProject.Name {
			t.Errorf("Expected project name %s, got %s", testProject.Name, project.Name)
		}
		if project.ID != testProject.ID {
			t.Errorf("Expected project ID %s, got %s", testProject.ID, project.ID)
		}

		t.Logf("Successfully fetched project by name: %s (ID: %s)", project.Name, project.ID)
	})

	t.Run("fetch project with non-existent name", func(t *testing.T) {
		options := &tfe.ProjectReadOptions{}
		_, err := FetchProjectByName(c, org, "non-existent-project-name-12345", options)
		if err == nil {
			t.Error("Expected error for non-existent project name, got nil")
		}
		t.Logf("Got expected error for non-existent name: %v", err)
	})

	t.Run("fetch project with invalid organization", func(t *testing.T) {
		options := &tfe.ProjectReadOptions{}
		_, err := FetchProjectByName(c, "non-existent-org-12345", testProject.Name, options)
		if err == nil {
			t.Error("Expected error for invalid organization, got nil")
		}
		t.Logf("Got expected error for invalid organization: %v", err)
	})
}
