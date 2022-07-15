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
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mmcdole/gofeed"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

// releaseCmd represents the release command
var (
	releaseCmd = &cobra.Command{
		Use:   "release",
		Short: "TFE release commands",
		Long:  "Work with binaries needed for TFE release installations.",
	}

	releaseTfeCmd = &cobra.Command{
		Use:   "tfe",
		Short: "TFE release commands",
		Long:  "Terraform Enterprise release commands.",
	}

	releaseTfeListCmd = &cobra.Command{
		Use:   "list",
		Short: "List TFE release",
		Long:  "List available Terraform Enterprise release.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return releaseTfeList()
		},
		PreRun: bindPFlags,
	}

	releaseTfeShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show TFE release",
		Long:  "Show a Terraform Enterprise release, including release notes.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return releaseTfeShow()
		},
		PreRun: bindPFlags,
	}

	releaseTfeDownloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download TFE release binary",
		Long:  "Download a Terraform Enterprise release binary.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return releaseTfeDownload()
		},
		PreRun: bindPFlags,
	}

	releaseReplicatedCmd = &cobra.Command{
		Use:   "replicated",
		Short: "Replicated release commands",
		Long:  "Replicated release commands.",
	}

	releaseReplicatedListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Replicated binaries",
		Long:  "List available Replicated release.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return releaseReplicatedList()
		},
		PreRun: bindPFlags,
	}

	releaseReplicatedDownloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download Replicated release binary",
		Long:  "Download a Replicated release binary.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return releaseReplicatedDownload()
		},
		PreRun: bindPFlags,
	}
)

func init() {
	// `tfx release tfe list`
	releaseTfeListCmd.Flags().StringP("licenseId", "l", "", "License Id for TFE/Replicated")
	releaseTfeListCmd.Flags().StringP("password", "p", "", "Password to authenticate")
	releaseTfeListCmd.Flags().IntP("maxResults", "r", 10, "The number of results to print (optional, defaults to 10)")
	releaseTfeListCmd.MarkFlagRequired("licenseId")
	releaseTfeListCmd.MarkFlagRequired("password")

	// `tfx release tfe show`
	releaseTfeShowCmd.Flags().StringP("licenseId", "l", "", "License Id for TFE/Replicated")
	releaseTfeShowCmd.Flags().StringP("password", "p", "", "Password to authenticate")
	releaseTfeShowCmd.Flags().StringP("release", "r", "", "Release Sequence (i.e. 610, 619, etc...)")
	releaseTfeShowCmd.MarkFlagRequired("licenseId")
	releaseTfeShowCmd.MarkFlagRequired("password")
	releaseTfeShowCmd.MarkFlagRequired("release")

	// `tfx release tfe download`
	releaseTfeDownloadCmd.Flags().StringP("licenseId", "l", "", "License Id for TFE/Replicated")
	releaseTfeDownloadCmd.Flags().StringP("password", "p", "", "Password to authenticate")
	releaseTfeDownloadCmd.Flags().StringP("release", "r", "", "Release Sequence (i.e. 610, 619, etc...)")
	releaseTfeDownloadCmd.Flags().StringP("directory", "d", "./", "Directory to save binary (optional, defaults to current directory)")
	releaseTfeDownloadCmd.MarkFlagRequired("licenseId")
	releaseTfeDownloadCmd.MarkFlagRequired("password")
	releaseTfeDownloadCmd.MarkFlagRequired("release")

	// `tfx release replicated list`
	releaseReplicatedListCmd.Flags().IntP("maxResults", "r", 10, "The number of results to print (optional, defaults to 10)")

	// `tfx release replicated download`
	releaseReplicatedDownloadCmd.Flags().StringP("directory", "d", "./", "Directory to save binary (optional, defaults to current directory)")
	releaseReplicatedDownloadCmd.Flags().StringP("version", "v", "", "Version of Replicated to Download (i.e. 0.0.1)")

	rootCmd.AddCommand(releaseCmd)
	releaseCmd.AddCommand(releaseTfeCmd)
	releaseTfeCmd.AddCommand(releaseTfeListCmd)
	releaseTfeCmd.AddCommand(releaseTfeShowCmd)
	releaseTfeCmd.AddCommand(releaseTfeDownloadCmd)

	releaseCmd.AddCommand(releaseReplicatedCmd)
	releaseReplicatedCmd.AddCommand(releaseReplicatedListCmd)
	releaseReplicatedCmd.AddCommand(releaseReplicatedDownloadCmd)
}

