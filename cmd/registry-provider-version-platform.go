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
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fatih/color"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	// `tfx registry provider version platform` commands
	registryProviderVersionPlatformCmd = &cobra.Command{
		Use:   "platform",
		Short: "Provider Version Platforms in Private Registry Commands",
		Long:  "Commands to work with Provider Version Platforms in a Private Registry of a TFx Organization.",
	}

	// `tfx registry provider version platform list` command
	registryProviderVersionPlatformListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Provider Version Platforms in a Private Registry",
		Long:  "List Provider Version Platforms for a Provider in a Private Registry of a TFx Organization. ",
		RunE: func(cmd *cobra.Command, args []string) error {
			return registryProviderVersionPlatformList()
		},
		PreRun: bindPFlags,
	}

	// `tfx registry provider version platform create` command
	registryProviderVersionPlatformCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create/Update a Provider Version Platform in a Private Registry",
		Long:  "Create/Update a Provider Version Platform for a Provider Version in a Private Registry of a TFx Organization. ",
		RunE: func(cmd *cobra.Command, args []string) error {
			return registryProviderVersionPlatformCreate()
		},
		PreRun: bindPFlags,
	}

	// `tfx registry provider version platform show` command
	registryProviderVersionPlatformShowCmd = &cobra.Command{
		Use:   "create",
		Short: "Show details about a Provider Version Platform in a Private Registry",
		Long:  "Show details about a Provider Version Platform for a Provider Version in a Private Registry of a TFx Organization. ",
		RunE: func(cmd *cobra.Command, args []string) error {
			return registryProviderVersionPlatformShow()
		},
		PreRun: bindPFlags,
	}

	// `tfx registry provider version platform delete` command
	registryProviderVersionPlatformDeleteCmd = &cobra.Command{
		Use:   "create",
		Short: "Delete a Provider Version Platform in a Private Registry",
		Long:  "Delete a Provider Version Platform for a Provider Version in a Private Registry of a TFx Organization. ",
		RunE: func(cmd *cobra.Command, args []string) error {
			return registryProviderVersionPlatformDelete()
		},
		PreRun: bindPFlags,
	}
)

func init() {
	// `tfx registry provider version platform list` arguments
	registryProviderVersionPlatformListCmd.Flags().StringP("name", "n", "", "Name of the Provider")
	registryProviderVersionPlatformListCmd.Flags().StringP("version", "v", "", "Version of Provider (i.e. 0.0.1)")
	registryProviderVersionPlatformListCmd.MarkFlagRequired("name")

	// `tfx registry provider version platform create` arguments
	registryProviderVersionPlatformCreateCmd.Flags().StringP("name", "n", "", "Name of the Provider")
	registryProviderVersionPlatformCreateCmd.Flags().StringP("version", "v", "", "Version of Provider (i.e. 0.0.1)")
	registryProviderVersionPlatformCreateCmd.Flags().StringP("os", "", "", "OS of the Provider Version Platform (linux, windows, darwin)")
	registryProviderVersionPlatformCreateCmd.Flags().StringP("arch", "", "", "ARCH of the Provider Version Platform (amd64, arm64)")
	registryProviderVersionPlatformCreateCmd.Flags().StringP("filename", "f", "", "Path to the filename that is the provider binary")
	registryProviderVersionPlatformCreateCmd.MarkFlagRequired("name")
	registryProviderVersionPlatformCreateCmd.MarkFlagRequired("version")
	registryProviderVersionPlatformCreateCmd.MarkFlagRequired("os")
	registryProviderVersionPlatformCreateCmd.MarkFlagRequired("arch")
	registryProviderVersionPlatformCreateCmd.MarkFlagRequired("filename")

	// `tfx registry provider version platform show` arguments
	// `tfx registry provider version platform delete` arguments

	registryProviderVersionCmd.AddCommand(registryProviderVersionPlatformCmd)
	registryProviderVersionPlatformCmd.AddCommand(registryProviderVersionPlatformListCmd)
	registryProviderVersionPlatformCmd.AddCommand(registryProviderVersionPlatformCreateCmd)
	registryProviderVersionPlatformCmd.AddCommand(registryProviderVersionPlatformShowCmd)
	registryProviderVersionPlatformCmd.AddCommand(registryProviderVersionPlatformDeleteCmd)
}

func registryProviderVersionPlatformList() error {
	// Validate flags
	orgName := *viperString("tfeOrganization")
	providerName := *viperString("name")
	// Attempt to prevent a non semantic version from being read
	providerVersion, err := viperSemanticVersionString("version")
	if err != nil {
		logError(err, "failed to parse semantic version")
	}

	client, ctx := getClientContext()

	// Read all providers in PMR
	fmt.Println("Reading Providers Platforms for Organization:", color.GreenString(orgName))
	fmt.Println("Provider Name:", color.GreenString(providerName))
	fmt.Println("Provider Version:", color.GreenString(providerVersion))
	platforms, err := client.RegistryProviderPlatforms.List(ctx, tfe.RegistryProviderVersionID{
		RegistryProviderID: tfe.RegistryProviderID{
			OrganizationName: orgName,
			Namespace:        orgName,
			RegistryName:     "private", // for some reason public doesn't work...
			Name:             providerName,
		},
		Version: providerVersion,
	}, &tfe.RegistryProviderPlatformListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})

	if err != nil {
		logError(err, "failed to read provider in PMR")
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"OS", "Arch", "ID", "Filename", "Shasum"})
	for _, i := range platforms.Items {

		// Is there any additional info we can get from reading directly? likely no
		// p, err := client.RegistryProviderPlatforms.Read(ctx, tfe.RegistryProviderPlatformID{
		// 	RegistryProviderVersionID: tfe.RegistryProviderVersionID{
		// 		RegistryProviderID: tfe.RegistryProviderID{
		// 			OrganizationName: orgName,
		// 			Namespace:        orgName,
		// 			RegistryName:     "private", // for some reason public doesn't work...
		// 			Name:             providerName,
		// 		},
		// 		Version: providerVersion,
		// 	},
		// 	OS:   i.OS,
		// 	Arch: i.Arch,
		// })
		// if err != nil {
		// 	logError(err, "failed to read platform")
		// }
		t.AppendRow(table.Row{i.OS, i.Arch, i.ID, i.Filename, i.Shasum})
	}
	t.SetStyle(table.StyleRounded)
	t.Render()

	//

	// fmt.Println(platform)
	return nil
}

func registryProviderVersionPlatformCreate() error {
	// Validate flags
	orgName := *viperString("tfeOrganization")
	providerName := *viperString("name")
	providerVersion := *viperString("version")
	providerOS := *viperString("os")
	providerARCH := *viperString("arch")
	providerFilename := *viperString("filename")

	if _, err := os.Stat(providerFilename); errors.Is(err, os.ErrNotExist) {
		logError(err, "Filename does not exist")
	}

	f, err := os.Open(providerFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		log.Fatal(err)
	}
	sum := hex.EncodeToString(hash.Sum(nil))

	filename := fmt.Sprintf("terraform-provider-%s_%s_%s_%s.zip",
		providerName,
		providerVersion,
		providerOS,
		providerARCH)

	fmt.Println(orgName)
	fmt.Println(providerFilename)
	fmt.Println(filename)
	fmt.Println(sum)
	return nil
}

func registryProviderVersionPlatformShow() error {
	fmt.Println(color.MagentaString("Function not implemented yet."))
	return nil
}

func registryProviderVersionPlatformDelete() error {
	fmt.Println(color.MagentaString("Function not implemented yet."))
	return nil
}
