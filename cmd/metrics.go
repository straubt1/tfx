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
	"time"

	"github.com/fatih/color"
	tfe "github.com/hashicorp/go-tfe"
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
	}
)

func init() {
	rootCmd.AddCommand(metricsCmd)
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

type MetricsAll struct {
	OrganizationCount int      `json:"OrganizationCount"`
	Organizations     []string `json:"Organizations"`
	WorkspaceCount    int      `json:"WorkspaceCount"`
	Workspaces        []string `json:"Workspaces"`
	RunCount          int      `json:"RunCount"`
	PolicyCheckCount  int      `json:"PolicyCheckCount"`
	PoliciesPassCount int      `json:"PoliciesPassCount"`
	PoliciesFailCount int      `json:"PoliciesFailCount"`
	QueryTime         string   `json:"QueryTime"`
}

type MetricsWorkspaces struct {
	WorkspaceCount    int      `json:"WorkspaceCount"`
	Workspaces        []string `json:"Workspaces"`
	RunCount          int      `json:"RunCount"`
	PolicyCheckCount  int      `json:"PolicyCheckCount"`
	PoliciesPassCount int      `json:"PoliciesPassCount"`
	PoliciesFailCount int      `json:"PoliciesFailCount"`
}

func getAllMetrics() (*MetricsAll, error) {
	result := &MetricsAll{}
	start := time.Now()
	client, ctx := getClientContext()

	orgs, err := getAllOrganizations(ctx, client)
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
		result.PolicyCheckCount += ws.PolicyCheckCount
		result.PoliciesPassCount += ws.PoliciesPassCount
		result.PoliciesFailCount += ws.PoliciesFailCount
	}

	elapsed := time.Since(start)
	result.QueryTime = elapsed.String()

	return result, nil
}

// func getAllAdminOrganizations(ctx context.Context, client *tfe.Client) ([]*tfe.AdminOrganization, error) {
// 	var err error
// 	var ol *tfe.AdminOrganizationList
// 	var organizationItems []*tfe.AdminOrganization
// 	pageNumber := 1
// 	for {
// 		ol, err = client.Admin.Organizations.List(ctx, tfe.AdminOrganizationListOptions{
// 			ListOptions: tfe.ListOptions{
// 				PageSize: 100,
// 			}})
// 		if err != nil {
// 			return nil, err
// 		}

// 		organizationItems = append(organizationItems, ol.Items...)
// 		if ol.NextPage == 0 {
// 			break
// 		}
// 		pageNumber++
// 	}

// 	return organizationItems, nil
// }

func getAllOrganizationWorkspaces(ctx context.Context, client *tfe.Client, orgName string) (*MetricsWorkspaces, error) {
	result := &MetricsWorkspaces{}
	workspaces, err := getAllWorkspaces(ctx, client, orgName, "")
	if err != nil {
		return nil, err
	}

	result.WorkspaceCount = len(workspaces)
	result.Workspaces = make([]string, len(workspaces))
	for i, ws := range workspaces {
		result.Workspaces[i] = ws.Name

		runs, err := client.Runs.List(ctx, ws.ID, tfe.RunListOptions{
			ListOptions: tfe.ListOptions{
				PageSize: 100,
			},
			Include: tfe.String(""),
		})
		if err != nil {
			return nil, err
		}
		result.RunCount += len(runs.Items)
		for _, r := range runs.Items {
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

func printMetricsTable(m *MetricsAll) error {
	fmt.Println(color.BlueString("Organization Count:       "), m.OrganizationCount)
	fmt.Println(color.BlueString("Workspace Count:          "), m.WorkspaceCount)
	fmt.Println(color.BlueString("Run Count:                "), m.RunCount)
	fmt.Println(color.BlueString("Policy Check Count:       "), m.PolicyCheckCount)
	fmt.Println(color.BlueString("Policies Passed Count:    "), m.PoliciesPassCount)
	fmt.Println(color.BlueString("Policies Failed Count:    "), m.PoliciesFailCount)
	fmt.Println()
	fmt.Println("Metrics Query Time:", color.YellowString(m.QueryTime))
	// fmt.Println(m)
	return nil
}
