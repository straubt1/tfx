// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package flags

import (
	"testing"

	"github.com/spf13/viper"
)

func TestParseProjectListFlags(t *testing.T) {
	tests := []struct {
		name  string
		setup func()
		want  ProjectListFlags
	}{
		{
			name:  "defaults when nothing set",
			setup: func() {},
			want:  ProjectListFlags{},
		},
		{
			name: "search and all set",
			setup: func() {
				viper.Set("search", "my-project")
				viper.Set("all", true)
			},
			want: ProjectListFlags{Search: "my-project", All: true},
		},
		{
			name: "search only",
			setup: func() {
				viper.Set("search", "prod")
			},
			want: ProjectListFlags{Search: "prod"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			tt.setup()

			got, err := ParseProjectListFlags(nil)
			if err != nil {
				t.Fatalf("ParseProjectListFlags() error = %v", err)
			}
			if *got != tt.want {
				t.Errorf("ParseProjectListFlags() = %+v, want %+v", *got, tt.want)
			}
		})
	}
}

func TestParseProjectShowFlags(t *testing.T) {
	tests := []struct {
		name  string
		setup func()
		want  ProjectShowFlags
	}{
		{
			name:  "defaults when nothing set",
			setup: func() {},
			want:  ProjectShowFlags{},
		},
		{
			name: "id set",
			setup: func() {
				viper.Set("id", "prj-abc123")
			},
			want: ProjectShowFlags{ID: "prj-abc123"},
		},
		{
			name: "name set",
			setup: func() {
				viper.Set("name", "my-project")
			},
			want: ProjectShowFlags{Name: "my-project"},
		},
		{
			name: "both id and name set",
			setup: func() {
				viper.Set("id", "prj-abc123")
				viper.Set("name", "my-project")
			},
			want: ProjectShowFlags{ID: "prj-abc123", Name: "my-project"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			tt.setup()

			got, err := ParseProjectShowFlags(nil)
			if err != nil {
				t.Fatalf("ParseProjectShowFlags() error = %v", err)
			}
			if *got != tt.want {
				t.Errorf("ParseProjectShowFlags() = %+v, want %+v", *got, tt.want)
			}
		})
	}
}
