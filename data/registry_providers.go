// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/output"
)

// ListRegistryProviders lists providers in the private registry for an organization, up to maxItems
func ListRegistryProviders(c *client.TfxClient, orgName string, maxItems int) ([]*tfe.RegistryProvider, error) {
	output.Get().Logger().Debug("Listing registry providers", "org", orgName, "maxItems", maxItems)

	items, err := client.FetchAll(c.Context, func(pageNumber int) ([]*tfe.RegistryProvider, *client.Pagination, error) {
		opts := &tfe.RegistryProviderListOptions{ListOptions: tfe.ListOptions{PageNumber: pageNumber, PageSize: 100}}
		res, err := c.Client.RegistryProviders.List(c.Context, orgName, opts)
		if err != nil {
			return nil, nil, err
		}
		return res.Items, client.NewPaginationFromTFE(res.Pagination), nil
	})
	if err != nil {
		return nil, err
	}

	if maxItems > 0 && len(items) > maxItems {
		items = items[:maxItems]
	}
	return items, nil
}

// CreateRegistryProvider creates a new provider in the private registry
func CreateRegistryProvider(c *client.TfxClient, orgName, name string) (*tfe.RegistryProvider, error) {
	output.Get().Logger().Debug("Creating registry provider", "org", orgName, "name", name)
	return c.Client.RegistryProviders.Create(c.Context, orgName, tfe.RegistryProviderCreateOptions{
		Name:         name,
		Namespace:    orgName,
		RegistryName: tfe.PrivateRegistry,
	})
}

// ReadRegistryProvider reads a provider
func ReadRegistryProvider(c *client.TfxClient, orgName, name string) (*tfe.RegistryProvider, error) {
	output.Get().Logger().Debug("Reading registry provider", "org", orgName, "name", name)
	return c.Client.RegistryProviders.Read(c.Context, tfe.RegistryProviderID{
		OrganizationName: orgName,
		Name:             name,
		Namespace:        orgName,
		RegistryName:     tfe.PrivateRegistry,
	}, &tfe.RegistryProviderReadOptions{Include: []tfe.RegistryProviderIncludeOps{}})
}

// DeleteRegistryProvider deletes a provider
func DeleteRegistryProvider(c *client.TfxClient, orgName, name string) error {
	output.Get().Logger().Debug("Deleting registry provider", "org", orgName, "name", name)
	return c.Client.RegistryProviders.Delete(c.Context, tfe.RegistryProviderID{
		OrganizationName: orgName,
		Name:             name,
		Namespace:        orgName,
		RegistryName:     tfe.PrivateRegistry,
	})
}

// ListRegistryProviderVersions lists versions for a provider
func ListRegistryProviderVersions(c *client.TfxClient, orgName, name string) ([]*tfe.RegistryProviderVersion, error) {
	output.Get().Logger().Debug("Listing provider versions", "org", orgName, "name", name)
	return client.FetchAll(c.Context, func(pageNumber int) ([]*tfe.RegistryProviderVersion, *client.Pagination, error) {
		opts := &tfe.RegistryProviderVersionListOptions{ListOptions: tfe.ListOptions{PageNumber: pageNumber, PageSize: 100}}
		res, err := c.Client.RegistryProviderVersions.List(c.Context, tfe.RegistryProviderID{
			OrganizationName: orgName,
			Namespace:        orgName,
			RegistryName:     tfe.PrivateRegistry,
			Name:             name,
		}, opts)
		if err != nil {
			return nil, nil, err
		}
		return res.Items, client.NewPaginationFromTFE(res.Pagination), nil
	})
}

// CreateRegistryProviderVersion creates a provider version
func CreateRegistryProviderVersion(c *client.TfxClient, orgName, name, version, keyID string) (*tfe.RegistryProviderVersion, error) {
	output.Get().Logger().Debug("Creating provider version", "org", orgName, "name", name, "version", version)
	return c.Client.RegistryProviderVersions.Create(c.Context, tfe.RegistryProviderID{
		OrganizationName: orgName,
		Namespace:        orgName,
		RegistryName:     tfe.PrivateRegistry,
		Name:             name,
	}, tfe.RegistryProviderVersionCreateOptions{
		Version: version,
		KeyID:   keyID,
	})
}

