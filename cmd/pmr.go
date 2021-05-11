/*
Copyright © 2021 Tom Straub <tstraub@hashicorp.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	pmrCmd = &cobra.Command{
		Use:   "pmr",
		Short: "Private Module Registry",
		Long:  "Work with Private Module Registry of a TFx Organization.",
	}

	pmrListCmd = &cobra.Command{
		Use: "list",
		// Aliases: []string{"ls"},
		Short: "List Private Module Registry",
		Long:  "List Private Module Registry of a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pmrList()
		},
		PreRun: bindPFlags,
	}

	pmrCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create Private Module",
		Long:  "Create Private Module for a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pmrCreate()
		},
		PreRun: bindPFlags,
	}

	pmrCreateVersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Create Private Module Version",
		Long:  "Create Private Module Version for a Private Module.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pmrCreateVersion()
		},
		PreRun: bindPFlags,
	}

	pmrShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show Private Module",
		Long:  "Show Private Module details for a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pmrShow()
		},
		PreRun: bindPFlags,
	}

	pmrShowVersionsCmd = &cobra.Command{
		Use:   "versions",
		Short: "Show Private Module Versions",
		Long:  "Show Private Module Versions for a Module.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pmrShowVersions()
		},
		PreRun: bindPFlags,
	}

	pmrDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete Private Module",
		Long:  "Delete Private Module details for a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pmrDelete()
		},
		PreRun: bindPFlags,
	}

	pmrDeleteVersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Delete Private Module Version",
		Long:  "Delete Private Module Version for a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pmrDeleteVersion()
		},
		PreRun: bindPFlags,
	}
)

func init() {
	// All `tfx pmr` commands
	// pmrCmd.PersistentFlags().StringP("workspaceName", "w", "", "Workspace name")

	// `tfx pmr create`
	pmrCreateCmd.Flags().StringP("name", "n", "", "Name of the Module (no spaces)")
	pmrCreateCmd.Flags().StringP("provider", "p", "", "Name of the provider (no spaces) (i.e. aws, azure, google)")

	// `tfx pmr create version`
	pmrCreateVersionCmd.Flags().StringP("name", "n", "", "Name of the Module (no spaces)")
	pmrCreateVersionCmd.Flags().StringP("provider", "p", "", "Name of the provider (no spaces) (i.e. aws, azure, google)")
	pmrCreateVersionCmd.Flags().String("moduleVersion", "", "Version of module (i.e. 0.0.1)")
	pmrCreateVersionCmd.Flags().StringP("directory", "d", "./", "Directory of Terraform (defaults to current directory)")

	// `tfx pmr show`
	pmrShowCmd.Flags().StringP("name", "n", "", "Name of the Module (no spaces)")
	pmrShowCmd.Flags().StringP("provider", "p", "", "Name of the provider (no spaces) (i.e. aws, azure, google)")

	// `tfx pmr show versions`
	pmrShowVersionsCmd.Flags().StringP("name", "n", "", "Name of the Module (no spaces)")
	pmrShowVersionsCmd.Flags().StringP("provider", "p", "", "Name of the provider (no spaces) (i.e. aws, azure, google)")

	// `tfx pmr delete`
	pmrDeleteCmd.Flags().StringP("name", "n", "", "Name of the Module (no spaces)")
	pmrDeleteCmd.Flags().StringP("provider", "p", "", "Name of the provider (no spaces) (i.e. aws, azure, google)")

	// `tfx pmr delete version`
	pmrDeleteVersionCmd.Flags().StringP("name", "n", "", "Name of the Module (no spaces)")
	pmrDeleteVersionCmd.Flags().StringP("provider", "p", "", "Name of the provider (no spaces) (i.e. aws, azure, google)")
	pmrDeleteVersionCmd.Flags().String("moduleVersion", "", "Version of module (i.e. 0.0.1)")

	rootCmd.AddCommand(pmrCmd)
	pmrCmd.AddCommand(pmrListCmd)
	pmrCmd.AddCommand(pmrCreateCmd)
	pmrCreateCmd.AddCommand(pmrCreateVersionCmd)
	pmrCmd.AddCommand(pmrShowCmd)
	pmrShowCmd.AddCommand(pmrShowVersionsCmd)
	pmrCmd.AddCommand(pmrDeleteCmd)
	pmrDeleteCmd.AddCommand(pmrDeleteVersionCmd)
}

func pmrList() error {
	// Validate flags
	orgName := *viperString("tfeOrganization")
	wsName := *viperString("workspaceName")
	client, ctx := getClientContext()

	// Read workspace
	fmt.Print("Reading Workspace ", color.GreenString(wsName), " for ID...")
	w, err := client.Workspaces.Read(ctx, orgName, wsName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(" Found:", w.ID)

	// Get all config versions and show the current config
	pmr, err := client.ConfigurationVersions.List(ctx, w.ID, tfe.ConfigurationVersionListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 10,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Id", "Speculative", "Status"})
	for _, i := range pmr.Items {
		t.AppendRow(table.Row{i.ID, i.Speculative, i.Status})
	}
	t.SetStyle(table.StyleRounded)
	t.Render()

	return nil
}

func pmrCreate() error {
	// Validate flags
	orgName := *viperString("tfeOrganization")
	moduleName := *viperString("name")
	providerName := *viperString("provider")
	client, ctx := getClientContext()

	// Create Module
	fmt.Print("Creating Private Module ...")
	pmr, err := client.RegistryModules.Create(ctx, orgName, tfe.RegistryModuleCreateOptions{
		Name:     tfe.String(moduleName),
		Provider: tfe.String(providerName),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(" ID:", pmr.ID)

	return nil
}

func pmrCreateVersion() error {
	// Validate flags
	hostname := *viperString("tfeHostname")
	token := *viperString("tfeToken")
	orgName := *viperString("tfeOrganization")
	moduleName := *viperString("name")
	providerName := *viperString("provider")
	moduleVersion := *viperString("moduleVersion")
	dir := *viperString("directory")
	// client, ctx := getClientContext()

	var err error

	// // Read module
	// var r *tfe.RegistryModule
	// r, err = client.RegistryModules.Read(ctx, orgName, moduleName, providerName)
	// _ = r
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _, i := range r.VersionStatuses {
	// 	fmt.Println(i.Version, i.Status, i.Error)
	// }

	// create module version to get URL
	var url *string
	url, err = RegistryModulesCreateVersion(token, hostname, orgName,
		moduleName, providerName, moduleVersion)
	if err != nil {
		return err
	}
	// fmt.Println(url)
	RegistryModulesUpload(token, url, dir)

	fmt.Println("Module Version Created", moduleVersion)
	return nil
}

func pmrShow() error {
	// Validate flags
	orgName := *viperString("tfeOrganization")
	moduleName := *viperString("name")
	providerName := *viperString("provider")
	client, ctx := getClientContext()

	// Read Config Version
	fmt.Print("Reading Module for ", color.GreenString(moduleName), "/", color.GreenString(providerName), "...")
	pmr, err := client.RegistryModules.Read(ctx, orgName, moduleName, providerName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(" pmr Found")
	fmt.Println(color.BlueString("ID:        "), pmr.ID)
	fmt.Println(color.BlueString("Status:    "), pmr.Status)
	fmt.Println(color.BlueString("Version Count:  "), len(pmr.VersionStatuses))
	fmt.Println(color.BlueString("Created:   "), pmr.CreatedAt)
	fmt.Println(color.BlueString("Updated:   "), pmr.UpdatedAt)

	return nil
}

func pmrShowVersions() error {
	// Validate flags
	orgName := *viperString("tfeOrganization")
	moduleName := *viperString("name")
	providerName := *viperString("provider")
	client, ctx := getClientContext()

	// Read Config Version
	fmt.Print("Reading Module for ", color.GreenString(moduleName), "/", color.GreenString(providerName), "...")
	pmr, err := client.RegistryModules.Read(ctx, orgName, moduleName, providerName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(" pmr Found")

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Version", "Status"})
	for _, i := range pmr.VersionStatuses {
		t.AppendRow(table.Row{i.Version, i.Status})
	}
	t.SetStyle(table.StyleRounded)
	t.Render()

	return nil
}

func pmrDelete() error {
	return nil
}

func pmrDeleteVersion() error {
	// Validate flags
	orgName := *viperString("tfeOrganization")
	moduleName := *viperString("name")
	providerName := *viperString("provider")
	moduleVersion := *viperString("moduleVersion")
	client, ctx := getClientContext()

	// Read Config Version
	fmt.Print("Deleting Module Version for ", color.GreenString(moduleName), "/", color.GreenString(providerName),
		":", color.GreenString(moduleVersion), "...")
	err := client.RegistryModules.DeleteVersion(ctx, orgName, moduleName, providerName, moduleVersion)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(" Deleted")

	return nil
}
