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
// If the viper config contains "tfxLogFile", HTTP logging will be enabled to that file
func NewFromViperWithContext(ctx context.Context) (*TfxClient, error) {
	hostname := viper.GetString("tfeHostname")
	token := viper.GetString("tfeToken")
	organization := viper.GetString("tfeOrganization")
	logFile := viper.GetString("tfxLogFile")

	return NewWithContext(ctx, hostname, token, organization, logFile)
}
