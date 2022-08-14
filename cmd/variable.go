// Copyright © 2021 Tom Straub <github.com/straubt1>

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package cmd

import (
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// `tfx variable` commands
	variableCmd = &cobra.Command{
		Use:   "variable",
		Short: "Variable Commands",
		Long:  "Commands to work with Workspace Variables.",
	}

	// `tfx variable list` command
	variableListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Variables",
		Long:  "List Variables in a Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return variableList(
				getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("workspace"))
		},
	}

	// `tfx variable create` command
	variableCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a Variable",
		Long:  "Create a Variable in a Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			value := *viperString("value")
			valueFile := *viperString("valueFile")
			if value == "" && valueFile == "" {
				return errors.New("required flag \"key\" or \"keyFile\" not set")
			}
			if value != "" && valueFile != "" {
				return errors.New("too many flags, only \"key\" or \"keyFile\" can be set, not both")
			}

			if valueFile != "" {
				if !isFile(valueFile) {
					return errors.New("valueFile does not exist")
				}
				o.AddMessageUserProvided("Variable Filename was contents will be used: ", valueFile)
				var err error
				value, err = readFile(valueFile) //override value
				if err != nil {
					return errors.Wrap(err, "unable to read the file passed")
				}
			}

			return variableCreate(
				getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("workspace"),
				*viperString("key"),
				value,
				*viperString("description"),
				*viperBool("env"),
				*viperBool("hcl"),
				*viperBool("sensitive"))
		},
	}

	// `tfx variable update` command
	variableUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update a Variable",
		Long:  "Update a Variable in a Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			value := *viperString("value")
			valueFile := *viperString("valueFile")
			if value == "" && valueFile == "" {
				return errors.New("required flag \"key\" or \"keyFile\" not set")
			}
			if value != "" && valueFile != "" {
				return errors.New("too many flags, only \"key\" or \"keyFile\" can be set, not both")
			}

			if valueFile != "" {
				if !isFile(valueFile) {
					return errors.New("valueFile does not exist")
				}
				o.AddMessageUserProvided("Variable Filename was contents will be used: ", valueFile)
				var err error
				value, err = readFile(valueFile) //override value
				if err != nil {
					return errors.Wrap(err, "unable to read the file passed")
				}
			}

			return variableUpdate(
				getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("workspace"),
				*viperString("key"),
				value,
				*viperString("description"),
				*viperBool("env"),
				*viperBool("hcl"),
				*viperBool("sensitive"))
		},
	}

	// `tfx variable show` command
	variableShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show details of a Variable",
		Long:  "Show details of a Variable in a Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return variableShow(
				getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("workspace"),
				*viperString("key"))
		},
	}

	// `tfx variable delete` command
	variableDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a Variable",
		Long:  "Delete a Variable in a Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return variableDelete(
				getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("workspace"),
				*viperString("key"))
		},
	}
)

