# tfx - HCP Terraform/Terraform Enterprise CLI

## Project Overview

tfx is a standalone CLI for HCP Terraform (Terraform Cloud) and Terraform Enterprise built with Go and Cobra. It provides API-driven workflows for managing workspaces, organizations, projects, runs, and the private registry.

**Key External Dependency:** `github.com/hashicorp/go-tfe` - Official TFE/TFC Go client library

## Architecture (Layered Design)

The codebase follows a strict **5-layer separation of concerns** pattern (see `Refactor.md` for detailed rationale):

```
cmd/ → flags/ → client/ → data/ → TFE API
  ↓      ↓         ↓        ↓
  └──────┴─────────┴────────┴──→ views/ → output/ → stdout
```

### Layer Responsibilities

1. **cmd/** - Command orchestration via Cobra
   - Wire layers together (typically 10-20 lines per command)
   - Example pattern: `cmd/project.go`, `cmd/workspace.go`
   - Each command file has an `init()` that registers flags and subcommands

2. **flags/** - Type-safe flag parsing
   - Convert cobra flags into structured config objects
   - Pattern: `*Flags` struct + `Parse*Flags(cmd)` function
   - Example: `flags/project.go` has `ProjectListFlags`, `ProjectShowFlags`

3. **client/** - TFE client wrapper with context
   - `TfxClient` struct wraps `*tfe.Client` with org/context/hostname
   - `NewFromViper()` - Primary constructor, reads from viper config
   - `FetchAll[T]()` - **Generic pagination helper** (use this for all list operations)

4. **data/** - API interaction & business logic
   - All TFE API calls happen here
   - Functions named `Fetch*()` - e.g., `FetchProjects()`, `FetchWorkspace()`
   - Uses `client.FetchAll()` for automatic pagination

5. **views/** - Output rendering (format-agnostic)
   - Each view has `Render()` method
   - Views use `output.Get()` singleton for terminal/JSON rendering
   - Pattern: Embed `*BaseView` for common functionality

6. **output/** - Singleton output system
   - `output.Get()` returns global instance (auto-initializes from viper)
   - Supports `--json` flag for machine-readable output
   - Manages spinner (disabled when `TFX_LOG` env var set)

## Command Implementation Pattern

When adding a new command, follow this exact sequence:

```go
// 1. cmd/mycommand.go - Orchestration
func myCommandList(cmdConfig *flags.MyCommandListFlags) error {
    v := view.NewMyCommandListView()
    c, err := client.NewFromViper()
    if err != nil { return v.RenderError(err) }
    
    data, err := data.FetchMyThings(c, cmdConfig.Search)
    if err != nil { return v.RenderError(err) }
    
    return v.Render(data)
}

// 2. flags/mycommand.go - Flag parsing
type MyCommandListFlags struct {
    Search string
}
func ParseMyCommandListFlags(cmd *cobra.Command) (*MyCommandListFlags, error) {
    return &MyCommandListFlags{
        Search: cmd.Flags().GetString("search"),
    }, nil
}

// 3. data/mycommand.go - API calls
func FetchMyThings(c *client.TfxClient, search string) ([]*tfe.MyThing, error) {
    return client.FetchAll(c.Context, func(pageNum int) ([]*tfe.MyThing, *client.Pagination, error) {
        opts := &tfe.MyThingListOptions{
            ListOptions: tfe.ListOptions{PageNumber: pageNum, PageSize: 100},
        }
        result, err := c.Client.MyThings.List(c.Context, c.OrganizationName, opts)
        return result.Items, client.NewPaginationFromTFE(result.Pagination), err
    })
}

// 4. views/mycommand_list.go - Rendering
type MyCommandListView struct{ *BaseView }
func NewMyCommandListView() *MyCommandListView { 
    return &MyCommandListView{NewBaseView()} 
}
func (v *MyCommandListView) Render(data []*tfe.MyThing) error {
    // Use v.renderer for output
}
```

## Testing

### Unit Tests
- Standard Go tests: `go test ./...`
- No special setup needed

### Integration Tests

Two types of integration tests, both using `//go:build integration` tag:

1. **Command-level tests** - `integration/` package (separate directory)
   - Tests full command execution via `cmd.GetRootCommand().Execute()`
   - Examples: `integration/organization_test.go`, `integration/workspace_test.go`
   - Run with: `task test-integration-cmd`

2. **Data-level tests** - `data/*_integration_test.go` (co-located with code)
   - Tests API layer directly (data.FetchProjects, etc.)
   - Examples: `data/organizations_integration_test.go`, `data/projects_integration_test.go`
   - Run with: `task test-integration-data`

**Setup:** Environment variables in `secrets/.env-int` (see `secrets/.env-int.example`)

**Environment variables:**
- Required: `TFE_HOSTNAME`, `TFE_TOKEN`, `TFE_ORGANIZATION`
- Optional: `TEST_WORKSPACE_NAME`, `TEST_PROJECT_NAME`, `TEST_PROJECT_ID`

**Integration test pattern:**
```go
//go:build integration

func TestMyCommand(t *testing.T) {
    hostname, token, org := setupTest(t)  // Helper in setup_test.go
    err := executeCommand(t, 
        []string{"mycommand", "list"}, 
        hostname, token, org,
    )
    if err != nil { t.Errorf("Command failed: %v", err) }
}
```

### Test Helpers
- **cmd/root_test_helpers.go** - Exports `GetRootCommand()` for integration tests only (not compiled in production builds)

## Build & Development

### Task Commands (via Taskfile.yml)
```bash
task go-build              # Build binary
task test                  # Unit tests
task test-integration-cmd  # Integration tests (commands)
task test-integration-data # Integration tests (data layer)
task serve-docs           # Serve mkdocs documentation
```

### Configuration System
- **Viper + Cobra:** All commands use viper for config management
- **Config file:** `.tfx.hcl` (HCL format, optional, searches current dir and home)
- **Persistent flags:** `--tfeHostname`, `--tfeToken`, `--tfeOrganization` (or via ENV)
- **ENV variables:** `TFE_HOSTNAME`, `TFE_TOKEN`, `TFE_ORGANIZATION`

### Debugging (Fully Implemented)
- `TFX_LOG=DEBUG` - Enable debug logging to terminal (auto-disables spinner)
- `TFX_LOG_PATH=/path/to/dir` - Save debug logs to directory
- `TFX_HTTP_LOG=/path/to/file` - Log all HTTP requests/responses (see `client/README.md`)
  - Logs include full request/response with redacted auth tokens
  - Useful for debugging API interactions

## Critical Patterns & Conventions

### DO Use Generic Pagination
```go
// ✅ CORRECT - Uses FetchAll generic helper
data, err := client.FetchAll(ctx, func(pageNum int) ([]T, *Pagination, error) { ... })

// ❌ WRONG - Manual pagination loop
for { /* manual page iteration */ }
```

### DO Use Views for All Output
```go
// ✅ CORRECT - View handles rendering
view := view.NewMyView()
return view.Render(data)

