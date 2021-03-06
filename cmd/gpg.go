/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"io/ioutil"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// gpgCmd represents the gpg command
var (
	gpgCmd = &cobra.Command{
		Use:   "gpg",
		Short: "GPG Keys",
		Long:  "Work with GPG Keys of a TFx Organization.",
	}

	gpgListCmd = &cobra.Command{
		Use:   "list",
		Short: "List GPG Keys",
		Long:  "List GPG Keys of a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return gpgList(
				getTfxClientContext())
		},
	}

	gpgCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create GPG Key",
		Long:  "Create GPG Key for a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !isFile(*viperString("publicKey")) {
				return errors.New("publicKey file does not exist")
			}

			return gpgCreate(
				getTfxClientContext(),
				*viperString("namespace"),
				*viperString("publicKey"),
				*viperString("registryName"))
		},
	}

	gpgShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show GPG Key",
		Long:  "Show GPG Key for a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return gpgShow(
				getTfxClientContext(),
				*viperString("namespace"),
				*viperString("keyId"))
		},
	}

	gpgDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete GPG Key",
		Long:  "Delete GPG Key for a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return gpgDelete(
				getTfxClientContext(),
				*viperString("namespace"),
				*viperString("keyId"))
		},
	}
)

func init() {
	// `tfx gpg list`

	// `tfx gpg create`
	gpgCreateCmd.Flags().StringP("namespace", "n", "", "Namespace (typically the organization name)")
	gpgCreateCmd.Flags().StringP("publicKey", "k", "", "File path to the public GPG key")
	gpgCreateCmd.Flags().StringP("registryName", "r", "private", "Registry name")
	gpgCreateCmd.MarkFlagRequired("namespace")
	gpgCreateCmd.MarkFlagRequired("publicKey")

	// `tfx gpg show`
	gpgShowCmd.Flags().StringP("namespace", "n", "", "Namespace (typically the organization name)")
	gpgShowCmd.Flags().StringP("keyId", "k", "", "GPG key Id")
	gpgShowCmd.Flags().StringP("registryName", "r", "private", "Registry name")
	gpgShowCmd.MarkFlagRequired("namespace")
	gpgShowCmd.MarkFlagRequired("keyId")

	// `tfx gpg delete`
	gpgDeleteCmd.Flags().StringP("namespace", "n", "", "Namespace (typically the organization name)")
	gpgDeleteCmd.Flags().StringP("keyId", "k", "", "GPG key Id")
	gpgDeleteCmd.Flags().StringP("registryName", "r", "private", "Registry name")
	gpgDeleteCmd.MarkFlagRequired("namespace")
	gpgDeleteCmd.MarkFlagRequired("keyId")

	rootCmd.AddCommand(gpgCmd)
	gpgCmd.AddCommand(gpgListCmd)
	gpgCmd.AddCommand(gpgCreateCmd)
	gpgCmd.AddCommand(gpgShowCmd)
	gpgCmd.AddCommand(gpgDeleteCmd)
}

func gpgList(c TfxClientContext) error {
	o.AddMessageUserProvided("List GPG Keys for Organization:", c.OrganizationName)
	gpg, err := ListGPGKeys(c)
	if err != nil {
		return errors.Wrap(err, "unable to list gpg keys")
	}

	o.AddTableHeader("Key Id", "Namespace", "Updated At", "Created At")
	for _, i := range gpg.Keys {
		o.AddTableRows(i.Attributes.KeyID, i.Attributes.Namespace, FormatDateTime(i.Attributes.UpdatedAt), FormatDateTime(i.Attributes.CreatedAt))
	}
	o.Close()

	return nil
}

func gpgCreate(c TfxClientContext, namespace string, publicKey string, registryName string) error {
	o.AddMessageUserProvided("Create GPG Key for Organization:", c.OrganizationName)
	b, err := ioutil.ReadFile(publicKey)
	if err != nil {
		return errors.Wrap(err, "failed to read publicKey file")
	}
	publicKeyContents := string(b)

	g, err := c.Client.GPGKeys.Create(c.Context, tfe.RegistryName(registryName), tfe.GPGKeyCreateOptions{
		Namespace:  namespace,
		AsciiArmor: publicKeyContents,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create gpg key")
	}

	o.AddMessageUserProvided("GPG Key Created", "")
	o.AddDeferredMessageRead("ID", g.ID)
	o.AddDeferredMessageRead("Created", FormatDateTime(g.CreatedAt))
	o.AddDeferredMessageRead("Updated", FormatDateTime(g.UpdatedAt))
	o.AddDeferredMessageRead("AsciiArmor", "\n"+g.AsciiArmor)
	o.Close()

	return nil
}

func gpgShow(c TfxClientContext, namespace string, keyId string) error {
	o.AddMessageUserProvided("Show a GPG Key for Organization:", c.OrganizationName)
	g, err := c.Client.GPGKeys.Read(c.Context, tfe.GPGKeyID{
		Namespace:    namespace,
		RegistryName: tfe.PrivateRegistry,
		KeyID:        keyId,
	})
	if err != nil {
		return errors.Wrap(err, "failed to read gpg key")
	}

	o.AddMessageUserProvided("GPG Key Found", "")
	o.AddDeferredMessageRead("ID", g.ID)
	o.AddDeferredMessageRead("Created", FormatDateTime(g.CreatedAt))
	o.AddDeferredMessageRead("Updated", FormatDateTime(g.UpdatedAt))
	o.AddDeferredMessageRead("AsciiArmor", "\n"+g.AsciiArmor)
	o.Close()

	return nil
}

func gpgDelete(c TfxClientContext, namespace string, keyId string) error {
	o.AddMessageUserProvided("Delete GPG Key for Organization:", c.OrganizationName)
	// TODO: verify GPG key is not in use before deleting

	err := c.Client.GPGKeys.Delete(c.Context, tfe.GPGKeyID{
		Namespace:    namespace,
		RegistryName: tfe.PrivateRegistry,
		KeyID:        keyId,
	})
	if err != nil {
		return errors.Wrap(err, "failed to delete gpg key")
	}
	o.AddMessageUserProvided("GPG Key Deleted", "")
	o.AddDeferredMessageRead("Status", "Success")
	o.Close()

	return nil
}
