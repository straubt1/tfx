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

// releasesCmd represents the releases command
var (
	releasesCmd = &cobra.Command{
		Use:   "releases",
		Short: "releases helper commands",
		Long:  "Work with binaries needed for TFE releases installations.",
	}

	releasesTfeCmd = &cobra.Command{
		Use:   "tfe",
		Short: "TFE releases commands",
		Long:  "TFE releases commands to work with TFE binaries.",
	}

	releasesTfeListCmd = &cobra.Command{
		Use:   "list",
		Short: "List TFE releases",
		Long:  "List available TFE releases for releases download.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return releasesTfeList()
		},
		PreRun: bindPFlags,
	}

	releasesTfeShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show TFE release",
		Long:  "Show a TFE release, including release notes.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return releasesTfeShow()
		},
		PreRun: bindPFlags,
	}

	releasesTfeDownloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download TFE release binary",
		Long:  "Download a TFE release binary.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return releasesTfeDownload()
		},
		PreRun: bindPFlags,
	}

	releasesReplicatedCmd = &cobra.Command{
		Use:   "replicated",
		Short: "Replicated releases commands",
		Long:  "Replicated releases commands to work with Replicated binaries.",
	}

	releasesReplicatedListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Replicated binaries",
		Long:  "List available Replicated binaries for download.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return releasesReplicatedList()
		},
		PreRun: bindPFlags,
	}

	releasesReplicatedDownloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download Replicated release binary",
		Long:  "Download a Replicated release binary.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return releasesReplicatedDownload()
		},
		PreRun: bindPFlags,
	}
)

func init() {
	// `tfx releases tfe list`
	releasesTfeListCmd.Flags().StringP("licenseId", "l", "", "License Id for TFE/Replicated")
	releasesTfeListCmd.Flags().StringP("password", "p", "", "Password to authenticate")
	releasesTfeListCmd.Flags().IntP("maxResults", "r", 10, "The number of results to print (optional, defaults to 10)")
	releasesTfeListCmd.MarkFlagRequired("licenseId")
	releasesTfeListCmd.MarkFlagRequired("password")

	// `tfx releases tfe show`
	releasesTfeShowCmd.Flags().StringP("licenseId", "l", "", "License Id for TFE/Replicated")
	releasesTfeShowCmd.Flags().StringP("password", "p", "", "Password to authenticate")
	releasesTfeShowCmd.Flags().StringP("release", "r", "", "Release Sequence (i.e. 610, 619, etc...)")
	releasesTfeShowCmd.MarkFlagRequired("licenseId")
	releasesTfeShowCmd.MarkFlagRequired("password")
	releasesTfeShowCmd.MarkFlagRequired("release")

	// `tfx releases tfe download`
	releasesTfeDownloadCmd.Flags().StringP("licenseId", "l", "", "License Id for TFE/Replicated")
	releasesTfeDownloadCmd.Flags().StringP("password", "p", "", "Password to authenticate")
	releasesTfeDownloadCmd.Flags().StringP("release", "r", "", "Release Sequence (i.e. 610, 619, etc...)")
	releasesTfeDownloadCmd.Flags().StringP("directory", "d", "./", "Directory to save binary (optional, defaults to current directory)")
	releasesTfeDownloadCmd.MarkFlagRequired("licenseId")
	releasesTfeDownloadCmd.MarkFlagRequired("password")
	releasesTfeDownloadCmd.MarkFlagRequired("release")

	// `tfx releases replicated list`
	releasesReplicatedListCmd.Flags().IntP("maxResults", "r", 10, "The number of results to print (optional, defaults to 10)")

	// `tfx releases replicated download`
	releasesReplicatedDownloadCmd.Flags().StringP("directory", "d", "./", "Directory to save binary (optional, defaults to current directory)")
	releasesReplicatedDownloadCmd.Flags().StringP("version", "v", "", "Version of Replicated to Download (i.e. 0.0.1)")

	rootCmd.AddCommand(releasesCmd)
	releasesCmd.AddCommand(releasesTfeCmd)
	releasesTfeCmd.AddCommand(releasesTfeListCmd)
	releasesTfeCmd.AddCommand(releasesTfeShowCmd)
	releasesTfeCmd.AddCommand(releasesTfeDownloadCmd)

	releasesCmd.AddCommand(releasesReplicatedCmd)
	releasesReplicatedCmd.AddCommand(releasesReplicatedListCmd)
	releasesReplicatedCmd.AddCommand(releasesReplicatedDownloadCmd)
}

func releasesTfeList() error {
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
	for index, i := range tfeBinaries.Releases {
		t.AppendRow(table.Row{i.ReleaseSequence, i.Label, i.Required, i.ReleaseDate.String()})
		if index+1 >= maxResults {
			break
		}
	}
	t.SetStyle(table.StyleRounded)
	t.Render()

	return nil
}

func releasesTfeShow() error {
	// Validate flags
	licenseId := *viperString("licenseId")
	password := *viperString("password")
	release := *viperInt("release")

	tfeBinaries, err := ListTFEBinaries(password, licenseId)
	if err != nil {
		return err
	}

	idx := slices.IndexFunc(tfeBinaries.Releases, func(c TFERelease) bool { return c.ReleaseSequence == release })
	if idx < 0 {
		fmt.Println(color.RedString("Error: "), "Unable to find release sequence: ", release)
		return nil
	}
	tfeRelease := tfeBinaries.Releases[idx]

	fmt.Println(" Found")
	fmt.Println(color.BlueString("Release Sequence: "), tfeRelease.ReleaseSequence)
	fmt.Println(color.BlueString("Label:            "), tfeRelease.Label)
	fmt.Println(color.BlueString("Release Date:     "), tfeRelease.ReleaseDate)
	fmt.Println(color.BlueString("Required:         "), tfeRelease.Required)
	fmt.Println(color.BlueString("Release Notes:    "))
	fmt.Println(tfeRelease.ReleaseNotes)

	return nil
}

func releasesTfeDownload() error {
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
	path := fmt.Sprintf("%stfe-%s.releases", directory, release)

	//Download file
	err = DownloadBinary(tfeUrl.URL, path)
	if err != nil {
		return err
	}

	return nil
}

func releasesReplicatedList() error {
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

func releasesReplicatedDownload() error {
	// Validate flags
	directory := *viperString("directory")
	// Attempt to prevent a non semantic version from being requested
	version, err := viperSemanticVersionString("version")
	if err != nil {
		logError(err, "failed to parse semantic version")
	}

	// Get url - escape "%2B" as "%%2B", + symbol
	url := fmt.Sprintf("https://s3.amazonaws.com/replicated-releases-work/stable/replicated-%s%%2B%s%%2B%s.tar.gz",
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
