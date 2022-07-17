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
	"strings"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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
			return releaseReplicatedList(
				*viperInt("maxResults"))
		},
	}

	// `tfx release replicated download` commands
	releaseReplicatedDownloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download a Replicated release binary",
		Long:  "Download a Replicated release binary.",
		RunE: func(cmd *cobra.Command, args []string) error {
			version, err := viperSemanticVersionString("version")
			if err != nil {
				return errors.New("failed to parse semantic version")
			}
			if !isDirectory(*viperString("directory")) {
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
	releaseReplicatedListCmd.Flags().IntP("maxResults", "r", 10, "The number of results to print")

	// `tfx release replicated download`
	releaseReplicatedDownloadCmd.Flags().StringP("directory", "d", "./", "Directory to save binary")
	releaseReplicatedDownloadCmd.Flags().StringP("version", "v", "", "Version of Replicated to Download (i.e. 0.0.1)")

	releaseCmd.AddCommand(releaseReplicatedCmd)
	releaseReplicatedCmd.AddCommand(releaseReplicatedListCmd)
	releaseReplicatedCmd.AddCommand(releaseReplicatedDownloadCmd)
}

func releaseReplicatedList(maxResults int) error {
	o.AddMessageUserProvided("List Available Replicated Releases", "")
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://release-notes.replicated.com/index.xml")
	if err != nil {
		return errors.Wrap(err, "failed to build replicated url")
	}

	o.AddTableHeader("Version", "Published Date")
	for index, i := range feed.Items {
		o.AddTableRows(i.Title, i.Published)
		if index+1 >= maxResults {
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
