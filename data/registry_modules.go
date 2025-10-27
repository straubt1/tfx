// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	tfe "github.com/hashicorp/go-tfe"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/output"
)

// ListRegistryModules lists modules in the private registry for an organization, up to maxItems
func ListRegistryModules(c *client.TfxClient, orgName string, maxItems int) ([]*tfe.RegistryModule, error) {
	output.Get().Logger().Debug("Listing registry modules", "org", orgName, "maxItems", maxItems)

	items, err := client.FetchAll(c.Context, func(pageNumber int) ([]*tfe.RegistryModule, *client.Pagination, error) {
		opts := &tfe.RegistryModuleListOptions{ListOptions: tfe.ListOptions{PageNumber: pageNumber, PageSize: 100}}
		res, err := c.Client.RegistryModules.List(c.Context, orgName, opts)
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

// CreateRegistryModule creates a new module in the private registry
func CreateRegistryModule(c *client.TfxClient, orgName, name, provider string) (*tfe.RegistryModule, error) {
	output.Get().Logger().Debug("Creating registry module", "org", orgName, "name", name, "provider", provider)
	return c.Client.RegistryModules.Create(c.Context, orgName, tfe.RegistryModuleCreateOptions{
		Name:     &name,
		Provider: &provider,
	})
}

// ReadRegistryModule reads a module by name and provider
func ReadRegistryModule(c *client.TfxClient, orgName, name, provider string) (*tfe.RegistryModule, error) {
	output.Get().Logger().Debug("Reading registry module", "org", orgName, "name", name, "provider", provider)
	return c.Client.RegistryModules.Read(c.Context, tfe.RegistryModuleID{
		Organization: orgName,
		Name:         name,
		Provider:     provider,
		Namespace:    orgName,
		RegistryName: tfe.PrivateRegistry,
	})
}

// DeleteRegistryModule deletes a module by name and provider
func DeleteRegistryModule(c *client.TfxClient, orgName, name, provider string) error {
	output.Get().Logger().Debug("Deleting registry module", "org", orgName, "name", name, "provider", provider)
	return c.Client.RegistryModules.DeleteProvider(c.Context, tfe.RegistryModuleID{
		Organization: orgName,
		Name:         name,
		Provider:     provider,
		RegistryName: tfe.PrivateRegistry,
	})
}

// ListRegistryModuleVersions returns the module and its VersionStatuses
func ListRegistryModuleVersions(c *client.TfxClient, orgName, name, provider string) (*tfe.RegistryModule, error) {
	output.Get().Logger().Debug("Listing module versions", "org", orgName, "name", name, "provider", provider)
	return ReadRegistryModule(c, orgName, name, provider)
}

// CreateRegistryModuleVersion creates a new module version and uploads the directory
func CreateRegistryModuleVersion(c *client.TfxClient, orgName, name, provider, version, directory string) (*tfe.RegistryModuleVersion, error) {
	output.Get().Logger().Debug("Creating module version", "org", orgName, "name", name, "provider", provider, "version", version)
	mv, err := c.Client.RegistryModules.CreateVersion(c.Context, tfe.RegistryModuleID{
		Organization: orgName,
		Name:         name,
		Provider:     provider,
		Namespace:    orgName,
		RegistryName: tfe.PrivateRegistry,
	}, tfe.RegistryModuleCreateVersionOptions{Version: &version})
	if err != nil {
		return nil, err
	}
	if err := c.Client.RegistryModules.Upload(c.Context, *mv, directory); err != nil {
		return nil, err
	}
	return mv, nil
}

// DeleteRegistryModuleVersion deletes a module version
func DeleteRegistryModuleVersion(c *client.TfxClient, orgName, name, provider, version string) error {
	output.Get().Logger().Debug("Deleting module version", "org", orgName, "name", name, "provider", provider, "version", version)
	return c.Client.RegistryModules.DeleteVersion(c.Context, tfe.RegistryModuleID{
		Organization: orgName,
		Name:         name,
		Provider:     provider,
		Namespace:    orgName,
		RegistryName: tfe.PrivateRegistry,
	}, version)
}
