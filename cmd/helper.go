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
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	tfe "github.com/hashicorp/go-tfe"
)

// Create or read the CV to prepare for plan
func createOrReadConfigurationVersion(ctx context.Context, client *tfe.Client, workspaceId string, cvId string, tfDirectory string, speculative bool) (*tfe.ConfigurationVersion, error) {
	var err error
	var cv *tfe.ConfigurationVersion

	if cvId == "" { // config id was not give, create a new one
		fmt.Print("Creating new Config Version ...")
		cv, err = client.ConfigurationVersions.Create(ctx, workspaceId, tfe.ConfigurationVersionCreateOptions{
			AutoQueueRuns: tfe.Bool(false), // wait for upload
			Speculative:   tfe.Bool(speculative),
		})
		if err != nil {
			return nil, err
		}
		fmt.Println(" ID:", color.BlueString(cv.ID))

		// upload code
		err = client.ConfigurationVersions.Upload(ctx, cv.UploadURL, tfDirectory)
		if err != nil {
			return nil, err
		}
	} else {
		// config id was given, read
		fmt.Println("Using existing Config Version ...", cvId)
		cv, err = client.ConfigurationVersions.Read(ctx, cvId)
		if err != nil {
			return nil, err
		}
		fmt.Println(" ID:", color.BlueString(cv.ID), " Status: ", color.BlueString(string(cv.Status)))
		if cv.Status != "uploaded" {
			return nil, errors.New("provider configuration version is not allowed")
		}
	}

	return cv, nil
}

// Get Run logs
func getRunLogs(ctx context.Context, client *tfe.Client, planId string) error {
	var err error
	var logs io.Reader
	logs, err = client.Plans.Logs(ctx, planId)
	if err != nil {
		return err
	}

	fmt.Println("------------------------------------------------------------------------")
	// mostly found from here: https://github.com/hashicorp/terraform/blob/89f986ded6fb07e7d5f27aaf340f69c353860c12/backend/remote/backend_plan.go#L332
	reader := bufio.NewReaderSize(logs, 64*1024)
	for next := true; next; {
		var l, line []byte

		for isPrefix := true; isPrefix; {
			l, isPrefix, err = reader.ReadLine()
			if err != nil {
				if err != io.EOF {
					return err
				}
				next = false
			}
			line = append(line, l...)
		}

		if next || len(line) > 0 {
			fmt.Println(string(line))
		}
	}

	// fmt.Println("------------------------------------------------------------------------")
	return nil
}

// Get Cost Estimation logs, if any
func getCostEstimationLogs(ctx context.Context, client *tfe.Client, r *tfe.Run) error {
	if r.CostEstimate == nil {
		return nil
	}
	msgPrefix := "Cost estimation"
	started := time.Now()
	updated := started
	for i := 0; ; i++ {
		// select {
		// case <-stopCtx.Done():
		// 	return stopCtx.Err()
		// case <-cancelCtx.Done():
		// 	return cancelCtx.Err()
		// case <-time.After(backoff(backoffMin, backoffMax, i)):
		// }

		// Retrieve the cost estimate to get its current status.
		ce, err := client.CostEstimates.Read(ctx, r.CostEstimate.ID)
		if err != nil {
			return err
		}

		// If the run is canceled or errored, but the cost-estimate still has
		// no result, there is nothing further to render.
		if ce.Status != tfe.CostEstimateFinished {
			if r.Status == tfe.RunCanceled || r.Status == tfe.RunErrored {
				return nil
			}
		}

		// checking if i == 0 so as to avoid printing this starting horizontal-rule
		// every retry, and that it only prints it on the first (i=0) attempt.
		if i == 0 {
			fmt.Println("------------------------------------------------------------------------")
		}

		switch ce.Status {
		case tfe.CostEstimateFinished:
			delta, err := strconv.ParseFloat(ce.DeltaMonthlyCost, 64)
			if err != nil {
				return err
			}

			sign := "+"
			if delta < 0 {
				sign = "-"
			}

			deltaRepr := strings.Replace(ce.DeltaMonthlyCost, "-", "", 1)

			fmt.Println(msgPrefix + ":")
			fmt.Printf("Resources: %d of %d estimated", ce.MatchedResourcesCount, ce.ResourcesCount)
			fmt.Printf("           $%s/mo %s$%s", ce.ProposedMonthlyCost, sign, deltaRepr)

			if len(r.PolicyChecks) == 0 && r.HasChanges {
				fmt.Println("------------------------------------------------------------------------")
			}

			// if b.CLI != nil {
			// 	b.CLI.Output(b.Colorize().Color(msgPrefix + ":"))
			// 	b.CLI.Output(b.Colorize().Color(fmt.Sprintf("Resources: %d of %d estimated", ce.MatchedResourcesCount, ce.ResourcesCount)))
			// 	b.CLI.Output(b.Colorize().Color(fmt.Sprintf("           $%s/mo %s$%s", ce.ProposedMonthlyCost, sign, deltaRepr)))

			// 	if len(r.PolicyChecks) == 0 && r.HasChanges && op.Type == backend.OperationTypeApply {
			// 		b.CLI.Output("------------------------------------------------------------------------")
			// 	}
			// }

			return nil
		case tfe.CostEstimatePending, tfe.CostEstimateQueued:
			// Check if 30 seconds have passed since the last update.
			current := time.Now()
			if i == 0 || current.Sub(updated).Seconds() > 30 {
				updated = current
				elapsed := ""

				// Calculate and set the elapsed time.
				if i > 0 {
					elapsed = fmt.Sprintf(
						" (%s elapsed)", current.Sub(started).Truncate(30*time.Second))
				}
				fmt.Println(msgPrefix + ":")
				fmt.Println("Waiting for cost estimate to complete..." + elapsed)
			}
			continue
		case tfe.CostEstimateSkippedDueToTargeting:
			fmt.Println(msgPrefix + ":")
			fmt.Println("Not available for this plan, because it was created with the -target option.")
			fmt.Println("------------------------------------------------------------------------")
			return nil
		case tfe.CostEstimateErrored:
			fmt.Println(msgPrefix + " errored.")
			fmt.Println("------------------------------------------------------------------------")
			return nil
		case tfe.CostEstimateCanceled:
			return errors.New(msgPrefix + " canceled.")
		default:
			return errors.New("Unknown or unexpected cost estimate state: " + string(ce.Status))
		}
	}
}

