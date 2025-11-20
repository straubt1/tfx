// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/cmd/flags"
	view "github.com/straubt1/tfx/cmd/views"
	"github.com/straubt1/tfx/data"
	pkgfile "github.com/straubt1/tfx/pkg/file"
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
		Long:  "List Provider Version Platforms for a Provider in a Private Registry of a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseRegistryProviderVersionPlatformListFlags(cmd)
			if err != nil {
				return err
			}
			if _, err := viperSemanticVersionString("version"); err != nil {
				return errors.Wrap(err, "Failed to Parse Semantic Version")
			}
			return registryProviderVersionPlatformList(cmdConfig)
		},
	}

	// `tfx registry provider version platform create` command
	registryProviderVersionPlatformCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create/Update a Provider Version Platform in a Private Registry",
		Long:  "Create/Update a Provider Version Platform for a Provider Version in a Private Registry of a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseRegistryProviderVersionPlatformCreateFlags(cmd)
			if err != nil {
				return err
			}
			if _, err := viperSemanticVersionString("version"); err != nil {
				return errors.Wrap(err, "Failed to Parse Semantic Version")
			}
			if !pkgfile.IsFile(cmdConfig.Filename) {
				return errors.New("filename does not exist")
			}
			return registryProviderVersionPlatformCreate(cmdConfig)
		},
	}

	// `tfx registry provider version platform show` command
	registryProviderVersionPlatformShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show details about a Provider Version Platform in a Private Registry",
		Long:  "Show details about a Provider Version Platform for a Provider Version in a Private Registry of a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseRegistryProviderVersionPlatformShowFlags(cmd)
			if err != nil {
				return err
			}
			if _, err := viperSemanticVersionString("version"); err != nil {
				return errors.Wrap(err, "Failed to Parse Semantic Version")
			}
			return registryProviderVersionPlatformShow(cmdConfig)
		},
	}

	// `tfx registry provider version platform delete` command
	registryProviderVersionPlatformDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a Provider Version Platform in a Private Registry",
		Long:  "Delete a Provider Version Platform for a Provider Version in a Private Registry of a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseRegistryProviderVersionPlatformDeleteFlags(cmd)
			if err != nil {
				return err
			}
			if _, err := viperSemanticVersionString("version"); err != nil {
				return errors.Wrap(err, "Failed to Parse Semantic Version")
			}
			return registryProviderVersionPlatformDelete(cmdConfig)
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

func registryProviderVersionPlatformList(cmdConfig *flags.RegistryProviderVersionPlatformListFlags) error {
	v := view.NewRegistryProviderPlatformListView()
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}
	v.PrintCommandHeader("List Provider Platforms in Registry for Organization: %s", c.OrganizationName)
	v.PrintCommandFilter("Provider Name: %s", cmdConfig.Name)
	v.PrintCommandFilter("Provider Version: %s", cmdConfig.Version)
	items, err := data.ListRegistryProviderPlatforms(c, c.OrganizationName, cmdConfig.Name, cmdConfig.Version)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to list provider version platforms"))
	}
	return v.Render(items)
}

func registryProviderVersionPlatformCreate(cmdConfig *flags.RegistryProviderVersionPlatformCreateFlags) error {
	v := view.NewRegistryProviderPlatformCreateView()
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}
	v.PrintCommandHeader("Create Provider Platform in Registry for Organization: %s", c.OrganizationName)
	f, err := os.Open(cmdConfig.Filename)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to open provider file"))
	}
	defer f.Close()

	v.Renderer().Message("Hashing Provider File")
	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return v.RenderError(errors.Wrap(err, "Failed to Hash File"))
	}
	sum := hex.EncodeToString(hash.Sum(nil))

	filename := fmt.Sprintf("terraform-provider-%s_%s_%s_%s.zip",
		cmdConfig.Name,
		cmdConfig.Version,
		cmdConfig.OS,
		cmdConfig.Arch)
	v.Renderer().Message("Building Provider Filename: %s", filename)

	rpp, err := data.CreateRegistryProviderPlatform(c, c.OrganizationName, cmdConfig.Name, cmdConfig.Version, cmdConfig.OS, cmdConfig.Arch, sum, filename)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to create provider version platform"))
	}

	v.Renderer().Message("Uploading Provider Version Platform...")
	err = UploadBinary(rpp.Links["provider-binary-upload"].(string), cmdConfig.Filename)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to upload binary to provider version platform"))
	}

	return v.Render(rpp)
}

func registryProviderVersionPlatformShow(cmdConfig *flags.RegistryProviderVersionPlatformShowFlags) error {
	v := view.NewRegistryProviderPlatformShowView()
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}
	v.PrintCommandHeader("Show Provider Platform in Registry for Organization: %s", c.OrganizationName)
	rpp, err := data.ReadRegistryProviderPlatform(c, c.OrganizationName, cmdConfig.Name, cmdConfig.Version, cmdConfig.OS, cmdConfig.Arch)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to read provider version platform"))
	}
	return v.Render(rpp)
}

func registryProviderVersionPlatformDelete(cmdConfig *flags.RegistryProviderVersionPlatformDeleteFlags) error {
	v := view.NewRegistryProviderPlatformDeleteView()
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}
	v.PrintCommandHeader("Delete Provider Platform in Registry for Organization: %s", c.OrganizationName)
	if err := data.DeleteRegistryProviderPlatform(c, c.OrganizationName, cmdConfig.Name, cmdConfig.Version, cmdConfig.OS, cmdConfig.Arch); err != nil {
		return v.RenderError(errors.Wrap(err, "failed to delete provider version platform"))
	}
	return v.Render(cmdConfig.Name)
}
