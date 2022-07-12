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
	"fmt"
	"os"

	"github.com/fatih/color"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/jedib0t/go-pretty/v6/table"
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
		Long:  "List Providers in a Private Registry of a TFx Organization. ",
		RunE: func(cmd *cobra.Command, args []string) error {
			return registryProviderList()
		},
		PreRun: bindPFlags,
	}

	// `tfx registry provider create` command
	registryProviderCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a Provider in a Private Registry",
		Long:  "Create a Provider in a Private Registry of a TFx Organization. ",
		RunE: func(cmd *cobra.Command, args []string) error {
			return registryProviderCreate()
		},
		PreRun: bindPFlags,
	}

	// `tfx registry provider show` command
	registryProviderShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show details about a Provider in a Private Registry",
		Long:  "Show details about a Provider in a Private Registry of a TFx Organization. ",
		RunE: func(cmd *cobra.Command, args []string) error {
			return registryProviderShow()
		},
		PreRun: bindPFlags,
	}

	// `tfx registry provider delete` command
	registryProviderDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a Provider in a Private Registry",
		Long:  "Delete a Provider in a Private Registry of a TFx Organization. ",
		RunE: func(cmd *cobra.Command, args []string) error {
			return registryProviderDelete()
		},
		PreRun: bindPFlags,
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

func registryProviderList() error {
	// Validate flags
	orgName := *viperString("tfeOrganization")

	client, ctx := getClientContext()

	// Read all providers in PMR
	fmt.Println("Reading Providers for Organization:", color.GreenString(orgName))
	modules, err := client.RegistryProviders.List(ctx, orgName, &tfe.RegistryProviderListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
		// Include: &[]tfe.RegistryProviderIncludeOps{"provider-versions"}, does not work, cant get provider versions from this call?
	})
	if err != nil {
		logError(err, "failed to read providers in PMR")
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Registry", "ID", "Published"})
	for _, i := range modules.Items {

		t.AppendRow(table.Row{i.Name, i.RegistryName, i.ID, i.UpdatedAt})
	}
	t.SetStyle(table.StyleRounded)
	t.Render()

	return nil
}

func registryProviderCreate() error {
	// Validate flags
	orgName := *viperString("tfeOrganization")
	name := *viperString("name")

	client, ctx := getClientContext()

	fmt.Println("Creating Provider in Registry for Organization:", color.GreenString(orgName))
	provider, err := client.RegistryProviders.Create(ctx, orgName, tfe.RegistryProviderCreateOptions{
		Name:         name,
		Namespace:    orgName, // always org name for RegistryName "private"
		RegistryName: tfe.PrivateRegistry,
	})
	if err != nil {
		logError(err, "failed to create Provider")
	}

	fmt.Println(color.BlueString("Name:      "), provider.Name)
	fmt.Println(color.BlueString("ID:        "), provider.ID)
	fmt.Println(color.BlueString("Namespace: "), provider.Namespace)
	fmt.Println(color.BlueString("Created:   "), provider.UpdatedAt)

	return nil
}

func registryProviderShow() error {
	// Validate flags
	orgName := *viperString("tfeOrganization")
	name := *viperString("name")

	client, ctx := getClientContext()

	fmt.Println("Creating Provider in Registry for Organization:", color.GreenString(orgName))
	provider, err := client.RegistryProviders.Read(ctx, tfe.RegistryProviderID{
		OrganizationName: orgName,
		Name:             name,
		Namespace:        orgName, // always org name for RegistryName "private"
		RegistryName:     tfe.PrivateRegistry,
	}, &tfe.RegistryProviderReadOptions{
		Include: []tfe.RegistryProviderIncludeOps{},
	})
	if err != nil {
		logError(err, "failed to create Provider")
	}

	fmt.Println(color.BlueString("Name:      "), provider.Name)
	fmt.Println(color.BlueString("ID:        "), provider.ID)
	fmt.Println(color.BlueString("Namespace: "), provider.Namespace)
	fmt.Println(color.BlueString("Created:   "), provider.UpdatedAt)

	return nil
}

func registryProviderDelete() error {
	// Validate flags
	orgName := *viperString("tfeOrganization")
	name := *viperString("name")

	client, ctx := getClientContext()

	fmt.Println("Delete Provider in Registry for Organization:", color.GreenString(orgName))
	err := client.RegistryProviders.Delete(ctx, tfe.RegistryProviderID{
		OrganizationName: orgName,
		Name:             name,
		Namespace:        orgName, // always org name for RegistryName "private"
		RegistryName:     tfe.PrivateRegistry,
	})
	if err != nil {
		logError(err, "failed to delete Provider")
	}

	fmt.Println(color.BlueString("Provider Delete: "), name)

	return nil
}