// Get Policy logs, if any
func getPolicyLogs(ctx context.Context, client *tfe.Client, r *tfe.Run) error {
	if r.PolicyChecks == nil {
		return nil
	}

	fmt.Println("------------------------------------------------------------------------")

	for i, pc := range r.PolicyChecks {
		// Read the policy check logs. This is a blocking call that will only
		// return once the policy check is complete.
		logs, err := client.PolicyChecks.Logs(ctx, pc.ID)
		if err != nil {
			return err
		}
		reader := bufio.NewReaderSize(logs, 64*1024)

		// Retrieve the policy check to get its current status.
		pc, err := client.PolicyChecks.Read(ctx, pc.ID)
		if err != nil {
			return err
		}

		// If the run is canceled or errored, but the policy check still has
		// no result, there is nothing further to render.
		if r.Status == tfe.RunCanceled || r.Status == tfe.RunErrored {
			switch pc.Status {
			case tfe.PolicyPending, tfe.PolicyQueued, tfe.PolicyUnreachable:
				continue
			}
		}

		var msgPrefix string
		switch pc.Scope {
		case tfe.PolicyScopeOrganization:
			msgPrefix = "Organization policy check"
		case tfe.PolicyScopeWorkspace:
			msgPrefix = "Workspace policy check"
		default:
			msgPrefix = fmt.Sprintf("Unknown policy check (%s)", pc.Scope)
		}

		fmt.Println(msgPrefix + ":")

		for next := true; next; {
			var l, line []byte

			for isPrefix := true; isPrefix; {
				l, isPrefix, err = reader.ReadLine()
				if err != nil {
					if err != io.EOF {
						return err
					}
					next = false
				}
				line = append(line, l...)
			}

			if next || len(line) > 0 {
				fmt.Println(string(line))
			}
		}

		switch pc.Status {
		case tfe.PolicyPasses:
			if r.HasChanges || i < len(r.PolicyChecks)-1 {
				fmt.Println("------------------------------------------------------------------------")
			}
			continue
		case tfe.PolicyErrored:
			return fmt.Errorf(msgPrefix + " errored.")
		case tfe.PolicyHardFailed:
			return fmt.Errorf(msgPrefix + " hard failed.")
		case tfe.PolicySoftFailed:
			fmt.Println("PolicySoftFailed")
			// runUrl := fmt.Sprintf(runHeader, b.hostname, b.organization, op.Workspace, r.ID)

			// if op.Type == backend.OperationTypePlan || op.UIOut == nil || op.UIIn == nil ||
			// 	!pc.Actions.IsOverridable || !pc.Permissions.CanOverride {
			// 	return fmt.Errorf(msgPrefix + " soft failed.\n" + runUrl)
			// }

			// if op.AutoApprove {
			// 	if _, err = b.client.PolicyChecks.Override(stopCtx, pc.ID); err != nil {
			// 		return generalError(fmt.Sprintf("Failed to override policy check.\n%s", runUrl), err)
			// 	}
			// } else {
			// 	opts := &terraform.InputOpts{
			// 		Id:          "override",
			// 		Query:       "\nDo you want to override the soft failed policy check?",
			// 		Description: "Only 'override' will be accepted to override.",
			// 	}
			// 	err = b.confirm(stopCtx, op, opts, r, "override")
			// 	if err != nil && err != errRunOverridden {
			// 		return fmt.Errorf(
			// 			fmt.Sprintf("Failed to override: %s\n%s\n", err.Error(), runUrl),
			// 		)
			// 	}

			// 	if err != errRunOverridden {
			// 		if _, err = b.client.PolicyChecks.Override(stopCtx, pc.ID); err != nil {
			// 			return generalError(fmt.Sprintf("Failed to override policy check.\n%s", runUrl), err)
			// 		}
			// 	} else {
			// 		b.CLI.Output(fmt.Sprintf("The run needs to be manually overridden or discarded.\n%s\n", runUrl))
			// 	}
			// }

			fmt.Println("------------------------------------------------------------------------")
		default:
			return fmt.Errorf("Unknown or unexpected policy state: %s", pc.Status)
		}
	}

	return nil
}

