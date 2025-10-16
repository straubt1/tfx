// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/cmd/flags"
	view "github.com/straubt1/tfx/cmd/views"
	"github.com/straubt1/tfx/data"
)

var (
	// `tfx variable` commands
	variableCmd = &cobra.Command{
		Aliases: []string{"var"},
		Use:     "variable",
		Short:   "Variable Commands",
		Long:    "Commands to work with Workspace Variables.",
	}

	// `tfx variable list` command
	variableListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Variables",
		Long:  "List Variables in a Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseVariableListFlags(cmd)
			if err != nil {
				return err
			}
			return variableList(cmdConfig)
		},
	}

	// `tfx variable create` command
	variableCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a Variable",
		Long:  "Create a Variable in a Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseVariableCreateFlags(cmd)
			if err != nil {
				return err
			}

			return variableCreate(cmdConfig)
		},
	}

	// `tfx variable update` command
	variableUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update a Variable",
		Long:  "Update a Variable in a Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseVariableUpdateFlags(cmd)
			if err != nil {
				return err
			}

			return variableUpdate(cmdConfig)
		},
	}

	// `tfx variable show` command
	variableShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show details of a Variable",
		Long:  "Show details of a Variable in a Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseVariableShowFlags(cmd)
			if err != nil {
				return err
			}
			return variableShow(cmdConfig)
		},
	}

	// `tfx variable delete` command
	variableDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a Variable",
		Long:  "Delete a Variable in a Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseVariableDeleteFlags(cmd)
			if err != nil {
				return err
			}
			return variableDelete(cmdConfig)
		},
	}
)

func init() {
	// `tfx variable list` command
	variableListCmd.Flags().StringP("workspace-name", "w", "", "Name of the Workspace")
	variableListCmd.MarkFlagRequired("workspace-name")

	// `tfx variable create` command
	variableCreateCmd.Flags().StringP("workspace-name", "w", "", "Name of the Workspace")
	variableCreateCmd.Flags().StringP("key", "k", "", "Key of the Variable")
	variableCreateCmd.Flags().StringP("value", "v", "", "Value of the Variable (value or valueFile must be set)")
	variableCreateCmd.Flags().StringP("value-file", "f", "", "Path to a variable text file, the contents of the file will be used (value or valueFile must be set)")
	variableCreateCmd.Flags().StringP("description", "d", "", "Description of the Variable (optional)")
	variableCreateCmd.Flags().BoolP("env", "e", false, "Variable is an Environment Variable (optional, defaults to false)")
	variableCreateCmd.Flags().BoolP("hcl", "", false, "Value of Variable is HCL (optional, defaults to false)")
	variableCreateCmd.Flags().BoolP("sensitive", "s", false, "Variable is Sensitive (optional, defaults to false)")
	variableCreateCmd.MarkFlagRequired("workspace-name")
	variableCreateCmd.MarkFlagRequired("key")
	variableCreateCmd.MarkFlagsMutuallyExclusive("value", "value-file")
	variableCreateCmd.MarkFlagsOneRequired("value", "value-file")

	// `tfx variable update` command
	variableUpdateCmd.Flags().StringP("workspace-name", "w", "", "Name of the Workspace")
	variableUpdateCmd.Flags().StringP("key", "k", "", "Key of the Variable")
	variableUpdateCmd.Flags().StringP("value", "v", "", "Value of the Variable (value or valueFile must be set)")
	variableUpdateCmd.Flags().StringP("value-file", "f", "", "Path to a variable text file, the contents of the file will be used (value or valueFile must be set)")
	variableUpdateCmd.Flags().StringP("description", "d", "", "Description of the Variable (optional)")
	variableUpdateCmd.Flags().BoolP("env", "e", false, "Variable is an Environment Variable (optional, defaults to false)")
	variableUpdateCmd.Flags().BoolP("hcl", "", false, "Value of Variable is HCL (optional, defaults to false)")
	variableUpdateCmd.Flags().BoolP("sensitive", "s", false, "Variable is Sensitive (optional, defaults to false)")
	variableUpdateCmd.MarkFlagRequired("workspace-name")
	variableUpdateCmd.MarkFlagRequired("key")
	variableUpdateCmd.MarkFlagsMutuallyExclusive("value", "value-file")
	variableUpdateCmd.MarkFlagsOneRequired("value", "value-file")

	// `tfx variable show` command
	variableShowCmd.Flags().StringP("workspace-name", "w", "", "Name of the Workspace")
	variableShowCmd.Flags().StringP("key", "k", "", "Key of the Variable")
	variableShowCmd.MarkFlagRequired("workspace-name")
	variableShowCmd.MarkFlagRequired("key")

	// `tfx variable delete` command
	variableDeleteCmd.Flags().StringP("workspace-name", "w", "", "Name of the Workspace")
	variableDeleteCmd.Flags().StringP("key", "k", "", "Key of the Variable")
	variableDeleteCmd.MarkFlagRequired("workspace-name")
	variableDeleteCmd.MarkFlagRequired("key")

	workspaceCmd.AddCommand(variableCmd)
	variableCmd.AddCommand(variableListCmd)
	variableCmd.AddCommand(variableCreateCmd)
	variableCmd.AddCommand(variableUpdateCmd)
	variableCmd.AddCommand(variableShowCmd)
	variableCmd.AddCommand(variableDeleteCmd)
}

