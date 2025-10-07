package client

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name         string
		hostname     string
		token        string
		organization string
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "valid configuration",
			hostname:     "app.terraform.io",
			token:        "test-token",
			organization: "test-org",
			wantErr:      false,
		},
		{
			name:         "missing hostname",
			hostname:     "",
			token:        "test-token",
			organization: "test-org",
			wantErr:      true,
			errMsg:       "hostname is required",
		},
		{
			name:         "missing token",
			hostname:     "app.terraform.io",
			token:        "",
			organization: "test-org",
			wantErr:      true,
			errMsg:       "token is required",
		},
		{
			name:         "organization is optional",
			hostname:     "app.terraform.io",
			token:        "test-token",
			organization: "",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(tt.hostname, tt.token, tt.organization)
			if tt.wantErr {
				if err == nil {
					t.Errorf("New() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("New() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("New() unexpected error = %v", err)
				return
			}
			if client == nil {
				t.Error("New() returned nil client")
				return
			}
			if client.Hostname != tt.hostname {
				t.Errorf("New() hostname = %v, want %v", client.Hostname, tt.hostname)
			}
			if client.OrganizationName != tt.organization {
				t.Errorf("New() organization = %v, want %v", client.OrganizationName, tt.organization)
			}
		})
	}
}

func TestNewWithContext(t *testing.T) {
	ctx := context.Background()
	client, err := NewWithContext(ctx, "app.terraform.io", "test-token", "test-org")
	if err != nil {
		t.Fatalf("NewWithContext() unexpected error = %v", err)
	}
	if client.Context != ctx {
		t.Error("NewWithContext() context not set correctly")
	}
}
