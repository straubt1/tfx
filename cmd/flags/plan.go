// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// PlanShowFlags holds flags for show plan
type PlanShowFlags struct {
	ID string
}

// PlanLogsFlags holds flags for logs plan
type PlanLogsFlags struct {
	ID string
}

// PlanJSONOutputFlags holds flags for jsonoutput plan
type PlanJSONOutputFlags struct {
	ID string
}

// PlanCreateFlags holds flags for create plan
type PlanCreateFlags struct {
	WorkspaceName   string
	Directory       string
	ConfigurationID string
	Message         string
	Speculative     bool
	Destroy         bool
	EnvVars         map[string]string
}

func ParsePlanShowFlags(cmd *cobra.Command) (*PlanShowFlags, error) {
	return &PlanShowFlags{ID: viper.GetString("id")}, nil
}

func ParsePlanLogsFlags(cmd *cobra.Command) (*PlanLogsFlags, error) {
	return &PlanLogsFlags{ID: viper.GetString("id")}, nil
}

func ParsePlanJSONOutputFlags(cmd *cobra.Command) (*PlanJSONOutputFlags, error) {
	return &PlanJSONOutputFlags{ID: viper.GetString("id")}, nil
}

func ParsePlanCreateFlags(cmd *cobra.Command) (*PlanCreateFlags, error) {
	return &PlanCreateFlags{
		WorkspaceName:   viper.GetString("name"),
		Directory:       viper.GetString("directory"),
		ConfigurationID: viper.GetString("configuration-id"),
		Message:         viper.GetString("message"),
		Speculative:     viper.GetBool("speculative"),
		Destroy:         viper.GetBool("destroy"),
		EnvVars:         make(map[string]string), // Will be populated from the env slice in the command
	}, nil
}
