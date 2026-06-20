// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package cmd

import (
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/cmd/flags"
	view "github.com/straubt1/tfx/cmd/views"
	"github.com/straubt1/tfx/data"
	pkgfile "github.com/straubt1/tfx/pkg/file"
)

var (
	// `tfx variable-set variable` command
	varsetVariableCmd = &cobra.Command{
		Use:     "variable",
		Aliases: []string{"var"},
		Short:   "Variable Set Variable Commands",
		Long:    "Commands to work with Variables in a Variable Set.",
	}

	// `tfx variable-set variable list` command
	varsetVariableListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Variables",
		Long:  "List Variables in a Variable Set.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseVariableSetVariableListFlags(cmd)
			if err != nil {
				return err
			}
			return variableSetVariableList(cmdConfig)
		},
	}

	// `tfx variable-set variable create` command
	varsetVariableCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a Variable",
		Long:  "Create a Variable in a Variable Set.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseVariableSetVariableCreateFlags(cmd)
			if err != nil {
				return err
			}
			return variableSetVariableCreate(cmdConfig)
		},
	}

	// `tfx variable-set variable update` command
	varsetVariableUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update a Variable",
		Long:  "Update a Variable in a Variable Set.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseVariableSetVariableUpdateFlags(cmd)
			if err != nil {
				return err
			}
			return variableSetVariableUpdate(cmdConfig)
		},
	}

	// `tfx variable-set variable show` command
	varsetVariableShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show details of a Variable",
		Long:  "Show details of a Variable in a Variable Set.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseVariableSetVariableShowFlags(cmd)
			if err != nil {
				return err
			}
			return variableSetVariableShow(cmdConfig)
		},
	}

	// `tfx variable-set variable delete` command
	varsetVariableDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a Variable",
		Long:  "Delete a Variable in a Variable Set.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseVariableSetVariableDeleteFlags(cmd)
			if err != nil {
				return err
			}
			return variableSetVariableDelete(cmdConfig)
		},
	}
)

func init() {
	addVarsetVariableCommonFlags(varsetVariableListCmd)
	varsetVariableListCmd.MarkFlagsOneRequired("varset-id", "varset-name")

	addVarsetVariableCommonFlags(varsetVariableCreateCmd)
	varsetVariableCreateCmd.Flags().StringP("key", "k", "", "Key of the Variable")
	varsetVariableCreateCmd.Flags().StringP("value", "v", "", "Value of the Variable (value or value-file must be set)")
	varsetVariableCreateCmd.Flags().StringP("value-file", "f", "", "Path to a variable text file, the contents of the file will be used (value or value-file must be set)")
	varsetVariableCreateCmd.Flags().StringP("description", "d", "", "Description of the Variable (optional)")
	varsetVariableCreateCmd.Flags().BoolP("env", "e", false, "Variable is an Environment Variable (optional, defaults to false)")
	varsetVariableCreateCmd.Flags().BoolP("hcl", "", false, "Value of Variable is HCL (optional, defaults to false)")
	varsetVariableCreateCmd.Flags().BoolP("sensitive", "s", false, "Variable is Sensitive (optional, defaults to false)")
	varsetVariableCreateCmd.MarkFlagsOneRequired("varset-id", "varset-name")
	varsetVariableCreateCmd.MarkFlagRequired("key")
	varsetVariableCreateCmd.MarkFlagsMutuallyExclusive("value", "value-file")
	varsetVariableCreateCmd.MarkFlagsOneRequired("value", "value-file")

	addVarsetVariableCommonFlags(varsetVariableUpdateCmd)
	varsetVariableUpdateCmd.Flags().StringP("key", "k", "", "Key of the Variable")
	varsetVariableUpdateCmd.Flags().StringP("value", "v", "", "Value of the Variable (value or value-file must be set)")
	varsetVariableUpdateCmd.Flags().StringP("value-file", "f", "", "Path to a variable text file, the contents of the file will be used (value or value-file must be set)")
	varsetVariableUpdateCmd.Flags().StringP("description", "d", "", "Description of the Variable (optional)")
	varsetVariableUpdateCmd.Flags().BoolP("env", "e", false, "Variable is an Environment Variable (optional, defaults to false)")
	varsetVariableUpdateCmd.Flags().BoolP("hcl", "", false, "Value of Variable is HCL (optional, defaults to false)")
	varsetVariableUpdateCmd.Flags().BoolP("sensitive", "s", false, "Variable is Sensitive (optional, defaults to false)")
	varsetVariableUpdateCmd.MarkFlagsOneRequired("varset-id", "varset-name")
	varsetVariableUpdateCmd.MarkFlagRequired("key")
	varsetVariableUpdateCmd.MarkFlagsMutuallyExclusive("value", "value-file")
	varsetVariableUpdateCmd.MarkFlagsOneRequired("value", "value-file")

	addVarsetVariableCommonFlags(varsetVariableShowCmd)
	varsetVariableShowCmd.Flags().StringP("key", "k", "", "Key of the Variable")
	varsetVariableShowCmd.MarkFlagsOneRequired("varset-id", "varset-name")
	varsetVariableShowCmd.MarkFlagRequired("key")

	addVarsetVariableCommonFlags(varsetVariableDeleteCmd)
	varsetVariableDeleteCmd.Flags().StringP("key", "k", "", "Key of the Variable")
	varsetVariableDeleteCmd.MarkFlagsOneRequired("varset-id", "varset-name")
	varsetVariableDeleteCmd.MarkFlagRequired("key")

	varsetCmd.AddCommand(varsetVariableCmd)
	varsetVariableCmd.AddCommand(varsetVariableListCmd)
	varsetVariableCmd.AddCommand(varsetVariableCreateCmd)
	varsetVariableCmd.AddCommand(varsetVariableUpdateCmd)
	varsetVariableCmd.AddCommand(varsetVariableShowCmd)
	varsetVariableCmd.AddCommand(varsetVariableDeleteCmd)
}

