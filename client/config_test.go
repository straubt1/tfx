package client

import (
	"context"
	"testing"

	"github.com/spf13/viper"
)

func TestNewFromViper(t *testing.T) {
	tests := []struct {
		name       string
		setupViper func()
		wantErr    bool
		errMsg     string
	}{
		{
			name: "valid viper configuration",
			setupViper: func() {
				viper.Set("tfeHostname", "app.terraform.io")
				viper.Set("tfeToken", "test-token")
				viper.Set("tfeOrganization", "test-org")
			},
			wantErr: false,
		},
		{
			name: "missing hostname in viper",
			setupViper: func() {
				viper.Set("tfeHostname", "")
				viper.Set("tfeToken", "test-token")
				viper.Set("tfeOrganization", "test-org")
			},
			wantErr: true,
			errMsg:  "hostname is required",
		},
		{
			name: "missing token in viper",
			setupViper: func() {
				viper.Set("tfeHostname", "app.terraform.io")
				viper.Set("tfeToken", "")
				viper.Set("tfeOrganization", "test-org")
			},
			wantErr: true,
			errMsg:  "token is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset viper before each test
			viper.Reset()
			tt.setupViper()

			client, err := NewFromViper()
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewFromViper() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("NewFromViper() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("NewFromViper() unexpected error = %v", err)
				return
			}
			if client == nil {
				t.Error("NewFromViper() returned nil client")
			}
		})
	}
}

func TestNewFromViperWithContext(t *testing.T) {
	viper.Reset()
	viper.Set("tfeHostname", "app.terraform.io")
	viper.Set("tfeToken", "test-token")
	viper.Set("tfeOrganization", "test-org")

	ctx := context.Background()
	client, err := NewFromViperWithContext(ctx)
	if err != nil {
		t.Fatalf("NewFromViperWithContext() unexpected error = %v", err)
	}
	if client.Context != ctx {
		t.Error("NewFromViperWithContext() context not set correctly")
	}
}
