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
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

// var workspaceName string
var directory string
var configId string

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("plan called", directory)
		var err error

		// terraformPath := "/Users/tstraub/Projects/straubt1.github.com/crispy-telegram/terraform/"
		var terraformPath string
		terraformPath, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(terraformPath)
		terraformPath = "/Users/tstraub/tfx/terraform" //debug

		config := &tfe.Config{
			Address: "https://" + tfeHostname,
			Token:   tfeToken,
		}

		var client *tfe.Client
		client, err = tfe.NewClient(config)
		if err != nil {
			log.Fatal(err)
		}
		// Create a context
		ctx := context.Background()

		var w *tfe.Workspace
		// Read workspace
		w, err = client.Workspaces.Read(ctx, tfeOrganization, workspaceName)
		if err != nil {
			log.Fatal(err)
		}

		// cvId := ""
		// cvId := "cv-3vW3Q3QbPvfYYcp8"
		var cv *tfe.ConfigurationVersion

		if configId == "" {
			fmt.Println("Creating new Config Version")
			// create config version
			cv, err = client.ConfigurationVersions.Create(ctx, w.ID, tfe.ConfigurationVersionCreateOptions{
				AutoQueueRuns: tfe.Bool(true),
				Speculative:   tfe.Bool(true),
			})
			if err != nil {
				log.Fatal(err)
			}

			// upload code
			err = client.ConfigurationVersions.Upload(ctx, cv.UploadURL, terraformPath)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("Using existing Config Version", configId)
			// read config version
			cv, err = client.ConfigurationVersions.Read(ctx, configId)
			if err != nil {
				log.Fatal(err)
			}
		}

		// create run
		var r *tfe.Run
		r, err = client.Runs.Create(ctx, tfe.RunCreateOptions{
			IsDestroy:            tfe.Bool(false),
			Message:              tfe.String("TFx is here"),
			ConfigurationVersion: cv,
			Workspace:            w,
			// TargetAddrs: [],
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Plan URL:", "https://"+tfeHostname+"/app/"+tfeOrganization+"/workspaces/"+workspaceName+"/runs/"+r.ID)
		fmt.Println(" run ID:", r.ID)
		fmt.Println(" cv ID:", r.ConfigurationVersion.ID)

		//
		// for {
		// 	fmt.Println(r.Status)
		// 	if r.Status == "planned_and_finished" {
		// 		break
		// 	}
		// 	// get current status
		// 	r, err = client.Runs.Read(ctx, r.ID)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 		break
		// 	}
		// 	time.Sleep(1 * time.Second)
		// }
		var logs io.Reader
		logs, err = client.Plans.Logs(ctx, r.Plan.ID)
		if err != nil {
			log.Fatal(err)
		}
		reader := bufio.NewReaderSize(logs, 64*1024)
		// https://github.com/hashicorp/terraform/blob/89f986ded6fb07e7d5f27aaf340f69c353860c12/backend/remote/backend_plan.go#L332
		for next := true; next; {
			var l, line []byte

			for isPrefix := true; isPrefix; {
				l, isPrefix, err = reader.ReadLine()
				if err != nil {
					if err != io.EOF {
						log.Fatal(err)
					}
					next = false
				}
				line = append(line, l...)
			}

			if next || len(line) > 0 {
				fmt.Println(string(line))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(planCmd)

	planCmd.PersistentFlags().StringVarP(&workspaceName, "workspaceName", "w", "", "Workspace Name")
	planCmd.PersistentFlags().StringVarP(&configId, "configId", "c", "", "Configuration Version Id (optional)")
	planCmd.PersistentFlags().StringVarP(&directory, "directory", "d", "", "Directory containing Terraform")
}
