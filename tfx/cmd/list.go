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
	"context"
	"fmt"
	"log"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
		tfeHostName := "https://firefly.tfe.rocks"
		tfeToken := "8WWd0B9Dqma2KQ.atlasv1.tIyHcx12dm4bWdLjlHbyuXWjNAopC4pquKxKzfm01Ez9SvL7JM1zAtTkPs3wQ98FIVg"
		// tfeOrganization := "firefly"
		tfeWorkspaceId := "ws-SNfnd5PLjEPPH2H6" // test-tom

		config := &tfe.Config{
			Address: tfeHostName,
			Token:   tfeToken,
		}

		client, err := tfe.NewClient(config)
		if err != nil {
			log.Fatal(err)
		}

		// Create a context
		ctx := context.Background()

		// Get all config versions and show the current config
		w, err := client.ConfigurationVersions.List(ctx, tfeWorkspaceId, tfe.ConfigurationVersionListOptions{})
		if err != nil {
			log.Fatal(err)
		}
		for _, i := range w.Items {
			fmt.Println(i.ID, i.Speculative, i.Status)
		}

		// // Create a new workspace
		// w, err := client.Workspaces.Create(ctx, "firefly", tfe.WorkspaceCreateOptions{
		// 	Name: tfe.String("tfx-test"),
		// })
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// // Update the workspace
		// w, err = client.Workspaces.Update(ctx, "firefly", w.Name, tfe.WorkspaceUpdateOptions{
		// 	AutoApply:        tfe.Bool(false),
		// 	TerraformVersion: tfe.String("0.11.1"),
		// 	WorkingDirectory: tfe.String("my-app/infra"),
		// })
		// if err != nil {
		// 	log.Fatal(err)
		// }

	},
}

func init() {
	cvCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
