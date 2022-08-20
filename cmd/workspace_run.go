// Copyright © 2021 Tom Straub <github.com/straubt1>

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// `tfx workspace run` commands
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Workspace Runs",
		Long:  "Work with Runs of a TFx Workspace.",
	}

	// `tfx workspace list` command
	runListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Runs",
		Long:  "List Runs of a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(
				getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("workspaceName"))
		},
	}

	// `tfx workspace create` command
	runCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create Run",
		Long:  "Create Run for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(
				getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("workspaceName"),
				*viperString("message"),
				*viperString("cvId"))
		},
	}

	// `tfx workspace show` command
	runShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show Run",
		Long:  "Show Run details for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runShow(
				getTfxClientContext(),
				*viperString("runId"))
		},
	}
)

func init() {
	// `tfx workspace run` commands

	// `tfx workspace run list` command
	runListCmd.Flags().StringP("workspaceName", "w", "", "Workspace name")
	runListCmd.MarkFlagRequired("workspaceName")

	// `tfx workspace run create` command
	runCreateCmd.Flags().StringP("workspaceName", "w", "", "Workspace name")
	// runCreateCmd.Flags().StringP("directory", "d", "./", "Directory of Terraform (defaults to current directory)")
	runCreateCmd.Flags().StringP("message", "m", "", "Run Message (optional)")
	runCreateCmd.Flags().StringP("cvId", "i", "", "Configuration Version (optional)")
	runCreateCmd.MarkFlagRequired("workspaceName")

	// `tfx workspace run show` command
	runShowCmd.Flags().StringP("runId", "i", "", "Run Id (i.e. run-*)")
	runShowCmd.MarkFlagRequired("runId")

	workspaceCmd.AddCommand(runCmd)
	runCmd.AddCommand(runListCmd)
	runCmd.AddCommand(runCreateCmd)
	runCmd.AddCommand(runShowCmd)
}

func workspaceRunListAll(c TfxClientContext, workspaceId string) ([]*tfe.Run, error) {
	allItems := []*tfe.Run{}
	opts := tfe.RunListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 1, PageSize: 100},
		// Include all the things - https://www.terraform.io/cloud-docs/api-docs/run#run-operations
		Operation: "plan_only,plan_and_apply,refresh_only,destroy,empty_apply",
		Include:   []tfe.RunIncludeOpt{},
	}
	for {
		items, err := c.Client.Runs.List(c.Context, workspaceId, &opts)
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

func runList(c TfxClientContext, orgName string, workspaceName string) error {
	o.AddMessageUserProvided("List Runs for Workspace:", workspaceName)
	workspaceId, err := getWorkspaceId(c, orgName, workspaceName)
	if err != nil {
		return errors.Wrap(err, "unable to read workspace id")
	}

	items, err := workspaceRunListAll(c, workspaceId)
	if err != nil {
		return errors.Wrap(err, "failed to list variables")
	}

	o.AddTableHeader("Id", "Configuration Version", "Status", "Plan Only", "Terraform Version", "Created", "Message")
	for _, i := range items {
		o.AddTableRows(i.ID, i.ConfigurationVersion.ID, i.Status, i.PlanOnly, i.TerraformVersion, FormatDateTime(i.CreatedAt), i.Message)
	}

	return nil
}

func runCreate(c TfxClientContext, orgName string, workspaceName string, message string, cvId string) error {
	o.AddMessageUserProvided("Create Run for Workspace:", workspaceName)
	var cv *tfe.ConfigurationVersion
	w, err := c.Client.Workspaces.Read(c.Context, orgName, workspaceName)
	if err != nil {
		return errors.Wrap(err, "failed to read workspace")
	}

	if cvId != "" {
		o.AddMessageUserProvided("Configuration Version Provided:", cvId)
		cv, err = c.Client.ConfigurationVersions.Read(c.Context, cvId)
		if err != nil {
			return errors.Wrap(err, "failed to read provider runId")
		}
	} else {
		o.AddMessageUserProvided("The run will be created using the workspace's latest configuration version", "")
	}

	run, err := c.Client.Runs.Create(c.Context, tfe.RunCreateOptions{
		Workspace:            w,
		IsDestroy:            tfe.Bool(false),
		Message:              tfe.String(message),
		ConfigurationVersion: cv, // will be nil if not provided
	})
	if err != nil {
		return errors.Wrap(err, "failed to create a run")
	}

	o.AddMessageUserProvided("Run Created", "")
	o.AddDeferredMessageRead("ID", run.ID)
	o.AddDeferredMessageRead("Configuration Version", run.ConfigurationVersion.ID)
	o.AddDeferredMessageRead("Terraform Version", run.TerraformVersion)
	o.AddDeferredMessageRead("Link",
		fmt.Sprintf("https://%s/app/%s/workspaces/%s/runs/%s", c.Hostname, c.OrganizationName, workspaceName, run.ID))
	o.Close()

	return nil
}

func runShow(c TfxClientContext, runId string) error {
	o.AddMessageUserProvided("Show Run for Workspace:", runId)
	run, err := c.Client.Runs.ReadWithOptions(c.Context, runId, &tfe.RunReadOptions{
		Include: []tfe.RunIncludeOpt{},
	})
	if err != nil {
		return errors.Wrap(err, "failed to read run from id")
	}

	o.AddDeferredMessageRead("ID", run.ID)
	o.AddDeferredMessageRead("Configuration Version", run.ConfigurationVersion.ID)
	o.AddDeferredMessageRead("Status", run.Status)
	o.AddDeferredMessageRead("Message", run.Message)
	o.AddDeferredMessageRead("Terraform Version", run.TerraformVersion)
	o.AddDeferredMessageRead("Created", FormatDateTime(run.CreatedAt))
	o.Close()

	return nil
}

func getWorkspaceId(c TfxClientContext, orgName string, workspaceName string) (string, error) {
	w, err := c.Client.Workspaces.Read(c.Context, orgName, workspaceName)
	if err != nil {
		return "", err
	}

	return w.ID, nil
}