// ❌ WRONG - Direct fmt.Print or json.Marshal
fmt.Println(data)  // Don't do this
```

### DO Use Type-Safe Flag Configs
```go
// ✅ CORRECT - Structured flags
cmdConfig, err := flags.ParseMyFlags(cmd)
value := cmdConfig.MyField

// ❌ WRONG - String-based flag access
value, _ := cmd.Flags().GetString("my-field")
```

### DO Handle Errors with Views
```go
// ✅ CORRECT - View handles error rendering
if err != nil { return v.RenderError(err) }

// ❌ WRONG - Direct error return
if err != nil { return err }  // Breaks JSON mode
```

### DO Use Viper for Config
```go
// ✅ CORRECT - Viper-aware client
c, err := client.NewFromViper()

// ❌ WRONG - Hardcoded values
c, err := client.New("app.terraform.io", "token", "org")
```

## Common Gotchas

1. **Output Singleton:** `output.Get()` initializes from viper on first call - don't call before cobra flag parsing
2. **Integration Tests:** Must have `//go:build integration` tag or they won't compile with `-tags=integration`
3. **Spinner:** Auto-disabled when `TFX_LOG` set or `--json` flag used
4. **Required Flags:** Set via `cmd.MarkFlagRequired()` in `init()`, but viper can satisfy from ENV/config
5. **BaseView:** Always embed `*BaseView` in custom views for standard methods

## File Naming Conventions

- Commands: `cmd/resource.go` (e.g., `workspace.go`, `project.go`)
- Flags: `flags/resource.go` (matches cmd file)
- Data: `data/resources.go` (plural, e.g., `workspaces.go`, `projects.go`)
- Views: `views/resource_operation.go` (e.g., `project_list.go`, `workspace_show.go`)
- Integration tests: `integration/resource_test.go`

## Documentation

- **Site:** `site/` directory (mkdocs) - `task serve-docs` to preview
- **Architecture:** See `Refactor.md` for detailed architectural decisions
- **Client:** See `client/README.md` for client usage patterns
