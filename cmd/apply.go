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
	runId string
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Applies a Workspace Run.",
	Long:  `Creates a Run Apply based on an existing Run Plan and displays its Apply logs.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("apply called")
		var err error

		client, ctx := getContext()

		// Verify run can be applied
		var r *tfe.Run
		r, err = client.Runs.Read(ctx, runId)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(r.Status)
		if !runCanBeApplied(string(r.Status)) {
			fmt.Println("Run Id ", r.ID, "Can not be applied. Status:", r.Status)
			return
		}

		// Create Apply
		err = client.Runs.Apply(ctx, r.ID, tfe.RunApplyOptions{
			Comment: tfe.String("TFx did the apply"),
		})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Workspace Apply Created, Apply Id:", r.Apply.ID)
		fmt.Println("Navigate:", "https://"+tfeHostname+"/app/"+tfeOrganization+"/workspaces/"+workspaceName+"/runs/"+r.ID)
		fmt.Println()

		getApplyLogs(ctx, client, r.Apply.ID)

		fmt.Println("Apply Complete:", r.Apply.ID)
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.PersistentFlags().StringVarP(&runId, "runId", "r", "", "Run Id to apply")
}