//go:build ignore

// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	pkgfile "github.com/straubt1/tfx/pkg/file"
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
			m := *viperInt("max-items")
			if *viperBool("all") {
				m = math.MaxInt
			}

			return releaseTfeList(
				*viperString("license-id"),
				*viperString("password"),
				m)
		},
	}

	// `tfx release tfe show` commands
	releaseTfeShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show a TFE release",
		Long:  "Show a Terraform Enterprise release, including release notes.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return releaseTfeShow(
				*viperString("license-id"),
				*viperString("password"),
				*viperInt("release-sequence"))
		},
	}

	// `tfx release tfe download` commands
	releaseTfeDownloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download TFE release binary",
		Long:  "Download a Terraform Enterprise release binary.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !pkgfile.IsDirectory(*viperString("directory")) {
				return errors.New("directory file does not exist")
			}

			return releaseTfeDownload(
				*viperString("license-id"),
				*viperString("password"),
				*viperInt("release-sequence"),
				*viperString("directory"))
		},
	}
)

func init() {
	// `tfx release tfe list`
	releaseTfeListCmd.Flags().StringP("license-id", "l", "", "License Id for TFE/Replicated")
	releaseTfeListCmd.Flags().StringP("password", "p", "", "Password to authenticate")
	releaseTfeListCmd.Flags().IntP("max-items", "m", 10, "The number of results to print")
	releaseTfeListCmd.Flags().BoolP("all", "a", false, "Retrieve all results regardless of maxItems flag (optional)")
	releaseTfeListCmd.MarkFlagRequired("license-id")
	releaseTfeListCmd.MarkFlagRequired("password")

	// `tfx release tfe show`
	releaseTfeShowCmd.Flags().StringP("license-id", "l", "", "License Id for TFE/Replicated")
	releaseTfeShowCmd.Flags().StringP("password", "p", "", "Password to authenticate")
	releaseTfeShowCmd.Flags().IntP("release-sequence", "r", 0, "Release Sequence (i.e. 610, 619, etc...)")
	releaseTfeShowCmd.MarkFlagRequired("license-id")
	releaseTfeShowCmd.MarkFlagRequired("password")
	releaseTfeShowCmd.MarkFlagRequired("release-sequence")

	// `tfx release tfe download`
	releaseTfeDownloadCmd.Flags().StringP("license-id", "l", "", "License Id for TFE/Replicated")
	releaseTfeDownloadCmd.Flags().StringP("password", "p", "", "Password to authenticate to TFE/Replicated")
	releaseTfeDownloadCmd.Flags().IntP("release-sequence", "r", 0, "Release Sequence (i.e. 610, 619, etc...)")
	releaseTfeDownloadCmd.Flags().StringP("directory", "d", "./", "Directory to save binary")
	releaseTfeDownloadCmd.MarkFlagRequired("license-id")
	releaseTfeDownloadCmd.MarkFlagRequired("password")
	releaseTfeDownloadCmd.MarkFlagRequired("release-sequence")

	releaseCmd.AddCommand(releaseTfeCmd)
	releaseTfeCmd.AddCommand(releaseTfeListCmd)
	releaseTfeCmd.AddCommand(releaseTfeShowCmd)
	releaseTfeCmd.AddCommand(releaseTfeDownloadCmd)
}

func releaseTfeList(licenseId string, password string, maxItems int) error {
	o.AddMessageUserProvided("List Available Terraform Enterprise Releases", "")
	tfeBinaries, err := ListTFEBinaries(password, licenseId)
	if err != nil {
		return errors.Wrap(err, "failed to list TFE releases")
	}

	o.AddTableHeader("Sequence", "Label", "Required", "Release Date")
	for index, i := range tfeBinaries.Releases {
		o.AddTableRows(i.ReleaseSequence, i.Label, i.Required, FormatDateTime(i.ReleaseDate))
		if index+1 >= maxItems {
			break
		}
	}

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
