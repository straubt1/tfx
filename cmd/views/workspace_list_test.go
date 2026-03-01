// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	"encoding/json"
	"testing"
	"time"

	tfe "github.com/hashicorp/go-tfe"
)

func TestWorkspaceListView_Render(t *testing.T) {
	t.Run("empty list produces empty JSON array", func(t *testing.T) {
		v := NewWorkspaceListView()
		out := captureOutput(t, func() error {
			return v.Render([]*tfe.Workspace{}, false)
		})

		var result []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			t.Fatalf("output is not valid JSON: %v\n%s", err, out)
		}
		if len(result) != 0 {
			t.Errorf("expected empty array, got %d items", len(result))
		}
	})

	t.Run("single workspace basic fields", func(t *testing.T) {
		v := NewWorkspaceListView()
		workspaces := []*tfe.Workspace{
			{Name: "my-workspace", ID: "ws-abc123", ResourceCount: 5, Locked: true},
		}

		out := captureOutput(t, func() error {
			return v.Render(workspaces, false)
		})

		var result []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			t.Fatalf("output is not valid JSON: %v\n%s", err, out)
		}
		if len(result) != 1 {
			t.Fatalf("expected 1 item, got %d", len(result))
		}
		ws := result[0]
		if ws["name"] != "my-workspace" {
			t.Errorf("name = %v, want my-workspace", ws["name"])
		}
		if ws["id"] != "ws-abc123" {
			t.Errorf("id = %v, want ws-abc123", ws["id"])
		}
		if ws["locked"] != true {
			t.Errorf("locked = %v, want true", ws["locked"])
		}
		if ws["resourceCount"] != float64(5) {
			t.Errorf("resourceCount = %v, want 5", ws["resourceCount"])
		}
	})

	t.Run("nil organization field is empty string", func(t *testing.T) {
		v := NewWorkspaceListView()
		workspaces := []*tfe.Workspace{
			{Name: "ws-no-org", ID: "ws-noorg", Organization: nil},
		}

		out := captureOutput(t, func() error {
			return v.Render(workspaces, true)
		})

		var result []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			t.Fatalf("output is not valid JSON: %v", err)
		}
		// organization omitempty: empty string omitted from JSON
		if org, ok := result[0]["organization"]; ok && org != "" {
			t.Errorf("expected organization to be empty or absent, got %v", org)
		}
	})

	t.Run("workspace with organization", func(t *testing.T) {
		v := NewWorkspaceListView()
		workspaces := []*tfe.Workspace{
			{
				Name:         "ws-with-org",
				ID:           "ws-org123",
				Organization: &tfe.Organization{Name: "my-org"},
			},
		}

		out := captureOutput(t, func() error {
			return v.Render(workspaces, true)
		})

		var result []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			t.Fatalf("output is not valid JSON: %v", err)
		}
		if result[0]["organization"] != "my-org" {
			t.Errorf("organization = %v, want my-org", result[0]["organization"])
		}
	})

	t.Run("workspace with current run", func(t *testing.T) {
		v := NewWorkspaceListView()
		runTime := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
		workspaces := []*tfe.Workspace{
			{
				Name: "ws-with-run",
				ID:   "ws-run123",
				CurrentRun: &tfe.Run{
					CreatedAt: runTime,
					Status:    tfe.RunApplied,
				},
			},
		}

		out := captureOutput(t, func() error {
			return v.Render(workspaces, false)
		})

		var result []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			t.Fatalf("output is not valid JSON: %v", err)
		}
		ws := result[0]
		if ws["currentRunStatus"] != "applied" {
			t.Errorf("currentRunStatus = %v, want applied", ws["currentRunStatus"])
		}
		if ws["currentRunCreated"] == "" || ws["currentRunCreated"] == nil {
			t.Errorf("currentRunCreated should be set, got %v", ws["currentRunCreated"])
		}
	})

	t.Run("workspace with VCS repository", func(t *testing.T) {
		v := NewWorkspaceListView()
		workspaces := []*tfe.Workspace{
			{
				Name:    "ws-with-vcs",
				ID:      "ws-vcs123",
				VCSRepo: &tfe.VCSRepo{DisplayIdentifier: "my-org/my-repo"},
			},
		}

		out := captureOutput(t, func() error {
			return v.Render(workspaces, false)
		})

		var result []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			t.Fatalf("output is not valid JSON: %v", err)
		}
		if result[0]["repository"] != "my-org/my-repo" {
			t.Errorf("repository = %v, want my-org/my-repo", result[0]["repository"])
		}
	})

	t.Run("multiple workspaces", func(t *testing.T) {
		v := NewWorkspaceListView()
		workspaces := []*tfe.Workspace{
			{Name: "ws-one", ID: "ws-001"},
			{Name: "ws-two", ID: "ws-002"},
			{Name: "ws-three", ID: "ws-003"},
		}

		out := captureOutput(t, func() error {
			return v.Render(workspaces, false)
		})

		var result []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			t.Fatalf("output is not valid JSON: %v", err)
		}
		if len(result) != 3 {
			t.Errorf("expected 3 items, got %d", len(result))
		}
	})
}
