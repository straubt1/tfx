/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
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

// airgapCmd represents the airgap command
var (
	airgapCmd = &cobra.Command{
		Use:   "airgap",
		Short: "Airgap helper commands",
		Long:  "Work with binaries needed for TFE airgap installations.",
	}

	airgapTfeCmd = &cobra.Command{
		Use:   "tfe",
		Short: "TFE airgap commands",
		Long:  "TFE airgap commands to work with TFE binaries.",
	}

	airgapTfeListCmd = &cobra.Command{
		Use:   "list",
		Short: "List TFE releases",
		Long:  "List available TFE releases for airgap download.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return airgapTfeList()
		},
		PreRun: bindPFlags,
	}

	airgapTfeShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show TFE release",
		Long:  "Show a TFE release, including release notes.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return airgapTfeShow()
		},
		PreRun: bindPFlags,
	}

	airgapTfeDownloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download TFE release binary",
		Long:  "Download a TFE release binary.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return airgapTfeDownload()
		},
		PreRun: bindPFlags,
	}

	airgapReplicatedCmd = &cobra.Command{
		Use:   "replicated",
		Short: "Replicated airgap commands",
		Long:  "Replicated airgap commands to work with Replicated binaries.",
	}

	airgapReplicatedListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Replicated binaries",
		Long:  "List available Replicated binaries for download.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return airgapReplicatedList()
		},
		PreRun: bindPFlags,
	}

	airgapReplicatedDownloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download Replicated release binary",
		Long:  "Download a Replicated release binary.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return airgapReplicatedDownload()
		},
		PreRun: bindPFlags,
	}
)

func init() {
	// `tfx airgap tfe list`
	airgapTfeListCmd.Flags().StringP("licenseId", "l", "", "License Id for TFE/Replicated")
	airgapTfeListCmd.Flags().StringP("password", "p", "", "Password to authenticate")
	airgapTfeListCmd.MarkFlagRequired("licenseId")
	airgapTfeListCmd.MarkFlagRequired("password")

	// `tfx airgap tfe show`
	airgapTfeShowCmd.Flags().StringP("licenseId", "l", "", "License Id for TFE/Replicated")
	airgapTfeShowCmd.Flags().StringP("password", "p", "", "Password to authenticate")
	airgapTfeShowCmd.Flags().StringP("release", "r", "", "Release Sequence (i.e. 610, 619, etc...)")
	airgapTfeShowCmd.MarkFlagRequired("licenseId")
	airgapTfeShowCmd.MarkFlagRequired("password")
	airgapTfeShowCmd.MarkFlagRequired("release")

	// `tfx airgap tfe download`
	airgapTfeDownloadCmd.Flags().StringP("licenseId", "l", "", "License Id for TFE/Replicated")
	airgapTfeDownloadCmd.Flags().StringP("password", "p", "", "Password to authenticate")
	airgapTfeDownloadCmd.Flags().StringP("release", "r", "", "Release Sequence (i.e. 610, 619, etc...)")
	airgapTfeDownloadCmd.Flags().StringP("directory", "d", "./", "Directory to save binary (optional, defaults to current directory)")
	airgapTfeDownloadCmd.MarkFlagRequired("licenseId")
	airgapTfeDownloadCmd.MarkFlagRequired("password")
	airgapTfeDownloadCmd.MarkFlagRequired("release")

	// `tfx airgap replicated list`
	airgapReplicatedListCmd.Flags().IntP("maxResults", "r", 10, "The number of results to print (optional, defaults to 10)")

	// `tfx airgap replicated download`
	airgapReplicatedDownloadCmd.Flags().StringP("directory", "d", "./", "Directory to save binary (optional, defaults to current directory)")
	airgapReplicatedDownloadCmd.Flags().StringP("version", "v", "", "Version of Replicated to Download (i.e. 0.0.1)")

	rootCmd.AddCommand(airgapCmd)
	airgapCmd.AddCommand(airgapTfeCmd)
	airgapTfeCmd.AddCommand(airgapTfeListCmd)
	airgapTfeCmd.AddCommand(airgapTfeShowCmd)
	airgapTfeCmd.AddCommand(airgapTfeDownloadCmd)

	airgapCmd.AddCommand(airgapReplicatedCmd)
	airgapReplicatedCmd.AddCommand(airgapReplicatedListCmd)
	airgapReplicatedCmd.AddCommand(airgapReplicatedDownloadCmd)
}

func airgapTfeList() error {
	// Validate flags
	licenseId := *viperString("licenseId")
	password := *viperString("password")

	tfeBinaries, err := ListTFEBinaries(password, licenseId)
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Sequence", "Label", "Required", "Release Date"})
	for _, i := range tfeBinaries.Releases {
		t.AppendRow(table.Row{i.ReleaseSequence, i.Label, i.Required, i.ReleaseDate.String()})
	}
	t.SetStyle(table.StyleRounded)
	t.Render()

	return nil
}

func airgapTfeShow() error {
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

func airgapTfeDownload() error {
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
	path := fmt.Sprintf("%stfe-%s.airgap", directory, release)

	//Download file
	err = DownloadBinary(tfeUrl.URL, path)
	if err != nil {
		return err
	}

	return nil
}

func airgapReplicatedList() error {
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
		if index >= maxResults {
			break
		}
	}
	t.SetStyle(table.StyleRounded)
	t.Render()

	return nil
}

func airgapReplicatedDownload() error {
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
