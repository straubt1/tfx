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
	// `tfx registry module` commands
	registryModuleCmd = &cobra.Command{
		Use:   "module",
		Short: "Modules in Private Registry Commands",
		Long:  "Work with Private Module Registry of a TFx Organization.",
	}

	// `tfx registry module list` command
	registryModuleListCmd = &cobra.Command{
		Use:   "list",
		Short: "List modules",
		Long:  "List modules in the Private Module Registry of a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return registryModuleList(
				getTfxClientContext())
		},
	}

	// `tfx registry module create` command
	registryModuleCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a module",
		Long:  "Create a module in the Private Module Registry for a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return registryModuleCreate(
				getTfxClientContext(),
				*viperString("name"),
				*viperString("provider"))
		},
	}

	// `tfx registry module show` command
	registryModuleShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show a module",
		Long:  "Show a module details of a module in the Private Module Registry for a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return registryModuleShow(
				getTfxClientContext(),
				*viperString("name"),
				*viperString("provider"))
		},
	}

	// `tfx registry module delete` command
	registryModuleDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a module",
		Long:  "Delete a module in the Private Module Registry for a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return registryModuleDelete(
				getTfxClientContext(),
				*viperString("name"),
				*viperString("provider"))
		},
	}
)

func init() {
	// `tfx registry module list` arguments

	// `tfx registry module create` arguments
	registryModuleCreateCmd.Flags().StringP("name", "n", "", "Name of the Module (no spaces)")
	registryModuleCreateCmd.Flags().StringP("provider", "p", "", "Name of the provider (no spaces) (i.e. aws, azure, google)")
	registryModuleCreateCmd.MarkFlagRequired("name")
	registryModuleCreateCmd.MarkFlagRequired("provider")

	// `tfx registry module show` arguments
	registryModuleShowCmd.Flags().StringP("name", "n", "", "Name of the Module (no spaces)")
	registryModuleShowCmd.Flags().StringP("provider", "p", "", "Name of the provider (no spaces) (i.e. aws, azure, google)")
	registryModuleShowCmd.MarkFlagRequired("name")
	registryModuleShowCmd.MarkFlagRequired("provider")

	// `tfx registry module delete` arguments
	registryModuleDeleteCmd.Flags().StringP("name", "n", "", "Name of the Module (no spaces)")
	registryModuleDeleteCmd.Flags().StringP("provider", "p", "", "Name of the provider (no spaces) (i.e. aws, azure, google)")
	registryModuleDeleteCmd.MarkFlagRequired("name")
	registryModuleDeleteCmd.MarkFlagRequired("provider")

	registryCmd.AddCommand(registryModuleCmd)
	registryModuleCmd.AddCommand(registryModuleListCmd)
	registryModuleCmd.AddCommand(registryModuleCreateCmd)
	registryModuleCmd.AddCommand(registryModuleShowCmd)
	registryModuleCmd.AddCommand(registryModuleDeleteCmd)
}

func registryModuleListAll(c TfxClientContext, orgName string) ([]*tfe.RegistryModule, error) {
	allItems := []*tfe.RegistryModule{}
	opts := tfe.RegistryModuleListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100,
		},
	}

	for {
		items, err := c.Client.RegistryModules.List(c.Context, orgName, &opts)
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

func registryModuleList(c TfxClientContext) error {
	o.AddMessageUserProvided("List Modules for Organization:", c.OrganizationName)
	items, err := registryModuleListAll(c, c.OrganizationName)
	if err != nil {
		return errors.Wrap(err, "failed to list modules")
	}

	o.AddTableHeader("Name", "Provider", "ID", "Status", "Published", "Versions")
	for _, i := range items {
		o.AddTableRows(i.Name, i.Provider, i.ID, i.Status, i.UpdatedAt, len(i.VersionStatuses))
	}

	return nil
}

func registryModuleCreate(c TfxClientContext, moduleName string, providerName string) error {
	o.AddMessageUserProvided("Create Module for Organization:", c.OrganizationName)
	module, err := c.Client.RegistryModules.Create(c.Context, c.OrganizationName, tfe.RegistryModuleCreateOptions{
		Name:     &moduleName,
		Provider: &providerName,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create module")
	}

	o.AddMessageUserProvided("Module Created:", module.Name)
	o.AddDeferredMessageRead("ID", module.ID)
	o.AddDeferredMessageRead("Namespace", module.Namespace)
	o.AddDeferredMessageRead("Created", module.CreatedAt)
	o.AddDeferredMessageRead("Updated", module.CreatedAt)

	return nil
}

func registryModuleShow(c TfxClientContext, moduleName string, providerName string) error {
	o.AddMessageUserProvided("Show Module for Organization:", c.OrganizationName)
	module, err := c.Client.RegistryModules.Read(c.Context, tfe.RegistryModuleID{
		Organization: c.OrganizationName,
		Name:         moduleName,
		Provider:     providerName,
		Namespace:    c.OrganizationName,
		RegistryName: tfe.PrivateRegistry,
	})
	if err != nil {
		logError(err, "failed to show module")
	}

	o.AddDeferredMessageRead("ID", module.ID)
	o.AddDeferredMessageRead("Status", module.Status)
	o.AddDeferredMessageRead("Created", module.CreatedAt)
	o.AddDeferredMessageRead("Updated", module.UpdatedAt)
	o.AddDeferredMessageRead("Versions", len(module.VersionStatuses))
	if len(module.VersionStatuses) > 0 {
		o.AddDeferredMessageRead("Latest Version", module.VersionStatuses[0].Version)
	}

	return nil
}

func registryModuleDelete(c TfxClientContext, moduleName string, providerName string) error {
	o.AddMessageUserProvided("Delete Module for Organization:", c.OrganizationName)
	// RegistryModules.DeleteProvider requires the provider as well (if just the name is used, multiple modules could be deleted)
	err := c.Client.RegistryModules.DeleteProvider(c.Context, tfe.RegistryModuleID{
		Organization: c.OrganizationName,
		Name:         moduleName,
		Provider:     providerName,
		RegistryName: tfe.PrivateRegistry,
	})
	if err != nil {
		logError(err, "failed to delete module")
	}

	o.AddMessageUserProvided("Module Deleted:", moduleName)
	o.AddDeferredMessageRead("Status", "Success")

	return nil
}
