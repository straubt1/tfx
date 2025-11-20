// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	"github.com/coreos/go-semver/semver"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/cmd/flags"
	view "github.com/straubt1/tfx/cmd/views"
	"github.com/straubt1/tfx/data"
	pkgfile "github.com/straubt1/tfx/pkg/file"
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
			cmdConfig, err := flags.ParseRegistryModuleVersionListFlags(cmd)
			if err != nil {
				return err
			}
			return registryModuleVersionList(cmdConfig)
		},
	}

	// `tfx registry module version create` command
	registryModuleVersionCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a Module Version",
		Long:  "Create a Module Version of a module in the Private Module Registry for a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseRegistryModuleVersionCreateFlags(cmd)
			if err != nil {
				return err
			}
			if _, err := semver.NewVersion(cmdConfig.Version); err != nil {
				return errors.New("failed to parse semantic version")
			}
			if !pkgfile.IsDirectory(cmdConfig.Directory) {
				return errors.New("directory file does not exist")
			}

			return registryModuleVersionCreate(cmdConfig)
		},
	}

	// `tfx registry module version delete` command
	registryModuleVersionDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a Module Version",
		Long:  "Delete a Module Version of a module in the Private Module Registry for a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseRegistryModuleVersionDeleteFlags(cmd)
			if err != nil {
				return err
			}
			if _, err := semver.NewVersion(cmdConfig.Version); err != nil {
				return errors.New("failed to parse semantic version")
			}
			return registryModuleVersionDelete(cmdConfig)
		},
	}

	// `tfx registry module version download` command
	registryModuleVersionDownloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download a Module Version",
		Long:  "Download the Terraform code of Module Version in the Private Module Registry for a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseRegistryModuleVersionDownloadFlags(cmd)
			if err != nil {
				return err
			}
			if _, err := semver.NewVersion(cmdConfig.Version); err != nil {
				return errors.New("failed to parse semantic version")
			}
			directory, err := pkgfile.GetDirectory(cmdConfig.Directory, cmdConfig.Name, cmdConfig.Provider, cmdConfig.Version)
			if err != nil {
				return err
			}
			cmdConfig.Directory = directory
			return registryModuleVersionDownload(cmdConfig)
		},
	}
)

func init() {
	// `tfx registry module version list` arguments
	registryModuleVersionListCmd.Flags().StringP("name", "n", "", "Name of the Module (no spaces)")
	registryModuleVersionListCmd.Flags().StringP("provider", "p", "", "Name of the provider (no spaces) (i.e. aws, azure, google)")
	registryModuleVersionListCmd.MarkFlagRequired("name")
	registryModuleVersionListCmd.MarkFlagRequired("provider")

	// `tfx registry module version create` arguments
	registryModuleVersionCreateCmd.Flags().StringP("name", "n", "", "Name of the Module (no spaces)")
	registryModuleVersionCreateCmd.Flags().StringP("provider", "p", "", "Name of the provider (no spaces) (i.e. aws, azure, google)")
	registryModuleVersionCreateCmd.Flags().StringP("version", "v", "", "Version of module (i.e. 0.0.1)")
	registryModuleVersionCreateCmd.Flags().StringP("directory", "d", "./", "Directory of Terraform (optional, defaults to current directory)")
	registryModuleVersionCreateCmd.MarkFlagRequired("name")
	registryModuleVersionCreateCmd.MarkFlagRequired("provider")
	registryModuleVersionCreateCmd.MarkFlagRequired("version")

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

func registryModuleVersionList(cmdConfig *flags.RegistryModuleVersionListFlags) error {
	v := view.NewRegistryModuleVersionListView()
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}
	v.PrintCommandHeader("List Module Versions for Organization: %s", c.OrganizationName)
	module, err := data.ListRegistryModuleVersions(c, c.OrganizationName, cmdConfig.Name, cmdConfig.Provider)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to list module versions"))
	}
	return v.Render(module)
}

func registryModuleVersionCreate(cmdConfig *flags.RegistryModuleVersionCreateFlags) error {
	v := view.NewRegistryModuleVersionCreateView()
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}
	v.PrintCommandHeader("Create Module Version for Organization: %s", c.OrganizationName)
	module, err := data.CreateRegistryModuleVersion(c, c.OrganizationName, cmdConfig.Name, cmdConfig.Provider, cmdConfig.Version, cmdConfig.Directory)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to create module version"))
	}
	return v.Render(module)
}

func registryModuleVersionDelete(cmdConfig *flags.RegistryModuleVersionDeleteFlags) error {
	v := view.NewRegistryModuleVersionDeleteView()
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}
	v.PrintCommandHeader("Delete Module Version for Organization: %s", c.OrganizationName)
	err = data.DeleteRegistryModuleVersion(c, c.OrganizationName, cmdConfig.Name, cmdConfig.Provider, cmdConfig.Version)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to delete module version"))
	}
	return v.Render(cmdConfig.Name)
}

func registryModuleVersionDownload(cmdConfig *flags.RegistryModuleVersionDownloadFlags) error {
	// Use existing REST helper for download; render simple output
	v := view.NewBaseView() // simple renderer sufficient
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}
	v.PrintCommandHeader("Downloading Module Version: %s", cmdConfig.Name)
	_, err = DownloadModule(c, cmdConfig.Name, cmdConfig.Provider, cmdConfig.Version, cmdConfig.Directory)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to download module"))
	}
	return v.Renderer().RenderProperties([]view.PropertyPair{{Key: "Status", Value: "Success"}, {Key: "Directory", Value: cmdConfig.Directory}})
}
