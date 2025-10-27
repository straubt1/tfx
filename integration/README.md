# Integration Tests

Integration tests for tfx CLI commands that execute against real TFE/TFC APIs.

## Setup

Set the following environment variables:

```bash
# Required for all tests
export TFE_HOSTNAME="app.terraform.io"  # or your TFE hostname
export TFE_TOKEN="your-token-here"
export TFE_ORGANIZATION="your-org-name"

# Optional: for workspace show tests
export TEST_WORKSPACE_NAME="existing-workspace-name"

# Optional: for project show tests (at least one required)
export TEST_PROJECT_NAME="existing-project-name"
export TEST_PROJECT_ID="prj-xxxxx"
```

## Running Tests

```bash
# Run all integration tests
go test -v -tags=integration ./integration/...

# Run specific test
go test -v -tags=integration ./integration/ -run TestOrganization

# Run with timeout
go test -v -tags=integration -timeout 30s ./integration/...
```

## Test Coverage

### Organization Commands

- [x] organization list
- [x] organization show

### Workspace Commands

- [x] workspace list
- [x] workspace list (with search filter)
- [x] workspace show

### Project Commands

- [x] project list
- [x] project list (with search filter)
- [x] project show (by name)
- [x] project show (by id)

## Notes

- Tests use real API calls - they require valid credentials
- Tests do not create/destroy resources (read-only operations)
- Success criteria: command executes without error
- Tests are skipped if required environment variables are not set
