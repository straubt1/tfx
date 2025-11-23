// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/straubt1/tfx/cmd/flags"
	view "github.com/straubt1/tfx/cmd/views"
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
			cmdConfig, err := flags.ParseReleaseTfeListFlags(cmd)
			if err != nil {
				return err
			}
			return releaseTfeList(cmdConfig)
		},
	}

	// `tfx release tfe show` commands
	releaseTfeShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show a TFE release",
		Long:  "Show detailed information about a Terraform Enterprise release.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseReleaseTfeShowFlags(cmd)
			if err != nil {
				return err
			}
			return releaseTfeShow(cmdConfig)
		},
	}

	// `tfx release tfe download` commands
	releaseTfeDownloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download a TFE release image",
		Long:  "Download a Terraform Enterprise release image as a tar file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseReleaseTfeDownloadFlags(cmd)
			if err != nil {
				return err
			}
			return releaseTfeDownload(cmdConfig)
		},
	}

	// // `tfx release tfe download` commands
	// releaseTfeDownloadCmd = &cobra.Command{
	// 	Use:   "download",
	// 	Short: "Download TFE release binary",
	// 	Long:  "Download a Terraform Enterprise release binary.",
	// 	RunE: func(cmd *cobra.Command, args []string) error {
	// 		if !pkgfile.IsDirectory(*viperString("directory")) {
	// 			return errors.New("directory file does not exist")
	// 		}

	// 		return releaseTfeDownload(
	// 			*viperString("license-id"),
	// 			*viperString("password"),
	// 			*viperInt("release-sequence"),
	// 			*viperString("directory"))
	// 	},
	// }
)

func init() {
	// `tfx release tfe list`
	releaseTfeListCmd.Flags().StringP("tfe-license-path", "l", "", "Path to TFE license file")
	releaseTfeListCmd.Flags().StringP("registry-url", "r", "images.releases.hashicorp.com/hashicorp/terraform-enterprise", "Docker registry URL (optional)")
	releaseTfeListCmd.Flags().IntP("max-items", "m", 5, "The number of results to print")
	releaseTfeListCmd.Flags().BoolP("all", "a", false, "Retrieve all results regardless of maxItems flag (optional)")
	releaseTfeListCmd.Flags().BoolP("stable-only", "s", true, "Show only stable releases (GA and versioned formats) (optional)")
	releaseTfeListCmd.MarkFlagRequired("tfe-license-path")

	// `tfx release tfe show`
	releaseTfeShowCmd.Flags().StringP("tfe-license-path", "l", "", "Path to TFE license file")
	releaseTfeShowCmd.Flags().StringP("registry-url", "r", "images.releases.hashicorp.com/hashicorp/terraform-enterprise", "Docker registry URL (optional)")
	releaseTfeShowCmd.Flags().StringP("tag", "t", "", "TFE release tag to show")
	releaseTfeShowCmd.MarkFlagRequired("tfe-license-path")
	releaseTfeShowCmd.MarkFlagRequired("tag")

	// `tfx release tfe download`
	releaseTfeDownloadCmd.Flags().StringP("tfe-license-path", "l", "", "Path to TFE license file")
	releaseTfeDownloadCmd.Flags().StringP("registry-url", "r", "images.releases.hashicorp.com/hashicorp/terraform-enterprise", "Docker registry URL (optional)")
	releaseTfeDownloadCmd.Flags().StringP("tag", "t", "", "TFE release tag to download")
	releaseTfeDownloadCmd.Flags().StringP("output", "o", "", "Output file path (default: terraform-enterprise-{tag}.tar)")
	releaseTfeDownloadCmd.MarkFlagRequired("tfe-license-path")
	releaseTfeDownloadCmd.MarkFlagRequired("tag")

	releaseCmd.AddCommand(releaseTfeCmd)
	releaseTfeCmd.AddCommand(releaseTfeListCmd)
	releaseTfeCmd.AddCommand(releaseTfeShowCmd)
	releaseTfeCmd.AddCommand(releaseTfeDownloadCmd)
}

