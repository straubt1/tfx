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

### HTTP Request/Response Logging

For debugging or auditing purposes, you can enable HTTP logging to capture all API requests and responses:

```go
import "github.com/straubt1/tfx/client"

// Create a client with HTTP logging enabled
tfxClient, closer, err := client.NewFromViperWithLogging("/tmp/tfx-http.log")
if err != nil {
    log.Fatal(err)
}
// Important: Always close the log file when done
defer closer.Close()

// All HTTP requests and responses will now be logged to the file
projects, err := tfxClient.FetchProjects("my-org", "")
```

**Log File Format:**
```
################################################################################
# TFX HTTP LOG - Started at 2025-10-07T10:30:45Z
################################################################################

================================================================================
REQUEST @ 2025-10-07T10:30:45Z
================================================================================
GET /api/v2/organizations/my-org/projects?page[number]=1&page[size]=100 HTTP/1.1
Host: app.terraform.io
Authorization: Bearer [REDACTED]
Accept: application/vnd.api+json
...

--------------------------------------------------------------------------------
RESPONSE @ 2025-10-07T10:30:46Z
--------------------------------------------------------------------------------
HTTP/2.0 200 OK
Content-Type: application/vnd.api+json
...
```

**Environment Variable Support:**

You can also enable logging via environment variable in your commands:

```go
// In your command RunE function:
logFile := os.Getenv("TFX_HTTP_LOG")
if logFile != "" {
    c, closer, err := client.NewFromViperWithLogging(logFile)
    if err != nil {
        return err
    }
    defer closer.Close()
    // ... use c
} else {
    c, err := client.NewFromViper()
    if err != nil {
        return err
    }
    // ... use c
}
```

Then run:
```bash
TFX_HTTP_LOG=/tmp/debug.log ./tfx project list
```

## TfxClient Structure

The `TfxClient` struct provides:

- `Client` - The underlying `*tfe.Client` from hashicorp/go-tfe
- `Context` - A context.Context for API calls
- `Hostname` - The TFE/TFC hostname
- `OrganizationName` - The default organization name

## Client Methods

### Organization Operations
- `FetchOrganizations() ([]*tfe.Organization, error)` - Fetch all organizations

### Project Operations
- `FetchProjects(orgName, searchString string) ([]*tfe.Project, error)` - Fetch projects for an organization
- `FetchProjectsAcrossOrgs(searchString string) ([]*tfe.Project, error)` - Fetch projects across all organizations
- `FetchProject(projectID string, options *tfe.ProjectReadOptions) (*tfe.Project, error)` - Fetch a single project

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
