/*
Copyright © 2021 Tom Straub <github.com/straubt1>

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
	"errors"
	"fmt"

	"github.com/fatih/color"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

var (
	applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "Applies a Workspace Run.",
		Long:  `Creates a Run Apply based on an existing Run Plan and displays its Apply logs.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runApply()
		},
		PreRun: bindPFlags,
	}
)

func init() {
	// All `tfx apply` commands
	applyCmd.PersistentFlags().StringP("runId", "i", "", "Run Id to apply")

	applyCmd.MarkPersistentFlagRequired("runId")

	rootCmd.AddCommand(applyCmd)
}

func runApply() error {
	// var err error

	// Validate flags
	hostname := *viperString("tfeHostname")
	orgName := *viperString("tfeOrganization")
	wsName := *viperString("workspaceName")
	runId := *viperString("runId")

	client, ctx := getClientContext()

	// Verify run can be applied
	var r *tfe.Run
	r, err := client.Runs.Read(ctx, runId)
	if err != nil {
		logError(err, "failed to read run")
	}

	if !runCanBeApplied(string(r.Status)) {
		logError(errors.New("run id "+r.ID+" can not be applied. status: "+string(r.Status)),
			"Unable to apply run")
	}

	// Create Apply
	err = client.Runs.Apply(ctx, r.ID, tfe.RunApplyOptions{
		Comment: tfe.String("TFx did the apply"),
	})
	if err != nil {
		logError(err, "failed to create apply")
	}

	fmt.Println("Workspace Apply Created, Apply Id:", color.BlueString(r.Apply.ID))
	fmt.Println("Navigate:", "https://"+hostname+"/app/"+orgName+"/workspaces/"+wsName+"/runs/"+r.ID)
	fmt.Println()

	getApplyLogs(ctx, client, r.Apply.ID)
	if err != nil {
		logError(err, "failed to read apply logs")
	}

	fmt.Println("Apply Complete:", r.Apply.ID)
	return nil
}
