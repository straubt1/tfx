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

	tfe "github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

var (
	directory     string
	configId      string
	isSpeculative bool
	// envs      []string
)

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Create a Workspace Run to Plan from your local Terraform files.",
	Long: `Creates a Configuration Version (unless supplied), uploads your Terraform files,
	creates a run, and display the Run logs.`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		client, ctx := getContext()
		fmt.Println(tfeHostname, workspaceName)

		var w *tfe.Workspace
		// Read workspace
		w, err = client.Workspaces.Read(ctx, tfeOrganization, workspaceName)
		if err != nil {
			log.Fatal(err)
		}

		// Create vars
		err = createOrUpdateVariable(ctx, client, w.ID, "test", "val")
		if err != nil {
			log.Fatal(err)
		}

		// Create config version
		cv, err := createOrReadConfigurationVersion(ctx, client, w.ID, configId, directory, isSpeculative)
		if err != nil {
			log.Fatal(err)
		}

		// Create run
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
		fmt.Println("Workspace Run Created, Run Id:", r.ID, "Config Version:", r.ConfigurationVersion.ID)
		fmt.Println("Navigate:", "https://"+tfeHostname+"/app/"+tfeOrganization+"/workspaces/"+workspaceName+"/runs/"+r.ID)
		fmt.Println()

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

		// Get Run Logs
		err = getRunLogs(ctx, client, r.Plan.ID)
		if err != nil {
			log.Fatal(err)
		}

		// Get Cost Estimation Logs
		// Get Policy Check Logs: https://github.com/hashicorp/terraform/blob/89f986ded6fb07e7d5f27aaf340f69c353860c12/backend/remote/backend_common.go#L344
		fmt.Println("Run Complete:", r.ID)
	},
}

func init() {
	rootCmd.AddCommand(planCmd)

	planCmd.PersistentFlags().StringVarP(&workspaceName, "workspaceName", "w", "", "Workspace Name")
	planCmd.PersistentFlags().StringVarP(&configId, "configId", "c", "", "Configuration Version Id (optional)")
	planCmd.PersistentFlags().StringVarP(&directory, "directory", "d", "./", "Directory containing Terraform, default to working directory.")
	planCmd.PersistentFlags().BoolVarP(&isSpeculative, "isSpeculative", "s", false, "Is this plan speculative, default to false")
	// planCmd.PersistentFlags().StringSliceVarP(&envs, "envs", "e", []string{}, "Array on ENV")
}
