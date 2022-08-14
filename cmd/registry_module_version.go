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
	"io/ioutil"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// `tfx registry module version` commands
	registryModuleVersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Module Versions in Private Registry Commands",
		Long:  "Commands to work with Module Versions in the Private Registry of a TFx Organization.",
	}

	// `tfx registry module version list` command
	registryModuleVersionListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Module Versions",
		Long:  "List Modules Versions of a module in the Private Module Registry for a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return registryModuleVersionList(
				getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("name"),
				*viperString("provider"))
		},
	}

	// `tfx registry module version create` command
	registryModuleVersionCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a Module Version",
		Long:  "Create a Module Version of a module in the Private Module Registry for a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			moduleVersion, err := viperSemanticVersionString("version")
			if err != nil {
				return errors.New("failed to parse semantic version")
			}
			if !isDirectory(*viperString("directory")) {
				return errors.New("directory file does not exist")
			}

			return registryModuleVersionCreate(
				getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("name"),
				*viperString("provider"),
				moduleVersion,
				*viperString("directory"))
		},
	}

	// `tfx registry module version delete` command
	registryModuleVersionDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a Module Version",
		Long:  "Delete a Module Version of a module in the Private Module Registry for a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			moduleVersion, err := viperSemanticVersionString("version")
			if err != nil {
				return errors.New("failed to parse semantic version")
			}

			return registryModuleVersionDelete(
				getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("name"),
				*viperString("provider"),
				moduleVersion)
		},
	}

	// `tfx registry module version download` command
	registryModuleVersionDownloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download a Module Version",
		Long:  "Download the Terraform code of Module Version in the Private Module Registry for a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			moduleVersion, err := viperSemanticVersionString("version")
			if err != nil {
				return errors.New("failed to parse semantic version")
			}

			return registryModuleVersionDownload(
				getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("name"),
				*viperString("provider"),
				moduleVersion,
				*viperString("directory"))
		},
	}
)

func init() {
	// `tfx registry module version create` arguments
	registryModuleVersionCreateCmd.Flags().StringP("name", "n", "", "Name of the Module (no spaces)")
	registryModuleVersionCreateCmd.Flags().StringP("provider", "p", "", "Name of the provider (no spaces) (i.e. aws, azure, google)")
	registryModuleVersionCreateCmd.Flags().StringP("version", "v", "", "Version of module (i.e. 0.0.1)")
	registryModuleVersionCreateCmd.Flags().StringP("directory", "d", "./", "Directory of Terraform (optional, defaults to current directory)")
	registryModuleVersionCreateCmd.MarkFlagRequired("name")
	registryModuleVersionCreateCmd.MarkFlagRequired("provider")
	registryModuleVersionCreateCmd.MarkFlagRequired("version")

	// `tfx registry module version list` arguments
	registryModuleVersionListCmd.Flags().StringP("name", "n", "", "Name of the Module (no spaces)")
	registryModuleVersionListCmd.Flags().StringP("provider", "p", "", "Name of the provider (no spaces) (i.e. aws, azure, google)")
	registryModuleVersionListCmd.MarkFlagRequired("name")
	registryModuleVersionListCmd.MarkFlagRequired("provider")

	// `tfx registry module version delete` arguments
	registryModuleVersionDeleteCmd.Flags().StringP("name", "n", "", "Name of the Module (no spaces)")
	registryModuleVersionDeleteCmd.Flags().StringP("provider", "p", "", "Name of the provider (no spaces) (i.e. aws, azure, google)")
	registryModuleVersionDeleteCmd.Flags().StringP("version", "v", "", "Version of module (i.e. 0.0.1)")
	registryModuleVersionDeleteCmd.MarkFlagRequired("name")
	registryModuleVersionDeleteCmd.MarkFlagRequired("provider")
	registryModuleVersionDeleteCmd.MarkFlagRequired("version")

	// `tfx registry module version download` arguments
	registryModuleVersionDownloadCmd.Flags().StringP("name", "n", "", "Name of the Module (no spaces)")
	registryModuleVersionDownloadCmd.Flags().StringP("provider", "p", "", "Name of the provider (no spaces) (i.e. aws, azure, google)")
	registryModuleVersionDownloadCmd.Flags().StringP("version", "v", "", "Version of module (i.e. 0.0.1)")
	registryModuleVersionDownloadCmd.Flags().StringP("directory", "d", "", "Directory to download module to (optional, defaults to a temp directory)")
	registryModuleVersionDownloadCmd.MarkFlagRequired("name")
	registryModuleVersionDownloadCmd.MarkFlagRequired("provider")
	registryModuleVersionDownloadCmd.MarkFlagRequired("version")

	registryModuleCmd.AddCommand(registryModuleVersionCmd)
	registryModuleVersionCmd.AddCommand(registryModuleVersionCreateCmd)
	registryModuleVersionCmd.AddCommand(registryModuleVersionListCmd)
	registryModuleVersionCmd.AddCommand(registryModuleVersionDeleteCmd)
	registryModuleVersionCmd.AddCommand(registryModuleVersionDownloadCmd)
}

