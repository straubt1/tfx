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
	"errors"
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
		Short: "Create a Provider Version in a Private Registry",
		Long:  "Create a Provider Version for a Provider in a Private Registry of a TFx Organization. ",
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
	registryProviderVersionCreateCmd.Flags().StringP("name", "n", "", "Name of the Provider")
	registryProviderVersionCreateCmd.Flags().StringP("version", "v", "", "Version of Provider (i.e. 0.0.1)")
	registryProviderVersionCreateCmd.Flags().StringP("keyId", "", "", "GPG Key Id")
	registryProviderVersionCreateCmd.Flags().StringP("shasums", "", "", "Path to shasums")
	registryProviderVersionCreateCmd.Flags().StringP("shasumssig", "", "", "Path to shasumssig")
	registryProviderVersionCreateCmd.MarkFlagRequired("name")
	registryProviderVersionCreateCmd.MarkFlagRequired("version")
	registryProviderVersionCreateCmd.MarkFlagRequired("keyId")
	registryProviderVersionCreateCmd.MarkFlagRequired("shasums")
	registryProviderVersionCreateCmd.MarkFlagRequired("shasumssig")

	// `tfx registry provider version show` arguments
	registryProviderVersionShowCmd.Flags().StringP("name", "n", "", "Name of the Provider")
	registryProviderVersionShowCmd.Flags().StringP("version", "v", "", "Version of Provider (i.e. 0.0.1)")
	registryProviderVersionShowCmd.MarkFlagRequired("name")
	registryProviderVersionShowCmd.MarkFlagRequired("version")

	// `tfx registry provider version delete` arguments
	registryProviderVersionDeleteCmd.Flags().StringP("name", "n", "", "Name of the Provider")
	registryProviderVersionDeleteCmd.Flags().StringP("version", "v", "", "Version of Provider (i.e. 0.0.1)")
	registryProviderVersionDeleteCmd.MarkFlagRequired("name")
	registryProviderVersionDeleteCmd.MarkFlagRequired("version")

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
	// Validate flags
	orgName := *viperString("tfeOrganization")
	providerName := *viperString("name")
	providerVersion := *viperString("version")
	keyId := *viperString("keyId")
	shasums := *viperString("shasums")
	shasumssig := *viperString("shasumssig")

	if _, err := os.Stat(shasums); errors.Is(err, os.ErrNotExist) {
		logError(err, "shasums file does not exist")
	}
	if _, err := os.Stat(shasumssig); errors.Is(err, os.ErrNotExist) {
		logError(err, "shasumssig file does not exist")
	}

	client, ctx := getClientContext()
	// existing, err := client.RegistryProviderVersions.Read(ctx, tfe.RegistryProviderVersionID{
	// 	RegistryProviderID: tfe.RegistryProviderID{
	// 		OrganizationName: orgName,
	// 		Namespace:        orgName, // always org name for RegistryName "private"
	// 		RegistryName:     tfe.PrivateRegistry,
	// 		Name:             providerName,
	// 	},
	// 	Version: providerVersion,
	// })
	// if err != nil {
	// 	logError(err, "failed to find provider in PMR")
	// }
	// fmt.Println(existing)

	// Create provider in Registry
	fmt.Println("Create Provider for Organization:", color.GreenString(orgName))
	fmt.Println("Provider Name:", color.GreenString(providerName))
	p, err := client.RegistryProviderVersions.Create(ctx, tfe.RegistryProviderID{
		OrganizationName: orgName,
		Namespace:        orgName, // always org name for RegistryName "private"
		RegistryName:     tfe.PrivateRegistry,
		Name:             providerName,
	}, tfe.RegistryProviderVersionCreateOptions{
		Version: providerVersion,
		KeyID:   keyId,
		// Protocols: []string{},
	})
	if err != nil {
		// logError(err, "failed to create provider in PMR")
		fmt.Println("failed to create provider in PMR")
	}

	err = UploadBinary(p.Links["shasums-upload"].(string), shasums)
	if err != nil {
		logError(err, "failed to upload shasums")
	}
	err = UploadBinary(p.Links["shasums-sig-upload"].(string), shasumssig)
	if err != nil {
		logError(err, "failed to upload shasums sig")
	}
	fmt.Println(p.Links["shasums-upload"], p.Links["shasums-sig-upload"], p.CreatedAt)
	fmt.Println(shasums, shasumssig, p.CreatedAt)
	return nil
}

func registryProviderVersionShow() error {
	// Validate flags
	orgName := *viperString("tfeOrganization")
	providerName := *viperString("name")
	providerVersion := *viperString("version")

	client, ctx := getClientContext()

	provider, err := client.RegistryProviderVersions.Read(ctx, tfe.RegistryProviderVersionID{
		RegistryProviderID: tfe.RegistryProviderID{
			OrganizationName: orgName,
			Namespace:        orgName,
			RegistryName:     "private", // for some reason public doesn't work...
			Name:             providerName,
		},
		Version: providerVersion,
	})
	if err != nil {
		logError(err, "failed to read provider in PMR")
	}

	fmt.Println(color.BlueString("Name:                 "), provider.RegistryProvider.Name)
	fmt.Println(color.BlueString("Version:              "), provider.Version)
	fmt.Println(color.BlueString("ID:                   "), provider.ID)
	fmt.Println(color.BlueString("Shasums Uploaded:     "), provider.ShasumsUploaded)
	fmt.Println(color.BlueString("Shasums Sig Uploaded: "), provider.ShasumsSigUploaded)

	// If the Shasums have been uploaded, display them (might be a better place for this?)
	if provider.ShasumsUploaded {
		sha, err := DownloadTextFile(provider.Links["shasums-download"].(string))
		if err != nil {
			logError(err, "failed to read shassum in PMR")
		}
		fmt.Println(color.BlueString("Shasums:"))
		fmt.Println(sha)
	}

	return nil
}

func registryProviderVersionDelete() error {
	client, ctx := getClientContext()
	orgName := *viperString("tfeOrganization")
	providerName := *viperString("name")
	providerVersion := *viperString("version")

	fmt.Println("Delete Provider Version in Registry for Organization:", color.GreenString(orgName))
	err := client.RegistryProviderVersions.Delete(ctx, tfe.RegistryProviderVersionID{
		RegistryProviderID: tfe.RegistryProviderID{
			OrganizationName: orgName,
			Name:             providerName,
			Namespace:        orgName, // always org name for RegistryName "private"
			RegistryName:     tfe.PrivateRegistry,
		},
		Version: providerVersion,
	})
	if err != nil {
		logError(err, "failed to delete Provider Version")
	}

	fmt.Println(color.BlueString("Provider Deleted: "), providerName, providerVersion)
	return nil
}
