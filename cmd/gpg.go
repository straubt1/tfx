/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/jedib0t/go-pretty/v6/table"
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
			return gpgList()
		},
		PreRun: bindPFlags,
	}

	gpgCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create GPG Key",
		Long:  "Create GPG Key for a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return gpgCreate()
		},
		PreRun: bindPFlags,
	}

	gpgShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show GPG Key",
		Long:  "Show GPG Key for a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return gpgShow()
		},
		PreRun: bindPFlags,
	}

	gpgDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete GPG Key",
		Long:  "Delete GPG Key for a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return gpgDelete()
		},
		PreRun: bindPFlags,
	}
)

func init() {
	// `tfx gpg list`
	// none

	// `tfx gpg create`
	gpgCreateCmd.Flags().StringP("namespace", "n", "", "Namespace (typically the organization name)")
	gpgCreateCmd.Flags().StringP("publicKey", "k", "", "File path to the public GPG key")
	gpgCreateCmd.Flags().StringP("registryName", "r", "private", "Registry name (optional, defaults to 'private')")
	gpgCreateCmd.MarkFlagRequired("namespace")
	gpgCreateCmd.MarkFlagRequired("publicKey")

	// `tfx gpg show`
	gpgShowCmd.Flags().StringP("namespace", "n", "", "Namespace (typically the organization name)")
	gpgShowCmd.Flags().StringP("keyId", "k", "", "GPG key Id")
	gpgShowCmd.Flags().StringP("registryName", "r", "private", "Registry name (optional, defaults to 'private')")
	gpgShowCmd.MarkFlagRequired("namespace")
	gpgShowCmd.MarkFlagRequired("keyId")

	// `tfx gpg delete`
	gpgDeleteCmd.Flags().StringP("namespace", "n", "", "Namespace (typically the organization name)")
	gpgDeleteCmd.Flags().StringP("keyId", "k", "", "GPG key Id")
	gpgDeleteCmd.Flags().StringP("registryName", "r", "private", "Registry name (optional, defaults to 'private')")
	gpgDeleteCmd.MarkFlagRequired("namespace")
	gpgDeleteCmd.MarkFlagRequired("keyId")

	rootCmd.AddCommand(gpgCmd)
	gpgCmd.AddCommand(gpgListCmd)
	gpgCmd.AddCommand(gpgCreateCmd)
	gpgCmd.AddCommand(gpgShowCmd)
	gpgCmd.AddCommand(gpgDeleteCmd)
}

func gpgList() error {
	// Validate flags
	hostname := *viperString("tfeHostname")
	token := *viperString("tfeToken")
	orgName := *viperString("tfeOrganization")

	// Read all GPG Keys in Org
	fmt.Println("Reading GPG Keys for Organization:", color.GreenString(orgName))
	gpg, err := ListGPGKeys(token, hostname, orgName)
	if err != nil {
		logError(err, "failed to read GPG Keys")
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Key Id", "Namespace", "Updated At", "Created At"})
	for _, i := range gpg.Keys {
		t.AppendRow(table.Row{i.Attributes.KeyID, i.Attributes.Namespace, i.Attributes.UpdatedAt.String(), i.Attributes.CreatedAt.String()})
	}
	t.SetStyle(table.StyleRounded)
	t.Render()

	return nil
}

func gpgCreate() error {
	// Validate flags
	namespace := *viperString("namespace")
	publicKey := *viperString("publicKey")
	registryName := *viperString("registryName")
	client, ctx := getClientContext()

	// Verify file exists
	b, err := ioutil.ReadFile(publicKey)
	if err != nil {
		fmt.Print(err)
	}
	publicKeyContents := string(b)

	// Create GPG Key
	fmt.Print("Creating GPG Key ", color.GreenString(namespace), "/", color.GreenString(registryName), " ... ")
	key, err := client.GPGKeys.Create(ctx, tfe.RegistryName(registryName), tfe.GPGKeyCreateOptions{
		Namespace:  namespace,
		AsciiArmor: publicKeyContents,
	})
	if err != nil {
		logError(err, "failed to create GPG Key")
	}
	fmt.Println(" Created with ID: ", color.BlueString(key.KeyID))

	return nil
}

func gpgShow() error {
	// Validate flags
	namespace := *viperString("namespace")
	keyId := *viperString("keyId")
	client, ctx := getClientContext()

	// Show GPG Key
	fmt.Print("Showing GPG Key ", color.GreenString(namespace), "/", color.GreenString(keyId), " ...")
	pmr, err := client.GPGKeys.Read(ctx, tfe.GPGKeyID{
		Namespace:    namespace,
		RegistryName: tfe.PrivateRegistry,
		KeyID:        keyId,
	})
	if err != nil {
		logError(err, "failed to show module")
	}
	fmt.Println(" Found")
	fmt.Println(color.BlueString("ID:        "), pmr.KeyID)
	fmt.Println(color.BlueString("Created:   "), pmr.CreatedAt)
	fmt.Println(color.BlueString("Updated:   "), pmr.UpdatedAt)
	fmt.Println(color.BlueString("Status:    "), pmr.AsciiArmor)

	return nil
}

func gpgDelete() error {
	// Validate flags
	namespace := *viperString("namespace")
	keyId := *viperString("keyId")
	client, ctx := getClientContext()

	// Delete GPG Key
	fmt.Print("Deleting GPG Key ", color.GreenString(namespace), "/", color.GreenString(keyId), " ... ")
	err := client.GPGKeys.Delete(ctx, tfe.GPGKeyID{
		Namespace:    namespace,
		RegistryName: tfe.PrivateRegistry,
		KeyID:        keyId,
	})
	if err != nil {
		logError(err, "failed to delete GPG Key")
	}
	fmt.Println(" Deleted")

	return nil
}