func init() {
	// `tfx variable list` command
	variableListCmd.Flags().StringP("workspace", "w", "", "Name of the Workspace")
	variableListCmd.MarkFlagRequired("workspace")

	// `tfx variable create` command
	variableCreateCmd.Flags().StringP("workspace", "w", "", "Name of the Workspace")
	variableCreateCmd.Flags().StringP("key", "k", "", "Key of the Variable")
	variableCreateCmd.Flags().StringP("value", "v", "", "Value of the Variable (value or valueFile must be set)")
	variableCreateCmd.Flags().StringP("valueFile", "f", "", "Path to a variable text file, the contents of the file will be used (value or valueFile must be set)")
	variableCreateCmd.Flags().StringP("description", "d", "", "Description of the Variable (optional)")
	variableCreateCmd.Flags().BoolP("env", "e", false, "Variable is an Environment Variable (optional, defaults to false)")
	variableCreateCmd.Flags().BoolP("hcl", "", false, "Value of Variable is HCL (optional, defaults to false)")
	variableCreateCmd.Flags().BoolP("sensitive", "", false, "Variable is Sensitive (optional, defaults to false)")
	variableCreateCmd.MarkFlagRequired("workspace")
	variableCreateCmd.MarkFlagRequired("key")
	// command code to ensure ONE of these is set
	// variableCreateCmd.MarkFlagRequired("value")
	// variableCreateCmd.MarkFlagRequired("valueFile")

	// `tfx variable update` command
	variableUpdateCmd.Flags().StringP("workspace", "w", "", "Name of the Workspace")
	variableUpdateCmd.Flags().StringP("key", "k", "", "Key of the Variable")
	variableUpdateCmd.Flags().StringP("value", "v", "", "Value of the Variable (value or valueFile must be set)")
	variableUpdateCmd.Flags().StringP("valueFile", "f", "", "Path to a variable text file, the contents of the file will be used (value or valueFile must be set)")
	variableUpdateCmd.Flags().StringP("description", "d", "", "Description of the Variable (optional)")
	variableUpdateCmd.Flags().BoolP("env", "e", false, "Variable is an Environment Variable (optional, defaults to false)")
	variableUpdateCmd.Flags().BoolP("hcl", "", false, "Value of Variable is HCL (optional, defaults to false)")
	variableUpdateCmd.Flags().BoolP("sensitive", "", false, "Variable is Sensitive (optional, defaults to false)")
	variableUpdateCmd.MarkFlagRequired("workspace")
	variableUpdateCmd.MarkFlagRequired("key")
	// command code to ensure ONE of these is set
	// variableUpdateCmd.MarkFlagRequired("value")
	// variableUpdateCmd.MarkFlagRequired("valueFile")

	// `tfx variable show` command
	variableShowCmd.Flags().StringP("workspace", "w", "", "Name of the Workspace")
	variableShowCmd.Flags().StringP("key", "k", "", "Key of the Variable")
	variableShowCmd.MarkFlagRequired("workspace")
	variableShowCmd.MarkFlagRequired("key")

	// `tfx variable delete` command
	variableDeleteCmd.Flags().StringP("workspace", "w", "", "Name of the Workspace")
	variableDeleteCmd.Flags().StringP("key", "k", "", "Key of the Variable")
	variableDeleteCmd.MarkFlagRequired("workspace")
	variableDeleteCmd.MarkFlagRequired("key")

	workspaceCmd.AddCommand(variableCmd)
	variableCmd.AddCommand(variableListCmd)
	variableCmd.AddCommand(variableCreateCmd)
	variableCmd.AddCommand(variableUpdateCmd)
	variableCmd.AddCommand(variableShowCmd)
	variableCmd.AddCommand(variableDeleteCmd)
}

func variablesListAll(c TfxClientContext, workspaceId string) ([]*tfe.Variable, error) {
	allItems := []*tfe.Variable{}
	opts := tfe.VariableListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100,
		},
	}
	for {
		items, err := c.Client.Variables.List(c.Context, workspaceId, &opts)
		if err != nil {
			return nil, err
		}

		allItems = append(allItems, items.Items...)
		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage
	}

	return allItems, nil
}

func variableList(c TfxClientContext, orgName string, workspaceName string) error {
	o.AddMessageUserProvided("List Variables for Workspace:", workspaceName)
	workspaceId, err := getWorkspaceId(c, orgName, workspaceName)
	if err != nil {
		return errors.Wrap(err, "unable to read workspace id")
	}

	items, err := variablesListAll(c, workspaceId)
	if err != nil {
		return errors.Wrap(err, "failed to list variables")
	}

	o.AddTableHeader("Id", "Key", "Value", "Sensitive", "HCL", "Category", "Description")
	for _, i := range items {
		o.AddTableRows(i.ID, i.Key, i.Value, i.Sensitive, i.HCL, i.Category, i.Description)
	}
	o.Close()

	return nil
}

