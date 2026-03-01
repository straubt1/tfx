// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	"os"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/output"
)

// FetchGPGKeys fetches all GPG keys for a given registry namespace
// Note: Currently only supports private registry
func FetchGPGKeys(c *client.TfxClient, namespace string) ([]*tfe.GPGKey, error) {
	output.Get().Logger().Debug("Fetching GPG keys", "namespace", namespace)

	// Use ListPrivate which returns all keys for the namespace
	opts := tfe.GPGKeyListOptions{
		Namespaces: []string{namespace},
	}

	keys, err := c.Client.GPGKeys.ListPrivate(c.Context, opts)
	if err != nil {
		output.Get().Logger().Error("Failed to fetch GPG keys", "namespace", namespace, "error", err)
		return nil, errors.Wrap(err, "failed to list GPG keys")
	}

	output.Get().Logger().Debug("GPG keys fetched successfully", "namespace", namespace, "count", len(keys.Items))
	return keys.Items, nil
}

// FetchGPGKey fetches a single GPG key by ID
func FetchGPGKey(c *client.TfxClient, namespace string, registryName tfe.RegistryName, keyID string) (*tfe.GPGKey, error) {
	output.Get().Logger().Debug("Fetching GPG key", "namespace", namespace, "keyID", keyID)

	gpgKeyID := tfe.GPGKeyID{
		RegistryName: registryName,
		Namespace:    namespace,
		KeyID:        keyID,
	}

	key, err := c.Client.GPGKeys.Read(c.Context, gpgKeyID)
	if err != nil {
		output.Get().Logger().Error("Failed to fetch GPG key", "namespace", namespace, "keyID", keyID, "error", err)
		return nil, errors.Wrap(err, "failed to read GPG key")
	}

	output.Get().Logger().Debug("GPG key fetched successfully", "namespace", namespace, "keyID", keyID)
	return key, nil
}

// CreateGPGKey creates a new GPG key
func CreateGPGKey(c *client.TfxClient, registryName tfe.RegistryName, namespace string, publicKeyPath string) (*tfe.GPGKey, error) {
	output.Get().Logger().Debug("Creating GPG key", "namespace", namespace, "publicKeyPath", publicKeyPath)

	// Read the public key file
	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		output.Get().Logger().Error("Failed to read public key file", "path", publicKeyPath, "error", err)
		return nil, errors.Wrap(err, "failed to read public key file")
	}

	opts := tfe.GPGKeyCreateOptions{
		Namespace:  namespace,
		AsciiArmor: string(publicKeyBytes),
	}

	key, err := c.Client.GPGKeys.Create(c.Context, registryName, opts)
	if err != nil {
		output.Get().Logger().Error("Failed to create GPG key", "namespace", namespace, "error", err)
		return nil, errors.Wrap(err, "failed to create GPG key")
	}

	output.Get().Logger().Debug("GPG key created successfully", "namespace", namespace, "keyID", key.KeyID)
	return key, nil
}

// DeleteGPGKey deletes a GPG key
func DeleteGPGKey(c *client.TfxClient, namespace string, registryName tfe.RegistryName, keyID string) error {
	output.Get().Logger().Debug("Deleting GPG key", "namespace", namespace, "keyID", keyID)

	gpgKeyID := tfe.GPGKeyID{
		RegistryName: registryName,
		Namespace:    namespace,
		KeyID:        keyID,
	}

	err := c.Client.GPGKeys.Delete(c.Context, gpgKeyID)
	if err != nil {
		output.Get().Logger().Error("Failed to delete GPG key", "namespace", namespace, "keyID", keyID, "error", err)
		return errors.Wrap(err, "failed to delete GPG key")
	}

	output.Get().Logger().Debug("GPG key deleted successfully", "namespace", namespace, "keyID", keyID)
	return nil
}
