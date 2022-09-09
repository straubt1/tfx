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
	"math"
	"strings"

	"github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// workspaceCmd represents the workspace command
var (
	// `tfx workspace` commands
	workspaceCmd = &cobra.Command{
		Use:     "workspace",
		Aliases: []string{"ws"},
		Short:   "Workspace Commands",
		Long:    "Work with TFx Workspaces",
	}

	// `tfx workspace list` command
	workspaceListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Workspaces",
		Long:  "List Workspaces in a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !validateRunStatus(*viperString("run-status")) {
				return errors.New("run status given is now allowed")
			}

			if *viperBool("all") {
				return workspaceListAll(
					getTfxClientContext(),
					*viperString("search"),
					*viperString("run-status"))
			} else {
				return workspaceList(
					getTfxClientContext(),
					*viperString("search"),
					*viperString("run-status"))
			}
		},
	}

	// `tfx workspace show` command
	workspaceShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show Workspace",
		Long:  "Show Workspace in a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return workspaceShow(
				getTfxClientContext(),
				*viperString("name"))
		},
	}
)

func init() {
	// `tfx workspace list`
	workspaceListCmd.Flags().StringP("search", "s", "", "Search string for Workspace Name (optional).")
	workspaceListCmd.Flags().String("run-status", "", "Filter on current run status (optional).")
	workspaceListCmd.Flags().BoolP("all", "a", false, "List All Organizations Workspaces (optional).")

	// `tfx workspace show`
	workspaceShowCmd.Flags().StringP("name", "n", "", "Name of the workspace.")
	workspaceShowCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(workspaceCmd)
	workspaceCmd.AddCommand(workspaceListCmd)
	workspaceCmd.AddCommand(workspaceShowCmd)
}

func workspaceListAllForOrganization(c TfxClientContext, orgName string, searchString string) ([]*tfe.Workspace, error) {
	allItems := []*tfe.Workspace{}
	opts := tfe.WorkspaceListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 1, PageSize: 100},
		Search:      searchString,
		// Tags:        "",
		// ExcludeTags: "",
		Include: []tfe.WSIncludeOpt{"organization", "current_run"},
	}
	for {
		items, err := c.Client.Workspaces.List(c.Context, orgName, &opts)
		if err != nil {
			return nil, err
		}

		allItems = append(allItems, items.Items...)
		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage
	}

	return allItems, nil
}

func organizationListAll(c TfxClientContext) ([]*tfe.Organization, error) {
	allItems := []*tfe.Organization{}
	opts := tfe.OrganizationListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100},
	}
	for {
		items, err := c.Client.Organizations.List(c.Context, &opts)
		if err != nil {
			return nil, err
		}

		allItems = append(allItems, items.Items...)
		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage
	}

	return allItems, nil
}

func workspaceList(c TfxClientContext, searchString string, runStatus string) error {
	o.AddMessageUserProvided("List Workspaces for Organization:", c.OrganizationName)
	items, err := workspaceListAllForOrganization(c, c.OrganizationName, searchString)
	if err != nil {
		return errors.Wrap(err, "failed to list variables")
	}

	if runStatus != "" {
		items, err = filterWorkspaces(items, runStatus)
		if err != nil {
			logError(err, "failed to filter workspaces by run status")
		}
		o.AddFormattedMessageCalculated("Found %d Filtered Workspaces", len(items))
	} else {
		o.AddFormattedMessageCalculated("Found %d Workspaces", len(items))
	}

	o.AddTableHeader("Name", "Id", "Current Run Created", "Status", "Repository", "Locked")
	for _, i := range items {
		cr_created_at := ""
		cr_status := ""
		if i.CurrentRun != nil {
			cr_created_at = FormatDateTime(i.CurrentRun.CreatedAt)
			cr_status = string(i.CurrentRun.Status)
		}
		ws_repo := ""
		if i.VCSRepo != nil {
			ws_repo = i.VCSRepo.DisplayIdentifier
		}

		o.AddTableRows(i.Name, i.ID, cr_created_at, cr_status, ws_repo, i.Locked)
	}

	return nil
}

