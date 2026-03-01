// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package flags

import (
	"testing"

	"github.com/spf13/viper"
)

func TestParseWorkspaceListFlags(t *testing.T) {
	tests := []struct {
		name  string
		setup func()
		want  WorkspaceListFlags
	}{
		{
			name:  "defaults when nothing set",
			setup: func() {},
			want:  WorkspaceListFlags{},
		},
		{
			name: "all fields set",
			setup: func() {
				viper.Set("search", "my-workspace")
				viper.Set("wildcard-name", "ws-*")
				viper.Set("run-status", "applied")
				viper.Set("project-id", "prj-123")
				viper.Set("tags", "env:prod")
				viper.Set("exclude-tags", "env:dev")
				viper.Set("all", true)
			},
			want: WorkspaceListFlags{
				Search:       "my-workspace",
				WildcardName: "ws-*",
				RunStatus:    "applied",
				ProjectID:    "prj-123",
				Tags:         "env:prod",
				ExcludeTags:  "env:dev",
				All:          true,
			},
		},
		{
			name: "search only",
			setup: func() {
				viper.Set("search", "prod")
			},
			want: WorkspaceListFlags{Search: "prod"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			tt.setup()

			got, err := ParseWorkspaceListFlags(nil)
			if err != nil {
				t.Fatalf("ParseWorkspaceListFlags() error = %v", err)
			}
			if *got != tt.want {
				t.Errorf("ParseWorkspaceListFlags() = %+v, want %+v", *got, tt.want)
			}
		})
	}
}

func TestParseWorkspaceShowFlags(t *testing.T) {
	tests := []struct {
		name  string
		setup func()
		want  WorkspaceShowFlags
	}{
		{
			name:  "default empty",
			setup: func() {},
			want:  WorkspaceShowFlags{},
		},
		{
			name: "name set",
			setup: func() {
				viper.Set("name", "my-workspace")
			},
			want: WorkspaceShowFlags{Name: "my-workspace"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			tt.setup()

			got, err := ParseWorkspaceShowFlags(nil)
			if err != nil {
				t.Fatalf("ParseWorkspaceShowFlags() error = %v", err)
			}
			if *got != tt.want {
				t.Errorf("ParseWorkspaceShowFlags() = %+v, want %+v", *got, tt.want)
			}
		})
	}
}
