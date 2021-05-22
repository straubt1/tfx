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
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/tidwall/sjson"

	"github.com/fatih/color"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	stateCmd = &cobra.Command{
		Use:   "state",
		Short: "State Versions",
		Long:  "Work with State Versions of a TFx Workspace.",
	}

	stateListCmd = &cobra.Command{
		Use:   "list",
		Short: "List State Versions",
		Long:  "List State Versions of a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return stateList()
		},
		PreRun: bindPFlags,
	}

	stateDownloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download State Version",
		Long:  "Download State Version for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return stateDownload()
		},
		PreRun: bindPFlags,
	}

	stateCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create State Version",
		Long:  "Create State Version for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return stateCreate()
		},
		PreRun: bindPFlags,
	}

	stateShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show State Version",
		Long:  "Show State Version details for a TFx Workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return stateShow()
		},
		PreRun: bindPFlags,
	}
)

func init() {
	// All `tfx state list` commands
	stateListCmd.Flags().StringP("workspaceName", "w", "", "Workspace name")
	stateListCmd.MarkFlagRequired("workspaceName")

	// `tfx state download`
	stateDownloadCmd.Flags().StringP("stateId", "i", "", "State Version Id (i.e. sv-*)")
	stateDownloadCmd.Flags().StringP("filename", "f", "", "File to save State Version to (i.e. terraform.tfstate)")
	stateDownloadCmd.MarkFlagRequired("stateId")
	stateDownloadCmd.MarkFlagRequired("filename")

	// `tfx state create`
	stateCreateCmd.Flags().StringP("workspaceName", "w", "", "Workspace name")
	stateCreateCmd.Flags().StringP("filename", "f", "", "File to create a State Version from (i.e. terraform.tfstate)")
	stateCreateCmd.MarkFlagRequired("workspaceName")
	stateCreateCmd.MarkFlagRequired("filename")

	// `tfx state show`
	stateShowCmd.Flags().StringP("stateId", "i", "", "State Version Id (i.e. sv-*)")
	stateShowCmd.MarkFlagRequired("stateId")

	rootCmd.AddCommand(stateCmd)
	stateCmd.AddCommand(stateListCmd)
	stateCmd.AddCommand(stateDownloadCmd)
	stateCmd.AddCommand(stateCreateCmd)
	stateCmd.AddCommand(stateShowCmd)
}