func variableList(cmdConfig *flags.VariableListFlags) error {
	// Create view for rendering
	v := view.NewVariableListView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Listing variables for workspace '%s'", cmdConfig.WorkspaceName)

	workspaceID, err := data.GetWorkspaceID(c, c.OrganizationName, cmdConfig.WorkspaceName)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "unable to read workspace id"))
	}

	variables, err := data.FetchVariables(c, workspaceID)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to list variables"))
	}

	return v.Render(cmdConfig.WorkspaceName, variables)
}

func variableCreate(cmdConfig *flags.VariableCreateFlags) error {
	// Create view for rendering
	v := view.NewVariableCreateView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Creating variable '%s' for workspace '%s'", cmdConfig.Key, cmdConfig.WorkspaceName)

	workspaceID, err := data.GetWorkspaceID(c, c.OrganizationName, cmdConfig.WorkspaceName)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "unable to read workspace id"))
	}

	// Handle value from file if specified
	value := cmdConfig.Value
	if cmdConfig.ValueFile != "" {
		if !isFile(cmdConfig.ValueFile) {
			return v.RenderError(errors.New("valueFile does not exist"))
		}
		v.PrintCommandFilter("Variable filename contents will be used: %s", cmdConfig.ValueFile)
		value, err = readFile(cmdConfig.ValueFile)
		if err != nil {
			return v.RenderError(errors.Wrap(err, "unable to read the file passed"))
		}
	}

	// Determine category
	var category *tfe.CategoryType
	if cmdConfig.Env {
		category = tfe.Category(tfe.CategoryEnv)
	} else {
		category = tfe.Category(tfe.CategoryTerraform)
	}

	// Create the variable
	variable, err := data.CreateVariable(c, workspaceID, tfe.VariableCreateOptions{
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

func variableUpdate(cmdConfig *flags.VariableUpdateFlags) error {
	// Create view for rendering
	v := view.NewVariableUpdateView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Updating variable '%s' for workspace '%s'", cmdConfig.Key, cmdConfig.WorkspaceName)

	workspaceID, err := data.GetWorkspaceID(c, c.OrganizationName, cmdConfig.WorkspaceName)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "unable to read workspace id"))
	}

	variableID, err := data.GetVariableID(c, workspaceID, cmdConfig.Key)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "unable to read variable id"))
	}

	// Handle value from file if specified
	value := cmdConfig.Value
	if cmdConfig.ValueFile != "" {
		if !isFile(cmdConfig.ValueFile) {
			return v.RenderError(errors.New("valueFile does not exist"))
		}
		v.PrintCommandFilter("Variable filename contents will be used: %s", cmdConfig.ValueFile)
		value, err = readFile(cmdConfig.ValueFile)
		if err != nil {
			return v.RenderError(errors.Wrap(err, "unable to read the file passed"))
		}
	}

	// Determine category
	var category *tfe.CategoryType
	if cmdConfig.Env {
		category = tfe.Category(tfe.CategoryEnv)
	} else {
		category = tfe.Category(tfe.CategoryTerraform)
	}

	// Update the variable
	variable, err := data.UpdateVariable(c, workspaceID, variableID, tfe.VariableUpdateOptions{
		Key:         &cmdConfig.Key,
		Value:       &value,
		Description: &cmdConfig.Description,
		Category:    category,
		HCL:         &cmdConfig.HCL,
		Sensitive:   &cmdConfig.Sensitive,
	})
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to update variable"))
	}

	return v.Render(variable)
}

func variableShow(cmdConfig *flags.VariableShowFlags) error {
	// Create view for rendering
	v := view.NewVariableShowView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Showing variable '%s' for workspace '%s'", cmdConfig.Key, cmdConfig.WorkspaceName)

	workspaceID, err := data.GetWorkspaceID(c, c.OrganizationName, cmdConfig.WorkspaceName)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "unable to read workspace id"))
	}

	variable, err := data.FetchVariable(c, workspaceID, cmdConfig.Key)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "unable to read variable"))
	}

	return v.Render(variable)
}

func variableDelete(cmdConfig *flags.VariableDeleteFlags) error {
	// Create view for rendering
	v := view.NewVariableDeleteView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	v.PrintCommandHeader("Deleting variable '%s' for workspace '%s'", cmdConfig.Key, cmdConfig.WorkspaceName)

	workspaceID, err := data.GetWorkspaceID(c, c.OrganizationName, cmdConfig.WorkspaceName)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "unable to read workspace id"))
	}

	variableID, err := data.GetVariableID(c, workspaceID, cmdConfig.Key)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "unable to read variable id"))
	}

	err = data.DeleteVariable(c, workspaceID, variableID)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to delete variable"))
	}

	return v.Render(cmdConfig.Key)
}
