---
name: new-data-function
description: Create a new function in the data/ package following project conventions
---

# Data Layer Conventions

The `data/` package is the **sole boundary** between TFx and the go-tfe SDK. All TFE/HCP Terraform API calls happen exclusively here. No other package should import `go-tfe` for API calls. This isolation makes the codebase easier to test and maintain.

## Pagination — Always use `client/pagination.go`

Use `client.FetchAll[T]()` for **all** paginated list operations. Never write manual pagination loops. Use `client.NewPaginationFromTFE()` to convert go-tfe pagination. PageSize should always be 100.

**Canonical example** — `data/projects.go:FetchProjects()`:

```go
func FetchThings(c *client.TfxClient, orgName string) ([]*tfe.Thing, error) {
    return client.FetchAll(c.Context, func(pageNumber int) ([]*tfe.Thing, *client.Pagination, error) {
        opts := &tfe.ThingListOptions{
            ListOptions: tfe.ListOptions{PageNumber: pageNumber, PageSize: 100},
        }
        result, err := c.Client.Things.List(c.Context, orgName, opts)
        if err != nil {
            return nil, nil, err
        }
        return result.Items, client.NewPaginationFromTFE(result.Pagination), nil
    })
}
```

## Function Naming

| Pattern | Purpose | Example |
|---------|---------|---------|
| `Fetch*` (plural) | Paginated list returning slice | `FetchProjects`, `FetchWorkspaces` |
| `Fetch*` (singular) | Read single item by ID | `FetchProject`, `FetchRun` |
| `Fetch*ByName` | Search + exact match by name | `FetchProjectByName` |
| `Get*ID` | Resolve name to ID | `GetWorkspaceID` |
| `Create*` / `Update*` / `Delete*` | State-changing operations | `CreateVariable`, `DeleteVariable` |
| `List*` | List with optional `maxItems` limit | `ListRegistryModules` |

## Function Signature

- First parameter is always `c *client.TfxClient`
- Return go-tfe types when possible (e.g., `[]*tfe.Project`)
- When aggregating across multiple API calls, define custom result types in `cmd/views/` to avoid import cycles (see `cmd/views/run_policy.go:RunPolicyResult`, `cmd/views/admin_metrics.go:MetricsWorkspace`)

## SDK Gaps — Direct HTTP Calls

When go-tfe doesn't expose a needed field, use a direct HTTP call via `c.Hostname` and `c.Token`. See `fetchEvaluationOutputs()` in `data/policy_checks.go` for the pattern. Always add a comment explaining why the raw call is necessary.

## Logging

Use `output.Get().Logger()` with structured key-value pairs:

- `Debug` — function entry/exit with parameters
- `Trace` — per-page pagination details
- `Error` — failures with full context
- `Info` — completion summaries with counts

## Error Handling

- Use `github.com/pkg/errors` — wrap with context: `errors.Wrap(err, "message")`
- Use `errors.Errorf()` for custom "not found" messages
- In multi-item loops, `log.Error` + `continue` for individual item failures (don't fail the whole batch)

## File Organization

- One file per resource: `data/projects.go`, `data/workspaces.go`
- File names use plural: `projects.go`, `variables.go`
- Integration tests co-located: `data/projects_integration_test.go`

## File Header

```go
// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT
```

## Reference Implementations

- **List with pagination**: `data/projects.go` — `FetchProjects()`
- **Single read**: `data/runs.go` — `FetchRun()`
- **Read by name**: `data/projects.go` — `FetchProjectByName()`
- **Create/Update/Delete**: `data/variables.go` — `CreateVariable()`, `UpdateVariable()`, `DeleteVariable()`
- **Complex aggregation with custom types**: `data/policy_checks.go` — `FetchRunPolicyDetails()`
- **Direct HTTP for SDK gaps**: `data/policy_checks.go` — `fetchEvaluationOutputs()`

## Task

Create the data layer function for: $ARGUMENTS

Follow all conventions above. Use `data/projects.go` as your primary reference.