func addVarsetVariableCommonFlags(cmd *cobra.Command) {
	cmd.Flags().String("varset-id", "", "ID of the Variable Set")
	cmd.Flags().String("varset-name", "", "Name of the Variable Set")
	cmd.Flags().String("organization-name", "", "Organization name (optional, defaults to configured organization).")
	cmd.Flags().String("project-name", "", "Scope for resolving the variable set by name (optional).")
	cmd.Flags().String("workspace-name", "", "Scope for resolving the variable set by name (optional).")
	cmd.MarkFlagsMutuallyExclusive("varset-id", "varset-name")
}

func resolveVariableSetForVariableCmd(c *client.TfxClient, varsetName, varsetID string, scope flags.VariableSetScopeFlags) (*tfe.VariableSet, error) {
	return data.ResolveVariableSet(c, c.OrganizationName, variableSetScopeFromFlags(scope), varsetName, varsetID)
}

func variableSetVariableList(cmdConfig *flags.VariableSetVariableListFlags) error {
	v := view.NewVariableSetVariableListView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	label := cmdConfig.VarsetID
	if label == "" {
		label = cmdConfig.VarsetName
	}
	v.PrintCommandHeader("Listing variables for variable set '%s'", label)

	vs, err := resolveVariableSetForVariableCmd(c, cmdConfig.VarsetName, cmdConfig.VarsetID, cmdConfig.VariableSetScopeFlags)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to resolve variable set"))
	}

	variables, err := data.FetchVariableSetVariables(c, vs.ID)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to list variables"))
	}

	return v.Render(variables)
}

