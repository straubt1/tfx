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
	varsetListCmd.Flags().BoolP("all", "a", false, "List variable sets across all organizations (optional).")
	addVarsetScopeFlags(varsetListCmd)
	varsetListCmd.MarkFlagsMutuallyExclusive("all", "project-name", "workspace-name")

	// `tfx variable-set show` flags
	varsetShowCmd.Flags().StringP("id", "i", "", "ID of the Variable Set.")
	varsetShowCmd.Flags().StringP("name", "n", "", "Name of the Variable Set.")
	addVarsetScopeFlags(varsetShowCmd)
	varsetShowCmd.MarkFlagsMutuallyExclusive("id", "name")
	varsetShowCmd.MarkFlagsOneRequired("id", "name")

	// `tfx variable-set create` flags
	varsetCreateCmd.Flags().StringP("name", "n", "", "Name of the Variable Set (required).")
	varsetCreateCmd.Flags().StringP("description", "d", "", "Description of the Variable Set (optional).")
	varsetCreateCmd.Flags().Bool("global", false, "Apply this Variable Set to all workspaces in the organization (optional).")
	varsetCreateCmd.Flags().Bool("priority", false, "Variable values in this set override workspace-level values (optional).")
	varsetCreateCmd.Flags().String("organization-name", "", "Organization name (optional, defaults to configured organization).")
	varsetCreateCmd.Flags().String("project-name", "", "Create as a project-owned variable set (optional).")
	varsetCreateCmd.Flags().String("workspace-name", "", "Apply the variable set to this workspace after creation (optional).")
	varsetCreateCmd.MarkFlagRequired("name")
	varsetCreateCmd.MarkFlagsMutuallyExclusive("global", "project-name")

	// `tfx variable-set delete` flags
	varsetDeleteCmd.Flags().StringP("id", "i", "", "ID of the Variable Set to delete.")
	varsetDeleteCmd.Flags().StringP("name", "n", "", "Name of the Variable Set to delete.")
	addVarsetScopeFlags(varsetDeleteCmd)
	varsetDeleteCmd.MarkFlagsMutuallyExclusive("id", "name")
	varsetDeleteCmd.MarkFlagsOneRequired("id", "name")

	rootCmd.AddCommand(varsetCmd)
	varsetCmd.AddCommand(varsetListCmd)
	varsetCmd.AddCommand(varsetShowCmd)
	varsetCmd.AddCommand(varsetCreateCmd)
	varsetCmd.AddCommand(varsetDeleteCmd)
}

func addVarsetScopeFlags(cmd *cobra.Command) {
	cmd.Flags().String("organization-name", "", "Organization name (optional, defaults to configured organization).")
	cmd.Flags().String("project-name", "", "Filter or scope to a project (optional).")
	cmd.Flags().String("workspace-name", "", "Filter or scope to a workspace (optional).")
}

func variableSetScopeFromFlags(f flags.VariableSetScopeFlags) data.VariableSetScope {
	return data.VariableSetScope{
		All:              f.All,
		OrganizationName: f.OrganizationName,
		ProjectName:      f.ProjectName,
		WorkspaceName:    f.WorkspaceName,
		Search:           f.Search,
	}
}

func varsetScopeFromListFlags(f *flags.VariableSetListFlags) data.VariableSetScope {
	return variableSetScopeFromFlags(f.VariableSetScopeFlags)
}

func varsetListHeader(v *view.VariableSetListView, c *client.TfxClient, f *flags.VariableSetListFlags) {
	if f.All {
		if f.Search != "" {
			v.PrintCommandHeader("Listing variable sets across all organizations matching '%s'", f.Search)
		} else {
			v.PrintCommandHeader("Listing variable sets across all organizations")
		}
		return
	}

	org := flags.ResolveVarsetOrg(f.OrganizationName, c.OrganizationName)

	if f.WorkspaceName != "" {
		if f.Search != "" {
			v.PrintCommandHeader("Listing variable sets for workspace '%s' in organization '%s' matching '%s'", f.WorkspaceName, org, f.Search)
		} else {
			v.PrintCommandHeader("Listing variable sets for workspace '%s' in organization '%s'", f.WorkspaceName, org)
		}
		return
	}

	if f.ProjectName != "" {
		if f.Search != "" {
			v.PrintCommandHeader("Listing variable sets for project '%s' in organization '%s' matching '%s'", f.ProjectName, org, f.Search)
		} else {
			v.PrintCommandHeader("Listing variable sets for project '%s' in organization '%s'", f.ProjectName, org)
		}
		return
	}

	if f.Search != "" {
		v.PrintCommandHeader("Listing variable sets in organization '%s' matching '%s'", org, f.Search)
	} else {
		v.PrintCommandHeader("Listing variable sets in organization '%s'", org)
	}
}

func variableSetList(cmdConfig *flags.VariableSetListFlags) error {
	v := view.NewVariableSetListView()
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	varsetListHeader(v, c, cmdConfig)

	scope := varsetScopeFromListFlags(cmdConfig)
	items, err := data.ListVariableSetsWithScope(c, c.OrganizationName, scope)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to list variable sets"))
	}
	return v.Render(items, cmdConfig.All)
}

func variableSetShow(cmdConfig *flags.VariableSetShowFlags) error {
	v := view.NewVariableSetShowView()
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	scope := variableSetScopeFromFlags(cmdConfig.VariableSetScopeFlags)
	if cmdConfig.ID != "" {
		v.PrintCommandHeader("Showing variable set '%s'", cmdConfig.ID)
	} else {
		v.PrintCommandHeader("Showing variable set '%s'", cmdConfig.Name)
	}

	vs, err := data.ResolveVariableSet(c, c.OrganizationName, scope, cmdConfig.Name, cmdConfig.ID)
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

	orgName := flags.ResolveVarsetOrg(cmdConfig.OrganizationName, c.OrganizationName)
	v.PrintCommandHeader("Creating variable set '%s' in organization '%s'", cmdConfig.Name, orgName)

	params := data.VariableSetCreateParams{
		Name:          cmdConfig.Name,
		Description:   cmdConfig.Description,
		Global:        cmdConfig.Global,
		Priority:      cmdConfig.Priority,
		ProjectName:   cmdConfig.ProjectName,
		WorkspaceName: cmdConfig.WorkspaceName,
	}
	vs, err := data.CreateVariableSet(c, orgName, params)
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

	scope := variableSetScopeFromFlags(cmdConfig.VariableSetScopeFlags)
	label := cmdConfig.ID
	if label == "" {
		label = cmdConfig.Name
	}
	v.PrintCommandHeader("Deleting variable set '%s'", label)

	vs, err := data.ResolveVariableSet(c, c.OrganizationName, scope, cmdConfig.Name, cmdConfig.ID)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to resolve variable set"))
	}

	err = data.DeleteVariableSet(c, vs.ID)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to delete variable set"))
	}
	return v.Render(vs.ID)
}
