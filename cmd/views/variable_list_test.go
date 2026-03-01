// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	"encoding/json"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
)

func TestVariableListView_Render(t *testing.T) {
	t.Run("empty list produces empty JSON array", func(t *testing.T) {
		v := NewVariableListView()
		out := captureOutput(t, func() error {
			return v.Render("my-workspace", []*tfe.Variable{})
		})

		var result []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			t.Fatalf("output is not valid JSON: %v\n%s", err, out)
		}
		if len(result) != 0 {
			t.Errorf("expected empty array, got %d items", len(result))
		}
	})

	t.Run("single env variable", func(t *testing.T) {
		v := NewVariableListView()
		variables := []*tfe.Variable{
			{
				ID:          "var-abc123",
				Key:         "AWS_REGION",
				Value:       "us-east-1",
				Category:    tfe.CategoryEnv,
				Sensitive:   false,
				HCL:         false,
				Description: "AWS region",
			},
		}

		out := captureOutput(t, func() error {
			return v.Render("my-workspace", variables)
		})

		var result []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			t.Fatalf("output is not valid JSON: %v\n%s", err, out)
		}
		if len(result) != 1 {
			t.Fatalf("expected 1 item, got %d", len(result))
		}
		variable := result[0]
		if variable["id"] != "var-abc123" {
			t.Errorf("id = %v, want var-abc123", variable["id"])
		}
		if variable["key"] != "AWS_REGION" {
			t.Errorf("key = %v, want AWS_REGION", variable["key"])
		}
		if variable["value"] != "us-east-1" {
			t.Errorf("value = %v, want us-east-1", variable["value"])
		}
		if variable["category"] != "env" {
			t.Errorf("category = %v, want env", variable["category"])
		}
		if variable["sensitive"] != false {
			t.Errorf("sensitive = %v, want false", variable["sensitive"])
		}
	})

	t.Run("sensitive variable has value in JSON output", func(t *testing.T) {
		v := NewVariableListView()
		variables := []*tfe.Variable{
			{
				ID:        "var-secret",
				Key:       "API_TOKEN",
				Value:     "",
				Category:  tfe.CategoryEnv,
				Sensitive: true,
			},
		}

		out := captureOutput(t, func() error {
			return v.Render("my-workspace", variables)
		})

		var result []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			t.Fatalf("output is not valid JSON: %v", err)
		}
		if result[0]["sensitive"] != true {
			t.Errorf("sensitive = %v, want true", result[0]["sensitive"])
		}
	})

	t.Run("HCL terraform variable", func(t *testing.T) {
		v := NewVariableListView()
		variables := []*tfe.Variable{
			{
				ID:       "var-hcl",
				Key:      "instance_type",
				Value:    `"t3.medium"`,
				Category: tfe.CategoryTerraform,
				HCL:      true,
			},
		}

		out := captureOutput(t, func() error {
			return v.Render("my-workspace", variables)
		})

		var result []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			t.Fatalf("output is not valid JSON: %v", err)
		}
		if result[0]["hcl"] != true {
			t.Errorf("hcl = %v, want true", result[0]["hcl"])
		}
		if result[0]["category"] != "terraform" {
			t.Errorf("category = %v, want terraform", result[0]["category"])
		}
	})

	t.Run("multiple variables preserve order", func(t *testing.T) {
		v := NewVariableListView()
		variables := []*tfe.Variable{
			{ID: "var-001", Key: "ALPHA"},
			{ID: "var-002", Key: "BETA"},
			{ID: "var-003", Key: "GAMMA"},
		}

		out := captureOutput(t, func() error {
			return v.Render("my-workspace", variables)
		})

		var result []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			t.Fatalf("output is not valid JSON: %v", err)
		}
		if len(result) != 3 {
			t.Fatalf("expected 3 items, got %d", len(result))
		}
		for i, expectedKey := range []string{"ALPHA", "BETA", "GAMMA"} {
			if result[i]["key"] != expectedKey {
				t.Errorf("item %d key = %v, want %s", i, result[i]["key"], expectedKey)
			}
		}
	})
}