// Get Apply logs
func getApplyLogs(ctx context.Context, client *tfe.Client, applyId string) error {
	var err error
	var logs io.Reader
	logs, err = client.Applies.Logs(ctx, applyId)
	if err != nil {
		return err
	}

	fmt.Println("------------------------------------------------------------------------")
	// mostly found from here: https://github.com/hashicorp/terraform/blob/89f986ded6fb07e7d5f27aaf340f69c353860c12/backend/remote/backend_plan.go#L332
	reader := bufio.NewReaderSize(logs, 64*1024)
	for next := true; next; {
		var l, line []byte

		for isPrefix := true; isPrefix; {
			l, isPrefix, err = reader.ReadLine()
			if err != nil {
				if err != io.EOF {
					return err
				}
				next = false
			}
			line = append(line, l...)
		}

		if next || len(line) > 0 {
			fmt.Println(string(line))
		}
	}

	return nil
}

// Ensure variable is up to date (upsert)
func createOrUpdateEnvVariables(ctx context.Context, client *tfe.Client, workspaceId string, variables map[string]string) error {
	var err error
	var allV *tfe.VariableList
	isSensitive := false

	// Read all variables and search
	// TODO: is there a better way? API doesnt expose a variable by name lookup
	allV, err = client.Variables.List(ctx, workspaceId, tfe.VariableListOptions{})
	if err != nil {
		return err
	}

	// loop over each env variable
	for key, val := range variables {
		// determine if variable already exists
		var found *tfe.Variable
		for i := range allV.Items {
			if allV.Items[i].Key == key {
				found = allV.Items[i]
			}
		}

		timestamp := time.Now()
		if found == nil {
			fmt.Print("Creating new Variable: ", color.GreenString(key), " ...")
			_, err = client.Variables.Create(ctx, workspaceId, tfe.VariableCreateOptions{
				Key:         &key,
				Value:       &val,
				Description: tfe.String(fmt.Sprintf("Written by TFx at %s", timestamp)),
				Category:    tfe.Category("env"),
				Sensitive:   tfe.Bool(isSensitive),
			})
			if err != nil {
				return err
			}
		} else {
			fmt.Print("Updating existing Variable: ", color.GreenString(key), " ...")
			_, err = client.Variables.Update(ctx, workspaceId, found.ID, tfe.VariableUpdateOptions{
				Key:         &key,
				Value:       &val,
				Description: tfe.String(fmt.Sprintf("Written by TFx at %s", timestamp)),
				// Category:    tfe.Category("env"),
				Sensitive: tfe.Bool(isSensitive),
			})
			if err != nil {
				return err
			}
		}
		fmt.Println(" Done")
	}

	return nil
}

// Determine if a run status can be applied
func runCanBeApplied(status string) bool {
	allowed := []string{"planned", "cost_estimated", "policy_checked"}

	for _, a := range allowed {
		if status == a {
			return true
		}
	}
	return false
}
