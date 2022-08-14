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

	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
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
		Long:  "List Provider Versions for a Provider in a Private Registry of a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return registryProviderVersionList(
				getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("name"))
		},
	}

	// `tfx registry provider version create` command
	registryProviderVersionCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a Provider Version in a Private Registry",
		Long:  "Create a Provider Version for a Provider in a Private Registry of a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !isFile(*viperString("shasums")) {
				return errors.New("shasums file does not exist")
			}
			if !isFile(*viperString("shasumssig")) {
				return errors.New("shasumssig file does not exist")
			}

			return registryProviderVersionCreate(
				getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("name"),
				*viperString("version"),
				*viperString("keyId"),
				*viperString("shasums"),
				*viperString("shasumssig"),
			)
		},
	}

	// `tfx registry provider version show` command
	registryProviderVersionShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show details of a Provider Version in a Private Registry",
		Long:  "Show details of a Provider Version for a Provider in a Private Registry of a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return registryProviderVersionShow(getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("name"),
				*viperString("version"))
		},
	}

	// `tfx registry provider version delete` command
	registryProviderVersionDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a Provider Version in a Private Registry",
		Long:  "Delete a Provider Version for a Provider in a Private Registry of a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return registryProviderVersionDelete(getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("name"),
				*viperString("version"))
		},
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

func registryProviderVersionsListAll(c TfxClientContext, orgName string, providerName string) ([]*tfe.RegistryProviderVersion, error) {
	allItems := []*tfe.RegistryProviderVersion{}
	opts := tfe.RegistryProviderVersionListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100,
		},
	}
	for {
		items, err := c.Client.RegistryProviderVersions.List(c.Context,
			tfe.RegistryProviderID{
				OrganizationName: orgName,
				Namespace:        orgName, // always org name for RegistryName "private"
				RegistryName:     tfe.PrivateRegistry,
				Name:             providerName,
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

func registryProviderVersionList(c TfxClientContext, orgName string, providerName string) error {
	o.AddMessageUserProvided("List Provider Versions in Registry for Organization:", orgName)
	o.AddMessageUserProvided("Provider Name:", providerName)
	items, err := registryProviderVersionsListAll(c, orgName, providerName)
	if err != nil {
		return errors.Wrap(err, "Failed to list provider versions")
	}

	o.AddTableHeader("Name", "Registry", "Published", "SHASUM Uploaded", "SHASUM Sig Uploaded")
	for _, i := range items {
		o.AddTableRows(i.ID, i.Version, i.UpdatedAt, i.ShasumsUploaded, i.ShasumsSigUploaded)
	}
	o.Close()

	return nil
}

func registryProviderVersionCreate(c TfxClientContext, orgName string, providerName string, providerVersion string, keyId string, shasums string, shasumssig string) error {
	o.AddMessageUserProvided("Create Provider Version in Registry for Organization:", orgName)
	o.AddMessageUserProvided("Provider Name:", providerName)
	p, err := c.Client.RegistryProviderVersions.Create(c.Context, tfe.RegistryProviderID{
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
		return errors.Wrap(err, "failed to create provider version")
	}

	o.AddMessageUserProvided("Uploading shasums and sig", "")
	err = UploadBinary(p.Links["shasums-upload"].(string), shasums)
	if err != nil {
		logError(err, "failed to upload shasums")
	}
	err = UploadBinary(p.Links["shasums-sig-upload"].(string), shasumssig)
	if err != nil {
		logError(err, "failed to upload shasums sig")
	}
	fmt.Println(shasums, shasumssig, p.CreatedAt)

	o.AddMessageUserProvided("Provider Version Created", "")
	o.AddDeferredMessageRead("Name", providerName)
	o.AddDeferredMessageRead("ID", p.ID)
	o.AddDeferredMessageRead("Version", p.Version)
	o.AddDeferredMessageRead("Created", p.UpdatedAt)
	o.Close()

	return nil
}

func registryProviderVersionShow(c TfxClientContext, orgName string, providerName string, providerVersion string) error {
	o.AddMessageUserProvided("Show Provider Version in Registry for Organization:", orgName)
	provider, err := c.Client.RegistryProviderVersions.Read(c.Context, tfe.RegistryProviderVersionID{
		RegistryProviderID: tfe.RegistryProviderID{
			OrganizationName: orgName,
			Namespace:        orgName, // always org name for RegistryName "private"
			RegistryName:     tfe.PrivateRegistry,
			Name:             providerName,
		},
		Version: providerVersion,
	})
	if err != nil {
		return errors.Wrap(err, "failed to read provider version")
	}

	o.AddDeferredMessageRead("Name", providerName)
	o.AddDeferredMessageRead("Version", provider.Version)
	o.AddDeferredMessageRead("ID", provider.ID)
	o.AddDeferredMessageRead("Shasums Uploaded", provider.ShasumsUploaded)
	o.AddDeferredMessageRead("Shasums Sig Uploaded", provider.ShasumsSigUploaded)
	// If the Shasums have been uploaded, display them (might be a better place for this?)
	if provider.ShasumsUploaded {
		sha, err := DownloadTextFile(provider.Links["shasums-download"].(string))
		if err != nil {
			return errors.Wrap(err, "Failed to read shasums download link")
		}
		o.AddDeferredMessageRead("Shasums", "\n"+sha)
	}

	o.Close()

	return nil
}

func registryProviderVersionDelete(c TfxClientContext, orgName string, providerName string, providerVersion string) error {
	o.AddMessageUserProvided("Delete Provider Version in Registry for Organization:", orgName)
	err := c.Client.RegistryProviderVersions.Delete(c.Context, tfe.RegistryProviderVersionID{
		RegistryProviderID: tfe.RegistryProviderID{
			OrganizationName: orgName,
			Name:             providerName,
			Namespace:        orgName, // always org name for RegistryName "private"
			RegistryName:     tfe.PrivateRegistry,
		},
		Version: providerVersion,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to Delete Provider Version")
	}

	o.AddMessageUserProvided("Provider Version Deleted:", providerName)
	o.AddDeferredMessageRead("Status", "Success")
	o.Close()

	return nil
}
