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

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")

		tfeHostName := "https://firefly.tfe.rocks"
		tfeToken := "8WWd0B9Dqma2KQ.atlasv1.tIyHcx12dm4bWdLjlHbyuXWjNAopC4pquKxKzfm01Ez9SvL7JM1zAtTkPs3wQ98FIVg"
		// tfeOrganization := "firefly"
		tfeWorkspaceId := "ws-oNdGGThW8gx7Lyg4" // test-tom
		terraformPath := "/Users/tstraub/Projects/straubt1.github.com/crispy-telegram/terraform"

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

		w, err := client.ConfigurationVersions.Create(ctx, tfeWorkspaceId, tfe.ConfigurationVersionCreateOptions{
			AutoQueueRuns: tfe.Bool(false),
			Speculative:   tfe.Bool(true),
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(w.ID, w.Status, w.UploadURL, w.AutoQueueRuns)
		upload := w.UploadURL

		err = client.ConfigurationVersions.Upload(ctx, upload, terraformPath)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(w)

	},
}

func init() {
	cvCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
