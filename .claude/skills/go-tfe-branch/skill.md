---
name: go-tfe-branch
description: Set up a new branch in the go-tfe fork and implement a code change. Use when starting work on a new contribution to go-tfe — syncs main with upstream, creates a branch, and implements the requested change for review. The go-tfe repo is at ../go-tfe/ relative to this project.
---

# go-tfe Branch Setup & Implementation

You are helping implement a contribution to the [hashicorp/go-tfe](https://github.com/hashicorp/go-tfe) repository. The fork lives at `../go-tfe/` relative to the current working directory (i.e., `/Users/tstraub/Projects/straubt1.github.com/go-tfe/`).

## Your Task

$ARGUMENTS

## Step-by-Step Process

### 1. Sync main with upstream

Work in the go-tfe directory. Ensure the upstream remote exists and main is current:

```bash
cd ~/Projects/straubt1.github.com/go-tfe

# Add upstream if not present
git remote get-url upstream 2>/dev/null || git remote add upstream https://github.com/hashicorp/go-tfe.git

# Fetch upstream and sync main
git fetch upstream
git checkout main
git merge upstream/main --ff-only
git push origin main
```

If `--ff-only` fails (diverged), report this to the user and stop — do not force-push or rebase without explicit instruction.

### 2. Create a feature branch

Branch names should be short, lowercase, hyphen-separated, descriptive of the change. Examples:
- `add-workspace-foo-field`
- `add-project-bar-method`
- `fix-run-list-pagination`

Derive the branch name from the task description. Create and check it out:

```bash
git checkout -b <branch-name>
```

### 3. Implement the change

Read existing code carefully before writing anything. Understand the patterns in place.

#### go-tfe Code Conventions

**File header** (all .go files):
```go
// Copyright IBM Corp. 2018, 2025
// SPDX-License-Identifier: MPL-2.0
```

**Package**: `package tfe`

**Interface + implementation pattern** (for resources):
```go
// Interface (exported, upper-camel)
type Things interface {
    List(ctx context.Context, org string, options *ThingListOptions) (*ThingList, error)
    Create(ctx context.Context, options ThingCreateOptions) (*Thing, error)
    Read(ctx context.Context, thingID string) (*Thing, error)
    Update(ctx context.Context, thingID string, options ThingUpdateOptions) (*Thing, error)
    Delete(ctx context.Context, thingID string) error
}

// Compile-time proof
var _ Things = (*things)(nil)

// Private implementation struct
type things struct {
    client *Client
}
```

**Struct tags**: Use `jsonapi` tags for all API-serialized fields:
- Required fields: `jsonapi:"attr,field-name"`
- Optional fields: `jsonapi:"attr,field-name,omitempty"`
- Relations: `jsonapi:"relation,relation-name"`
- Primary ID: `jsonapi:"primary,type-name"`

**Optional fields** use pointers: `*string`, `*bool`, `*int`

**Options structs** for create/update:
```go
type ThingCreateOptions struct {
    Type string `jsonapi:"primary,things"`          // required for jsonapi
    Name *string `jsonapi:"attr,name"`              // required fields still use pointer in options
    Description *string `jsonapi:"attr,description,omitempty"` // optional
}
```

**Method implementation pattern**:
```go
func (s *things) Create(ctx context.Context, options ThingCreateOptions) (*Thing, error) {
    if err := options.valid(); err != nil {
        return nil, err
    }
    req, err := s.client.NewRequest("POST", "things", &options)
    if err != nil {
        return nil, err
    }
    t := &Thing{}
    err = s.client.Do(ctx, req, t)
    if err != nil {
        return nil, err
    }
    return t, nil
}
```

**Validation**:
```go
func (o ThingCreateOptions) valid() error {
    if !validString(o.Name) {
        return ErrRequiredName
    }
    return nil
}
```

**List options** always embed `ListOptions`:
```go
type ThingListOptions struct {
    ListOptions
    Search string `url:"search[name],omitempty"`
}
```

**API docs comment** on interfaces:
```go
// Things describes all the thing related methods that the Terraform
// Enterprise API supports.
//
// TFE API docs: https://developer.hashicorp.com/terraform/cloud-docs/api-docs/things
type Things interface {
```

**Wire up** in `tfe.go`: The new resource must be added to the `Client` struct and initialized in `NewClient()`. Look at how existing resources are wired — follow the exact same pattern.

#### When adding a field to an existing struct

- Add the field to the struct with proper `jsonapi` tag
- Add the field to CreateOptions/UpdateOptions if it can be set
- Validate required fields in `valid()`
- Read existing tests in `*_integration_test.go` to understand what's tested

#### Tests

Always add or update tests for every change. There are two patterns:

**Unit tests** (preferred for new fields/structs — no live TFE instance needed):
Use `httptest.NewServer` to serve a crafted JSON-API response and assert the field is correctly deserialized. Create a `*_test.go` file (e.g., `policy_evaluation_test.go`) using `package tfe`. Two subtests per field:
1. Field present — serve a response with the field populated, assert the value is correct
2. Field absent — serve a response without the field, assert it is nil/empty (omitempty works)

```go
func TestThing_FieldDeserialization(t *testing.T) {
    t.Parallel()
    testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/vnd.api+json")
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, `{ "data": [...] }`)
    }))
    defer testServer.Close()

    client, err := NewClient(&Config{Address: testServer.URL, Token: "fake-token"})
    require.NoError(t, err)
    // call the relevant client method and assert
}
```

**Integration test updates** (for `*_integration_test.go` or `*_beta_test.go`):
- Add assertions for the new field where the resource is already being read/listed
- If the field cannot be populated in the test environment (e.g., Sentinel print output when OPA is used), add a comment explaining why and assert it is empty

#### CHANGELOG.md

Add an entry under `# Unreleased` at the top of CHANGELOG.md. Format:
```
* Adds `FieldName` field to `StructName` by @straubt1 [#<PR_NUMBER>](https://github.com/hashicorp/go-tfe/pull/<PR_NUMBER>)
```

Use a placeholder like `#XXXX` for the PR number — it will be filled in after the PR is created. Pick the correct category:
- `## Enhancements` — new fields, new methods, new options
- `## Bug Fixes` — fixing incorrect behavior

**CHANGELOG rules:**
- Describe the change in terms of the go-tfe API only — what struct, field, or method changed
- Never mention tfx, the tfx project, or any external consumer of go-tfe
- Do not explain why a caller might want the field; just describe what was added

### 4. Run tests (if possible)

```bash
cd /Users/tstraub/Projects/straubt1.github.com/go-tfe
go vet ./...
go build ./...
```

Unit tests (non-integration):
```bash
go test ./... -run "^Test[^I]" -timeout=5m
```

Report any failures. Do not silently skip test failures.

### 5. Present for review

Do NOT commit. Display:
1. A summary of all files changed and what was changed in each
2. The intended commit message (so the user can copy it)
3. Any open questions or things to verify

**Intended commit message format**:
```
<type>: <short description>

<optional body with more detail if needed>
```

Types: `feat`, `fix`, `refactor`, `test`, `docs`

Example:
```
feat: add SettingOverwrites field to ProjectUpdateOptions

Exposes the setting-overwrites attribute in the update options struct
so callers can modify project-level setting overrides via the API.
```

## Important Rules

- Never commit unless the user explicitly asks
- Never push unless the user explicitly asks
- Always read existing files before writing — match existing patterns exactly
- If the task is ambiguous, ask for clarification before writing code
- Do not mention this tool or its name in commit messages or PR descriptions