// ReadRegistryProviderVersion reads a provider version
func ReadRegistryProviderVersion(c *client.TfxClient, orgName, name, version string) (*tfe.RegistryProviderVersion, error) {
	output.Get().Logger().Debug("Reading provider version", "org", orgName, "name", name, "version", version)
	return c.Client.RegistryProviderVersions.Read(c.Context, tfe.RegistryProviderVersionID{
		RegistryProviderID: tfe.RegistryProviderID{
			OrganizationName: orgName,
			Namespace:        orgName,
			RegistryName:     tfe.PrivateRegistry,
			Name:             name,
		},
		Version: version,
	})
}

// DeleteRegistryProviderVersion deletes a provider version
func DeleteRegistryProviderVersion(c *client.TfxClient, orgName, name, version string) error {
	output.Get().Logger().Debug("Deleting provider version", "org", orgName, "name", name, "version", version)
	return c.Client.RegistryProviderVersions.Delete(c.Context, tfe.RegistryProviderVersionID{
		RegistryProviderID: tfe.RegistryProviderID{
			OrganizationName: orgName,
			Name:             name,
			Namespace:        orgName,
			RegistryName:     tfe.PrivateRegistry,
		},
		Version: version,
	})
}

// ListRegistryProviderPlatforms lists platforms for a provider version
func ListRegistryProviderPlatforms(c *client.TfxClient, orgName, name, version string) ([]*tfe.RegistryProviderPlatform, error) {
	output.Get().Logger().Debug("Listing provider platforms", "org", orgName, "name", name, "version", version)
	return client.FetchAll(c.Context, func(pageNumber int) ([]*tfe.RegistryProviderPlatform, *client.Pagination, error) {
		opts := &tfe.RegistryProviderPlatformListOptions{ListOptions: tfe.ListOptions{PageNumber: pageNumber, PageSize: 100}}
		res, err := c.Client.RegistryProviderPlatforms.List(c.Context, tfe.RegistryProviderVersionID{
			RegistryProviderID: tfe.RegistryProviderID{
				OrganizationName: orgName,
				Namespace:        orgName,
				RegistryName:     tfe.PrivateRegistry,
				Name:             name,
			},
			Version: version,
		}, opts)
		if err != nil {
			return nil, nil, err
		}
		return res.Items, client.NewPaginationFromTFE(res.Pagination), nil
	})
}

// CreateRegistryProviderPlatform creates a provider platform record
func CreateRegistryProviderPlatform(c *client.TfxClient, orgName, name, version, os, arch, shasum, filename string) (*tfe.RegistryProviderPlatform, error) {
	output.Get().Logger().Debug("Creating provider platform", "org", orgName, "name", name, "version", version, "os", os, "arch", arch)
	return c.Client.RegistryProviderPlatforms.Create(c.Context, tfe.RegistryProviderVersionID{
		RegistryProviderID: tfe.RegistryProviderID{
			OrganizationName: orgName,
			Namespace:        orgName,
			RegistryName:     tfe.PrivateRegistry,
			Name:             name,
		},
		Version: version,
	}, tfe.RegistryProviderPlatformCreateOptions{
		OS:       os,
		Arch:     arch,
		Shasum:   shasum,
		Filename: filename,
	})
}

