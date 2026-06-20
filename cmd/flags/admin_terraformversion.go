// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package flags

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AdminTerraformVersionListFlags holds all flags for the admin terraform-version list command
type AdminTerraformVersionListFlags struct {
	Search string
}

// AdminTerraformVersionShowFlags holds all flags for the admin terraform-version show command
type AdminTerraformVersionShowFlags struct {
	Version string
}

// AdminTerraformVersionCreateFlags holds all flags for the admin terraform-version create command
type AdminTerraformVersionCreateFlags struct {
	Version  string
	URL      string
	SHA      string
	Official bool
	Enabled  bool
	Beta     bool
}

// AdminTerraformVersionCreateOfficialFlags holds all flags for the admin terraform-version create official command
type AdminTerraformVersionCreateOfficialFlags struct {
	Version string
	Enabled bool
	Beta    bool
}

// AdminTerraformVersionDeleteFlags holds all flags for the admin terraform-version delete command
type AdminTerraformVersionDeleteFlags struct {
	Version string
}

// AdminTerraformVersionEnableDisableFlags holds all flags for the admin terraform-version enable/disable commands
type AdminTerraformVersionEnableDisableFlags struct {
	Versions []string
	All      bool
}

// AdminTerraformVersionDisableAllFlags holds all flags for the admin terraform-version disable all command
type AdminTerraformVersionDisableAllFlags struct {
	Except     []string
	Before     string
	NotInUse   bool
	Deprecated bool
	Unofficial bool
	Official   bool
	Beta       bool
}

// AdminTerraformVersionEnableAllFlags holds all flags for the admin terraform-version enable all command
type AdminTerraformVersionEnableAllFlags struct {
	Include    []string
	Except     []string
	Unofficial bool
	Official   bool
	Beta       bool
}

// ParseAdminTerraformVersionListFlags creates AdminTerraformVersionListFlags from the current command context
func ParseAdminTerraformVersionListFlags(cmd *cobra.Command) (*AdminTerraformVersionListFlags, error) {
	return &AdminTerraformVersionListFlags{
		Search: viper.GetString("search"),
	}, nil
}

// ParseAdminTerraformVersionShowFlags creates AdminTerraformVersionShowFlags from the current command context
func ParseAdminTerraformVersionShowFlags(cmd *cobra.Command) (*AdminTerraformVersionShowFlags, error) {
	version := viper.GetString("version")
	if err := validateSemanticVersion(version); err != nil {
		return nil, err
	}
	return &AdminTerraformVersionShowFlags{
		Version: version,
	}, nil
}

// ParseAdminTerraformVersionCreateFlags creates AdminTerraformVersionCreateFlags from the current command context
func ParseAdminTerraformVersionCreateFlags(cmd *cobra.Command) (*AdminTerraformVersionCreateFlags, error) {
	version := viper.GetString("version")
	if err := validateSemanticVersion(version); err != nil {
		return nil, err
	}

	sha := viper.GetString("sha")
	if err := validateSHA(sha); err != nil {
		return nil, err
	}

	// Note: disable flag is inverted to enabled
	return &AdminTerraformVersionCreateFlags{
		Version:  version,
		URL:      viper.GetString("url"),
		SHA:      sha,
		Official: viper.GetBool("official"),
		Enabled:  !viper.GetBool("disable"),
		Beta:     viper.GetBool("beta"),
	}, nil
}

// ParseAdminTerraformVersionCreateOfficialFlags creates AdminTerraformVersionCreateOfficialFlags from the current command context
func ParseAdminTerraformVersionCreateOfficialFlags(cmd *cobra.Command) (*AdminTerraformVersionCreateOfficialFlags, error) {
	version := viper.GetString("version")
	if err := validateSemanticVersion(version); err != nil {
		return nil, err
	}

	// Note: disable flag is inverted to enabled
	return &AdminTerraformVersionCreateOfficialFlags{
		Version: version,
		Enabled: !viper.GetBool("disable"),
		Beta:    viper.GetBool("beta"),
	}, nil
}

// ParseAdminTerraformVersionDeleteFlags creates AdminTerraformVersionDeleteFlags from the current command context
func ParseAdminTerraformVersionDeleteFlags(cmd *cobra.Command) (*AdminTerraformVersionDeleteFlags, error) {
	return &AdminTerraformVersionDeleteFlags{
		Version: viper.GetString("version"),
	}, nil
}

