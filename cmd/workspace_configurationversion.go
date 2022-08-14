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
	"bytes"
	"io/ioutil"
	"math"

	"github.com/hashicorp/go-slug"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// `tfx workspace configuration-version` commands
	cvCmd = &cobra.Command{
		Use:     "configuration-version",
		Aliases: []string{"cv"},
		Short:   "Configuration Version Commands",
		Long:    "Work with Configuration Versions of a TFx Workspace.",
	}

	// `tfx workspace configuration-version list` command
	cvListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Configuration Versions",
		Long:  "List Configuration Versions of a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			m := *viperInt("maxItems")
			if *viperBool("all") {
				m = math.MaxInt
			}
			return cvList(
				getTfxClientContext(),
				*viperString("workspaceName"),
				m)
		},
	}

	// `tfx workspace configuration-version create` command
	cvCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create Configuration Version",
		Long:  "Create Configuration Version for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !isDirectory(*viperString("directory")) {
				return errors.New("directory file does not exist")
			}

			return cvCreate(
				getTfxClientContext(),
				*viperString("workspaceName"),
				*viperString("directory"),
				*viperBool("speculative"))
		},
	}

	// `tfx workspace configuration-version show` command
	cvShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show Configuration Version",
		Long:  "Show Configuration Version details for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cvShow(
				getTfxClientContext(),
				*viperString("configurationId"))
		},
	}

	// `tfx workspace configuration-version download` command
	cvDownloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download the Configuration Version",
		Long:  "Download the Configuration Version code for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cvDownload(
				getTfxClientContext(),
				*viperString("configurationId"),
				*viperString("directory"))
		},
	}
)

func init() {
	// `tfx workspace configuration-version` commands

	// `tfx workspace configuration-version list` command
	cvListCmd.Flags().StringP("workspaceName", "w", "", "Workspace name")
	cvListCmd.Flags().IntP("maxItems", "", 10, "Max number of results (optional)")
	cvListCmd.Flags().BoolP("all", "a", false, "Retrieve all results regardless of maxItems flag (optional)")
	cvListCmd.MarkFlagRequired("workspaceName")

	// `tfx cv create`
	cvCreateCmd.Flags().StringP("workspaceName", "w", "", "Workspace name")
	cvCreateCmd.Flags().StringP("directory", "d", "./", "Directory of Terraform (optional, defaults to current directory)")
	cvCreateCmd.Flags().BoolP("speculative", "s", false, "Perform a Speculative Plan (optional, defaults to false)")
	cvCreateCmd.MarkFlagRequired("workspaceName")

	// `tfx cv show`
	cvShowCmd.Flags().StringP("configurationId", "i", "", "Configuration Version Id (i.e. cv-*)")
	cvShowCmd.MarkFlagRequired("configurationId")

	// `tfx cv download`
	cvDownloadCmd.Flags().StringP("configurationId", "i", "", "Configuration Version Id (i.e. cv-*)")
	cvDownloadCmd.Flags().StringP("directory", "d", "", "Directory to download Configuration Version to (optional, defaults to a temp directory)")
	cvDownloadCmd.MarkFlagRequired("configurationId")

	workspaceCmd.AddCommand(cvCmd)
	cvCmd.AddCommand(cvListCmd)
	cvCmd.AddCommand(cvCreateCmd)
	cvCmd.AddCommand(cvShowCmd)
	cvCmd.AddCommand(cvDownloadCmd)
}

func cvListAll(c TfxClientContext, workspaceId string, maxItems int) ([]*tfe.ConfigurationVersion, error) {
	pageSize := 100

	if maxItems < 100 {
		pageSize = maxItems // Only get what we need in one page
	}
	allItems := []*tfe.ConfigurationVersion{}
	opts := tfe.ConfigurationVersionListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 1, PageSize: pageSize},
		Include:     []tfe.ConfigVerIncludeOpt{"ingress_attributes"},
	}
	for {
		items, err := c.Client.ConfigurationVersions.List(c.Context, workspaceId, &opts)
		if err != nil {
			return nil, err
		}

		allItems = append(allItems, items.Items...)
		if len(allItems) >= maxItems {
			break // Hit the max, break. For maxItems > 100 it is possible to return more than max in this approach
		}

		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage
	}

	return allItems, nil
}

