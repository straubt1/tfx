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
	"bufio"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tidwall/sjson"
)

var (
	// `tfx workspace state` commands
	stateCmd = &cobra.Command{
		Use:     "state-version",
		Aliases: []string{"sv"},
		Short:   "State Version Commands",
		Long:    "Work with State Versions of a TFx Workspace.",
	}

	// `tfx workspace state list` command
	stateListCmd = &cobra.Command{
		Use:   "list",
		Short: "List State Versions",
		Long:  "List State Versions of a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			m := *viperInt("max-items")
			if *viperBool("all") {
				m = math.MaxInt
			}
			return stateList(
				getTfxClientContext(),
				*viperString("workspace-name"),
				m)
		},
	}

	// `tfx workspace state create` command
	stateCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create State Version",
		Long:  "Create State Version for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !isFile(*viperString("filename")) {
				return errors.New("state file does not exist")
			}

			return stateCreate(
				getTfxClientContext(),
				*viperString("workspace-name"),
				*viperString("filename"))
		},
	}

	// `tfx workspace state show` command
	stateShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show State Version",
		Long:  "Show State Version details for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return stateShow(
				getTfxClientContext(),
				*viperString("state-id"))
		},
	}

	// `tfx workspace state download` command
	stateDownloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download State Version",
		Long:  "Download State Version for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			directory, err := getDirectory(*viperString("directory"))
			if err != nil {
				return err
			}

			return stateDownload(
				getTfxClientContext(),
				*viperString("state-id"),
				directory)
		},
	}
)

func init() {
	// `tfx workspace state list` command
	stateListCmd.Flags().StringP("workspace-name", "w", "", "Workspace name")
	stateListCmd.Flags().IntP("max-items", "", 10, "Max number of results (optional)")
	stateListCmd.Flags().BoolP("all", "a", false, "Retrieve all results regardless of maxItems flag (optional)")
	stateListCmd.MarkFlagRequired("workspace-name")

	// `tfx workspace state create` command
	stateCreateCmd.Flags().StringP("workspace-name", "w", "", "Workspace name")
	stateCreateCmd.Flags().StringP("filename", "f", "", "Filename of the state file to create")
	stateCreateCmd.MarkFlagRequired("workspace-name")
	stateCreateCmd.MarkFlagRequired("filename")

	// `tfx workspace state show` command
	stateShowCmd.Flags().StringP("state-id", "i", "", "State Version Id (i.e. sv-*)")
	stateShowCmd.MarkFlagRequired("state-id")

	// `tfx workspace state download` command
	stateDownloadCmd.Flags().StringP("state-id", "i", "", "State Version Id (i.e. sv-*)")
	stateDownloadCmd.Flags().StringP("directory", "d", "", "Directory of download state version (optional, defaults to a temp directory)")
	stateDownloadCmd.Flags().StringP("filename", "f", "", "Filename to save State Version as (optional)")
	stateDownloadCmd.MarkFlagRequired("state-id")

	workspaceCmd.AddCommand(stateCmd)
	stateCmd.AddCommand(stateListCmd)
	stateCmd.AddCommand(stateDownloadCmd)
	stateCmd.AddCommand(stateCreateCmd)
	stateCmd.AddCommand(stateShowCmd)
}

