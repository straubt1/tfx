// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package client

import (
	"context"

	"github.com/spf13/viper"
)

// NewFromViper creates a TfxClient using configuration from viper
func NewFromViper() (*TfxClient, error) {
	return NewFromViperWithContext(context.Background())
}

// NewFromViperWithContext creates a TfxClient using viper configuration with a parent context
func NewFromViperWithContext(ctx context.Context) (*TfxClient, error) {
	hostname := viper.GetString("hostname")
	token := viper.GetString("token")
	organization := viper.GetString("organization")

	return NewWithContext(ctx, hostname, token, organization)
}

// NewFromViperForTUI creates a TfxClient for TUI mode.
// It always installs a LoggingTransport that publishes HTTP events to bus,
// enabling the API Inspector panel regardless of whether TFX_LOG is set.
func NewFromViperForTUI(bus *APIEventBus) (*TfxClient, error) {
	hostname := viper.GetString("hostname")
	token := viper.GetString("token")
	organization := viper.GetString("organization")

	return NewWithContextAndBus(context.Background(), hostname, token, organization, bus)
}
