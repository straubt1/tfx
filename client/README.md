# TFX Client Package

This package provides a clean interface for creating and managing Terraform Enterprise/Cloud API clients.

## Usage

### Basic Usage

```go
import "github.com/straubt1/tfx/client"

// Create a client with explicit configuration
tfxClient, err := client.New("app.terraform.io", "your-token", "your-org")
if err != nil {
    log.Fatal(err)
}

// Use the client
workspaces, err := tfxClient.Client.Workspaces.List(
    tfxClient.Context,
    tfxClient.OrganizationName,
    nil,
)
```

### Using Viper Configuration

```go
import "github.com/straubt1/tfx/client"

// Create a client from viper configuration (tfeHostname, tfeToken, tfeOrganization)
tfxClient, err := client.NewFromViper()
if err != nil {
    log.Fatal(err)
}
```

### Using Custom Context

```go
import (
    "context"
    "time"
    "github.com/straubt1/tfx/client"
)

// Create a client with timeout context
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

tfxClient, err := client.NewWithContext(ctx, "app.terraform.io", "your-token", "your-org")
if err != nil {
    log.Fatal(err)
}
```

## TfxClient Structure

The `TfxClient` struct provides:

- `Client` - The underlying `*tfe.Client` from hashicorp/go-tfe
- `Context` - A context.Context for API calls
- `Hostname` - The TFE/TFC hostname
- `OrganizationName` - The default organization name

## Error Handling

All functions return proper errors that can be wrapped and handled gracefully:

```go
tfxClient, err := client.NewFromViper()
if err != nil {
    return fmt.Errorf("failed to create TFE client: %w", err)
}
```

## Testing

The package includes unit tests. Run them with:

```bash
go test ./client/...
```
