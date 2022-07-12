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
	// `tfx registry provider version` commands
	registryProviderVersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Provider Versions in Private Registry Commands",
		Long:  "Commands to work with Provider Versions in a Private Registry of a TFx Organization.",
	}

	// `tfx registry provider version list` command
	registryProviderVersionListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Provider Versions in a Private Registry",
		Long:  "List Provider Versions for a Provider in a Private Registry of a TFx Organization. ",
		RunE: func(cmd *cobra.Command, args []string) error {
			return registryProviderVersionList()
		},
		PreRun: bindPFlags,
	}

	// `tfx registry provider version create` command
	registryProviderVersionCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create/Update a Provider Version in a Private Registry",
		Long:  "Create/Update a Provider Version for a Provider in a Private Registry of a TFx Organization. ",
		RunE: func(cmd *cobra.Command, args []string) error {
			return registryProviderVersionCreate()
		},
		PreRun: bindPFlags,
	}

	// `tfx registry provider version show` command
	registryProviderVersionShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show details of a Provider Version in a Private Registry",
		Long:  "Show details of a Provider Version for a Provider in a Private Registry of a TFx Organization. ",
		RunE: func(cmd *cobra.Command, args []string) error {
			return registryProviderVersionShow()
		},
		PreRun: bindPFlags,
	}

	// `tfx registry provider version delete` command
	registryProviderVersionDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a Provider Version in a Private Registry",
		Long:  "Delete a Provider Version for a Provider in a Private Registry of a TFx Organization. ",
		RunE: func(cmd *cobra.Command, args []string) error {
			return registryProviderVersionDelete()
		},
		PreRun: bindPFlags,
	}
)

func init() {
	// `tfx registry provider version list` arguments
	registryProviderVersionListCmd.Flags().StringP("name", "n", "", "Name of the Provider")
	registryProviderVersionListCmd.MarkFlagRequired("name")

	// `tfx registry provider version create` arguments
	// `tfx registry provider version show` arguments
	// `tfx registry provider version delete` arguments

	registryProviderCmd.AddCommand(registryProviderVersionCmd)
	registryProviderVersionCmd.AddCommand(registryProviderVersionListCmd)
	registryProviderVersionCmd.AddCommand(registryProviderVersionCreateCmd)
	registryProviderVersionCmd.AddCommand(registryProviderVersionShowCmd)
	registryProviderVersionCmd.AddCommand(registryProviderVersionDeleteCmd)
}

func registryProviderVersionList() error {
	// Validate flags
	orgName := *viperString("tfeOrganization")
	providerName := *viperString("name")

	client, ctx := getClientContext()

	// Read all providers in PMR
	fmt.Println("Reading Providers for Organization:", color.GreenString(orgName))
	fmt.Println("Provider Name:", color.GreenString(providerName))
	provider, err := client.RegistryProviderVersions.List(ctx, tfe.RegistryProviderID{
		OrganizationName: orgName,
		Namespace:        orgName,
		RegistryName:     "private", // for some reason public doesn't work...
		Name:             providerName,
	}, &tfe.RegistryProviderVersionListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})

	if err != nil {
		logError(err, "failed to read provider in PMR")
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Registry", "Published", "SHASUM Uploaded", "SHASUM Sig Uploaded"})
	for _, i := range provider.Items {

		t.AppendRow(table.Row{i.ID, i.Version, i.UpdatedAt, i.ShasumsUploaded, i.ShasumsSigUploaded})
	}
	t.SetStyle(table.StyleRounded)
	t.Render()

	return nil
}

func registryProviderVersionCreate() error {
	fmt.Println(color.MagentaString("Function not implemented yet."))
	return nil
}

func registryProviderVersionShow() error {
	fmt.Println(color.MagentaString("Function not implemented yet."))
	return nil
}

func registryProviderVersionDelete() error {
	fmt.Println(color.MagentaString("Function not implemented yet."))
	return nil
}
