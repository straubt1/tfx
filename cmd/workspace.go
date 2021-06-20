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
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/hashicorp/go-tfe"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// workspaceCmd represents the workspace command
var (
	workspaceCmd = &cobra.Command{
		Use:     "workspace",
		Aliases: []string{"ws"},
		Short:   "Workspaces",
		Long:    "Work with TFx Workspaces",
	}

	workspaceListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Workspaces",
		Long:  "List Workspaces of a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return workspaceList()
		},
		PreRun: bindPFlags,
	}

	workspaceListAllCmd = &cobra.Command{
		Use:   "all",
		Short: "List All Workspaces",
		Long:  "List Workspaces of all TFx Organizations.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return workspaceListAll()
		},
		PreRun: bindPFlags,
	}

	workspaceShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show Workspace",
		Long:  "Show Workspace of a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return workspaceShow()
		},
		PreRun: bindPFlags,
	}
)

func init() {
	// `tfx workspace list`
	workspaceListCmd.Flags().StringP("search", "s", "", "Search string for Workspace Name (optional).")
	workspaceListCmd.Flags().String("run-status", "", "Filter on current run status (optional).")

	// `tfx workspace list all`
	workspaceListAllCmd.Flags().StringP("search", "s", "", "Search string for Workspace Name (optional).")

	// `tfx workspace show`
	workspaceShowCmd.Flags().StringP("workspaceName", "w", "", "Workspace name")
	workspaceShowCmd.MarkFlagRequired("workspaceName")

	rootCmd.AddCommand(workspaceCmd)
	workspaceCmd.AddCommand(workspaceListCmd)
	workspaceListCmd.AddCommand(workspaceListAllCmd)
	workspaceCmd.AddCommand(workspaceShowCmd)
}

