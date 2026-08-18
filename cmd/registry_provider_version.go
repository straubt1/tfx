// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	"github.com/coreos/go-semver/semver"
	tfe "github.com/hashicorp/go-tfe"
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
		Long:  "Create a Provider Version for a Provider in a Private Registry of a TFx Organization. Pass --directory to infer name, version, GPG key, checksums, and platforms from a folder staged by tfx registry provider download.",
		Example: `
tfx registry provider version create --directory ./providers/azurerm/5.0.0

tfx registry provider version create --name azurerm --version 5.0.0 --key-id 34365D9472D7468F --shasums ./SHA256SUMS --shasums-sig ./SHA256SUMS.sig`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseRegistryProviderVersionCreateFlags(cmd)
			if err != nil {
				return err
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
			if _, err := semver.NewVersion(cmdConfig.Version); err != nil {
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
			if _, err := semver.NewVersion(cmdConfig.Version); err != nil {
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
	registryProviderVersionCreateCmd.Flags().StringP("directory", "d", "", "Staged provider directory from tfx registry provider download (infers name, version, GPG key, and platforms)")

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
	if useDirectoryMode(cmdConfig) {
		return registryProviderVersionCreateFromDirectory(cmdConfig)
	}
	if err := validateExplicitVersionCreate(cmdConfig); err != nil {
		return err
	}
	return registryProviderVersionCreateExplicit(cmdConfig)
}

func useDirectoryMode(cmdConfig *flags.RegistryProviderVersionCreateFlags) bool {
	// postInitCommands binds every command's flags into one Viper instance, so a
	// sibling --directory default (e.g. download's ./providers) can leak onto this
	// flag. Prefer the explicit five-flag path when it is fully specified.
	explicit := cmdConfig.Name != "" && cmdConfig.Version != "" && cmdConfig.KeyID != "" &&
		cmdConfig.Shasums != "" && cmdConfig.ShasumsSig != ""
	if explicit {
		return false
	}
	return cmdConfig.Directory != ""
}

func validateExplicitVersionCreate(cmdConfig *flags.RegistryProviderVersionCreateFlags) error {
	if cmdConfig.Name == "" {
		return errors.New("name is required (or pass --directory)")
	}
	if cmdConfig.Version == "" {
		return errors.New("version is required (or pass --directory)")
	}
	if _, err := semver.NewVersion(cmdConfig.Version); err != nil {
		return errors.New("invalid semantic version")
	}
	if cmdConfig.KeyID == "" {
		return errors.New("key-id is required (or pass --directory)")
	}
	if cmdConfig.Shasums == "" {
		return errors.New("shasums is required (or pass --directory)")
	}
	if cmdConfig.ShasumsSig == "" {
		return errors.New("shasums-sig is required (or pass --directory)")
	}
	if !pkgfile.IsFile(cmdConfig.Shasums) {
		return errors.New("shasums file does not exist")
	}
	if !pkgfile.IsFile(cmdConfig.ShasumsSig) {
		return errors.New("shasumssig file does not exist")
	}
	return nil
}

func registryProviderVersionCreateExplicit(cmdConfig *flags.RegistryProviderVersionCreateFlags) error {
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
	v.Output().Message("Uploading shasums and sig")
	if err := uploadVersionChecksums(p, cmdConfig.Shasums, cmdConfig.ShasumsSig); err != nil {
		return v.RenderError(err)
	}
	return v.Render(p)
}

func registryProviderVersionCreateFromDirectory(cmdConfig *flags.RegistryProviderVersionCreateFlags) error {
	v := view.NewRegistryProviderVersionCreateView()
	staged, err := data.ReadStagedProviderDirectory(cmdConfig.Directory)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to read staged provider directory"))
	}
	keyID := staged.KeyID
	if cmdConfig.KeyID != "" {
		keyID = cmdConfig.KeyID
	}
	if keyID == "" {
		return v.RenderError(errors.New("could not infer GPG key id; pass --key-id"))
	}

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}
	v.PrintCommandHeader("Create Provider Version in Registry for Organization: %s", c.OrganizationName)
	v.PrintCommandFilter("Directory: %s", cmdConfig.Directory)
	v.PrintCommandFilter("Provider Name: %s", staged.Name)
	v.PrintCommandFilter("Version: %s", staged.Version)

	_, providerCreated, err := data.EnsureRegistryProvider(c, c.OrganizationName, staged.Name)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to ensure provider"))
	}
	if providerCreated {
		v.Output().Message("Created provider %s", staged.Name)
	}

	p, err := data.CreateRegistryProviderVersion(c, c.OrganizationName, staged.Name, staged.Version, keyID)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to create provider version"))
	}
	v.Output().Message("Uploading shasums and sig")
	if err := uploadVersionChecksums(p, staged.Shasums, staged.ShasumsSig); err != nil {
		return v.RenderError(err)
	}

	var platforms []*tfe.RegistryProviderPlatform
	for _, plat := range staged.Platforms {
		v.Output().Message("Uploading %s_%s", plat.OS, plat.Arch)
		rpp, err := data.UploadRegistryProviderPlatform(c, c.OrganizationName, staged.Name, staged.Version, plat.OS, plat.Arch, plat.Path)
		if err != nil {
			return v.RenderError(errors.Wrapf(err, "failed to upload platform %s_%s", plat.OS, plat.Arch))
		}
		platforms = append(platforms, rpp)
	}

	return v.RenderFromDirectory(&view.RegistryProviderVersionCreateFromDirectoryResult{
		Name:            staged.Name,
		Version:         staged.Version,
		KeyID:           keyID,
		ProviderCreated: providerCreated,
		ProviderVersion: p,
		Platforms:       platforms,
	})
}

func uploadVersionChecksums(p *tfe.RegistryProviderVersion, shasums, shasumsSig string) error {
	shasumsURL, ok := p.Links["shasums-upload"].(string)
	if !ok || shasumsURL == "" {
		return errors.New("provider version is missing shasums-upload link")
	}
	sigURL, ok := p.Links["shasums-sig-upload"].(string)
	if !ok || sigURL == "" {
		return errors.New("provider version is missing shasums-sig-upload link")
	}
	if err := data.UploadBinary(shasumsURL, shasums); err != nil {
		return errors.Wrap(err, "failed to upload shasums")
	}
	if err := data.UploadBinary(sigURL, shasumsSig); err != nil {
		return errors.Wrap(err, "failed to upload shasums sig")
	}
	return nil
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
		sha, err := data.DownloadTextFile(provider.Links["shasums-download"].(string))
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