// filterStableTags filters tags to only include stable releases:
// - Semantic versioning (e.g., 1.0.0, 1.1.1)
// - Older version format (e.g., v202504-1, v202312-1)
// All other tags are filtered out as intermediate builds
func filterStableTags(tags []string) []string {
	// Regex for semantic versioning (e.g., 1.0.0, 1.1.1)
	semanticVersionRegex := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	// Regex for versioned format (e.g., v202504-1, v202312-1)
	versionedFormatRegex := regexp.MustCompile(`^v\d{6}-\d+$`)

	var filtered []string
	for _, tag := range tags {
		if semanticVersionRegex.MatchString(tag) || versionedFormatRegex.MatchString(tag) {
			filtered = append(filtered, tag)
		}
	}
	return filtered
}

// sortReleases sorts tags with latest releases first
// Sorts by semantic versions first, then by versioned format (by date descending)
func sortReleases(tags []string) {
	semanticVersionRegex := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	versionedFormatRegex := regexp.MustCompile(`^v(\d{6})-(\d+)$`)

	sort.SliceStable(tags, func(i, j int) bool {
		tagI := tags[i]
		tagJ := tags[j]

		iIsSemantic := semanticVersionRegex.MatchString(tagI)
		jIsSemantic := semanticVersionRegex.MatchString(tagJ)

		// If both are semantic versions, compare them numerically
		if iIsSemantic && jIsSemantic {
			iParts := strings.Split(tagI, ".")
			jParts := strings.Split(tagJ, ".")

			for k := 0; k < 3 && k < len(iParts) && k < len(jParts); k++ {
				iNum, _ := strconv.Atoi(iParts[k])
				jNum, _ := strconv.Atoi(jParts[k])
				if iNum != jNum {
					return iNum > jNum // Higher version first
				}
			}
			return false
		}

		// If both are versioned format, compare by date and sequence
		iMatch := versionedFormatRegex.FindStringSubmatch(tagI)
		jMatch := versionedFormatRegex.FindStringSubmatch(tagJ)

		if len(iMatch) > 0 && len(jMatch) > 0 {
			iDate := iMatch[1]
			jDate := jMatch[1]

			if iDate != jDate {
				return iDate > jDate // Later date first
			}

			// If dates are the same, compare sequence numbers
			iSeq, _ := strconv.Atoi(iMatch[2])
			jSeq, _ := strconv.Atoi(jMatch[2])
			return iSeq > jSeq // Higher sequence first
		}

		// Semantic versions come before versioned format
		if iIsSemantic {
			return true
		}
		return false
	})
}

// getImageCreatedDate fetches the creation date for a tag from the registry
func getImageCreatedDate(tag string, ref name.Repository, auth authn.Authenticator) (string, error) {
	// Build the full image reference with the tag
	imageRef, err := name.NewTag(ref.String() + ":" + tag)
	if err != nil {
		return "", err
	}

	// Fetch the image to get metadata
	img, err := remote.Image(imageRef, remote.WithAuth(auth))
	if err != nil {
		return "", err
	}

	// Get the config file which contains the creation date
	configFile, err := img.ConfigFile()
	if err != nil {
		return "", err
	}

	// Format the created time as a readable string
	if !configFile.Created.IsZero() {
		return configFile.Created.Time.Format("01-02-2006"), nil
	}

	return "N/A", nil
}

func releaseTfeList(cmdConfig *flags.ReleaseTfeListFlags) error {
	// Create view for rendering
	v := view.NewReleaseTfeListView()

	v.PrintCommandHeader("Listing Terraform Enterprise releases")

	// Read the license file
	licenseData, err := os.ReadFile(cmdConfig.TfeLicensePath)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to read TFE license file"))
	}

	// Create authentication from license file
	auth := &authn.Basic{
		Username: "terraform",
		Password: string(licenseData),
	}

	// Parse the registry URL
	ref, err := name.NewRepository(cmdConfig.RegistryURL)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "invalid registry URL"))
	}

	// List tags from the registry
	tags, err := remote.List(ref, remote.WithAuth(auth))
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to list TFE releases from registry"))
	}

	// Build releases data for the view
	releases := make([]map[string]interface{}, 0)
	maxItems := cmdConfig.MaxItems
	if cmdConfig.All {
		maxItems = len(tags)
	}

	// Filter tags if stable-only is enabled
	var filteredTags []string
	if cmdConfig.StableOnly {
		filteredTags = filterStableTags(tags)
	} else {
		filteredTags = tags
	}

	// Sort filtered tags with latest first
	sortReleases(filteredTags)

	// Limit filtered tags based on maxItems
	if len(filteredTags) > maxItems && !cmdConfig.All {
		filteredTags = filteredTags[:maxItems]
	}

	// Fetch creation date for each tag that will be displayed
	for _, tag := range filteredTags {
		createdDate, err := getImageCreatedDate(tag, ref, auth)
		if err != nil {
			// If we can't fetch the date, use N/A
			createdDate = "N/A"
		}

		releases = append(releases, map[string]interface{}{
			"Tag":     tag,
			"Created": createdDate,
		})
	}

	return v.Render(releases)
}

