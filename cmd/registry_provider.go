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
	// `tfx registry provider` commands
	registryProviderCmd = &cobra.Command{
		Use:   "provider",
		Short: "Providers in Private Registry Commands",
		Long:  "Commands to work with Providers in a Private Registry of a TFx Organization.",
	}

	// `tfx registry provider list` command
	registryProviderListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Providers in a Private Registry",
		Long:  "List Providers in a Private Registry of a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return registryProviderList(
				getTfxClientContext(),
				*viperString("tfeOrganization"))
		},
	}

	// `tfx registry provider create` command
	registryProviderCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a Provider in a Private Registry",
		Long:  "Create a Provider in a Private Registry of a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return registryProviderCreate(
				getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("name"))
		},
	}

	// `tfx registry provider show` command
	registryProviderShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show details about a Provider in a Private Registry",
		Long:  "Show details about a Provider in a Private Registry of a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return registryProviderShow(
				getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("name"))
		},
	}

	// `tfx registry provider delete` command
	registryProviderDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a Provider in a Private Registry",
		Long:  "Delete a Provider in a Private Registry of a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return registryProviderDelete(
				getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("name"))
		},
	}
)

func init() {
	// `tfx registry provider list` arguments

	// `tfx registry provider create` arguments
	registryProviderCreateCmd.Flags().StringP("name", "n", "", "Name of the Provider")
	registryProviderCreateCmd.MarkFlagRequired("name")

	// `tfx registry provider show` arguments
	registryProviderShowCmd.Flags().StringP("name", "n", "", "Name of the Provider")
	registryProviderShowCmd.MarkFlagRequired("name")

	// `tfx registry provider delete` arguments
	registryProviderDeleteCmd.Flags().StringP("name", "n", "", "Name of the Provider")
	registryProviderDeleteCmd.MarkFlagRequired("name")

	registryCmd.AddCommand(registryProviderCmd)
	registryProviderCmd.AddCommand(registryProviderListCmd)
	registryProviderCmd.AddCommand(registryProviderCreateCmd)
	registryProviderCmd.AddCommand(registryProviderShowCmd)
	registryProviderCmd.AddCommand(registryProviderDeleteCmd)
}

func registryProviderListAll(c TfxClientContext, orgName string) ([]*tfe.RegistryProvider, error) {
	allItems := []*tfe.RegistryProvider{}
	opts := tfe.RegistryProviderListOptions{
		// RegistryName: tfe.PrivateRegistry, // Can restrict to just private
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100,
		},
		// Include: &[]tfe.RegistryProviderIncludeOps{"provider-versions"}, does not work, cant get provider versions from this call?
	}
	for {
		items, err := c.Client.RegistryProviders.List(c.Context, orgName, &opts)
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

func registryProviderList(c TfxClientContext, orgName string) error {
	o.AddMessageUserProvided("List Providers in Registry for Organization:", orgName)
	items, err := registryProviderListAll(c, orgName)
	if err != nil {
		return errors.Wrap(err, "failed to list providers")
	}

	o.AddTableHeader("Name", "Registry", "ID", "Published")
	for _, i := range items {
		o.AddTableRows(i.Name, i.RegistryName, i.ID, i.UpdatedAt)
	}
	o.Close()

	return nil
}

func registryProviderCreate(c TfxClientContext, orgName string, providerName string) error {
	o.AddMessageUserProvided("Create Provider in Registry for Organization:", orgName)
	provider, err := c.Client.RegistryProviders.Create(c.Context, orgName, tfe.RegistryProviderCreateOptions{
		Name:         providerName,
		Namespace:    orgName, // always org name for RegistryName "private"
		RegistryName: tfe.PrivateRegistry,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create provider")
	}

	o.AddMessageUserProvided("Provider Created:", provider.Name)
	o.AddDeferredMessageRead("ID", provider.ID)
	o.AddDeferredMessageRead("Namespace", provider.Namespace)
	o.AddDeferredMessageRead("Created", provider.UpdatedAt)
	o.Close()

	return nil
}

func registryProviderShow(c TfxClientContext, orgName string, providerName string) error {
	o.AddMessageUserProvided("Show Provider in Registry for Organization:", orgName)
	provider, err := c.Client.RegistryProviders.Read(c.Context, tfe.RegistryProviderID{
		OrganizationName: orgName,
		Name:             providerName,
		Namespace:        orgName, // always org name for RegistryName "private"
		RegistryName:     tfe.PrivateRegistry,
	}, &tfe.RegistryProviderReadOptions{
		Include: []tfe.RegistryProviderIncludeOps{},
	})
	if err != nil {
		return errors.Wrap(err, "failed to read provider")
	}

	o.AddDeferredMessageRead("Name", provider.Name)
	o.AddDeferredMessageRead("ID", provider.ID)
	o.AddDeferredMessageRead("Namespace", provider.Namespace)
	o.AddDeferredMessageRead("Created", provider.UpdatedAt)
	o.Close()

	return nil
}

func registryProviderDelete(c TfxClientContext, orgName string, providerName string) error {
	o.AddMessageUserProvided("Delete Provider in Registry for Organization:", orgName)
	err := c.Client.RegistryProviders.Delete(c.Context, tfe.RegistryProviderID{
		OrganizationName: orgName,
		Name:             providerName,
		Namespace:        orgName, // always org name for RegistryName "private"
		RegistryName:     tfe.PrivateRegistry,
	})
	if err != nil {
		return errors.Wrap(err, "failed to delete provider")
	}

	o.AddMessageUserProvided("Provider Deleted:", providerName)
	o.AddDeferredMessageRead("Status", "Success")
	o.Close()

	return nil
}