// ReadRegistryProviderPlatform reads a provider platform
func ReadRegistryProviderPlatform(c *client.TfxClient, orgName, name, version, os, arch string) (*tfe.RegistryProviderPlatform, error) {
	output.Get().Logger().Debug("Reading provider platform", "org", orgName, "name", name, "version", version, "os", os, "arch", arch)
	return c.Client.RegistryProviderPlatforms.Read(c.Context, tfe.RegistryProviderPlatformID{
		RegistryProviderVersionID: tfe.RegistryProviderVersionID{
			RegistryProviderID: tfe.RegistryProviderID{
				OrganizationName: orgName,
				Namespace:        orgName,
				RegistryName:     tfe.PrivateRegistry,
				Name:             name,
			},
			Version: version,
		},
		OS:   os,
		Arch: arch,
	})
}

// DeleteRegistryProviderPlatform deletes a provider platform
func DeleteRegistryProviderPlatform(c *client.TfxClient, orgName, name, version, os, arch string) error {
	output.Get().Logger().Debug("Deleting provider platform", "org", orgName, "name", name, "version", version, "os", os, "arch", arch)
	return c.Client.RegistryProviderPlatforms.Delete(c.Context, tfe.RegistryProviderPlatformID{
		RegistryProviderVersionID: tfe.RegistryProviderVersionID{
			RegistryProviderID: tfe.RegistryProviderID{
				OrganizationName: orgName,
				Namespace:        orgName,
				RegistryName:     tfe.PrivateRegistry,
				Name:             name,
			},
			Version: version,
		},
		OS:   os,
		Arch: arch,
	})
}

// EnsureRegistryProvider returns the private-registry provider, creating it if missing.
func EnsureRegistryProvider(c *client.TfxClient, orgName, name string) (p *tfe.RegistryProvider, created bool, err error) {
	p, err = ReadRegistryProvider(c, orgName, name)
	if err == nil {
		return p, false, nil
	}
	if !isNotFound(err) {
		return nil, false, err
	}
	output.Get().Logger().Info("Provider not found, creating", "org", orgName, "name", name)
	p, err = CreateRegistryProvider(c, orgName, name)
	if err != nil {
		return nil, false, err
	}
	return p, true, nil
}

func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, tfe.ErrResourceNotFound) {
		return true
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "not found") || strings.Contains(s, "resource not found")
}

// UploadRegistryProviderPlatform hashes zipPath, creates the platform record, and PUTs the zip.
func UploadRegistryProviderPlatform(c *client.TfxClient, orgName, name, version, osName, arch, zipPath string) (*tfe.RegistryProviderPlatform, error) {
	output.Get().Logger().Debug("Uploading provider platform", "org", orgName, "name", name, "version", version, "os", osName, "arch", arch, "path", zipPath)

	sum, err := fileSHA256(zipPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to hash file")
	}

	filename := fmt.Sprintf("terraform-provider-%s_%s_%s_%s.zip", name, version, osName, arch)
	rpp, err := CreateRegistryProviderPlatform(c, orgName, name, version, osName, arch, sum, filename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create provider version platform")
	}

	uploadURL, ok := rpp.Links["provider-binary-upload"].(string)
	if !ok || uploadURL == "" {
		return nil, errors.New("provider platform is missing provider-binary-upload link")
	}
	if err := UploadBinary(uploadURL, zipPath); err != nil {
		return nil, errors.Wrap(err, "failed to upload binary to provider version platform")
	}
	return rpp, nil
}

// UploadRegistryProviderPlatforms uploads each staged zip concurrently, limited
// by concurrency. Results are in input order. After the first error, no new
// uploads are started; in-flight PUTs are allowed to finish.
func UploadRegistryProviderPlatforms(c *client.TfxClient, orgName, name, version string, platforms []StagedProviderPlatform, concurrency int) ([]*tfe.RegistryProviderPlatform, error) {
	return mapConcurrent(concurrency, platforms, func(plat StagedProviderPlatform) (*tfe.RegistryProviderPlatform, error) {
		rpp, err := UploadRegistryProviderPlatform(c, orgName, name, version, plat.OS, plat.Arch, plat.Path)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to upload platform %s_%s", plat.OS, plat.Arch)
		}
		return rpp, nil
	})
}
