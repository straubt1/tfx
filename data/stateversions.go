// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	"bufio"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"os"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/output"
	"github.com/tidwall/sjson"
)

// StateFile mirrors key metadata from state for creation
type StateFile struct {
	Version          int64  `json:"version"`
	TerraformVersion string `json:"terraform_version"`
	Serial           int64  `json:"serial"`
	Lineage          string `json:"lineage"`
}

// FetchStateVersions lists state versions for a workspace name with max-items
func FetchStateVersions(c *client.TfxClient, orgName, workspaceName string, maxItems int) ([]*tfe.StateVersion, error) {
	output.Get().Logger().Debug("Fetching state versions", "organization", orgName, "workspace", workspaceName, "maxItems", maxItems)

	pageSize := 100
	if maxItems > 0 && maxItems < 100 {
		pageSize = maxItems
	}

	var all []*tfe.StateVersion
	opts := &tfe.StateVersionListOptions{
		ListOptions:  tfe.ListOptions{PageNumber: 1, PageSize: pageSize},
		Organization: orgName,
		Workspace:    workspaceName,
	}

	for {
		res, err := c.Client.StateVersions.List(c.Context, opts)
		if err != nil {
			output.Get().Logger().Error("Failed to list state versions", "workspace", workspaceName, "page", opts.PageNumber, "error", err)
			return nil, err
		}

		all = append(all, res.Items...)
		if maxItems > 0 && len(all) >= maxItems {
			break
		}

		if res.CurrentPage >= res.TotalPages {
			break
		}
		opts.PageNumber = res.NextPage
	}

	if maxItems > 0 && len(all) > maxItems {
		all = all[:maxItems]
	}

	output.Get().Logger().Debug("State versions fetched", "count", len(all))
	return all, nil
}

// CreateStateVersionFromFile reads a state file and creates a new state version
func CreateStateVersionFromFile(c *client.TfxClient, orgName, workspaceName, filename string) (*tfe.StateVersion, error) {
	if filename == "" {
		return nil, errors.New("state file does not exist")
	}

	// Read file
	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open state file")
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	content, err := ioReadAll(reader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read state file")
	}

	// Parse state file metadata
	var st StateFile
	if err := json.Unmarshal(content, &st); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal state file")
	}

	// Get workspace ID
	workspaceID, err := GetWorkspaceID(c, orgName, workspaceName)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read workspace id")
	}

	// Get current state to increment serial
	newSerial := int64(0)
	currentState, _ := c.Client.StateVersions.ReadCurrent(c.Context, workspaceID)
	if currentState != nil {
		newSerial = currentState.Serial + 1
	}

	// Prepare state payload
	opts, err := updateStateFile(string(content), newSerial, st.Lineage)
	if err != nil {
		return nil, errors.Wrap(err, "unable to update state file")
	}

	// Lock workspace
	if _, err := c.Client.Workspaces.Lock(c.Context, workspaceID, tfe.WorkspaceLockOptions{Reason: tfe.String("TFx locking to create new State File")}); err != nil {
		return nil, errors.Wrap(err, "failed to lock workspace")
	}

	// Create state version
	sv, err := c.Client.StateVersions.Create(c.Context, workspaceID, *opts)
	if err != nil {
		// attempt unlock on error
		c.Client.Workspaces.Unlock(c.Context, workspaceID)
		return nil, errors.Wrap(err, "failed to create state version")
	}

	// Unlock workspace (best effort)
	if _, err := c.Client.Workspaces.Unlock(c.Context, workspaceID); err != nil {
		output.Get().Logger().Warn("Failed to unlock workspace after state create", "workspaceID", workspaceID, "error", err)
	}

	return sv, nil
}

// FetchStateVersion reads a state version with includes
func FetchStateVersion(c *client.TfxClient, stateID string) (*tfe.StateVersion, error) {
	sv, err := c.Client.StateVersions.ReadWithOptions(c.Context, stateID, &tfe.StateVersionReadOptions{
		Include: []tfe.StateVersionIncludeOpt{"created_by", "run", "run.created_by", "run.configuration_version", "outputs"},
	})
	if err != nil {
		return nil, err
	}
	return sv, nil
}

// DownloadStateVersion downloads the state version bytes
func DownloadStateVersion(c *client.TfxClient, stateID string) ([]byte, error) {
	sv, err := c.Client.StateVersions.Read(c.Context, stateID)
	if err != nil {
		return nil, err
	}
	return c.Client.StateVersions.Download(c.Context, sv.DownloadURL)
}

// updateStateFile prepares the state create options based on new serial and lineage
func updateStateFile(stateContents string, newSerial int64, lineage string) (*tfe.StateVersionCreateOptions, error) {
	// take current json file, and update serial to the new serial
	newContentString, err := sjson.Set(stateContents, "serial", newSerial)
	newContent := []byte(newContentString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set new serial")
	}

	// base64 encode and calculate checksum
	encoded := base64.StdEncoding.EncodeToString(newContent)
	checksum := md5.Sum(newContent)
	checksumStr := make([]byte, 0, 32)
	for _, b := range checksum {
		checksumStr = append(checksumStr, "0123456789abcdef"[b>>4], "0123456789abcdef"[b&0x0f])
	}

	svCreateOpts := tfe.StateVersionCreateOptions{
		Lineage: tfe.String(lineage),
		MD5:     tfe.String(string(checksumStr)),
		Serial:  tfe.Int64(newSerial),
		State:   tfe.String(encoded),
	}
	return &svCreateOpts, nil
}

// ioReadAll isolates ioutil.ReadAll deprecation concerns
func ioReadAll(r *bufio.Reader) ([]byte, error) {
	return r.ReadBytes(0)
}
