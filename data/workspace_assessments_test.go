// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/straubt1/tfx/client"
)

// newAssessmentsTestServer returns a fake TFE API that serves a single workspace
// and records PATCH requests made against it.
func newAssessmentsTestServer(t *testing.T, assessmentsEnabled bool, patchStatus int) (*client.TfxClient, *[]map[string]interface{}) {
	t.Helper()
	var patches []map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		switch {
		case r.URL.Path == "/api/v2/ping":
			w.Header().Set("TFP-API-Version", "2.6")
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/organizations/my-org/workspaces/my-ws":
			fmt.Fprintf(w, `{"data":{"id":"ws-123","type":"workspaces","attributes":{"name":"my-ws","assessments-enabled":%t}}}`, assessmentsEnabled)
		case r.Method == http.MethodPatch && r.URL.Path == "/api/v2/workspaces/ws-123":
			body, _ := io.ReadAll(r.Body)
			var payload map[string]interface{}
			if err := json.Unmarshal(body, &payload); err != nil {
				t.Errorf("invalid PATCH body: %v", err)
			}
			patches = append(patches, payload)
			if patchStatus != http.StatusOK {
				w.WriteHeader(patchStatus)
				fmt.Fprint(w, `{"errors":[{"status":"422","title":"invalid attribute"}]}`)
				return
			}
			fmt.Fprint(w, `{"data":{"id":"ws-123","type":"workspaces","attributes":{"name":"my-ws"}}}`)
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/organizations/my-org/workspaces/missing-ws":
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"errors":[{"status":"404","title":"not found"}]}`)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	tfeClient, err := tfe.NewClient(&tfe.Config{Address: server.URL, Token: "test-token", HTTPClient: server.Client()})
	if err != nil {
		t.Fatalf("failed to create tfe client: %v", err)
	}
	return &client.TfxClient{Client: tfeClient, Context: context.Background(), OrganizationName: "my-org"}, &patches
}

func assessmentsAttr(t *testing.T, payload map[string]interface{}) interface{} {
	t.Helper()
	d, _ := payload["data"].(map[string]interface{})
	attrs, _ := d["attributes"].(map[string]interface{})
	v, ok := attrs["assessments-enabled"]
	if !ok {
		t.Fatalf("PATCH body missing assessments-enabled: %v", payload)
	}
	return v
}

func TestSetWorkspaceAssessmentsByName(t *testing.T) {
	tests := []struct {
		name        string
		current     bool
		enabled     bool
		wantStatus  string
		wantPatched bool
	}{
		{name: "enable when disabled", current: false, enabled: true, wantStatus: "Enabled", wantPatched: true},
		{name: "disable when enabled", current: true, enabled: false, wantStatus: "Disabled", wantPatched: true},
		{name: "enable when already enabled", current: true, enabled: true, wantStatus: "Already enabled"},
		{name: "disable when already disabled", current: false, enabled: false, wantStatus: "Already disabled"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, patches := newAssessmentsTestServer(t, tt.current, http.StatusOK)

			status, err := SetWorkspaceAssessmentsByName(c, "my-org", "my-ws", tt.enabled)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if status != tt.wantStatus {
				t.Errorf("status = %q, want %q", status, tt.wantStatus)
			}

			if !tt.wantPatched {
				if len(*patches) != 0 {
					t.Fatalf("expected no PATCH, got %d", len(*patches))
				}
				return
			}
			if len(*patches) != 1 {
				t.Fatalf("expected 1 PATCH, got %d", len(*patches))
			}
			if got := assessmentsAttr(t, (*patches)[0]); got != tt.enabled {
				t.Errorf("assessments-enabled = %v, want %v", got, tt.enabled)
			}
		})
	}
}

func TestSetWorkspaceAssessmentsUpdateError(t *testing.T) {
	c, _ := newAssessmentsTestServer(t, false, http.StatusUnprocessableEntity)

	_, err := SetWorkspaceAssessmentsByName(c, "my-org", "my-ws", true)
	if err == nil {
		t.Fatal("expected error from failed update")
	}
}

func TestSetWorkspaceAssessmentsByNameNotFound(t *testing.T) {
	c, patches := newAssessmentsTestServer(t, false, http.StatusOK)

	_, err := SetWorkspaceAssessmentsByName(c, "my-org", "missing-ws", true)
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected not found error, got %v", err)
	}
	if len(*patches) != 0 {
		t.Fatalf("expected no PATCH, got %d", len(*patches))
	}
}

func TestSetWorkspaceAssessmentsSkipsReadForFetchedWorkspace(t *testing.T) {
	c, patches := newAssessmentsTestServer(t, false, http.StatusOK)

	status, err := SetWorkspaceAssessments(c, &tfe.Workspace{ID: "ws-123", Name: "my-ws"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != "Enabled" {
		t.Errorf("status = %q, want Enabled", status)
	}
	if len(*patches) != 1 {
		t.Fatalf("expected 1 PATCH, got %d", len(*patches))
	}
}