func stateListAll(c TfxClientContext, workspaceName string, maxItems int) ([]*tfe.StateVersion, error) {
	pageSize := 100

	if maxItems < 100 {
		pageSize = maxItems // Only get what we need in one page
	}
	allItems := []*tfe.StateVersion{}
	opts := tfe.StateVersionListOptions{
		ListOptions:  tfe.ListOptions{PageNumber: 1, PageSize: pageSize},
		Organization: c.OrganizationName,
		Workspace:    workspaceName,
	}
	for {
		items, err := c.Client.StateVersions.List(c.Context, &opts)
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

func stateList(c TfxClientContext, workspaceName string, maxItems int) error {
	o.AddMessageUserProvided("List State Versions for Workspace:", workspaceName)
	items, err := stateListAll(c, workspaceName, maxItems)
	if err != nil {
		return errors.Wrap(err, "failed to list state versions")
	}

	o.AddTableHeader("Id", "Terraform Version", "Serial", "Run Id", "Created")
	for _, i := range items {
		runId := ""
		if i.Run != nil {
			runId = i.Run.ID
		}
		o.AddTableRows(i.ID, i.TerraformVersion, i.Serial, runId, FormatDateTime(i.CreatedAt))
	}

	return nil
}

type StateFile struct {
	Version          int64  `json:"version"`
	TerraformVersion string `json:"terraform_version"`
	Serial           int64  `json:"serial"`
	Lineage          string `json:"lineage"`
}

// Reads a state file by filename
// Reads the workspace current state
// If provided lineage is not equal to the current state of workspace, API will error
// Will increment Serial to one more than last state (or zero if none)
func stateCreate(c TfxClientContext, workspaceName string, filename string) error {
	o.AddMessageUserProvided("Create State Version for Workspace:", workspaceName)

	o.AddMessageUserProvided("Read state file and Parse:", filename)
	f, err := os.Open(filename)
	if err != nil {
		return errors.Wrap(err, "failed to open state file")
	}
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	contentString := string(content)
	defer f.Close()
	// Read state file so we can get info (serial, lineage, etc...)
	var st StateFile
	err = json.Unmarshal(content, &st)
	if err != nil {
		return errors.Wrap(err, "failed to marshal state file")
	}
	o.AddDeferredMessageRead("Provided Lineage", st.Lineage)
	o.AddDeferredMessageRead("Provided Serial", st.Serial)

	workspaceId, err := getWorkspaceId(c, workspaceName)
	if err != nil {
		return errors.Wrap(err, "unable to read workspace id")
	}

	// Get current State Version, so we can increment the serial
	// fmt.Print("Reading Current State Version ...")
	newSerial := int64(0)
	currentState, _ := c.Client.StateVersions.ReadCurrent(c.Context, workspaceId)
	if currentState != nil {
		newSerial = currentState.Serial + 1
		o.AddDeferredMessageRead("Existing Serial", currentState.Serial)
	} else {
		o.AddDeferredMessageRead("Existing Serial", "")
	}
	// else -> Workspace not found, workspace does not have a current state version, or user unauthorized to perform action.
	// If currentState == nil, then assume no current state, serial will be zero
	opts, err := updateStateFile(contentString, newSerial, st.Lineage)
	if err != nil {
		return errors.Wrap(err, "unable to update state file")
	}

	o.AddMessageUserProvided("Locking Workspace...", "")
	_, err = c.Client.Workspaces.Lock(c.Context, workspaceId, tfe.WorkspaceLockOptions{
		Reason: tfe.String("TFx locking to create new State File"),
	})
	if err != nil {
		return errors.Wrap(err, "failed to lock workspace")
	}

	o.AddMessageUserProvided("Creating State Version...", "")
	state, err := c.Client.StateVersions.Create(c.Context, workspaceId, *opts)
	if err != nil {
		// if error, attempt to unlock workspace
		o.AddMessageUserProvided("Unlocking Workspace...", "")
		_, err2 := c.Client.Workspaces.Unlock(c.Context, workspaceId)
		if err2 != nil {
			return errors.Wrap(err2, "failed to unlock workspace")
		}
		return errors.Wrap(err, "failed to create state version")
	}
	o.AddDeferredMessageRead("Created Serial", state.Serial)

	o.AddMessageUserProvided("Unlocking Workspace...", "")
	_, err = c.Client.Workspaces.Unlock(c.Context, workspaceId)
	if err != nil {
		return errors.Wrap(err, "failed to unlock workspace")
	}

	return nil
}

func stateShow(c TfxClientContext, stateId string) error {
	o.AddMessageUserProvided("Show State Version for Workspace from Id:", stateId)
	state, err := c.Client.StateVersions.ReadWithOptions(c.Context, stateId, &tfe.StateVersionReadOptions{
		Include: []tfe.StateVersionIncludeOpt{"created_by", "run", "run.created_by", "run.configuration_version", "outputs"},
	})
	if err != nil {
		return errors.Wrap(err, "failed to read state version from provided id")
	}

	o.AddDeferredMessageRead("ID", state.ID)
	o.AddDeferredMessageRead("Created", FormatDateTime(state.CreatedAt))
	o.AddDeferredMessageRead("Terraform Version", state.TerraformVersion)
	o.AddDeferredMessageRead("Serial", state.Serial)
	o.AddDeferredMessageRead("State Version", state.StateVersion)
	o.AddDeferredMessageRead("Run Id", state.Run.ID)
	for _, i := range state.Outputs {
		o.AddDeferredMessageRead("output_"+i.Name, i.Value)
	}

	return nil
}

func stateDownload(c TfxClientContext, stateId string, directory string) error {
	o.AddMessageUserProvided("Downloading State Version from Id:", stateId)
	state, err := c.Client.StateVersions.Read(c.Context, stateId)
	if err != nil {
		return errors.Wrap(err, "failed to read state version with provided id")
	}

	o.AddMessageUserProvided("State Version Found, download started...", "")
	buff, err := c.Client.StateVersions.Download(c.Context, state.DownloadURL)
	if err != nil {
		return errors.Wrap(err, "failed to download state version")
	}

	fullPath := filepath.Join(directory, fmt.Sprintf("%s.state", stateId))
	err = ioutil.WriteFile(fullPath, buff, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to save state version")
	}

	o.AddDeferredMessageRead("Status", "Success")
	o.AddDeferredMessageRead("File", fullPath)

	return nil
}

func updateStateFile(stateContents string, newSerial int64, lineage string) (*tfe.StateVersionCreateOptions, error) {
	// take current json file, and update serial to the new serial
	newContentString, err := sjson.Set(stateContents, "serial", newSerial)
	newContent := []byte(newContentString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set new serial")
	}

	// base64 encode and calculate checksum
	encoded := base64.StdEncoding.EncodeToString(newContent)
	checksum := fmt.Sprintf("%x", md5.Sum(newContent))

	svCreateOpts := tfe.StateVersionCreateOptions{
		Lineage: tfe.String(lineage),
		MD5:     tfe.String(checksum),
		Serial:  tfe.Int64(newSerial),
		State:   tfe.String(encoded),
	}

	return &svCreateOpts, nil
}
