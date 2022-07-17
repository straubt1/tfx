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
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pkg/errors"
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
			providerVersion, err := viperSemanticVersionString("version")
			if err != nil {
				return errors.Wrap(err, "Failed to Parse Semantic Version")
			}

			return registryProviderVersionPlatformList(
				getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("name"),
				providerVersion)
		},
	}

	// `tfx registry provider version platform create` command
	registryProviderVersionPlatformCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create/Update a Provider Version Platform in a Private Registry",
		Long:  "Create/Update a Provider Version Platform for a Provider Version in a Private Registry of a TFx Organization. ",
		RunE: func(cmd *cobra.Command, args []string) error {
			providerVersion, err := viperSemanticVersionString("version")
			if err != nil {
				return errors.Wrap(err, "Failed to Parse Semantic Version")
			}

			providerFilename := *viperString("filename")
			if !isFile(providerFilename) {
				return errors.New("filename does not exist")
			}

			return registryProviderVersionPlatformCreate(
				getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("name"),
				providerVersion,
				*viperString("os"),
				*viperString("arch"),
				providerFilename)
		},
	}

	// `tfx registry provider version platform show` command
	registryProviderVersionPlatformShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show details about a Provider Version Platform in a Private Registry",
		Long:  "Show details about a Provider Version Platform for a Provider Version in a Private Registry of a TFx Organization. ",
		RunE: func(cmd *cobra.Command, args []string) error {
			providerVersion, err := viperSemanticVersionString("version")
			if err != nil {
				return errors.Wrap(err, "Failed to Parse Semantic Version")
			}

			return registryProviderVersionPlatformShow(
				getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("name"),
				providerVersion,
				*viperString("os"),
				*viperString("arch"))
		},
	}

	// `tfx registry provider version platform delete` command
	registryProviderVersionPlatformDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a Provider Version Platform in a Private Registry",
		Long:  "Delete a Provider Version Platform for a Provider Version in a Private Registry of a TFx Organization. ",
		RunE: func(cmd *cobra.Command, args []string) error {
			providerVersion, err := viperSemanticVersionString("version")
			if err != nil {
				return errors.Wrap(err, "Failed to Parse Semantic Version")
			}

			return registryProviderVersionPlatformDelete(
				getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("name"),
				providerVersion,
				*viperString("os"),
				*viperString("arch"))
		},
	}
)

