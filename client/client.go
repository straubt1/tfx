package client

import (
	"context"
	"fmt"

	tfe "github.com/hashicorp/go-tfe"
)

// TfxClient encapsulates the TFE client and context for API operations
type TfxClient struct {
	Client           *tfe.Client
	Context          context.Context
	Hostname         string
	OrganizationName string
}

// New creates a new TFE client with the provided configuration
func New(hostname, token, organization string) (*TfxClient, error) {
	return NewWithContext(context.Background(), hostname, token, organization, "")
}

// NewWithContext creates a new TFE client with a parent context
// If logFile is provided (non-empty), HTTP request/response logging will be enabled
func NewWithContext(ctx context.Context, hostname, token, organization, logFile string) (*TfxClient, error) {
	if hostname == "" {
		return nil, fmt.Errorf("hostname is required")
	}
	if token == "" {
		return nil, fmt.Errorf("token is required")
	}

	var config *tfe.Config

	// Conditionally enable HTTP logging if logFile is provided
	if logFile != "" {
		httpClient, _, err := NewHTTPClientWithLogging(logFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create HTTP client with logging: %w", err)
		}

		config = &tfe.Config{
			Address:    fmt.Sprintf("https://%s", hostname),
			Token:      token,
			HTTPClient: httpClient,
		}
	} else {
		// No logging - use default HTTP client
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
		OrganizationName: organization,
	}, nil
}
