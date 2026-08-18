// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	"os"
	"strings"

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

// EnsureGPGKey returns the private-registry GPG key for orgName, creating it from
// publicKeyPath when missing. orgName is the private-registry GPG namespace (the
// TFE/HCP organization), not the public-registry publisher.
func EnsureGPGKey(c *client.TfxClient, orgName, keyID, publicKeyPath string) (*tfe.GPGKey, bool, error) {
	output.Get().Logger().Debug("Ensuring GPG key", "namespace", orgName, "keyID", keyID)

	gpgKeyID := tfe.GPGKeyID{
		RegistryName: tfe.PrivateRegistry,
		Namespace:    orgName,
		KeyID:        keyID,
	}
	key, err := c.Client.GPGKeys.Read(c.Context, gpgKeyID)
	if err == nil {
		if err := requireGPGKeyID(key, keyID); err != nil {
			return nil, false, err
		}
		return key, false, nil
	}
	if !isNotFound(err) {
		return nil, false, errors.Wrap(err, "failed to read GPG key")
	}
	if publicKeyPath == "" {
		return nil, false, errors.Errorf("GPG key %s is not in the private registry and no .asc public key was found in the staged directory", keyID)
	}

	output.Get().Logger().Info("GPG key not found, creating", "namespace", orgName, "keyID", keyID)
	key, err = CreateGPGKey(c, tfe.PrivateRegistry, orgName, publicKeyPath)
	if err != nil {
		return nil, false, err
	}
	if err := requireGPGKeyID(key, keyID); err != nil {
		return nil, false, err
	}
	return key, true, nil
}

func requireGPGKeyID(key *tfe.GPGKey, want string) error {
	if key == nil {
		return errors.New("GPG key is nil")
	}
	got := strings.TrimSpace(key.KeyID)
	want = strings.TrimSpace(want)
	if !strings.EqualFold(got, want) {
		return errors.Errorf("GPG key id mismatch: registry has %s, expected %s", got, want)
	}
	return nil
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
