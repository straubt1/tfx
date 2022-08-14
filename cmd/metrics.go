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
	"os"
	"time"

	"github.com/araddon/dateparse"
	"github.com/fatih/color"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// metricsCmd represents the metrics command
var (
	metricsCmd = &cobra.Command{
		Use:   "metrics",
		Short: "Read metrics about TFx Usage",
		Long:  "Read details about how TFx is being used. This command can take a while to execute.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return metricsRun()
		},
		PreRun: bindPFlags,
		Hidden: true, // hide until this is better defined
	}

	metricsWorkspaceCmd = &cobra.Command{
		Use:   "workspace",
		Short: "Read metrics about TFx Workspace Usage",
		Long:  "Read details about how TFx Workspaces are being used. This command can take a while to execute.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return metricsWorkspaceRun()
		},
		PreRun: bindPFlags,
	}
)

func init() {
	// `tfx metrics workspace`
	metricsWorkspaceCmd.Flags().StringP("since", "s", "", "Start time when querying runs in the format MM/DD/YYYY hh:mm:ss. Examples: ['01/31/2021 10:30', '02/28/2021 10:30 AM', '03/20/2021'] (optional).")

	adminCmd.AddCommand(metricsCmd)
	metricsCmd.AddCommand(metricsWorkspaceCmd)
}

func metricsRun() error {
	fmt.Println("Getting metrics for all organizations and workspaces...", color.HiYellowString("this can take some time to complete..."))
	results, err := getAllMetrics()
	if err != nil {
		logError(err, "failed to get all metrics")
	}

	err = printMetricsTable(results)
	return err
}

func metricsWorkspaceRun() error {
	orgName := *viperString("tfeOrganization")
	sinceString := *viperString("since")
	since, err := parseTime(sinceString)
	if err != nil {
		logError(err, "failed to parse given since string")
	}

	fmt.Println("Getting metrics for all workspaces in", color.BlueString(orgName), "...", color.HiYellowString("this can take some time to complete..."))
	results, err := getAllMetricsWorkspace(orgName, since)
	if err != nil {
		logError(err, "failed to get all metrics for workspaces")
	}

	printMetricsWorkspaceTable(results)

	return nil
}

type MetricsAll struct {
	OrganizationCount int      `json:"OrganizationCount"`
	Organizations     []string `json:"Organizations"`
	WorkspaceCount    int      `json:"WorkspaceCount"`
	Workspaces        []string `json:"Workspaces"`
	RunCount          int      `json:"RunCount"`
	RunErroredCount   int      `json:"RunErroredCount"`
	RunDiscardedCount int      `json:"RunDiscardedCount"`
	RunCancelledCount int      `json:"RunCancelledCount"`
	PolicyCheckCount  int      `json:"PolicyCheckCount"`
	PoliciesPassCount int      `json:"PoliciesPassCount"`
	PoliciesFailCount int      `json:"PoliciesFailCount"`
	QueryTime         string   `json:"QueryTime"`
}

type MetricsWorkspaces struct {
	WorkspaceCount    int      `json:"WorkspaceCount"`
	Workspaces        []string `json:"Workspaces"`
	RunCount          int      `json:"RunCount"`
	RunErroredCount   int      `json:"RunErroredCount"`
	RunDiscardedCount int      `json:"RunDiscardedCount"`
	RunCancelledCount int      `json:"RunCancelledCount"`
	PolicyCheckCount  int      `json:"PolicyCheckCount"`
	PoliciesPassCount int      `json:"PoliciesPassCount"`
	PoliciesFailCount int      `json:"PoliciesFailCount"`
}

type MetricsWorkspaceResult struct {
	Workspaces []MetricsWorkspace
	Since      time.Time
	QueryTime  string `json:"QueryTime"`
}

type MetricsWorkspace struct {
	Name              string `json:"Name"`
	ID                string `json:"ID"`
	RunCount          int    `json:"RunCount"`
	RunErroredCount   int    `json:"RunErroredCount"`
	RunDiscardedCount int    `json:"RunDiscardedCount"`
	RunCancelledCount int    `json:"RunCancelledCount"`
	PolicyCheckCount  int    `json:"PolicyCheckCount"`
	PoliciesPassCount int    `json:"PoliciesPassCount"`
	PoliciesFailCount int    `json:"PoliciesFailCount"`
}

func getAllMetrics() (*MetricsAll, error) {
	result := &MetricsAll{}
	start := time.Now()
	client, ctx := getClientContext()

	orgs, err := organizationListAll(getTfxClientContext())
	if err != nil {
		return nil, err
	}

	// Organizations
	result.OrganizationCount = len(orgs)
	result.Organizations = make([]string, len(orgs))
	for i, v := range orgs {
		result.Organizations[i] = v.Name
	}

	// Loop on organizations
	for _, o := range orgs {
		ws, err := getAllOrganizationWorkspaces(ctx, client, o.Name)
		if err != nil {
			return nil, err
		}
		result.WorkspaceCount += ws.WorkspaceCount
		result.RunCount += ws.RunCount
		result.RunErroredCount += ws.RunErroredCount
		result.RunCancelledCount += ws.RunCancelledCount
		result.RunDiscardedCount += ws.RunDiscardedCount
		result.PolicyCheckCount += ws.PolicyCheckCount
		result.PoliciesPassCount += ws.PoliciesPassCount
		result.PoliciesFailCount += ws.PoliciesFailCount
	}

	elapsed := time.Since(start)
	result.QueryTime = elapsed.String()

	return result, nil
}