func stateList() error {
	// Validate flags
	orgName := *viperString("tfeOrganization")
	wsName := *viperString("workspaceName")
	client, ctx := getClientContext()

	// Read workspace
	fmt.Print("Reading Workspace ID for Name: ", color.GreenString(wsName), " ...")
	w, err := client.Workspaces.Read(ctx, orgName, wsName)
	if err != nil {
		logError(err, "failed to read workspace id")
	}
	fmt.Println(" Found:", color.BlueString(w.ID))

	// Get current state version
	curOpts := &tfe.StateVersionCurrentOptions{
		Include: "outputs",
	}
	state, err := client.StateVersions.CurrentWithOptions(ctx, w.ID, curOpts)
	if err != nil {
		logError(err, "failed to get current state version")
	}
	fmt.Println(color.BlueString("ID:     "), state.ID)
	fmt.Println(color.BlueString("Create: "), state.CreatedAt)
	// fmt.Println(color.BlueString("Run:    "), state.Run.ID) // Run can by nil

	// Get all state versions
	stateList, err := client.StateVersions.List(ctx, tfe.StateVersionListOptions{
		Organization: &orgName,
		Workspace:    &wsName,
	})
	if err != nil {
		logError(err, "failed to list state versions")
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Id", "Serial", "Run", "Created"})
	for _, i := range stateList.Items {
		t.AppendRow(table.Row{i.ID, i.Serial, i.Serial, i.CreatedAt})
	}
	t.SetStyle(table.StyleRounded)
	t.Render()

	return nil
}

func stateDownload() error {
	stateId := *viperString("stateId")
	filename := *viperString("filename")
	client, ctx := getClientContext()

	fmt.Print("Reading State Version for ID ", color.GreenString(stateId), " ...")
	state, err := client.StateVersions.Read(ctx, stateId)
	if err != nil {
		logError(err, "failed to read state version")
	}
	fmt.Println(" Found")

	buff, err := client.StateVersions.Download(ctx, state.DownloadURL)
	if err != nil {
		logError(err, "failed to download state version")
	}

	err = ioutil.WriteFile(filename, buff, 0644)
	if err != nil {
		logError(err, "failed to save state version")
	}

	return nil
}

type StateFile struct {
	Version          int64  `json:"version"`
	TerraformVersion string `json:"terraform_version"`
	Serial           int64  `json:"serial"`
	Lineage          string `json:"lineage"`
}

// to create a state file, the serial needs incremented and the lineage must match.
// reads the Workspace to increment serial and maintain lineage
func stateCreate() error {
	// Validate flags
	orgName := *viperString("tfeOrganization")
	wsName := *viperString("workspaceName")
	filename := *viperString("filename")
	client, ctx := getClientContext()

	// Open file and get bytes[]
	f, err := os.Open(filename)
	if err != nil {
		logError(err, "failed to open file")
	}
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	contentString := string(content)
	defer f.Close()

	// Read workspace
	fmt.Print("Reading Workspace ", color.GreenString(wsName), " for ID...")
	w, err := client.Workspaces.Read(ctx, orgName, wsName)
	if err != nil {
		logError(err, "failed to read workspace id")
	}
	fmt.Println(" Found:", w.ID)

	// // Read latest state file
	// currentState, err := client.StateVersions.Current(ctx, w.ID)
	// if err != nil {
	// 	logError(err, "failed to get current state version")
	// }
	// // if currentState == nil {

	// // }
	// currentStateContents, err := client.StateVersions.Download(ctx, currentState.DownloadURL)
	// if err != nil {
	// 	logError(err, "failed to get download current state file")
	// }
	// fmt.Println(string(currentStateContents))
	// var currSt StateFile
	// err = json.Unmarshal(content, &currentStateContents)
	// if err != nil {
	// 	logError(err, "failed to marshal state file")
	// }
	// fmt.Println("Current State file ...")
	// fmt.Println(color.BlueString("Lineage:    "), currSt.Lineage)
	// fmt.Println(color.BlueString("Serial:     "), currSt.Serial)
	// fmt.Println(color.BlueString("Version:    "), currSt.Version)
	// fmt.Println(color.BlueString("Terraform:  "), currSt.TerraformVersion)
	// fmt.Println()

	// Read state file so we can get info (serial, lineage, etc...)
	var st StateFile
	err = json.Unmarshal(content, &st)
	if err != nil {
		logError(err, "failed to marshal state file")
	}
	fmt.Println("Reading State file ", color.GreenString(filename), " ...")
	fmt.Println(color.BlueString("Lineage:    "), st.Lineage)
	fmt.Println(color.BlueString("Serial:     "), st.Serial)
	fmt.Println(color.BlueString("Version:    "), st.Version)
	fmt.Println(color.BlueString("Terraform:  "), st.TerraformVersion)
	fmt.Println()

	newSerial := st.Serial + 1
	newContentString, err := sjson.Set(contentString, "serial", newSerial)
	newContent := []byte(newContentString)
	if err != nil {
		logError(err, "failed to set new serial")
	}
	fmt.Println("New State Serial ", color.BlueString(strconv.FormatInt(newSerial, 10)), " ...")

	encoded := base64.StdEncoding.EncodeToString(newContent)
	checksum := fmt.Sprintf("%x", md5.Sum(newContent))
	fmt.Println("ENCODED: " + encoded)
	fmt.Println("CHECKSUM: " + checksum)

	_, err = client.Workspaces.Lock(ctx, w.ID, tfe.WorkspaceLockOptions{
		Reason: tfe.String("TFx locking to create new State File"),
	})
	if err != nil {
		logError(err, "failed to lock workspace")
	}

	// Create State Version
	fmt.Print("Creating State Version ...")
	state, err := client.StateVersions.Create(ctx, w.ID, tfe.StateVersionCreateOptions{
		Lineage: tfe.String(st.Lineage),
		MD5:     tfe.String(checksum),
		Serial:  tfe.Int64(newSerial),
		State:   tfe.String(encoded),
		// Force:   new(bool),
		// Run:     &tfe.Run{},
	})
	if err != nil {
		logError(err, "failed to create state version")
	}
	_ = state

	_, err = client.Workspaces.Unlock(ctx, w.ID)
	if err != nil {
		logError(err, "failed to unlock workspace")
	}
	return nil
}

func stateShow() error {
	// Validate flags
	stateId := *viperString("stateId")
	client, ctx := getClientContext()

	// Read Config Version
	fmt.Print("Reading State Version for ID ", color.GreenString(stateId), " ...")
	state, err := client.StateVersions.ReadWithOptions(ctx, stateId, &tfe.StateVersionReadOptions{
		Include: "outputs",
	})
	if err != nil {
		logError(err, "failed to read state version")
	}
	fmt.Println(" Found")
	fmt.Println(color.BlueString("ID:     "), state.ID)
	fmt.Println(color.BlueString("Create: "), state.CreatedAt)
	fmt.Println(color.BlueString("Run:    "), state.Run.ID)

	return nil
}
