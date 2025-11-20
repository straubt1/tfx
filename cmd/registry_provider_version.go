// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/cmd/flags"
	view "github.com/straubt1/tfx/cmd/views"
	"github.com/straubt1/tfx/data"
	pkgfile "github.com/straubt1/tfx/pkg/file"
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
			cmdConfig, err := flags.ParseRegistryProviderVersionListFlags(cmd)
			if err != nil {
				return err
			}
			return registryProviderVersionList(cmdConfig)
		},
	}

	// `tfx registry provider version create` command
	registryProviderVersionCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a Provider Version in a Private Registry",
		Long:  "Create a Provider Version for a Provider in a Private Registry of a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseRegistryProviderVersionCreateFlags(cmd)
			if err != nil {
				return err
			}
			if _, err := viperSemanticVersionString("version"); err != nil {
				return errors.New("invalid semantic version")
			}
			if !pkgfile.IsFile(cmdConfig.Shasums) {
				return errors.New("shasums file does not exist")
			}
			if !pkgfile.IsFile(cmdConfig.ShasumsSig) {
				return errors.New("shasumssig file does not exist")
			}
			return registryProviderVersionCreate(cmdConfig)
		},
	}

	// `tfx registry provider version show` command
	registryProviderVersionShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show details of a Provider Version in a Private Registry",
		Long:  "Show details of a Provider Version for a Provider in a Private Registry of a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseRegistryProviderVersionShowFlags(cmd)
			if err != nil {
				return err
			}
			if _, err := viperSemanticVersionString("version"); err != nil {
				return errors.New("invalid semantic version")
			}
			return registryProviderVersionShow(cmdConfig)
		},
	}

	// `tfx registry provider version delete` command
	registryProviderVersionDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a Provider Version in a Private Registry",
		Long:  "Delete a Provider Version for a Provider in a Private Registry of a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseRegistryProviderVersionDeleteFlags(cmd)
			if err != nil {
				return err
			}
			if _, err := viperSemanticVersionString("version"); err != nil {
				return errors.New("invalid semantic version")
			}
			return registryProviderVersionDelete(cmdConfig)
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
	registryProviderVersionCreateCmd.Flags().StringP("key-id", "", "", "GPG Key Id")
	registryProviderVersionCreateCmd.Flags().StringP("shasums", "", "", "Path to shasums")
	registryProviderVersionCreateCmd.Flags().StringP("shasums-sig", "", "", "Path to shasumssig")
	registryProviderVersionCreateCmd.MarkFlagRequired("name")
	registryProviderVersionCreateCmd.MarkFlagRequired("version")
	registryProviderVersionCreateCmd.MarkFlagRequired("key-id")
	registryProviderVersionCreateCmd.MarkFlagRequired("shasums")
	registryProviderVersionCreateCmd.MarkFlagRequired("shasums-sig")

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

func registryProviderVersionList(cmdConfig *flags.RegistryProviderVersionListFlags) error {
	v := view.NewRegistryProviderVersionListView()
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}
	v.PrintCommandHeader("List Provider Versions in Registry for Organization: %s", c.OrganizationName)
	v.PrintCommandFilter("Provider Name: %s", cmdConfig.Name)
	items, err := data.ListRegistryProviderVersions(c, c.OrganizationName, cmdConfig.Name)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "Failed to list provider versions"))
	}
	return v.Render(items)
}

func registryProviderVersionCreate(cmdConfig *flags.RegistryProviderVersionCreateFlags) error {
	v := view.NewRegistryProviderVersionCreateView()
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}
	v.PrintCommandHeader("Create Provider Version in Registry for Organization: %s", c.OrganizationName)
	v.PrintCommandFilter("Provider Name: %s", cmdConfig.Name)
	p, err := data.CreateRegistryProviderVersion(c, c.OrganizationName, cmdConfig.Name, cmdConfig.Version, cmdConfig.KeyID)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to create provider version"))
	}
	v.Renderer().Message("Uploading shasums and sig")
	if err := UploadBinary(p.Links["shasums-upload"].(string), cmdConfig.Shasums); err != nil {
		return v.RenderError(errors.Wrap(err, "failed to upload shasums"))
	}
	if err := UploadBinary(p.Links["shasums-sig-upload"].(string), cmdConfig.ShasumsSig); err != nil {
		return v.RenderError(errors.Wrap(err, "failed to upload shasums sig"))
	}
	fmt.Println(cmdConfig.Shasums, cmdConfig.ShasumsSig, p.CreatedAt)
	return v.Render(p)
}

func registryProviderVersionShow(cmdConfig *flags.RegistryProviderVersionShowFlags) error {
	v := view.NewRegistryProviderVersionShowView()
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}
	v.PrintCommandHeader("Show Provider Version in Registry for Organization: %s", c.OrganizationName)
	provider, err := data.ReadRegistryProviderVersion(c, c.OrganizationName, cmdConfig.Name, cmdConfig.Version)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to read provider version"))
	}
	var shas string
	if provider.ShasumsUploaded {
		sha, err := DownloadTextFile(provider.Links["shasums-download"].(string))
		if err != nil {
			return v.RenderError(errors.Wrap(err, "Failed to read shasums download link"))
		}
		shas = sha
	}
	return v.Render(provider, shas)
}

func registryProviderVersionDelete(cmdConfig *flags.RegistryProviderVersionDeleteFlags) error {
	v := view.NewRegistryProviderVersionDeleteView()
	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}
	v.PrintCommandHeader("Delete Provider Version in Registry for Organization: %s", c.OrganizationName)
	if err := data.DeleteRegistryProviderVersion(c, c.OrganizationName, cmdConfig.Name, cmdConfig.Version); err != nil {
		return v.RenderError(errors.Wrap(err, "Failed to Delete Provider Version"))
	}
	return v.Render(cmdConfig.Name)
}
