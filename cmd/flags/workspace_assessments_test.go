// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package flags

import (
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

func TestParseWorkspaceAssessmentsFlags(t *testing.T) {
	tests := []struct {
		name  string
		setup func()
		want  []string
	}{
		{
			name:  "no names set",
			setup: func() {},
			want:  []string{},
		},
		{
			name:  "single name",
			setup: func() { viper.Set("name", []string{"ws-a"}) },
			want:  []string{"ws-a"},
		},
		{
			name:  "multiple names",
			setup: func() { viper.Set("name", []string{"ws-a", "ws-b"}) },
			want:  []string{"ws-a", "ws-b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			tt.setup()

			got, err := ParseWorkspaceAssessmentsFlags(nil)
			if err != nil {
				t.Fatalf("ParseWorkspaceAssessmentsFlags() error = %v", err)
			}
			if len(tt.want) == 0 {
				if len(got.Names) != 0 {
					t.Errorf("Names = %#v, want empty", got.Names)
				}
				return
			}
			if !reflect.DeepEqual(got.Names, tt.want) {
				t.Errorf("Names = %#v, want %#v", got.Names, tt.want)
			}
		})
	}
}

func TestParseWorkspaceAssessmentsAllFlags(t *testing.T) {
	viper.Reset()
	viper.Set("search", "prod")
	viper.Set("wildcard-name", "app-*")
	viper.Set("project-id", "prj-123")
	viper.Set("tags", "team:a")
	viper.Set("exclude-tags", "legacy")

	got, err := ParseWorkspaceAssessmentsAllFlags(nil)
	if err != nil {
		t.Fatalf("ParseWorkspaceAssessmentsAllFlags() error = %v", err)
	}

	want := WorkspaceAssessmentsAllFlags{
		Search:       "prod",
		WildcardName: "app-*",
		ProjectID:    "prj-123",
		Tags:         "team:a",
		ExcludeTags:  "legacy",
	}
	if *got != want {
		t.Errorf("got %+v, want %+v", *got, want)
	}

	listFlags := got.ListFlags()
	wantList := WorkspaceListFlags{
		Search:       "prod",
		WildcardName: "app-*",
		ProjectID:    "prj-123",
		Tags:         "team:a",
		ExcludeTags:  "legacy",
	}
	if *listFlags != wantList {
		t.Errorf("ListFlags() = %+v, want %+v", *listFlags, wantList)
	}
}
