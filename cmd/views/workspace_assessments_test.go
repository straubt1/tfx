// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	"encoding/json"
	"testing"
)

func TestWorkspaceAssessmentsView_Render(t *testing.T) {
	t.Run("empty results produce empty JSON array", func(t *testing.T) {
		v := NewWorkspaceAssessmentsView()
		out := captureOutput(t, func() error {
			return v.Render([]WorkspaceAssessmentsResult{})
		})

		var result []WorkspaceAssessmentsResult
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			t.Fatalf("output is not valid JSON: %v\n%s", err, out)
		}
		if len(result) != 0 {
			t.Errorf("expected empty array, got %d items", len(result))
		}
	})

	t.Run("results keep order and fields", func(t *testing.T) {
		v := NewWorkspaceAssessmentsView()
		in := []WorkspaceAssessmentsResult{
			{Name: "ws-a", Status: "Enabled"},
			{Name: "ws-b", Status: "Already enabled"},
			{Name: "ws-c", Status: "resource not found"},
		}
		out := captureOutput(t, func() error {
			return v.Render(in)
		})

		var result []WorkspaceAssessmentsResult
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			t.Fatalf("output is not valid JSON: %v\n%s", err, out)
		}
		if len(result) != len(in) {
			t.Fatalf("expected %d items, got %d", len(in), len(result))
		}
		for i := range in {
			if result[i] != in[i] {
				t.Errorf("item %d = %+v, want %+v", i, result[i], in[i])
			}
		}
	})
}
