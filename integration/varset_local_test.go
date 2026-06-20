// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

//go:build integration
// +build integration

package integration

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

type varsetFixture struct {
	name       string
	createArgs []string
}

type varsetVariableFixture struct {
	varsetName string
	key        string
	deleteArgs []string
}

// TestVariableSetLocalProfileLifecycle exercises varset and variable CRUD against a
// named TFX profile (local / local.tfe.rocks). Creates several variable sets with
// mixed variable types, lists them, then deletes everything so end state matches start.
//
// Requires TFX_INTEGRATION_PROFILE=local and a matching profile block in ~/.tfx.hcl:
//
//	profile "local" {
//	  hostname     = "local.tfe.rocks"
//	  organization = "your-org"
//	  token        = "..."
//	}
//
// Set TFX_INTEGRATION_NO_CLEANUP=1 to skip delete steps and leave resources for inspection.
func TestVariableSetLocalProfileLifecycle(t *testing.T) {
	setupProfileIntegrationTest(t, localIntegrationProfile)

	prefix := fmt.Sprintf("tfx-varset-%s", uniquePetName())
	keyPrefix := strings.ReplaceAll(prefix, "-", "_")
	profile := localIntegrationProfile

	varsets := []varsetFixture{
		{
			name: prefix + "-org-a",
			createArgs: []string{
				"varset", "create",
				"--name", prefix + "-org-a",
				"--description", "TFx local integration test (org-owned A)",
			},
		},
		{
			name: prefix + "-org-b",
			createArgs: []string{
				"varset", "create",
				"--name", prefix + "-org-b",
				"--description", "TFx local integration test (org-owned B)",
				"--priority",
			},
		},
		{
			name: prefix + "-org-c",
			createArgs: []string{
				"varset", "create",
				"--name", prefix + "-org-c",
				"--description", "TFx local integration test (org-owned C)",
			},
		},
	}

	if projectName := os.Getenv("TEST_PROJECT_NAME"); projectName != "" {
		varsets = append(varsets, varsetFixture{
			name: prefix + "-project",
			createArgs: []string{
				"varset", "create",
				"--name", prefix + "-project",
				"--description", "TFx local integration test (project-owned)",
				"--project-name", projectName,
			},
		})
	}

	if workspaceName := os.Getenv("TEST_WORKSPACE_NAME"); workspaceName != "" {
		varsets = append(varsets, varsetFixture{
			name: prefix + "-workspace",
			createArgs: []string{
				"varset", "create",
				"--name", prefix + "-workspace",
				"--description", "TFx local integration test (workspace-assigned)",
				"--workspace-name", workspaceName,
			},
		})
	}

	variables := []varsetVariableFixture{
		{
			varsetName: prefix + "-org-a",
			key:        keyPrefix + "_tf_plain",
			deleteArgs: []string{"varset", "variable", "delete", "--varset-name", prefix + "-org-a", "--key", keyPrefix + "_tf_plain"},
		},
		{
			varsetName: prefix + "-org-a",
			key:        keyPrefix + "_env_var",
			deleteArgs: []string{"varset", "variable", "delete", "--varset-name", prefix + "-org-a", "--key", keyPrefix + "_env_var"},
		},
		{
			varsetName: prefix + "-org-a",
			key:        keyPrefix + "_hcl_var",
			deleteArgs: []string{"varset", "variable", "delete", "--varset-name", prefix + "-org-a", "--key", keyPrefix + "_hcl_var"},
		},
		{
			varsetName: prefix + "-org-a",
			key:        keyPrefix + "_sensitive",
			deleteArgs: []string{"varset", "variable", "delete", "--varset-name", prefix + "-org-a", "--key", keyPrefix + "_sensitive"},
		},
		{
			varsetName: prefix + "-org-b",
			key:        keyPrefix + "_b_tf",
			deleteArgs: []string{"varset", "variable", "delete", "--varset-name", prefix + "-org-b", "--key", keyPrefix + "_b_tf"},
		},
		{
			varsetName: prefix + "-org-b",
			key:        keyPrefix + "_b_env",
			deleteArgs: []string{"varset", "variable", "delete", "--varset-name", prefix + "-org-b", "--key", keyPrefix + "_b_env"},
		},
		{
			varsetName: prefix + "-org-c",
			key:        keyPrefix + "_c_tf",
			deleteArgs: []string{"varset", "variable", "delete", "--varset-name", prefix + "-org-c", "--key", keyPrefix + "_c_tf"},
		},
	}

	skipCleanup := integrationSkipCleanup()
	if skipCleanup {
		t.Logf("TFX_INTEGRATION_NO_CLEANUP set; resources will be left in place (search prefix %q)", prefix)
	} else {
		t.Cleanup(func() {
			for _, v := range variables {
				if err := executeCommandWithProfile(t, profile, v.deleteArgs); err != nil {
					t.Logf("cleanup: delete variable %s in %s: %v", v.key, v.varsetName, err)
				}
			}
			for i := len(varsets) - 1; i >= 0; i-- {
				vs := varsets[i]
				if err := executeCommandWithProfile(t, profile, []string{
					"varset", "delete", "--name", vs.name,
				}); err != nil {
					t.Logf("cleanup: delete varset %s: %v", vs.name, err)
				}
			}
		})
	}

	t.Run("create variable sets", func(t *testing.T) {
		for _, vs := range varsets {
			vs := vs
			t.Run(vs.name, func(t *testing.T) {
				if err := executeCommandWithProfile(t, profile, vs.createArgs); err != nil {
					t.Fatalf("create variable set %s: %v", vs.name, err)
				}
			})
		}
	})

	t.Run("create variables", func(t *testing.T) {
		creates := []struct {
			name string
			args []string
		}{
			{
				name: "terraform plain",
				args: []string{
					"varset", "variable", "create",
					"--varset-name", prefix + "-org-a",
					"--key", keyPrefix + "_tf_plain",
					"--value", "plain-value",
					"--description", "terraform category",
				},
			},
			{
				name: "environment",
				args: []string{
					"varset", "variable", "create",
					"--varset-name", prefix + "-org-a",
					"--key", keyPrefix + "_env_var",
					"--value", "env-value",
					"--env",
					"--description", "environment category",
				},
			},
			{
				name: "hcl",
				args: []string{
					"varset", "variable", "create",
					"--varset-name", prefix + "-org-a",
					"--key", keyPrefix + "_hcl_var",
					"--value", `{"region"="us-east-1"}`,
					"--hcl",
					"--description", "hcl category",
				},
			},
			{
				name: "sensitive",
				args: []string{
					"varset", "variable", "create",
					"--varset-name", prefix + "-org-a",
					"--key", keyPrefix + "_sensitive",
					"--value", "super-secret",
					"--sensitive",
					"--description", "sensitive variable",
				},
			},
			{
				name: "org-b terraform",
				args: []string{
					"varset", "variable", "create",
					"--varset-name", prefix + "-org-b",
					"--key", keyPrefix + "_b_tf",
					"--value", "beta-tf",
				},
			},
			{
				name: "org-b environment",
				args: []string{
					"varset", "variable", "create",
					"--varset-name", prefix + "-org-b",
					"--key", keyPrefix + "_b_env",
					"--value", "beta-env",
					"--env",
				},
			},
			{
				name: "org-c terraform",
				args: []string{
					"varset", "variable", "create",
					"--varset-name", prefix + "-org-c",
					"--key", keyPrefix + "_c_tf",
					"--value", "gamma-tf",
				},
			},
		}

		for _, tc := range creates {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				if err := executeCommandWithProfile(t, profile, tc.args); err != nil {
					t.Fatalf("create variable (%s): %v", tc.name, err)
				}
			})
		}
	})

	t.Run("list variable sets by search prefix", func(t *testing.T) {
		if err := executeCommandWithProfile(t, profile, []string{
			"varset", "list", "--search", prefix,
		}); err != nil {
			t.Fatalf("list variable sets: %v", err)
		}
	})

	t.Run("list variables in primary varset", func(t *testing.T) {
		if err := executeCommandWithProfile(t, profile, []string{
			"varset", "variable", "list", "--varset-name", prefix + "-org-a",
		}); err != nil {
			t.Fatalf("list variables: %v", err)
		}
	})

	t.Run("show variable set by name", func(t *testing.T) {
		if err := executeCommandWithProfile(t, profile, []string{
			"varset", "show", "--name", prefix + "-org-a",
		}); err != nil {
			t.Fatalf("show variable set: %v", err)
		}
	})

	if skipCleanup {
		t.Run("list retained varsets", func(t *testing.T) {
			for _, vs := range varsets {
				t.Logf("retained varset: %s", vs.name)
			}
			if err := executeCommandWithProfile(t, profile, []string{
				"varset", "list", "--search", prefix,
			}); err != nil {
				t.Fatalf("list retained varsets: %v", err)
			}
		})
	} else {
		t.Run("delete variables", func(t *testing.T) {
			for _, v := range variables {
				v := v
				t.Run(v.key, func(t *testing.T) {
					if err := executeCommandWithProfile(t, profile, v.deleteArgs); err != nil {
						t.Fatalf("delete variable %s: %v", v.key, err)
					}
				})
			}
			// Clear cleanup list so t.Cleanup does not retry successful deletes.
			variables = nil
		})

		t.Run("delete variable sets", func(t *testing.T) {
			for _, vs := range varsets {
				vs := vs
				t.Run(vs.name, func(t *testing.T) {
					if err := executeCommandWithProfile(t, profile, []string{
						"varset", "delete", "--name", vs.name,
					}); err != nil {
						t.Fatalf("delete variable set %s: %v", vs.name, err)
					}
				})
			}
			varsets = nil
		})

		t.Run("verify no leftover varsets for prefix", func(t *testing.T) {
			if err := executeCommandWithProfile(t, profile, []string{
				"varset", "list", "--search", prefix,
			}); err != nil {
				t.Fatalf("list after cleanup: %v", err)
			}
		})
	}
}
