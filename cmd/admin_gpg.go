// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	"os"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/cmd/flags"
	view "github.com/straubt1/tfx/cmd/views"
	"github.com/straubt1/tfx/data"
)

// gpgCmd represents the gpg command
var (
	// `tfx admin gpg` command
	gpgCmd = &cobra.Command{
		Use:   "gpg",
		Short: "GPG Keys",
		Long:  "Work with GPG Keys in the Private Registry",
		Example: `
List all GPG keys for a namespace:
tfx admin gpg list --namespace myorg

Create a GPG key from a file:
tfx admin gpg create --namespace myorg --public-key /path/to/key.asc

Show a GPG key:
tfx admin gpg show --namespace myorg --id 37AD5AEF6A5D6D5C`,
	}

	// `tfx admin gpg list` command
	gpgListCmd = &cobra.Command{
		Use:   "list",
		Short: "List GPG Keys",
		Long:  "List GPG Keys in the Private Registry for a namespace.",
		Example: `
tfx admin gpg list --namespace myorg`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseAdminGPGListFlags(cmd)
			if err != nil {
				return err
			}
			return gpgList(cmdConfig)
		},
	}

	// `tfx admin gpg create` command
	gpgCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create GPG Key",
		Long:  "Create GPG Key in the Private Registry for a namespace.",
		Example: `
tfx admin gpg create --namespace myorg --public-key /path/to/key.asc`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseAdminGPGCreateFlags(cmd)
			if err != nil {
				return err
			}
			return gpgCreate(cmdConfig)
		},
	}

	// `tfx admin gpg show` command
	gpgShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show GPG Key",
		Long:  "Show GPG Key details from the Private Registry.",
		Example: `
tfx admin gpg show --namespace myorg --id 37AD5AEF6A5D6D5C`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseAdminGPGShowFlags(cmd)
			if err != nil {
				return err
			}
			return gpgShow(cmdConfig)
		},
	}

	// `tfx admin gpg delete` command
	gpgDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete GPG Key",
		Long:  "Delete GPG Key from the Private Registry.",
		Example: `
tfx admin gpg delete --namespace myorg --id 37AD5AEF6A5D6D5C`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseAdminGPGDeleteFlags(cmd)
			if err != nil {
				return err
			}
			return gpgDelete(cmdConfig)
		},
	}
)

func init() {
	// `tfx admin gpg list` flags
	gpgListCmd.Flags().StringP("namespace", "n", "", "Namespace (typically the organization name)")
	gpgListCmd.Flags().StringP("registry-name", "r", "private", "Registry name (default: private)")

	// `tfx admin gpg create` flags
	gpgCreateCmd.Flags().StringP("namespace", "n", "", "Namespace (typically the organization name)")
	gpgCreateCmd.Flags().StringP("public-key", "k", "", "File path to the public GPG key")
	gpgCreateCmd.Flags().StringP("registry-name", "r", "private", "Registry name (default: private)")
	gpgCreateCmd.MarkFlagRequired("namespace")
	gpgCreateCmd.MarkFlagRequired("public-key")

	// `tfx admin gpg show` flags
	gpgShowCmd.Flags().StringP("namespace", "n", "", "Namespace (typically the organization name)")
	gpgShowCmd.Flags().StringP("id", "i", "", "GPG key ID")
	gpgShowCmd.Flags().StringP("registry-name", "r", "private", "Registry name (default: private)")
	gpgShowCmd.MarkFlagRequired("namespace")
	gpgShowCmd.MarkFlagRequired("id")

	// `tfx admin gpg delete` flags
	gpgDeleteCmd.Flags().StringP("namespace", "n", "", "Namespace (typically the organization name)")
	gpgDeleteCmd.Flags().StringP("id", "i", "", "GPG key ID")
	gpgDeleteCmd.Flags().StringP("registry-name", "r", "private", "Registry name (default: private)")
	gpgDeleteCmd.MarkFlagRequired("namespace")
	gpgDeleteCmd.MarkFlagRequired("id")

	adminCmd.AddCommand(gpgCmd)
	gpgCmd.AddCommand(gpgListCmd)
	gpgCmd.AddCommand(gpgCreateCmd)
	gpgCmd.AddCommand(gpgShowCmd)
	gpgCmd.AddCommand(gpgDeleteCmd)
}

func gpgList(cmdConfig *flags.AdminGPGListFlags) error {
	// Create view for rendering
	v := view.NewAdminGPGListView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	// Use namespace from flags, fallback to organization if not provided
	namespace := cmdConfig.Namespace
	if namespace == "" {
		namespace = c.OrganizationName
	}

	// Print command header
	v.PrintCommandHeader("Listing GPG keys for namespace '%s'", namespace)

	// Fetch GPG keys
	keys, err := data.FetchGPGKeys(c, namespace)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to list GPG keys"))
	}

	return v.Render(keys)
}

func gpgCreate(cmdConfig *flags.AdminGPGCreateFlags) error {
	// Create view for rendering
	v := view.NewAdminGPGCreateView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	// Validate public key file exists
	if _, err := os.Stat(cmdConfig.PublicKey); err != nil {
		if os.IsNotExist(err) {
			return v.RenderError(errors.New("public key file does not exist"))
		}
		return v.RenderError(errors.Wrap(err, "failed to access public key file"))
	}

	// Parse registry name
	registryName := tfe.RegistryName(cmdConfig.RegistryName)

	// Print command header
	v.PrintCommandHeader("Creating GPG key for namespace '%s' in registry '%s'", cmdConfig.Namespace, cmdConfig.RegistryName)

	// Create GPG key
	key, err := data.CreateGPGKey(c, registryName, cmdConfig.Namespace, cmdConfig.PublicKey)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to create GPG key"))
	}

	return v.Render(key)
}

func gpgShow(cmdConfig *flags.AdminGPGShowFlags) error {
	// Create view for rendering
	v := view.NewAdminGPGShowView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	// Parse registry name
	registryName := tfe.RegistryName(cmdConfig.RegistryName)

	// Print command header
	v.PrintCommandHeader("Showing GPG key '%s' for namespace '%s'", cmdConfig.ID, cmdConfig.Namespace)

	// Fetch GPG key
	key, err := data.FetchGPGKey(c, cmdConfig.Namespace, registryName, cmdConfig.ID)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to read GPG key"))
	}

	return v.Render(key)
}

func gpgDelete(cmdConfig *flags.AdminGPGDeleteFlags) error {
	// Create view for rendering
	v := view.NewAdminGPGDeleteView()

	c, err := client.NewFromViper()
	if err != nil {
		return v.RenderError(err)
	}

	// Parse registry name
	registryName := tfe.RegistryName(cmdConfig.RegistryName)

	// Print command header
	v.PrintCommandHeader("Deleting GPG key '%s' for namespace '%s'", cmdConfig.ID, cmdConfig.Namespace)
	// Note: The TFE API does not expose whether a GPG key is in use by a provider.
	// Deleting a key that is actively referenced by a provider version will break those operations.

	// Delete GPG key
	err = data.DeleteGPGKey(c, cmdConfig.Namespace, registryName, cmdConfig.ID)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to delete GPG key"))
	}

	return v.Render()
}