func init() {
	// `tfx registry provider version platform list` arguments
	registryProviderVersionPlatformListCmd.Flags().StringP("name", "n", "", "Name of the Provider")
	registryProviderVersionPlatformListCmd.Flags().StringP("version", "v", "", "Version of Provider (i.e. 0.0.1)")
	registryProviderVersionPlatformListCmd.MarkFlagRequired("name")
	registryProviderVersionPlatformListCmd.MarkFlagRequired("version")

	// `tfx registry provider version platform create` arguments
	registryProviderVersionPlatformCreateCmd.Flags().StringP("name", "n", "", "Name of the Provider")
	registryProviderVersionPlatformCreateCmd.Flags().StringP("version", "v", "", "Version of Provider (i.e. 0.0.1)")
	registryProviderVersionPlatformCreateCmd.Flags().StringP("os", "", "", "OS of the Provider Version Platform (linux, windows, darwin)")
	registryProviderVersionPlatformCreateCmd.Flags().StringP("arch", "", "", "ARCH of the Provider Version Platform (amd64, arm64)")
	registryProviderVersionPlatformCreateCmd.Flags().StringP("filename", "f", "", "Path to the file that is the provider binary. Must be a zip file. Actual filename does not matter.")
	registryProviderVersionPlatformCreateCmd.MarkFlagRequired("name")
	registryProviderVersionPlatformCreateCmd.MarkFlagRequired("version")
	registryProviderVersionPlatformCreateCmd.MarkFlagRequired("os")
	registryProviderVersionPlatformCreateCmd.MarkFlagRequired("arch")
	registryProviderVersionPlatformCreateCmd.MarkFlagRequired("filename")

	// `tfx registry provider version platform show` arguments
	registryProviderVersionPlatformShowCmd.Flags().StringP("name", "n", "", "Name of the Provider")
	// registryProviderVersionPlatformShowCmd.Flags().StringP("name", "n", "", "Name of the Provider")
	registryProviderVersionPlatformShowCmd.Flags().StringP("version", "v", "", "Version of Provider (i.e. 0.0.1)")
	registryProviderVersionPlatformShowCmd.Flags().StringP("os", "", "", "OS of the Provider Version Platform (linux, windows, darwin)")
	registryProviderVersionPlatformShowCmd.Flags().StringP("arch", "", "", "ARCH of the Provider Version Platform (amd64, arm64)")
	registryProviderVersionPlatformShowCmd.MarkFlagRequired("name")
	registryProviderVersionPlatformShowCmd.MarkFlagRequired("version")
	registryProviderVersionPlatformShowCmd.MarkFlagRequired("os")
	registryProviderVersionPlatformShowCmd.MarkFlagRequired("arch")

	// `tfx registry provider version platform delete` arguments
	registryProviderVersionPlatformDeleteCmd.Flags().StringP("name", "n", "", "Name of the Provider")
	registryProviderVersionPlatformDeleteCmd.Flags().StringP("version", "v", "", "Version of Provider (i.e. 0.0.1)")
	registryProviderVersionPlatformDeleteCmd.Flags().StringP("os", "", "", "OS of the Provider Version Platform (linux, windows, darwin)")
	registryProviderVersionPlatformDeleteCmd.Flags().StringP("arch", "", "", "ARCH of the Provider Version Platform (amd64, arm64)")
	registryProviderVersionPlatformDeleteCmd.MarkFlagRequired("name")
	registryProviderVersionPlatformDeleteCmd.MarkFlagRequired("version")
	registryProviderVersionPlatformDeleteCmd.MarkFlagRequired("os")
	registryProviderVersionPlatformDeleteCmd.MarkFlagRequired("arch")

	registryProviderVersionCmd.AddCommand(registryProviderVersionPlatformCmd)
	registryProviderVersionPlatformCmd.AddCommand(registryProviderVersionPlatformListCmd)
	registryProviderVersionPlatformCmd.AddCommand(registryProviderVersionPlatformCreateCmd)
	registryProviderVersionPlatformCmd.AddCommand(registryProviderVersionPlatformShowCmd)
	registryProviderVersionPlatformCmd.AddCommand(registryProviderVersionPlatformDeleteCmd)
}

func registryProviderVersionPlatformsListAll(c TfxClientContext, orgName string, providerName string, providerVersion string) ([]*tfe.RegistryProviderPlatform, error) {
	allItems := []*tfe.RegistryProviderPlatform{}
	opts := tfe.RegistryProviderPlatformListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100,
		},
	}
	for {
		items, err := c.Client.RegistryProviderPlatforms.List(c.Context,
			tfe.RegistryProviderVersionID{
				RegistryProviderID: tfe.RegistryProviderID{
					OrganizationName: orgName,
					Namespace:        orgName,
					RegistryName:     "private", // for some reason public doesn't work...
					Name:             providerName,
				},
				Version: providerVersion,
			}, &opts)

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

func registryProviderVersionPlatformList(c TfxClientContext, orgName string, providerName string, providerVersion string) error {
	fmt.Println("Reading Providers Platforms for Organization:", color.GreenString(orgName))
	fmt.Println("Provider Name:", color.GreenString(providerName))
	fmt.Println("Provider Version:", color.GreenString(providerVersion))
	items, err := registryProviderVersionPlatformsListAll(c, orgName, providerName, providerVersion)
	if err != nil {
		return errors.Wrap(err, "Failed to List Provider Version Platforms")
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"OS", "Arch", "ID", "Filename", "Shasum"})
	for _, i := range items {
		t.AppendRow(table.Row{i.OS, i.Arch, i.ID, i.Filename, i.Shasum})
	}
	t.SetStyle(table.StyleRounded)
	t.Render()

	return nil
}

