// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package client

import (
	"context"
	"fmt"
	"net/http"

	tfe "github.com/hashicorp/go-tfe"
)

// TfxClient encapsulates the TFE client and context for API operations
type TfxClient struct {
	Client           *tfe.Client
	Context          context.Context
	Hostname         string
	Token            string
	OrganizationName string
	EventBus         *APIEventBus // non-nil in TUI mode; publishes HTTP events to the inspector panel
}

// New creates a new TFE client with the provided configuration
func New(hostname, token, organization string) (*TfxClient, error) {
	return NewWithContext(context.Background(), hostname, token, organization)
}

// NewWithContext creates a new TFE client with a parent context.
// HTTP logging is enabled when TFX_LOG or TFX_LOG_PATH is set.
// To also attach an APIEventBus (TUI inspector), use NewWithContextAndBus.
func NewWithContext(ctx context.Context, hostname, token, organization string) (*TfxClient, error) {
	return NewWithContextAndBus(ctx, hostname, token, organization, nil)
}

// NewWithContextAndBus creates a new TFE client with optional event bus support.
// When bus is non-nil a LoggingTransport is always installed (regardless of TFX_LOG)
// so the TUI inspector panel receives every API call.
func NewWithContextAndBus(ctx context.Context, hostname, token, organization string, bus *APIEventBus) (*TfxClient, error) {
	if hostname == "" {
		return nil, fmt.Errorf("hostname is required")
	}
	if token == "" {
		return nil, fmt.Errorf("token is required")
	}

	var config *tfe.Config

	// Install the logging transport when TFX_LOG/TFX_LOG_PATH is set OR when an
	// event bus is provided (TUI mode).  A single transport handles both channels.
	if IsTFXLogEnabled() || bus != nil {
		httpClient, _, err := NewHTTPClientWithLogging(bus)
		if err != nil {
			return nil, fmt.Errorf("failed to create HTTP client with logging: %w", err)
		}
		config = &tfe.Config{
			Address:    fmt.Sprintf("https://%s", hostname),
			Token:      token,
			HTTPClient: httpClient,
		}
	} else {
		// No logging and no event bus — use the default HTTP client.
		config = &tfe.Config{
			Address: fmt.Sprintf("https://%s", hostname),
			Token:   token,
		}
	}

	client, err := tfe.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create TFE client: %w", err)
	}

	return &TfxClient{
		Client:           client,
		Context:          ctx,
		Hostname:         hostname,
		Token:            token,
		OrganizationName: organization,
		EventBus:         bus,
	}, nil
}

// newHTTPClientNoLogging returns a plain http.Client with no transport wrapping.
// Used when neither logging nor event bus is needed.
func newHTTPClientNoLogging() *http.Client {
	return &http.Client{}
}
