// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	"encoding/json"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
)

func TestProjectListView_Render(t *testing.T) {
	t.Run("empty list produces empty JSON array", func(t *testing.T) {
		v := NewProjectListView()
		out := captureOutput(t, func() error {
			return v.Render([]*tfe.Project{}, false)
		})

		var result []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			t.Fatalf("output is not valid JSON: %v\n%s", err, out)
		}
		if len(result) != 0 {
			t.Errorf("expected empty array, got %d items", len(result))
		}
	})

	t.Run("single project basic fields", func(t *testing.T) {
		v := NewProjectListView()
		projects := []*tfe.Project{
			{
				Name:        "my-project",
				ID:          "prj-abc123",
				Description: "A test project",
			},
		}

		out := captureOutput(t, func() error {
			return v.Render(projects, false)
		})

		var result []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			t.Fatalf("output is not valid JSON: %v\n%s", err, out)
		}
		if len(result) != 1 {
			t.Fatalf("expected 1 item, got %d", len(result))
		}
		p := result[0]
		if p["name"] != "my-project" {
			t.Errorf("name = %v, want my-project", p["name"])
		}
		if p["id"] != "prj-abc123" {
			t.Errorf("id = %v, want prj-abc123", p["id"])
		}
		if p["description"] != "A test project" {
			t.Errorf("description = %v, want 'A test project'", p["description"])
		}
	})

	t.Run("nil organization field is empty string", func(t *testing.T) {
		v := NewProjectListView()
		projects := []*tfe.Project{
			{Name: "no-org-project", ID: "prj-noorg", Organization: nil},
		}

		out := captureOutput(t, func() error {
			return v.Render(projects, true)
		})

		var result []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			t.Fatalf("output is not valid JSON: %v", err)
		}
		// organization omitempty: empty string may be absent from JSON
		if org, ok := result[0]["organization"]; ok && org != "" {
			t.Errorf("expected organization to be empty or absent, got %v", org)
		}
	})

	t.Run("project with organization", func(t *testing.T) {
		v := NewProjectListView()
		projects := []*tfe.Project{
			{
				Name:         "org-project",
				ID:           "prj-org123",
				Organization: &tfe.Organization{Name: "my-org"},
			},
		}

		out := captureOutput(t, func() error {
			return v.Render(projects, true)
		})

		var result []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			t.Fatalf("output is not valid JSON: %v", err)
		}
		if result[0]["organization"] != "my-org" {
			t.Errorf("organization = %v, want my-org", result[0]["organization"])
		}
	})

	t.Run("multiple projects", func(t *testing.T) {
		v := NewProjectListView()
		projects := []*tfe.Project{
			{Name: "project-one", ID: "prj-001"},
			{Name: "project-two", ID: "prj-002"},
		}

		out := captureOutput(t, func() error {
			return v.Render(projects, false)
		})

		var result []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			t.Fatalf("output is not valid JSON: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 items, got %d", len(result))
		}
		if result[0]["name"] != "project-one" {
			t.Errorf("first project name = %v, want project-one", result[0]["name"])
		}
		if result[1]["name"] != "project-two" {
			t.Errorf("second project name = %v, want project-two", result[1]["name"])
		}
	})
}
