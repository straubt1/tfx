// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	"math"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/cmd/flags"
	view "github.com/straubt1/tfx/cmd/views"
	"github.com/straubt1/tfx/data"
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
			cmdConfig, err := flags.ParseRegistryProviderListFlags(cmd)
			if err != nil {
				return err
			}
			if cmdConfig.All {
				cmdConfig.MaxItems = math.MaxInt
			}
			return registryProviderList(cmdConfig)
		},
	}

	// `tfx registry provider create` command
	registryProviderCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a Provider in a Private Registry",
		Long:  "Create a Provider in a Private Registry of a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseRegistryProviderCreateFlags(cmd)
			if err != nil {
				return err
			}
			return registryProviderCreate(cmdConfig)
		},
	}

	// `tfx registry provider show` command
	registryProviderShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show details about a Provider in a Private Registry",
		Long:  "Show details about a Provider in a Private Registry of a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseRegistryProviderShowFlags(cmd)
			if err != nil {
				return err
			}
			return registryProviderShow(cmdConfig)
		},
	}

	// `tfx registry provider delete` command
	registryProviderDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a Provider in a Private Registry",
		Long:  "Delete a Provider in a Private Registry of a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseRegistryProviderDeleteFlags(cmd)
			if err != nil {
				return err
			}
			return registryProviderDelete(cmdConfig)
		},
	}

	// `tfx registry provider download` command
	registryProviderDownloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download a Provider from the Public Registry",
		Long:  "Download provider binaries, SHA256SUMS, and signature from the public Terraform Registry and stage them for upload to a private registry.",
		Example: `
tfx registry provider download --name azurerm --version 5.0.0

tfx registry provider download --name azurerm --version 5.0.0 --platforms linux_amd64

tfx registry provider download --namespace hashicorp --name azurerm --version 5.0.0 --directory ./providers`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseRegistryProviderDownloadFlags(cmd)
			if err != nil {
				return err
			}
			if _, err := semver.NewVersion(cmdConfig.Version); err != nil {
				return errors.New("invalid semantic version")
			}
			return registryProviderDownload(cmdConfig)
		},
	}
)

func init() {
	// `tfx registry provider list` arguments
	registryProviderListCmd.Flags().IntP("max-items", "m", 10, "Max number of results (optional)")
	registryProviderListCmd.Flags().BoolP("all", "a", false, "Retrieve all results regardless of maxItems flag (optional)")

	// `tfx registry provider create` arguments
	registryProviderCreateCmd.Flags().StringP("name", "n", "", "Name of the Provider")
	registryProviderCreateCmd.MarkFlagRequired("name")

	// `tfx registry provider show` arguments
	registryProviderShowCmd.Flags().StringP("name", "n", "", "Name of the Provider")
	registryProviderShowCmd.MarkFlagRequired("name")

	// `tfx registry provider delete` arguments
	registryProviderDeleteCmd.Flags().StringP("name", "n", "", "Name of the Provider")
	registryProviderDeleteCmd.MarkFlagRequired("name")

	// `tfx registry provider download` arguments
	registryProviderDownloadCmd.Flags().String("namespace", "hashicorp", "Public registry namespace")
	registryProviderDownloadCmd.Flags().StringP("name", "n", "", "Name of the Provider")
	registryProviderDownloadCmd.Flags().StringP("version", "v", "", "Version of Provider (i.e. 5.0.0)")
	registryProviderDownloadCmd.Flags().StringSlice("platforms", flags.DefaultPublicProviderPlatforms, "Platforms to download as os_arch (comma separated)")
	registryProviderDownloadCmd.Flags().Bool("all-platforms", false, "Download every platform published for this version")
	registryProviderDownloadCmd.Flags().StringP("directory", "d", flags.DefaultPublicProviderDirectory, "Base directory to stage files into (<directory>/<name>/<version>)")
	registryProviderDownloadCmd.MarkFlagRequired("name")
	registryProviderDownloadCmd.MarkFlagRequired("version")

	registryCmd.AddCommand(registryProviderCmd)
	registryProviderCmd.AddCommand(registryProviderListCmd)
	registryProviderCmd.AddCommand(registryProviderCreateCmd)
	registryProviderCmd.AddCommand(registryProviderShowCmd)
	registryProviderCmd.AddCommand(registryProviderDeleteCmd)
	registryProviderCmd.AddCommand(registryProviderDownloadCmd)
}

func registryProviderList(cmdConfig *flags.RegistryProviderListFlags) error {
	v := view.NewRegistryProviderListView()
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}
	v.PrintCommandHeader("List Providers in Registry for Organization: %s", c.OrganizationName)
	items, err := data.ListRegistryProviders(c, c.OrganizationName, cmdConfig.MaxItems)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to list providers"))
	}
	return v.Render(items)
}

func registryProviderCreate(cmdConfig *flags.RegistryProviderCreateFlags) error {
	v := view.NewRegistryProviderCreateView()
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}
	v.PrintCommandHeader("Create Provider in Registry for Organization: %s", c.OrganizationName)
	provider, err := data.CreateRegistryProvider(c, c.OrganizationName, cmdConfig.Name)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to create provider"))
	}
	return v.Render(provider)
}

func registryProviderShow(cmdConfig *flags.RegistryProviderShowFlags) error {
	v := view.NewRegistryProviderShowView()
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}
	v.PrintCommandHeader("Show Provider in Registry for Organization: %s", c.OrganizationName)
	provider, err := data.ReadRegistryProvider(c, c.OrganizationName, cmdConfig.Name)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to read provider"))
	}
	return v.Render(provider)
}

func registryProviderDelete(cmdConfig *flags.RegistryProviderDeleteFlags) error {
	v := view.NewRegistryProviderDeleteView()
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}
	v.PrintCommandHeader("Delete Provider in Registry for Organization: %s", c.OrganizationName)
	if err := data.DeleteRegistryProvider(c, c.OrganizationName, cmdConfig.Name); err != nil {
		return v.RenderError(errors.Wrap(err, "failed to delete provider"))
	}
	return v.Render(cmdConfig.Name)
}

func registryProviderDownload(cmdConfig *flags.RegistryProviderDownloadFlags) error {
	v := view.NewRegistryProviderDownloadView()
	v.PrintCommandHeader("Download Provider from Public Registry")
	v.PrintCommandFilter("Provider: %s/%s %s", cmdConfig.Namespace, cmdConfig.Name, cmdConfig.Version)
	if cmdConfig.AllPlatforms {
		v.PrintCommandFilter("Platforms: all")
	} else {
		v.PrintCommandFilter("Platforms: %s", strings.Join(cmdConfig.Platforms, ", "))
	}
	v.Renderer().Message("Downloading provider artifacts...")

	result, err := data.DownloadPublicProvider(data.PublicProviderDownloadConfig{
		Namespace:    cmdConfig.Namespace,
		Name:         cmdConfig.Name,
		Version:      cmdConfig.Version,
		Directory:    cmdConfig.Directory,
		Platforms:    cmdConfig.Platforms,
		AllPlatforms: cmdConfig.AllPlatforms,
	})
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to download public provider"))
	}
	return v.Render(result)
}
