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

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

// pmrCmd represents the pmr command
var pmrCmd = &cobra.Command{
	Use:   "pmr",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		client, ctx := getContext()
		fmt.Println(tfeHostname, workspaceName)

		// modId := "mod-BAYTcvZTjMZ5dSzk"
		modName := "name"
		providerName := "aws"
		var r *tfe.RegistryModule
		// Create module
		// r, err = client.RegistryModules.Create(ctx, tfeOrganization, tfe.RegistryModuleCreateOptions{
		// 	Name:     modName,
		// 	Provider: providerName,
		// })

		// Delete module - debug
		// err = client.RegistryModules.Delete(ctx, tfeOrganization, modName)

		// Read module
		r, err = client.RegistryModules.Read(ctx, tfeOrganization, modName, providerName)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(r)

		var url *string
		url, err = RegistryModulesCreateVersion(tfeToken, tfeHostname, tfeOrganization,
			modName, providerName, "0.0.9")
		fmt.Println(url)
		// var mv *tfe.RegistryModuleVersion
		// mv, err = client.RegistryModules.CreateVersion(ctx, tfeOrganization, modName, providerName, tfe.RegistryModuleCreateVersionOptions{
		// 	Version: tfe.String("0.0.4"),
		// })
		if err != nil {
			log.Fatal(err)
		}
		err = RegistryModulesUpload(tfeToken, *url)
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Println(mv)
	},
}

func init() {
	rootCmd.AddCommand(pmrCmd)

	pmrCmd.PersistentFlags().StringVarP(&directory, "directory", "d", "./", "Directory containing Terraform, default to working directory.")
}