func variableCreate(c TfxClientContext, orgName string, workspaceName string,
	variableKey string, variableValue string, description string, isEnvironment bool, isHcl bool, isSensitive bool) error {
	o.AddMessageUserProvided("Create Variable for Workspace:", workspaceName)
	workspaceId, err := getWorkspaceId(c, orgName, workspaceName)
	if err != nil {
		return errors.Wrap(err, "unable to read workspace id")
	}

	// check if value is a file, if so use the contents
	// TODO: consider moving this to a different arg/command?
	if isFile(variableValue) {
		o.AddMessageUserProvided("Value passed as a filename, contents will be used: ", variableValue)
		variableValue, err = readFile(variableValue)
		if err != nil {
			return errors.Wrap(err, "unable to read the file passed")
		}
	}

	var category *tfe.CategoryType
	if isEnvironment {
		category = tfe.Category(tfe.CategoryEnv)
	} else {
		category = tfe.Category(tfe.CategoryTerraform)
	}
	variable, err := c.Client.Variables.Create(c.Context, workspaceId, tfe.VariableCreateOptions{
		Key:         &variableKey,
		Value:       &variableValue,
		Description: &description,
		Category:    category,
		HCL:         &isHcl,
		Sensitive:   &isSensitive,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to Create Variable")
	}

	o.AddMessageUserProvided("Variable Created:", variableKey)
	o.AddDeferredMessageRead("ID", variable.ID)
	o.AddDeferredMessageRead("Key", variable.Key)
	o.AddDeferredMessageRead("Value", variable.Value)
	o.AddDeferredMessageRead("Sensitive", variable.Sensitive)
	o.AddDeferredMessageRead("HCL", variable.HCL)
	o.AddDeferredMessageRead("Category", variable.Category)
	o.AddDeferredMessageRead("Description", variable.Description)
	o.Close()

	return nil
}

func variableUpdate(c TfxClientContext, orgName string, workspaceName string,
	variableKey string, variableValue string, description string, isEnvironment bool, isHcl bool, isSensitive bool) error {
	o.AddMessageUserProvided("Update Variable for Workspace:", workspaceName)
	workspaceId, err := getWorkspaceId(c, orgName, workspaceName)
	if err != nil {
		return errors.Wrap(err, "unable to read workspace id")
	}

	variableId, err := getVariableId(c, workspaceId, variableKey)
	if err != nil {
		return errors.Wrap(err, "unable to read variable id")
	}

	var category *tfe.CategoryType
	if isEnvironment {
		category = tfe.Category(tfe.CategoryEnv)
	} else {
		category = tfe.Category(tfe.CategoryTerraform)
	}
	variable, err := c.Client.Variables.Update(c.Context, workspaceId, variableId, tfe.VariableUpdateOptions{
		Key:         &variableKey,
		Value:       &variableValue,
		Description: &description,
		Category:    category,
		HCL:         &isHcl,
		Sensitive:   &isSensitive,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to Update Variable")
	}

	o.AddMessageUserProvided("Variable Updated", workspaceName)
	o.AddDeferredMessageRead("ID", variable.ID)
	o.AddDeferredMessageRead("Key", variable.Key)
	o.AddDeferredMessageRead("Value", variable.Value)
	o.AddDeferredMessageRead("Sensitive", variable.Sensitive)
	o.AddDeferredMessageRead("HCL", variable.HCL)
	o.AddDeferredMessageRead("Category", variable.Category)
	o.AddDeferredMessageRead("Description", variable.Description)
	o.Close()

	return nil
}

func variableShow(c TfxClientContext, orgName string, workspaceName string, variableKey string) error {
	o.AddMessageUserProvided("Show Variable for Workspace:", workspaceName)
	workspaceId, err := getWorkspaceId(c, orgName, workspaceName)
	if err != nil {
		return errors.Wrap(err, "unable to read workspace id")
	}

	variableId, err := getVariableId(c, workspaceId, variableKey)
	if err != nil {
		return errors.Wrap(err, "unable to read variable id")
	}

	variable, err := c.Client.Variables.Read(c.Context, workspaceId, variableId)
	if err != nil {
		return errors.Wrap(err, "unable to read variable")
	}

	o.AddDeferredMessageRead("ID", variable.ID)
	o.AddDeferredMessageRead("Key", variable.Key)
	o.AddDeferredMessageRead("Value", variable.Value)
	o.AddDeferredMessageRead("Sensitive", variable.Sensitive)
	o.AddDeferredMessageRead("HCL", variable.HCL)
	o.AddDeferredMessageRead("Category", variable.Category)
	o.AddDeferredMessageRead("Description", variable.Description)
	o.Close()

	return nil
}

func variableDelete(c TfxClientContext, orgName string, workspaceName string, variableKey string) error {
	// TODO: Add ability to delete multiple keys at once: https://github.com/spf13/cobra/issues/661
	o.AddMessageUserProvided("Delete Variable for Workspace:", workspaceName)
	workspaceId, err := getWorkspaceId(c, orgName, workspaceName)
	if err != nil {
		return errors.Wrap(err, "unable to read workspace id")
	}

	variableId, err := getVariableId(c, workspaceId, variableKey)
	if err != nil {
		return errors.Wrap(err, "unable to read variable id")
	}

	err = c.Client.Variables.Delete(c.Context, workspaceId, variableId)
	if err != nil {
		return errors.Wrap(err, "failed to delete variable")
	}

	o.AddMessageUserProvided("Variable Deleted:", variableKey)
	o.AddDeferredMessageRead("Status", "Success")
	o.Close()

	return nil
}

func getVariableId(c TfxClientContext, workspaceId string, variableKey string) (string, error) {
	vars, err := variablesListAll(c, workspaceId)
	if err != nil {
		return "", err
	}

	for _, v := range vars {
		if v.Key == variableKey {
			return v.ID, nil
		}
	}

	return "", errors.New("variable key not found")
}