func variableSetVariableCreate(cmdConfig *flags.VariableSetVariableCreateFlags) error {
	v := view.NewVariableSetVariableCreateView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	label := cmdConfig.VarsetID
	if label == "" {
		label = cmdConfig.VarsetName
	}
	v.PrintCommandHeader("Creating variable '%s' for variable set '%s'", cmdConfig.Key, label)

	vs, err := resolveVariableSetForVariableCmd(c, cmdConfig.VarsetName, cmdConfig.VarsetID, cmdConfig.VariableSetScopeFlags)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to resolve variable set"))
	}

	value := cmdConfig.Value
	if cmdConfig.ValueFile != "" {
		if !pkgfile.IsFile(cmdConfig.ValueFile) {
			return v.RenderError(errors.New("valueFile does not exist"))
		}
		v.PrintCommandFilter("Variable filename contents will be used: %s", cmdConfig.ValueFile)
		value, err = pkgfile.ReadFile(cmdConfig.ValueFile)
		if err != nil {
			return v.RenderError(errors.Wrap(err, "unable to read the file passed"))
		}
	}

	var category *tfe.CategoryType
	if cmdConfig.Env {
		category = tfe.Category(tfe.CategoryEnv)
	} else {
		category = tfe.Category(tfe.CategoryTerraform)
	}

	variable, err := data.CreateVariableSetVariable(c, vs.ID, tfe.VariableSetVariableCreateOptions{
		Key:         &cmdConfig.Key,
		Value:       &value,
		Description: &cmdConfig.Description,
		Category:    category,
		HCL:         &cmdConfig.HCL,
		Sensitive:   &cmdConfig.Sensitive,
	})
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to create variable"))
	}

	return v.Render(variable)
}

func variableSetVariableUpdate(cmdConfig *flags.VariableSetVariableUpdateFlags) error {
	v := view.NewVariableSetVariableUpdateView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	label := cmdConfig.VarsetID
	if label == "" {
		label = cmdConfig.VarsetName
	}
	v.PrintCommandHeader("Updating variable '%s' for variable set '%s'", cmdConfig.Key, label)

	vs, err := resolveVariableSetForVariableCmd(c, cmdConfig.VarsetName, cmdConfig.VarsetID, cmdConfig.VariableSetScopeFlags)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to resolve variable set"))
	}

	variableID, err := data.GetVariableSetVariableID(c, vs.ID, cmdConfig.Key)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "unable to read variable id"))
	}

	value := cmdConfig.Value
	if cmdConfig.ValueFile != "" {
		if !pkgfile.IsFile(cmdConfig.ValueFile) {
			return v.RenderError(errors.New("valueFile does not exist"))
		}
		v.PrintCommandFilter("Variable filename contents will be used: %s", cmdConfig.ValueFile)
		value, err = pkgfile.ReadFile(cmdConfig.ValueFile)
		if err != nil {
			return v.RenderError(errors.Wrap(err, "unable to read the file passed"))
		}
	}

	variable, err := data.UpdateVariableSetVariable(c, vs.ID, variableID, tfe.VariableSetVariableUpdateOptions{
		Key:         &cmdConfig.Key,
		Value:       &value,
		Description: &cmdConfig.Description,
		HCL:         &cmdConfig.HCL,
		Sensitive:   &cmdConfig.Sensitive,
	})
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to update variable"))
	}

	return v.Render(variable)
}

func variableSetVariableShow(cmdConfig *flags.VariableSetVariableShowFlags) error {
	v := view.NewVariableSetVariableShowView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	label := cmdConfig.VarsetID
	if label == "" {
		label = cmdConfig.VarsetName
	}
	v.PrintCommandHeader("Showing variable '%s' for variable set '%s'", cmdConfig.Key, label)

	vs, err := resolveVariableSetForVariableCmd(c, cmdConfig.VarsetName, cmdConfig.VarsetID, cmdConfig.VariableSetScopeFlags)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to resolve variable set"))
	}

	variable, err := data.FetchVariableSetVariable(c, vs.ID, cmdConfig.Key)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "unable to read variable"))
	}

	return v.Render(variable)
}

func variableSetVariableDelete(cmdConfig *flags.VariableSetVariableDeleteFlags) error {
	v := view.NewVariableSetVariableDeleteView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	label := cmdConfig.VarsetID
	if label == "" {
		label = cmdConfig.VarsetName
	}
	v.PrintCommandHeader("Deleting variable '%s' for variable set '%s'", cmdConfig.Key, label)

	vs, err := resolveVariableSetForVariableCmd(c, cmdConfig.VarsetName, cmdConfig.VarsetID, cmdConfig.VariableSetScopeFlags)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to resolve variable set"))
	}

	variableID, err := data.GetVariableSetVariableID(c, vs.ID, cmdConfig.Key)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "unable to read variable id"))
	}

	err = data.DeleteVariableSetVariable(c, vs.ID, variableID)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to delete variable"))
	}

	return v.Render(cmdConfig.Key)
}