func getAllOrganizationWorkspaces(ctx context.Context, client *tfe.Client, orgName string) (*MetricsWorkspaces, error) {
	result := &MetricsWorkspaces{}
	workspaces, err := workspaceListAllForOrganization(getTfxClientContext(), orgName, "")
	if err != nil {
		return nil, err
	}

	result.WorkspaceCount = len(workspaces)
	result.Workspaces = make([]string, len(workspaces))
	for i, ws := range workspaces {
		result.Workspaces[i] = ws.Name

		runs, err := client.Runs.List(ctx, ws.ID, &tfe.RunListOptions{
			ListOptions: tfe.ListOptions{
				PageSize: 100,
			},
			Include: []tfe.RunIncludeOpt{},
		})
		if err != nil {
			return nil, err
		}
		result.RunCount += len(runs.Items)
		for _, r := range runs.Items {
			if r.Status == "errored" {
				result.RunErroredCount++
			} else if r.Status == "canceled" || r.Status == "force_canceled" {
				result.RunCancelledCount++
			} else if r.Status == "discarded" {
				result.RunDiscardedCount++
			}

			result.PolicyCheckCount += len(r.PolicyChecks)
			for _, p := range r.PolicyChecks {
				pFull, _ := client.PolicyChecks.Read(ctx, p.ID)
				if pFull != nil {
					if pFull.Result != nil {
						result.PoliciesPassCount += pFull.Result.Passed
						result.PoliciesFailCount += pFull.Result.TotalFailed
					}
				}
			}
		}
	}

	return result, nil
}

func getAllMetricsWorkspace(orgName string, since time.Time) (*MetricsWorkspaceResult, error) {
	result := &MetricsWorkspaceResult{}
	start := time.Now()
	client, ctx := getClientContext()

	result.Since = since

	wsResults, err := getOrganizationMetricWorkspaces(ctx, client, orgName, since)
	if err != nil {
		return nil, err
	}
	result.Workspaces = *wsResults

	elapsed := time.Since(start)
	result.QueryTime = elapsed.String()
	return result, nil
}

func getOrganizationMetricWorkspaces(ctx context.Context, client *tfe.Client, orgName string, runSinceTime time.Time) (*[]MetricsWorkspace, error) {
	result := []MetricsWorkspace{}
	workspaces, err := workspaceListAllForOrganization(getTfxClientContext(), orgName, "")
	if err != nil {
		return nil, err
	}

	for _, ws := range workspaces {
		wsResult := MetricsWorkspace{}
		wsResult.Name = ws.Name
		wsResult.ID = ws.ID

		// TODO: pagination
		runs, err := client.Runs.List(ctx, ws.ID, &tfe.RunListOptions{
			ListOptions: tfe.ListOptions{
				PageSize: 100},
			Include: []tfe.RunIncludeOpt{},
		})
		if err != nil {
			return nil, err
		}
		// wsResult.RunCount += len(runs.Items)
		for _, r := range runs.Items {
			if runSinceTime.After(r.CreatedAt) {
				continue // Run outside time frame, ignore
			}
			wsResult.RunCount++

			if r.Status == "errored" {
				wsResult.RunErroredCount++
			} else if r.Status == "canceled" || r.Status == "force_canceled" {
				wsResult.RunCancelledCount++
			} else if r.Status == "discarded" {
				wsResult.RunDiscardedCount++
			}

			wsResult.PolicyCheckCount += len(r.PolicyChecks)
			for _, p := range r.PolicyChecks {
				pFull, _ := client.PolicyChecks.Read(ctx, p.ID)
				if pFull != nil {
					if pFull.Result != nil {
						wsResult.PoliciesPassCount += pFull.Result.Passed
						wsResult.PoliciesFailCount += pFull.Result.TotalFailed
					}
				}
			}
		}

		result = append(result, wsResult)
	}

	return &result, nil
}

func printMetricsTable(m *MetricsAll) error {
	fmt.Println(color.BlueString("Organization Count:       "), m.OrganizationCount)
	fmt.Println(color.BlueString("Workspace Count:          "), m.WorkspaceCount)
	fmt.Println(color.BlueString("Run Count:                "), m.RunCount)
	fmt.Println(color.BlueString("Run Errored Count:        "), m.RunErroredCount)
	fmt.Println(color.BlueString("Run Cancelled Count:      "), m.RunCancelledCount)
	fmt.Println(color.BlueString("Run Discarded Count:      "), m.RunDiscardedCount)
	fmt.Println(color.BlueString("Policy Check Count:       "), m.PolicyCheckCount)
	fmt.Println(color.BlueString("Policies Passed Count:    "), m.PoliciesPassCount)
	fmt.Println(color.BlueString("Policies Failed Count:    "), m.PoliciesFailCount)
	fmt.Println()
	fmt.Println("Metrics Query Time:", color.YellowString(m.QueryTime))
	// fmt.Println(m)
	return nil
}

func printMetricsWorkspaceTable(result *MetricsWorkspaceResult) {
	fmt.Println("Workspace Results since:", color.BlueString(timestamp(result.Since)))
	fmt.Println("Workspaces with no Runs found will be ommitted")
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Total Runs", "Errored Runs", "Discarded Runs", "Cancelled Runs"})
	for _, i := range result.Workspaces {
		if i.RunCount == 0 {
			continue //dont print if there are no runs
		}
		t.AppendRow(table.Row{i.Name, i.RunCount, i.RunErroredCount, i.RunDiscardedCount, i.RunCancelledCount})
	}
	t.SetStyle(table.StyleRounded)
	t.Render()

	fmt.Println()
	fmt.Println("Metrics Query Time:", color.YellowString(result.QueryTime))
}

func parseTime(s string) (time.Time, error) {
	zeroTime := time.Time{}
	// nothing passed, include all time
	if s == "" {
		return zeroTime, nil
	}
	// "3/1/2014 10:25 PM"
	t, err := dateparse.ParseLocal(s)
	if err != nil {
		return zeroTime, err
	}

	return t, nil
}
