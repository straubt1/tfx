// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RegistryProviderListFlags holds flags for provider list
type RegistryProviderListFlags struct {
	MaxItems int
	All      bool
}

// RegistryProviderCreateFlags holds flags for provider create
type RegistryProviderCreateFlags struct {
	Name string
}

// RegistryProviderShowFlags holds flags for provider show
type RegistryProviderShowFlags struct {
	Name string
}

// RegistryProviderDeleteFlags holds flags for provider delete
type RegistryProviderDeleteFlags struct {
	Name string
}

func ParseRegistryProviderListFlags(cmd *cobra.Command) (*RegistryProviderListFlags, error) {
	return &RegistryProviderListFlags{
		MaxItems: viper.GetInt("max-items"),
		All:      viper.GetBool("all"),
	}, nil
}

func ParseRegistryProviderCreateFlags(cmd *cobra.Command) (*RegistryProviderCreateFlags, error) {
	return &RegistryProviderCreateFlags{
		Name: viper.GetString("name"),
	}, nil
}

func ParseRegistryProviderShowFlags(cmd *cobra.Command) (*RegistryProviderShowFlags, error) {
	return &RegistryProviderShowFlags{
		Name: viper.GetString("name"),
	}, nil
}

func ParseRegistryProviderDeleteFlags(cmd *cobra.Command) (*RegistryProviderDeleteFlags, error) {
	return &RegistryProviderDeleteFlags{
		Name: viper.GetString("name"),
	}, nil
}

// DefaultPublicProviderDirectory is the staging base for public provider downloads.
const DefaultPublicProviderDirectory = "./providers"

// PublicRegistryHashiCorpNamespace is the public Terraform Registry namespace
// for official HashiCorp providers (download --namespace default).
const PublicRegistryHashiCorpNamespace = "hashicorp"

// DefaultPublicProviderPlatforms are the os_arch values downloaded when --platforms is omitted.
var DefaultPublicProviderPlatforms = []string{
	"linux_amd64",
	"darwin_arm64",
	"darwin_amd64",
	"windows_amd64",
}

// RegistryProviderDownloadFlags holds flags for public registry provider download
type RegistryProviderDownloadFlags struct {
	Namespace    string
	Name         string
	Version      string
	Directory    string
	Platforms    []string
	AllPlatforms bool
}

func ParseRegistryProviderDownloadFlags(cmd *cobra.Command) (*RegistryProviderDownloadFlags, error) {
	platforms := viper.GetStringSlice("platforms")
	if len(platforms) == 0 {
		platforms = append([]string{}, DefaultPublicProviderPlatforms...)
	}
	directory := DefaultPublicProviderDirectory
	if UserPassed(cmd, "directory") {
		if d, err := cmd.Flags().GetString("directory"); err == nil && d != "" {
			directory = d
		}
	}
	namespace := viper.GetString("namespace")
	if namespace == "" {
		namespace = PublicRegistryHashiCorpNamespace
	}
	return &RegistryProviderDownloadFlags{
		Namespace:    namespace,
		Name:         viper.GetString("name"),
		Version:      viper.GetString("version"),
		Directory:    directory,
		Platforms:    platforms,
		AllPlatforms: viper.GetBool("all-platforms"),
	}, nil
}
