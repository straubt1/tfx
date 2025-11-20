//go:build ignore

// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	"fmt"
	"math"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	pkgfile "github.com/straubt1/tfx/pkg/file"
)

var (
	// `tfx release replicated` commands
	releaseReplicatedCmd = &cobra.Command{
		Use:   "replicated",
		Short: "Replicated release commands",
		Long:  "Work with releases for Replicated.",
	}

	// `tfx release replicated list` commands
	releaseReplicatedListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Replicated binaries",
		Long:  "List available Replicated releases.",
		RunE: func(cmd *cobra.Command, args []string) error {
			m := *viperInt("max-items")
			if *viperBool("all") {
				m = math.MaxInt
			}

			return releaseReplicatedList(
				m)
		},
	}

	// `tfx release replicated download` commands
	releaseReplicatedDownloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download a Replicated release binary",
		Long:  "Download a Replicated release binary.",
		RunE: func(cmd *cobra.Command, args []string) error {
			version := *viperString("version")
			if _, err := semver.NewVersion(version); err != nil {
				return errors.New("failed to parse semantic version")
			}
			if !pkgfile.IsDirectory(*viperString("directory")) {
				return errors.New("directory file does not exist")
			}

			return releaseReplicatedDownload(
				*viperString("directory"),
				version)
		},
	}
)

func init() {
	// `tfx release replicated list`
	releaseReplicatedListCmd.Flags().IntP("max-items", "m", 10, "The number of results to print")
	releaseReplicatedListCmd.Flags().BoolP("all", "a", false, "Retrieve all results regardless of maxItems flag (optional)")

	// `tfx release replicated download`
	releaseReplicatedDownloadCmd.Flags().StringP("directory", "d", "./", "Directory to save binary")
	releaseReplicatedDownloadCmd.Flags().StringP("version", "v", "", "Version of Replicated to Download (i.e. 0.0.1)")

	releaseCmd.AddCommand(releaseReplicatedCmd)
	releaseReplicatedCmd.AddCommand(releaseReplicatedListCmd)
	releaseReplicatedCmd.AddCommand(releaseReplicatedDownloadCmd)
}

func releaseReplicatedList(maxItems int) error {
	o.AddMessageUserProvided("List Available Replicated Releases", "")
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://release-notes.replicated.com/index.xml")
	if err != nil {
		return errors.Wrap(err, "failed to build replicated url")
	}

	o.AddTableHeader("Version", "Published Date")
	for index, i := range feed.Items {
		o.AddTableRows(i.Title, i.Published)
		if index+1 >= maxItems {
			break
		}
	}
	o.Close()

	return nil
}

// TODO: fix JSON output
func releaseReplicatedDownload(directory string, version string) error {
	o.AddMessageUserProvided("Download Release binary for Replicated:", version)
	// Get url - escape "%2B" as "%%2B", + symbol
	url := fmt.Sprintf("https://s3.amazonaws.com/replicated-airgap-work/stable/replicated-%s%%2B%s%%2B%s.tar.gz",
		version,
		version,
		version)

	// Verify trailing, if not add it
	if !strings.HasSuffix(directory, "/") {
		directory += "/"
	}
	path := fmt.Sprintf("%sreplicated-%s.targz", directory, version)

	//Download file
	err := DownloadBinary(url, path)
	if err != nil {
		return errors.Wrap(err, "failed to download replicated binary")
	}

	o.AddMessageUserProvided("Release Downloaded!", "")

	return nil
}
