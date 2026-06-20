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

# Optional: profile-based varset lifecycle test (local.tfe.rocks)
export TFX_INTEGRATION_PROFILE="local"
# export TFX_CONFIG_FILE="$HOME/.tfx.hcl"  # when not using default paths

# Optional: keep varsets/variables after TestVariableSetLocalProfileLifecycle (for manual inspection)
export TFX_INTEGRATION_NO_CLEANUP=1
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

### Variable Set Commands

- [x] variable-set list / CRUD (env-based credentials)
- [x] variable-set lifecycle with `--profile local` (`TestVariableSetLocalProfileLifecycle`)

## Profile-based local testing

For a TFE instance at `local.tfe.rocks`, add a profile to `~/.tfx.hcl`:

```hcl
profile "local" {
  hostname     = "local.tfe.rocks"
  organization = "your-org"
  token        = "your-token"
}
```

Run the lifecycle test:

```bash
export TFX_INTEGRATION_PROFILE=local
go test -v -tags=integration -count=1 ./integration/ -run TestVariableSetLocalProfileLifecycle
```

Optional: set `TEST_PROJECT_NAME` and/or `TEST_WORKSPACE_NAME` to also create project-owned and workspace-assigned variable sets in the same run.

To leave created variable sets and variables in place (skip delete steps and `t.Cleanup`):

```bash
export TFX_INTEGRATION_NO_CLEANUP=1
go test -v -tags=integration -count=1 ./integration/ -run TestVariableSetLocalProfileLifecycle
```

The test logs the search prefix and varset names so you can find them in the UI or with `tfx varset list --search <prefix>`.

## Notes

- Tests use real API calls — they require valid credentials
- Most tests are read-only; `TestVariableSetLocalProfileLifecycle` creates and deletes variable sets (unless `TFX_INTEGRATION_NO_CLEANUP=1`)
- Success criteria: command executes without error (non-zero exit on CLI failure)
- Tests are skipped if required environment variables are not set