func releaseTfeList() error {
	// Validate flags
	licenseId := *viperString("licenseId")
	password := *viperString("password")
	maxResults := *viperInt("maxResults")

	tfeBinaries, err := ListTFEBinaries(password, licenseId)
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Sequence", "Label", "Required", "Release Date"})
	for index, i := range tfeBinaries.release {
		t.AppendRow(table.Row{i.releaseequence, i.Label, i.Required, i.ReleaseDate.String()})
		if index+1 >= maxResults {
			break
		}
	}
	t.SetStyle(table.StyleRounded)
	t.Render()

	return nil
}

func releaseTfeShow() error {
	// Validate flags
	licenseId := *viperString("licenseId")
	password := *viperString("password")
	release := *viperInt("release")

	tfeBinaries, err := ListTFEBinaries(password, licenseId)
	if err != nil {
		return err
	}

	idx := slices.IndexFunc(tfeBinaries.release, func(c TFERelease) bool { return c.releaseequence == release })
	if idx < 0 {
		fmt.Println(color.RedString("Error: "), "Unable to find release sequence: ", release)
		return nil
	}
	tfeRelease := tfeBinaries.release[idx]

	fmt.Println(" Found")
	fmt.Println(color.BlueString("Release Sequence: "), tfeRelease.releaseequence)
	fmt.Println(color.BlueString("Label:            "), tfeRelease.Label)
	fmt.Println(color.BlueString("Release Date:     "), tfeRelease.ReleaseDate)
	fmt.Println(color.BlueString("Required:         "), tfeRelease.Required)
	fmt.Println(color.BlueString("Release Notes:    "))
	fmt.Println(tfeRelease.ReleaseNotes)

	return nil
}

func releaseTfeDownload() error {
	// Validate flags
	licenseId := *viperString("licenseId")
	password := *viperString("password")
	release := *viperString("release")
	directory := *viperString("directory")

	// Get url
	tfeUrl, err := GetTFEBinary(password, licenseId, release)
	if err != nil {
		return err
	}

	// Verify directory
	_, err = os.Stat(directory)
	if err != nil {
		fmt.Println(color.RedString("Error: Invalid directory "), directory)
		return err
	}

	// Verify trailing, if not add it
	if !strings.HasSuffix(directory, "/") {
		directory += "/"
	}
	path := fmt.Sprintf("%stfe-%s.release", directory, release)

	//Download file
	err = DownloadBinary(tfeUrl.URL, path)
	if err != nil {
		return err
	}

	return nil
}

func releaseReplicatedList() error {
	// Validate flags
	maxResults := *viperInt("maxResults")

	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://release-notes.replicated.com/index.xml")
	fmt.Println(feed.Title)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Version", "Published Date"})
	for index, i := range feed.Items {
		t.AppendRow(table.Row{i.Title, i.Published})
		if index+1 >= maxResults {
			break
		}
	}
	t.SetStyle(table.StyleRounded)
	t.Render()

	return nil
}

func releaseReplicatedDownload() error {
	// Validate flags
	directory := *viperString("directory")
	// Attempt to prevent a non semantic version from being requested
	version, err := viperSemanticVersionString("version")
	if err != nil {
		logError(err, "failed to parse semantic version")
	}

	// Get url - escape "%2B" as "%%2B", + symbol
	url := fmt.Sprintf("https://s3.amazonaws.com/replicated-airgap-work/stable/replicated-%s%%2B%s%%2B%s.tar.gz",
		version,
		version,
		version)

	// Verify directory
	_, err = os.Stat(directory)
	if err != nil {
		fmt.Println(color.RedString("Error: Invalid directory "), directory)
		return err
	}

	// Verify trailing, if not add it
	if !strings.HasSuffix(directory, "/") {
		directory += "/"
	}
	path := fmt.Sprintf("%sreplicated-%s.targz", directory, version)

	//Download file
	err = DownloadBinary(url, path)
	if err != nil {
		return err
	}

	return nil
}
