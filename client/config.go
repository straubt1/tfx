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
	hostname := viper.GetString("tfeHostname")
	token := viper.GetString("tfeToken")
	organization := viper.GetString("tfeOrganization")

	return NewWithContext(ctx, hostname, token, organization)
}
