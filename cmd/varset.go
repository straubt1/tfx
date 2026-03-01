// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/cmd/flags"
	view "github.com/straubt1/tfx/cmd/views"
	"github.com/straubt1/tfx/data"
)

var (
	// `tfx variable-set` command
	varsetCmd = &cobra.Command{
		Use:     "variable-set",
		Aliases: []string{"varset"},
		Short:   "Variable Set Commands",
		Long:    "Work with TFx Variable Sets.",
	}

	// `tfx variable-set list` command
	varsetListCmd = &cobra.Command{
		Use:   "list",
		Short: "List variable sets",
		Long:  "List Variable Sets in a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseVariableSetListFlags(cmd)
			if err != nil {
				return err
			}
			return variableSetList(cmdConfig)
		},
	}

	// `tfx variable-set show` command
	varsetShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show variable set details",
		Long:  "Show details of a Variable Set, including assigned workspaces, projects, and variables.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseVariableSetShowFlags(cmd)
			if err != nil {
				return err
			}
			return variableSetShow(cmdConfig)
		},
	}

	// `tfx variable-set create` command
	varsetCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a variable set",
		Long:  "Create a new Variable Set in a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseVariableSetCreateFlags(cmd)
			if err != nil {
				return err
			}
			return variableSetCreate(cmdConfig)
		},
	}

	// `tfx variable-set delete` command
	varsetDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a variable set",
		Long:  "Delete a Variable Set from a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseVariableSetDeleteFlags(cmd)
			if err != nil {
				return err
			}
			return variableSetDelete(cmdConfig)
		},
	}
)

func init() {
	// `tfx variable-set list` flags
	varsetListCmd.Flags().StringP("search", "s", "", "Search string to filter variable sets by name (optional).")

	// `tfx variable-set show` flags
	varsetShowCmd.Flags().StringP("id", "i", "", "ID of the Variable Set (required).")
	varsetShowCmd.MarkFlagRequired("id")

	// `tfx variable-set create` flags
	varsetCreateCmd.Flags().StringP("name", "n", "", "Name of the Variable Set (required).")
	varsetCreateCmd.Flags().StringP("description", "d", "", "Description of the Variable Set (optional).")
	varsetCreateCmd.Flags().Bool("global", false, "Apply this Variable Set to all workspaces in the organization (optional).")
	varsetCreateCmd.Flags().Bool("priority", false, "Variable values in this set override workspace-level values (optional).")
	varsetCreateCmd.MarkFlagRequired("name")

	// `tfx variable-set delete` flags
	varsetDeleteCmd.Flags().StringP("id", "i", "", "ID of the Variable Set to delete (required).")
	varsetDeleteCmd.MarkFlagRequired("id")

	rootCmd.AddCommand(varsetCmd)
	varsetCmd.AddCommand(varsetListCmd)
	varsetCmd.AddCommand(varsetShowCmd)
	varsetCmd.AddCommand(varsetCreateCmd)
	varsetCmd.AddCommand(varsetDeleteCmd)
}

func variableSetList(cmdConfig *flags.VariableSetListFlags) error {
	v := view.NewVariableSetListView()
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}
	v.PrintCommandHeader("List Variable Sets for Organization: %s", c.OrganizationName)
	items, err := data.ListVariableSets(c, c.OrganizationName, cmdConfig.Search)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to list variable sets"))
	}
	return v.Render(items)
}

func variableSetShow(cmdConfig *flags.VariableSetShowFlags) error {
	v := view.NewVariableSetShowView()
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}
	v.PrintCommandHeader("Show Variable Set: %s", cmdConfig.ID)
	vs, err := data.ReadVariableSet(c, cmdConfig.ID)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to read variable set"))
	}
	return v.Render(vs)
}

func variableSetCreate(cmdConfig *flags.VariableSetCreateFlags) error {
	v := view.NewVariableSetCreateView()
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}
	v.PrintCommandHeader("Create Variable Set for Organization: %s", c.OrganizationName)
	vs, err := data.CreateVariableSet(c, c.OrganizationName, cmdConfig.Name, cmdConfig.Description, cmdConfig.Global, cmdConfig.Priority)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to create variable set"))
	}
	return v.Render(vs)
}

func variableSetDelete(cmdConfig *flags.VariableSetDeleteFlags) error {
	v := view.NewVariableSetDeleteView()
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}
	v.PrintCommandHeader("Delete Variable Set: %s", cmdConfig.ID)
	err = data.DeleteVariableSet(c, cmdConfig.ID)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to delete variable set"))
	}
	return v.Render(cmdConfig.ID)
}