func releaseTfeShow(cmdConfig *flags.ReleaseTfeShowFlags) error {
	// Create view for rendering
	v := view.NewReleaseTfeShowView()

	v.PrintCommandHeader("Terraform Enterprise Release Details")

	// Read the license file
	licenseData, err := os.ReadFile(cmdConfig.TfeLicensePath)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to read TFE license file"))
	}

	// Create authentication from license file
	auth := &authn.Basic{
		Username: "terraform",
		Password: string(licenseData),
	}

	// Parse the registry URL and create the full image reference
	ref, err := name.NewRepository(cmdConfig.RegistryURL)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "invalid registry URL"))
	}

	imageRef, err := name.NewTag(ref.String() + ":" + cmdConfig.Tag)
	if err != nil {
		return v.RenderError(errors.Wrap(err, "invalid tag"))
	}

	// Fetch the image
	img, err := remote.Image(imageRef, remote.WithAuth(auth))
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to fetch release image"))
	}

	// Get image digest
	digest, err := img.Digest()
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to get image digest"))
	}

	// Get config file with metadata
	configFile, err := img.ConfigFile()
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to get image config"))
	}

	// Get layers
	layers, err := img.Layers()
	if err != nil {
		return v.RenderError(errors.Wrap(err, "failed to get layers"))
	}

	// Calculate total image size by summing all layer sizes
	var totalSize int64
	for _, layer := range layers {
		layerSize, err := layer.Size()
		if err != nil {
			continue
		}
		totalSize += layerSize
	}

	// Build release data
	releaseData := map[string]interface{}{
		"Tag":     cmdConfig.Tag,
		"Digest":  digest.String(),
		"Size":    totalSize,
		"Created": configFile.Created.Time.Format("01-02-2006"),
		"OS":      configFile.OS,
		"Arch":    configFile.Architecture,
		"Layers":  layers,
		"Config":  configFile.Config,
		"History": configFile.History,
	}

	return v.Render(releaseData)
}

// TODO: Fix download spinner
func releaseTfeDownload(cmdConfig *flags.ReleaseTfeDownloadFlags) error {
	// Determine output filename
	outputPath := cmdConfig.Output
	if outputPath == "" {
		outputPath = "terraform-enterprise-" + cmdConfig.Tag + ".tar"
	}

	// Read the license file
	licenseData, err := os.ReadFile(cmdConfig.TfeLicensePath)
	if err != nil {
		return errors.Wrap(err, "failed to read TFE license file")
	}

	// Create authentication from license file
	auth := &authn.Basic{
		Username: "terraform",
		Password: string(licenseData),
	}

	// Parse the registry URL and create the full image reference
	ref, err := name.NewRepository(cmdConfig.RegistryURL)
	if err != nil {
		return errors.Wrap(err, "invalid registry URL")
	}

	imageRef, err := name.NewTag(ref.String() + ":" + cmdConfig.Tag)
	if err != nil {
		return errors.Wrap(err, "invalid tag")
	}

	fmt.Println("Downloading image:", imageRef.String())

	// Fetch the image
	img, err := remote.Image(imageRef, remote.WithAuth(auth))
	if err != nil {
		return errors.Wrap(err, "failed to fetch release image")
	}

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return errors.Wrap(err, "failed to create output file")
	}
	defer outputFile.Close()

	fmt.Println("Writing to:", outputPath)

	// Write image to tar file
	if err := tarball.Write(imageRef, img, outputFile); err != nil {
		os.Remove(outputPath) // Clean up on error
		return errors.Wrap(err, "failed to write image to tar file")
	}

	fmt.Println("Download complete:", outputPath)
	return nil
}
