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
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var (
	// `tfx release tfe` commands
	releaseTfeCmd = &cobra.Command{
		Use:   "tfe",
		Short: "TFE release commands",
		Long:  "Work with releases for Terraform Enterprise.",
	}

	// `tfx release tfe list` commands
	releaseTfeListCmd = &cobra.Command{
		Use:   "list",
		Short: "List TFE Releases",
		Long:  "List available Terraform Enterprise releases.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return releaseTfeList(
				*viperString("licenseId"),
				*viperString("password"),
				*viperInt("maxResults"))
		},
	}

	// `tfx release tfe show` commands
	releaseTfeShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show a TFE release",
		Long:  "Show a Terraform Enterprise release, including release notes.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return releaseTfeShow(
				*viperString("licenseId"),
				*viperString("password"),
				*viperInt("releaseSequence"))
		},
	}

	// `tfx release tfe download` commands
	releaseTfeDownloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download TFE release binary",
		Long:  "Download a Terraform Enterprise release binary.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !isDirectory(*viperString("directory")) {
				return errors.New("directory file does not exist")
			}

			return releaseTfeDownload(
				*viperString("licenseId"),
				*viperString("password"),
				*viperInt("releaseSequence"),
				*viperString("directory"))
		},
	}
)

func init() {
	// `tfx release tfe list`
	releaseTfeListCmd.Flags().StringP("licenseId", "l", "", "License Id for TFE/Replicated")
	releaseTfeListCmd.Flags().StringP("password", "p", "", "Password to authenticate")
	releaseTfeListCmd.Flags().IntP("maxResults", "r", 10, "The number of results to print")
	releaseTfeListCmd.MarkFlagRequired("licenseId")
	releaseTfeListCmd.MarkFlagRequired("password")

	// `tfx release tfe show`
	releaseTfeShowCmd.Flags().StringP("licenseId", "l", "", "License Id for TFE/Replicated")
	releaseTfeShowCmd.Flags().StringP("password", "p", "", "Password to authenticate")
	releaseTfeShowCmd.Flags().IntP("releaseSequence", "r", 0, "Release Sequence (i.e. 610, 619, etc...)")
	releaseTfeShowCmd.MarkFlagRequired("licenseId")
	releaseTfeShowCmd.MarkFlagRequired("password")
	releaseTfeShowCmd.MarkFlagRequired("releaseSequence")

	// `tfx release tfe download`
	releaseTfeDownloadCmd.Flags().StringP("licenseId", "l", "", "License Id for TFE/Replicated")
	releaseTfeDownloadCmd.Flags().StringP("password", "p", "", "Password to authenticate to TFE/Replicated")
	releaseTfeDownloadCmd.Flags().IntP("releaseSequence", "r", 0, "Release Sequence (i.e. 610, 619, etc...)")
	releaseTfeDownloadCmd.Flags().StringP("directory", "d", "./", "Directory to save binary")
	releaseTfeDownloadCmd.MarkFlagRequired("licenseId")
	releaseTfeDownloadCmd.MarkFlagRequired("password")
	releaseTfeDownloadCmd.MarkFlagRequired("releaseSequence")

	releaseCmd.AddCommand(releaseTfeCmd)
	releaseTfeCmd.AddCommand(releaseTfeListCmd)
	releaseTfeCmd.AddCommand(releaseTfeShowCmd)
	releaseTfeCmd.AddCommand(releaseTfeDownloadCmd)
}

func releaseTfeList(licenseId string, password string, maxResults int) error {
	o.AddMessageUserProvided("List Available Terraform Enterprise Releases", "")
	tfeBinaries, err := ListTFEBinaries(password, licenseId)
	if err != nil {
		return errors.Wrap(err, "failed to list TFE releases")
	}

	o.AddTableHeader("Sequence", "Label", "Required", "Release Date")
	for index, i := range tfeBinaries.Releases {
		o.AddTableRows(i.ReleaseSequence, i.Label, i.Required, FormatDateTime(i.ReleaseDate))
		if index+1 >= maxResults {
			break
		}
	}
	o.Close()

	return nil
}

func releaseTfeShow(licenseId string, password string, releaseSequence int) error {
	o.AddMessageUserProvided("Show Release details for Terraform Enterprise:", strconv.Itoa(releaseSequence))
	tfeBinaries, err := ListTFEBinaries(password, licenseId)
	if err != nil {
		return errors.Wrap(err, "failed to list TFE releases")
	}

	idx := slices.IndexFunc(tfeBinaries.Releases, func(c TFERelease) bool { return c.ReleaseSequence == releaseSequence })
	if idx < 0 {
		return errors.Wrap(err, "unable to fine release sequence in available TFE releases")
	}
	tfeRelease := tfeBinaries.Releases[idx]

	o.AddDeferredMessageRead("Release Sequence", tfeRelease.ReleaseSequence)
	o.AddDeferredMessageRead("Label", tfeRelease.Label)
	o.AddDeferredMessageRead("Release Date", FormatDateTime(tfeRelease.ReleaseDate))
	o.AddDeferredMessageRead("Required", tfeRelease.Required)
	o.AddDeferredMessageRead("Release Notes", "\n"+tfeRelease.ReleaseNotes)
	o.Close()

	return nil
}

// TODO: fix JSON output
func releaseTfeDownload(licenseId string, password string, releaseSequence int, directory string) error {
	o.AddMessageUserProvided("Download Release binary for Terraform Enterprise:", strconv.Itoa(releaseSequence))
	tfeUrl, err := GetTFEBinary(password, licenseId, releaseSequence)
	if err != nil {
		return errors.Wrap(err, "failed to get TFE releases")
	}

	// Verify trailing slash, if not add it
	if !strings.HasSuffix(directory, "/") {
		directory += "/"
	}
	path := fmt.Sprintf("%stfe-%d.release", directory, releaseSequence)

	//Download file
	err = DownloadBinary(tfeUrl.URL, path)
	if err != nil {
		return errors.Wrap(err, "failed to list download TFE binary")
	}

	o.AddMessageUserProvided("Release Downloaded!", "")

	return nil
}