// ParseAdminTerraformVersionEnableDisableFlags creates AdminTerraformVersionEnableDisableFlags from the current command context
func ParseAdminTerraformVersionEnableDisableFlags(cmd *cobra.Command) (*AdminTerraformVersionEnableDisableFlags, error) {
	versions := viper.GetStringSlice("versions")
	all := cmd.CalledAs() == "all" // Check if the "all" subcommand was called

	return &AdminTerraformVersionEnableDisableFlags{
		Versions: versions,
		All:      all,
	}, nil
}

// ParseAdminTerraformVersionDisableAllFlags creates AdminTerraformVersionDisableAllFlags from the current command context
func ParseAdminTerraformVersionDisableAllFlags(cmd *cobra.Command) (*AdminTerraformVersionDisableAllFlags, error) {
	except := viper.GetStringSlice("except")
	before := viper.GetString("before")
	notInUse := viper.GetBool("not-in-use")
	deprecated := viper.GetBool("deprecated")
	unofficial := viper.GetBool("unofficial")
	official := viper.GetBool("official")
	beta := viper.GetBool("beta")

	filterCount := 0
	if len(except) > 0 {
		filterCount++
	}
	if before != "" {
		filterCount++
	}
	if notInUse {
		filterCount++
	}
	if deprecated {
		filterCount++
	}
	if unofficial {
		filterCount++
	}
	if official {
		filterCount++
	}
	if beta {
		filterCount++
	}

	if filterCount > 1 {
		return nil, errors.New("only one filter flag may be set: --except, --before, --not-in-use, --beta, --deprecated, --official, or --unofficial")
	}

	if len(except) > 0 {
		for _, v := range except {
			if err := validateSemanticVersion(v); err != nil {
				return nil, errors.Wrapf(err, "invalid version in --except: %s", v)
			}
		}
	}

	if before != "" {
		if err := validateSemanticVersion(before); err != nil {
			return nil, errors.Wrapf(err, "invalid version in --before: %s", before)
		}
	}

	return &AdminTerraformVersionDisableAllFlags{
		Except:     except,
		Before:     before,
		NotInUse:   notInUse,
		Deprecated: deprecated,
		Unofficial: unofficial,
		Official:   official,
		Beta:       beta,
	}, nil
}

// ParseAdminTerraformVersionEnableAllFlags creates AdminTerraformVersionEnableAllFlags from the current command context
func ParseAdminTerraformVersionEnableAllFlags(cmd *cobra.Command) (*AdminTerraformVersionEnableAllFlags, error) {
	include := viper.GetStringSlice("include")
	except := viper.GetStringSlice("except")
	unofficial := viper.GetBool("unofficial")
	official := viper.GetBool("official")
	beta := viper.GetBool("beta")

	filterCount := 0
	if len(include) > 0 {
		filterCount++
	}
	if len(except) > 0 {
		filterCount++
	}
	if unofficial {
		filterCount++
	}
	if official {
		filterCount++
	}
	if beta {
		filterCount++
	}

	if filterCount > 1 {
		return nil, errors.New("only one filter flag may be set: --include, --except, --beta, --official, or --unofficial")
	}

	for _, v := range include {
		if err := validateSemanticVersion(v); err != nil {
			return nil, errors.Wrapf(err, "invalid version in --include: %s", v)
		}
	}

	for _, v := range except {
		if err := validateSemanticVersion(v); err != nil {
			return nil, errors.Wrapf(err, "invalid version in --except: %s", v)
		}
	}

	return &AdminTerraformVersionEnableAllFlags{
		Include:    include,
		Except:     except,
		Unofficial: unofficial,
		Official:   official,
		Beta:       beta,
	}, nil
}

// validateSemanticVersion validates that a version string is in semantic version format
func validateSemanticVersion(version string) error {
	// Simple regex for semantic versioning (major.minor.patch with optional pre-release and metadata)
	semanticVersionRegex := regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	if !semanticVersionRegex.MatchString(version) {
		return errors.New("failed to parse semantic version")
	}
	return nil
}

// validateSHA validates that a SHA string is 64 characters (SHA-256)
func validateSHA(sha string) error {
	if len(sha) != 64 {
		return fmt.Errorf("SHA checksum must be 64 characters long, got %d", len(sha))
	}
	// Validate it's hexadecimal
	shaRegex := regexp.MustCompile(`^[a-fA-F0-9]+$`)
	if !shaRegex.MatchString(sha) {
		return errors.New("SHA checksum must be hexadecimal")
	}
	return nil
}