func workspaceListAll(c TfxClientContext, searchString string, runStatus string) error {
	o.AddMessageUserProvided("List Workspaces for all available Organizations", "")
	orgs, err := organizationListAll(c)
	if err != nil {
		logError(err, "failed to list organizations")
	}

	var allWorkspaceList []*tfe.Workspace
	for _, v := range orgs {
		workspaceList, err := workspaceListAllForOrganization(c, v.Name, searchString)
		if err != nil {
			logError(err, "failed to list workspaces for organization")
		}
		allWorkspaceList = append(allWorkspaceList, workspaceList...)
	}

	if runStatus != "" {
		allWorkspaceList, err = filterWorkspaces(allWorkspaceList, runStatus)
		if err != nil {
			logError(err, "failed to filter workspaces by run status")
		}
		o.AddFormattedMessageCalculated("Found %d Filtered Workspaces", len(allWorkspaceList))
	} else {
		o.AddFormattedMessageCalculated("Found %d Workspaces", len(allWorkspaceList))
	}

	o.AddTableHeader("Organization", "Name", "Id", "Current Run Created", "Status", "Repository", "Locked")
	for _, i := range allWorkspaceList {
		cr_created_at := ""
		cr_status := ""
		if i.CurrentRun != nil {
			cr_created_at = FormatDateTime(i.CurrentRun.CreatedAt)
			cr_status = string(i.CurrentRun.Status)
		}
		ws_repo := ""
		if i.VCSRepo != nil {
			ws_repo = i.VCSRepo.DisplayIdentifier
		}

		o.AddTableRows(i.Organization.Name, i.Name, i.ID, cr_created_at, cr_status, ws_repo, i.Locked)
	}

	return nil
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

func workspaceShow(c TfxClientContext, workspaceName string) error {
	o.AddMessageUserProvided("Show Workspace:", workspaceName)
	w, err := c.Client.Workspaces.Read(c.Context, c.OrganizationName, workspaceName)
	if err != nil {
		logError(err, "failed to read workspace")
	}

	ta, err := workspaceTeamListAll(c, w.ID, math.MaxInt)
	if err != nil {
		return errors.Wrap(err, "failed to list teams")
	}

	o.AddDeferredMessageRead("ID", w.ID)
	o.AddDeferredMessageRead("Terraform Version", w.TerraformVersion)
	o.AddDeferredMessageRead("Execution Mode", w.ExecutionMode)
	o.AddDeferredMessageRead("Auto Apply", w.AutoApply)
	o.AddDeferredMessageRead("Working Directory", w.WorkingDirectory)
	o.AddDeferredMessageRead("Locked", w.Locked)

	if w.CurrentRun == nil {
		o.AddDeferredMessageRead("Current Run", "none")
	} else {
		run, err := c.Client.Runs.ReadWithOptions(c.Context, w.CurrentRun.ID, &tfe.RunReadOptions{
			Include: []tfe.RunIncludeOpt{},
		})
		if err != nil {
			logError(err, "failed to read workspace current run")
		}

		o.AddDeferredMessageRead("Current Run Id", run.ID)
		o.AddDeferredMessageRead("Current Run Status", run.Status)
		o.AddDeferredMessageRead("Current Run Created", FormatDateTime(run.CreatedAt))
	}

	// loop through team access and get team names (requires an additional API call)
	teamNames := []string{}
	for _, i := range ta {
		t, err := c.Client.Teams.Read(c.Context, i.Team.ID)
		if err != nil {
			return errors.Wrap(err, "failed to find team name")
		}
		teamNames = append(teamNames, t.Name)
	}
	o.AddDeferredMessageRead("Team Access", strings.Join(teamNames, ","))

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
