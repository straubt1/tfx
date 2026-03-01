// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package flags

import (
	"testing"

	"github.com/spf13/viper"
)

func TestParseVariableListFlags(t *testing.T) {
	tests := []struct {
		name  string
		setup func()
		want  VariableListFlags
	}{
		{
			name:  "defaults when nothing set",
			setup: func() {},
			want:  VariableListFlags{},
		},
		{
			name: "workspace-name set",
			setup: func() {
				viper.Set("name", "my-workspace")
			},
			want: VariableListFlags{WorkspaceName: "my-workspace"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			tt.setup()

			got, err := ParseVariableListFlags(nil)
			if err != nil {
				t.Fatalf("ParseVariableListFlags() error = %v", err)
			}
			if *got != tt.want {
				t.Errorf("ParseVariableListFlags() = %+v, want %+v", *got, tt.want)
			}
		})
	}
}

func TestParseVariableShowFlags(t *testing.T) {
	tests := []struct {
		name  string
		setup func()
		want  VariableShowFlags
	}{
		{
			name:  "defaults when nothing set",
			setup: func() {},
			want:  VariableShowFlags{},
		},
		{
			name: "all fields set",
			setup: func() {
				viper.Set("name", "my-workspace")
				viper.Set("key", "MY_VAR")
			},
			want: VariableShowFlags{WorkspaceName: "my-workspace", Key: "MY_VAR"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			tt.setup()

			got, err := ParseVariableShowFlags(nil)
			if err != nil {
				t.Fatalf("ParseVariableShowFlags() error = %v", err)
			}
			if *got != tt.want {
				t.Errorf("ParseVariableShowFlags() = %+v, want %+v", *got, tt.want)
			}
		})
	}
}

func TestParseVariableCreateFlags(t *testing.T) {
	tests := []struct {
		name  string
		setup func()
		want  VariableCreateFlags
	}{
		{
			name:  "defaults when nothing set",
			setup: func() {},
			want:  VariableCreateFlags{},
		},
		{
			name: "all fields set",
			setup: func() {
				viper.Set("name", "my-workspace")
				viper.Set("key", "MY_VAR")
				viper.Set("value", "my-value")
				viper.Set("value-file", "/tmp/val.txt")
				viper.Set("description", "a test variable")
				viper.Set("env", true)
				viper.Set("hcl", true)
				viper.Set("sensitive", true)
			},
			want: VariableCreateFlags{
				WorkspaceName: "my-workspace",
				Key:           "MY_VAR",
				Value:         "my-value",
				ValueFile:     "/tmp/val.txt",
				Description:   "a test variable",
				Env:           true,
				HCL:           true,
				Sensitive:     true,
			},
		},
		{
			name: "sensitive env variable",
			setup: func() {
				viper.Set("name", "ws1")
				viper.Set("key", "SECRET")
				viper.Set("value", "topsecret")
				viper.Set("env", true)
				viper.Set("sensitive", true)
			},
			want: VariableCreateFlags{
				WorkspaceName: "ws1",
				Key:           "SECRET",
				Value:         "topsecret",
				Env:           true,
				Sensitive:     true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			tt.setup()

			got, err := ParseVariableCreateFlags(nil)
			if err != nil {
				t.Fatalf("ParseVariableCreateFlags() error = %v", err)
			}
			if *got != tt.want {
				t.Errorf("ParseVariableCreateFlags() = %+v, want %+v", *got, tt.want)
			}
		})
	}
}

func TestParseVariableUpdateFlags(t *testing.T) {
	tests := []struct {
		name  string
		setup func()
		want  VariableUpdateFlags
	}{
		{
			name:  "defaults when nothing set",
			setup: func() {},
			want:  VariableUpdateFlags{},
		},
		{
			name: "all fields set",
			setup: func() {
				viper.Set("name", "my-workspace")
				viper.Set("key", "MY_VAR")
				viper.Set("value", "new-value")
				viper.Set("value-file", "")
				viper.Set("description", "updated description")
				viper.Set("env", false)
				viper.Set("hcl", true)
				viper.Set("sensitive", false)
			},
			want: VariableUpdateFlags{
				WorkspaceName: "my-workspace",
				Key:           "MY_VAR",
				Value:         "new-value",
				Description:   "updated description",
				HCL:           true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			tt.setup()

			got, err := ParseVariableUpdateFlags(nil)
			if err != nil {
				t.Fatalf("ParseVariableUpdateFlags() error = %v", err)
			}
			if *got != tt.want {
				t.Errorf("ParseVariableUpdateFlags() = %+v, want %+v", *got, tt.want)
			}
		})
	}
}

func TestParseVariableDeleteFlags(t *testing.T) {
	tests := []struct {
		name  string
		setup func()
		want  VariableDeleteFlags
	}{
		{
			name:  "defaults when nothing set",
			setup: func() {},
			want:  VariableDeleteFlags{},
		},
		{
			name: "all fields set",
			setup: func() {
				viper.Set("name", "my-workspace")
				viper.Set("key", "MY_VAR")
			},
			want: VariableDeleteFlags{WorkspaceName: "my-workspace", Key: "MY_VAR"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			tt.setup()

			got, err := ParseVariableDeleteFlags(nil)
			if err != nil {
				t.Fatalf("ParseVariableDeleteFlags() error = %v", err)
			}
			if *got != tt.want {
				t.Errorf("ParseVariableDeleteFlags() = %+v, want %+v", *got, tt.want)
			}
		})
	}
}