func workspaceList() error {
	orgName := *viperString("tfeOrganization")
	searchString := *viperString("search")
	runStatus := *viperString("run-status")
	if !validateRunStatus(runStatus) {
		logError(errors.New("run status given is now allowed"), "failed to supply a valid run status")
	}

	client, ctx := getClientContext()

	if searchString == "" {
		fmt.Println("Reading Workspaces for Organization:", color.GreenString(orgName))
	} else {
		fmt.Println("Reading Workspaces for Organization:", color.GreenString(orgName), "with workspace search string:", color.GreenString(searchString))
	}
	workspaceList, err := getAllWorkspaces(ctx, client, orgName, searchString)
	if err != nil {
		logError(err, "failed to list workspaces")
	}

	if runStatus != "" {
		workspaceList, err = filterWorkspaces(workspaceList, runStatus)
		if err != nil {
			logError(err, "failed to filter workspaces by run status")
		}
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Id", "Current Run Created", "Status"})
	for _, i := range workspaceList {
		cr_created_at := ""
		cr_status := ""
		if i.CurrentRun != nil {
			cr_created_at = timestamp(i.CurrentRun.CreatedAt)
			cr_status = string(i.CurrentRun.Status)
		}
		t.AppendRow(table.Row{i.Name, i.ID, cr_created_at, cr_status})
	}
	t.SetStyle(table.StyleRounded)
	t.Render()

	if searchString == "" {
		fmt.Println("Workspaces Found:", color.BlueString(strconv.Itoa(len(workspaceList))))
	} else {
		fmt.Println("Workspaces Found:", color.BlueString(strconv.Itoa(len(workspaceList))), "with workspace search string:", color.GreenString(searchString))
	}

	return nil
}

func workspaceListAll() error {
	searchString := *viperString("search")
	client, ctx := getClientContext()

	//
	orgs, err := getAllOrganizations(ctx, client)
	if err != nil {
		logError(err, "failed to list organizations")
	}

	aString := make([]string, len(orgs))
	for i, v := range orgs {
		aString[i] = v.Name
	}
	if searchString == "" {
		fmt.Println("Reading Workspaces for Organizations:", color.BlueString(strings.Join(aString, ", ")))
	} else {
		fmt.Println("Reading Workspaces for Organizations:", color.BlueString(strings.Join(aString, ", ")), "with workspace search string:", color.GreenString(searchString))
	}

	var allWorkspaceList []*tfe.Workspace
	for _, v := range orgs {
		workspaceList, err := getAllWorkspaces(ctx, client, v.Name, searchString)
		if err != nil {
			logError(err, "failed to list workspaces")
		}
		allWorkspaceList = append(allWorkspaceList, workspaceList...)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Organization", "Name", "Id", "Current Run Created", "Status"})
	for _, i := range allWorkspaceList {
		cr_created_at := ""
		cr_status := ""
		if i.CurrentRun != nil {
			cr_created_at = timestamp(i.CurrentRun.CreatedAt)
			cr_status = string(i.CurrentRun.Status)
		}
		t.AppendRow(table.Row{i.Organization.Name, i.Name, i.ID, cr_created_at, cr_status})
	}
	t.SetStyle(table.StyleRounded)
	t.Render()

	if searchString == "" {
		fmt.Println("Workspaces Found:", color.BlueString(strconv.Itoa(len(allWorkspaceList))))
	} else {
		fmt.Println("Workspaces Found:", color.BlueString(strconv.Itoa(len(allWorkspaceList))), "with workspace search string:", color.GreenString(searchString))
	}

	return nil
}

func getAllOrganizations(ctx context.Context, client *tfe.Client) ([]*tfe.Organization, error) {
	var err error
	var ol *tfe.OrganizationList
	var organizationItems []*tfe.Organization
	pageNumber := 1
	for {
		ol, err = client.Organizations.List(ctx, tfe.OrganizationListOptions{
			ListOptions: tfe.ListOptions{
				PageSize: 100,
			}})
		if err != nil {
			return nil, err
		}

		organizationItems = append(organizationItems, ol.Items...)
		if ol.NextPage == 0 {
			break
		}
		pageNumber++
	}

	return organizationItems, nil
}

func getAllWorkspaces(ctx context.Context, client *tfe.Client, orgName string, search string) ([]*tfe.Workspace, error) {
	var err error
	var wsl *tfe.WorkspaceList
	var workspaceItems []*tfe.Workspace
	pageNumber := 1
	for {
		wsl, err = client.Workspaces.List(ctx, orgName, tfe.WorkspaceListOptions{
			ListOptions: tfe.ListOptions{
				PageSize:   100,
				PageNumber: pageNumber,
			},
			Include: tfe.String("organization,current_run"),
			Search:  tfe.String(search),
			// A search string (partial workspace name) used to filter the results.
			// Search *string `url:"search[name],omitempty"`

			// A list of relations to include. See available resources https://www.terraform.io/docs/cloud/api/workspaces.html#available-related-resources
			// Include *string `url:"include"`
		})
		if err != nil {
			return nil, err
		}

		workspaceItems = append(workspaceItems, wsl.Items...)
		if wsl.NextPage == 0 {
			break
		}
		pageNumber++
	}

	return workspaceItems, nil
}

func filterWorkspaces(list []*tfe.Workspace, runStatus string) ([]*tfe.Workspace, error) {
	var result []*tfe.Workspace
	for _, w := range list {
		if w.CurrentRun != nil {
			if w.CurrentRun.Status == tfe.RunStatus(runStatus) {
				result = append(result, w)
			}
		}
	}

	return result, nil
}

func workspaceShow() error {
	orgName := *viperString("tfeOrganization")
	wsName := *viperString("workspaceName")
	client, ctx := getClientContext()

	fmt.Print("Reading Workspace ", color.GreenString(wsName), "...")
	w, err := client.Workspaces.Read(ctx, orgName, wsName)
	if err != nil {
		logError(err, "failed to read workspace id")
	}
	fmt.Println(" Found")

	fmt.Println(color.BlueString("Id:                 "), w.ID)
	fmt.Println(color.BlueString("Terraform Version:  "), w.TerraformVersion)
	fmt.Println(color.BlueString("Execution Mode:     "), w.ExecutionMode)
	// TODO: not populating
	// fmt.Println(color.BlueString("Last Updated:       "), timestamp(w.UpdatedAt))
	fmt.Println(color.BlueString("Auto Apply:         "), w.AutoApply)
	fmt.Println(color.BlueString("Working Directory:  "), w.WorkingDirectory)
	fmt.Println(color.BlueString("Locked:             "), w.Locked)
	if w.CurrentRun == nil {
		fmt.Println(color.BlueString("Current Run:         "), "")
	} else {
		run, err := client.Runs.ReadWithOptions(ctx, w.CurrentRun.ID, &tfe.RunReadOptions{
			Include: "",
		})
		if err != nil {
			logError(err, "failed to read workspace current run")
		}

		fmt.Println(color.BlueString("Current Run"))
		fmt.Println(color.BlueString("  Id:         "), run.ID)
		fmt.Println(color.BlueString("  Created At: "), timestamp(run.CreatedAt))
		fmt.Println(color.BlueString("  Status:     "), run.Status)
	}

	return nil
}

func validateRunStatus(s string) bool {
	if s == "" {
		return true
	}
	var runStatuses = [...]string{"pending", "plan_queued", "planning", "planned", "cost_estimating", "cost_estimated", "policy_checking", "policy_override", "policy_soft_failed", "policy_checked", "confirmed", "planned_and_finished", "apply_queued", "applying", "applied", "discarded", "errored", "canceled", "force_canceled"}
	for _, status := range runStatuses {
		if status == s {
			return true
		}
	}
	return false
}