func registryModuleVersionList(c TfxClientContext, orgName string, moduleName string, providerName string) error {
	o.AddMessageUserProvided("List Module Versions for Organization:", orgName)
	module, err := c.Client.RegistryModules.Read(c.Context, tfe.RegistryModuleID{
		Organization: orgName,
		Name:         moduleName,
		Provider:     providerName,
		Namespace:    orgName,
		RegistryName: tfe.PrivateRegistry,
	})
	if err != nil {
		return errors.Wrap(err, "failed to list module versions")
	}

	o.AddTableHeader("Version", "Status")
	for _, i := range module.VersionStatuses {
		o.AddTableRows(i.Version, i.Status)
	}
	o.Close()

	return nil
}

func registryModuleVersionCreate(c TfxClientContext, orgName string, moduleName string, providerName string,
	moduleVersion string, directory string) error {
	o.AddMessageUserProvided("Create Module Version for Organization:", orgName)
	module, err := c.Client.RegistryModules.CreateVersion(c.Context, tfe.RegistryModuleID{
		Organization: orgName,
		Name:         moduleName,
		Provider:     providerName,
		Namespace:    orgName,
		RegistryName: tfe.PrivateRegistry,
	}, tfe.RegistryModuleCreateVersionOptions{
		Version: &moduleVersion,
	})
	if err != nil {
		errors.Wrap(err, "failed to create module version")
	}
	o.AddMessageUserProvided("Module Created, Uploading...", "")
	err = c.Client.RegistryModules.Upload(c.Context, *module, directory)
	if err != nil {
		errors.Wrap(err, "failed to upload module version")
	}

	o.AddMessageUserProvided("Module Created:", module.RegistryModule.Name)
	o.AddDeferredMessageRead("ID", module.RegistryModule.ID)
	o.AddDeferredMessageRead("Created", module.CreatedAt)
	o.Close()

	return nil
}

func registryModuleVersionDelete(c TfxClientContext, orgName string, moduleName string, providerName string,
	moduleVersion string) error {
	o.AddMessageUserProvided("Delete Module Version for Organization:", orgName)
	err := c.Client.RegistryModules.DeleteVersion(c.Context, tfe.RegistryModuleID{
		Organization: orgName,
		Name:         moduleName,
		Provider:     providerName,
		Namespace:    orgName,
		RegistryName: tfe.PrivateRegistry,
	}, moduleVersion)
	if err != nil {
		return errors.Wrap(err, "failed to delete module version")
	}

	o.AddMessageUserProvided("Module Version Deleted:", moduleName)
	o.AddDeferredMessageRead("Status", "Success")
	o.Close()

	return nil
}

func registryModuleVersionDownload(c TfxClientContext, orgName string, moduleName string, providerName string,
	moduleVersion string, directory string) error {
	o.AddMessageUserProvided("Downloading Module Version:", moduleName)
	var err error
	// Determine a directory to unpack the slug contents into.
	if directory != "" {
		if !isDirectory(directory) {
			return errors.Wrap(err, "provider directory is not valid")
		}
	} else {
		o.AddMessageUserProvided("Directory not supplied, creating a temp directory", "")
		dst, err := ioutil.TempDir("", "slug")
		if err != nil {
			return errors.Wrap(err, "failed to create temp directory")
		}
		directory = dst
	}

	o.AddMessageUserProvided("Module Version Found, download started...", "")
	_, err = DownloadModule(c.Token, c.Hostname, orgName, moduleName, providerName, moduleVersion, directory)
	if err != nil {
		return errors.Wrap(err, "failed to download module")
	}

	o.AddDeferredMessageRead("Status", "Success")
	o.AddDeferredMessageRead("Directory", directory)
	o.Close()

	return nil
}
