// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type StateListFlags struct {
	WorkspaceName string
	MaxItems      int
}

type StateCreateFlags struct {
	WorkspaceName string
	Filename      string
}

type StateShowFlags struct {
	StateID string
}

type StateDownloadFlags struct {
	StateID   string
	Directory string
	Filename  string
}

func ParseStateListFlags(cmd *cobra.Command) (*StateListFlags, error) {
	return &StateListFlags{
		WorkspaceName: viper.GetString("name"),
		MaxItems:      viper.GetInt("max-items"),
	}, nil
}

func ParseStateCreateFlags(cmd *cobra.Command) (*StateCreateFlags, error) {
	return &StateCreateFlags{
		WorkspaceName: viper.GetString("name"),
		Filename:      viper.GetString("filename"),
	}, nil
}

func ParseStateShowFlags(cmd *cobra.Command) (*StateShowFlags, error) {
	return &StateShowFlags{StateID: viper.GetString("state-id")}, nil
}

func ParseStateDownloadFlags(cmd *cobra.Command) (*StateDownloadFlags, error) {
	return &StateDownloadFlags{
		StateID:   viper.GetString("state-id"),
		Directory: viper.GetString("directory"),
		Filename:  viper.GetString("filename"),
	}, nil
}