func cvList(c TfxClientContext, workspaceName string, maxItems int) error {
	o.AddMessageUserProvided("List Configuration Versions for Workspace:", workspaceName)
	workspaceId, err := getWorkspaceId(c, c.OrganizationName, workspaceName)
	if err != nil {
		return errors.Wrap(err, "unable to read workspace id")
	}

	items, err := cvListAll(c, workspaceId, maxItems)
	if err != nil {
		return errors.Wrap(err, "failed to list variables")
	}

	o.AddTableHeader("Id", "Speculative", "Status", "Repo", "Branch", "Commit", "Message")
	for _, i := range items {
		identifier, branch, commit, message := "", "", "", ""
		if i.IngressAttributes != nil { // only valid for VCS-driven workflow
			identifier = i.IngressAttributes.Identifier
			branch = i.IngressAttributes.Branch
			commit = i.IngressAttributes.CommitSHA
			if len(commit) > 7 {
				commit = commit[0:7] // shrink hash, nice to have
			}
			message = i.IngressAttributes.CommitMessage
		}
		o.AddTableRows(i.ID, i.Speculative, i.Status, identifier, branch, commit, message)
	}
	o.Close()

	return nil
}

func cvCreate(c TfxClientContext, workspaceName string, directory string, isSpeculative bool) error {
	o.AddMessageUserProvided("Create Configuration Version for Workspace:", workspaceName)
	o.AddMessageUserProvided("Code Directory:", directory)
	workspaceId, err := getWorkspaceId(c, c.OrganizationName, workspaceName)
	if err != nil {
		return errors.Wrap(err, "unable to read workspace id")
	}

	cv, err := c.Client.ConfigurationVersions.Create(c.Context, workspaceId, tfe.ConfigurationVersionCreateOptions{
		AutoQueueRuns: tfe.Bool(false),
		Speculative:   tfe.Bool(isSpeculative),
	})
	if err != nil {
		return errors.Wrap(err, "failed to create configuration version")
	}

	o.AddMessageUserProvided("Upload code to Configuration Version...", "")
	err = c.Client.ConfigurationVersions.Upload(c.Context, cv.UploadURL, directory)
	if err != nil {
		return errors.Wrap(err, "failed to upload code to the configuration version")
	}

	o.AddMessageUserProvided("Configuration Version Created", "")
	o.AddDeferredMessageRead("ID", cv.ID)
	o.AddDeferredMessageRead("Speculative", cv.Speculative)
	o.Close()

	return nil
}

func cvShow(c TfxClientContext, configurationId string) error {
	o.AddMessageUserProvided("Show Configuration Version for Workspace from Id:", configurationId)
	cv, err := c.Client.ConfigurationVersions.ReadWithOptions(c.Context, configurationId, &tfe.ConfigurationVersionReadOptions{
		Include: []tfe.ConfigVerIncludeOpt{"ingress_attributes"},
	})
	if err != nil {
		return errors.Wrap(err, "failed to read configuration version from provided id")
	}

	o.AddDeferredMessageRead("ID", cv.ID)
	o.AddDeferredMessageRead("Status", cv.Status)
	o.AddDeferredMessageRead("Speculative", cv.Speculative)
	if cv.ErrorMessage != "" {
		o.AddDeferredMessageRead("Error Message", cv.ErrorMessage)
	}
	if cv.IngressAttributes != nil { // only valid for VCS-driven workflow
		o.AddDeferredMessageRead("Repo", cv.IngressAttributes.Identifier)
		o.AddDeferredMessageRead("Branch", cv.IngressAttributes.Branch)
		o.AddDeferredMessageRead("Commit", cv.IngressAttributes.CommitSHA)
		o.AddDeferredMessageRead("Message", cv.IngressAttributes.CommitMessage)
		o.AddDeferredMessageRead("Link", cv.IngressAttributes.CommitURL)
	}
	o.Close()

	return nil
}

func cvDownload(c TfxClientContext, configurationId string, directory string) error {
	o.AddMessageUserProvided("Downloading Configuration Version from Id:", configurationId)
	var err error
	// Determine a directory to unpack the slug contents into.
	if directory != "" {
		if !isDirectory(directory) {
			return errors.Wrap(err, "configuration version directory is not valid")
		}
	} else {
		o.AddMessageUserProvided("Directory not supplied, creating a temp directory", "")
		dst, err := ioutil.TempDir("", "slug")
		if err != nil {
			return errors.Wrap(err, "failed to create temp directory")
		}
		directory = dst
	}

	cv, err := c.Client.ConfigurationVersions.Download(c.Context, configurationId)
	if err != nil {
		return errors.Wrap(err, "failed to download configuration version")
	}

	o.AddMessageUserProvided("Configuration Version Found, download started...", "")
	// convert byte slice to io.Reader
	reader := bytes.NewReader(cv)
	if err := slug.Unpack(reader, directory); err != nil {
		return errors.Wrap(err, "failed to unpack configuration version slug")
	}

	o.AddDeferredMessageRead("Status", "Success")
	o.AddDeferredMessageRead("Directory", directory)
	o.Close()

	return nil
}