func registryProviderVersionPlatformCreate(c TfxClientContext, orgName string, providerName string, providerVersion string, providerOS string, providerARCH string, providerFilename string) error {
	f, err := os.Open(providerFilename)
	if err != nil {
		return errors.Wrap(err, "Failed to Open File")
	}
	defer f.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return errors.Wrap(err, "Failed to Hash File")
	}
	sum := hex.EncodeToString(hash.Sum(nil))

	filename := fmt.Sprintf("terraform-provider-%s_%s_%s_%s.zip",
		providerName,
		providerVersion,
		providerOS,
		providerARCH)

	fmt.Println("Create Provider Platforms for Organization:", color.GreenString(orgName))
	rpp, err := c.Client.RegistryProviderPlatforms.Create(c.Context, tfe.RegistryProviderVersionID{
		RegistryProviderID: tfe.RegistryProviderID{
			OrganizationName: orgName,
			Namespace:        orgName, // always org name for RegistryName "private
			RegistryName:     tfe.PrivateRegistry,
			Name:             providerName,
		},
		Version: providerVersion,
	}, tfe.RegistryProviderPlatformCreateOptions{
		OS:       providerOS,
		Arch:     providerARCH,
		Shasum:   sum,
		Filename: filename,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to Create Provider Version Platform")
	}

	fmt.Println("Uploading Provider Version Platform...")
	err = UploadBinary(rpp.Links["provider-binary-upload"].(string), providerFilename)
	if err != nil {
		return errors.Wrap(err, "Failed to Upload Binary to Provider Version Platform")
	}
	fmt.Println("Provider Version Platform Uploaded!")
	return nil
}

func registryProviderVersionPlatformShow(c TfxClientContext, orgName string, providerName string, providerVersion string, providerOS string, providerARCH string) error {
	fmt.Println("Delete Provider Version Platform in Registry for Organization:", color.GreenString(orgName))
	rpp, err := c.Client.RegistryProviderPlatforms.Read(c.Context, tfe.RegistryProviderPlatformID{
		RegistryProviderVersionID: tfe.RegistryProviderVersionID{
			RegistryProviderID: tfe.RegistryProviderID{
				OrganizationName: orgName,
				Namespace:        orgName, // always org name for RegistryName "private
				RegistryName:     tfe.PrivateRegistry,
				Name:             providerName,
			},
			Version: providerVersion,
		},
		OS:   providerOS,
		Arch: providerARCH,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to Read Provider Version Platform")
	}

	fmt.Println(color.BlueString("Name:     "), providerName)
	fmt.Println(color.BlueString("ID:       "), rpp.ID)
	fmt.Println(color.BlueString("Version:  "), providerVersion)
	fmt.Println(color.BlueString("OS:       "), rpp.OS)
	fmt.Println(color.BlueString("ARCH:     "), rpp.Arch)
	fmt.Println(color.BlueString("Filename: "), rpp.Filename)
	fmt.Println(color.BlueString("Shasum:   "), rpp.Shasum)

	return nil
}

func registryProviderVersionPlatformDelete(c TfxClientContext, orgName string, providerName string, providerVersion string, providerOS string, providerARCH string) error {
	fmt.Println("Delete Provider Version Platform in Registry for Organization:", color.GreenString(orgName))
	err := c.Client.RegistryProviderPlatforms.Delete(c.Context, tfe.RegistryProviderPlatformID{
		RegistryProviderVersionID: tfe.RegistryProviderVersionID{
			RegistryProviderID: tfe.RegistryProviderID{
				OrganizationName: orgName,
				Namespace:        orgName, // always org name for RegistryName "private
				RegistryName:     tfe.PrivateRegistry,
				Name:             providerName,
			},
			Version: providerVersion,
		},
		OS:   providerOS,
		Arch: providerARCH,
	})

	if err != nil {
		return errors.Wrap(err, "Failed to Delete Provider Version Platform")
	}

	fmt.Println(color.BlueString("Provider Version Platform Deleted: "), providerName, providerVersion)
	return nil
}
