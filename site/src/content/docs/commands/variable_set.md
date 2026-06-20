---
title: Variable Set Commands
---

General commands to manage Variable Sets and their variables.

:::note
All commands below can be used with a `varset` alias. Variable subcommands also accept a `var` alias (`tfx varset var`).
:::

## Ownership vs assignment

Variable sets have two independent concepts:

| Concept | What it means | How to set it |
|---|---|---|
| **Ownership** | Who owns the variable set (organization or project) | Default create = org-owned. Use `--project-name` on create for a project-owned set. |
| **Assignment** | Where the variable set applies (workspaces/projects) | `--global` applies to all workspaces in the org. `--workspace-name` on create applies to one workspace after creation. Project/workspace list scopes show sets *assigned* to that target. |

Ownership and assignment are separate: an org-owned set can be assigned to specific workspaces; a project-owned set is owned by a project but may still be assigned elsewhere via the API.

### Scope flags (list, show, delete, and variable commands)

These flags select which variable sets to list or how to resolve a variable set by name:

| Flag | Used by | Description |
|---|---|---|
| `--organization-name` | list, show, delete, variable | Organization to use (defaults to configured organization). |
| `--project-name` | list, show, delete, variable | List sets assigned to a project, or narrow name lookup to that project scope. |
| `--workspace-name` | list, show, delete, variable | List sets assigned to a workspace, or narrow name lookup to that workspace scope. |
| `--all` / `-a` | list only | List variable sets across all organizations you can access. |
| `--search` / `-s` | list only | Filter results by variable set name (composes with any list scope). |

## `tfx variable-set list`

List variable sets. Default scope is the configured organization. Use scope flags to list by project, workspace, or all organizations.

| Flag | Description |
|---|---|
| `--search` / `-s` | Filter by variable set name (optional). |
| `--all` / `-a` | List across all organizations (mutually exclusive with `--project-name` and `--workspace-name`). |
| `--organization-name` | Organization to list (optional, defaults to configured org). |
| `--project-name` | List variable sets assigned to this project (optional). |
| `--workspace-name` | List variable sets assigned to this workspace (optional). |

**Organization Example (default)**

```sh
$ tfx varset list
Using config file: /Users/tstraub/.tfx.hcl
Listing variable sets in organization 'firefly'
╭──────────────────────┬─────────────────────┬────────┬──────────┬───────────────────────────────╮
│ NAME                 │ ID                  │ GLOBAL │ PRIORITY │ PARENT                        │
├──────────────────────┼─────────────────────┼────────┼──────────┼───────────────────────────────┤
│ aws-credentials      │ varset-abc123XYZ456 │ false  │ false    │ organization:firefly          │
│ common-tags          │ varset-def789GHI012 │ true   │ false    │ organization:firefly          │
│ production-overrides │ varset-jkl345MNO678 │ false  │ true     │ project:prj-ABC123defGHI789   │
╰──────────────────────┴─────────────────────┴────────┴──────────┴───────────────────────────────╯
```

**Search Example**

```sh
$ tfx varset list --search aws
Using config file: /Users/tstraub/.tfx.hcl
Listing variable sets in organization 'firefly' matching 'aws'
╭─────────────────┬─────────────────────┬────────┬──────────┬──────────────────────╮
│ NAME            │ ID                  │ GLOBAL │ PRIORITY │ PARENT               │
├─────────────────┼─────────────────────┼────────┼──────────┼──────────────────────┤
│ aws-credentials │ varset-abc123XYZ456 │ false  │ false    │ organization:firefly │
╰─────────────────┴─────────────────────┴────────┴──────────┴──────────────────────╯
```

**Project-scoped Example**

```sh
$ tfx varset list --project-name my-project
Using config file: /Users/tstraub/.tfx.hcl
Listing variable sets for project 'my-project' in organization 'firefly'
╭──────────────────────┬─────────────────────┬────────┬──────────┬─────────────────────────────╮
│ NAME                 │ ID                  │ GLOBAL │ PRIORITY │ PARENT                      │
├──────────────────────┼─────────────────────┼────────┼──────────┼─────────────────────────────┤
│ production-overrides │ varset-jkl345MNO678 │ false  │ true     │ project:prj-ABC123defGHI789 │
╰──────────────────────┴─────────────────────┴────────┴──────────┴─────────────────────────────╯
```

**Workspace-scoped Example**

```sh
$ tfx varset list --workspace-name prod-us-east-1
Using config file: /Users/tstraub/.tfx.hcl
Listing variable sets for workspace 'prod-us-east-1' in organization 'firefly'
╭─────────────────┬─────────────────────┬────────┬──────────┬──────────────────────╮
│ NAME            │ ID                  │ GLOBAL │ PRIORITY │ PARENT               │
├─────────────────┼─────────────────────┼────────┼──────────┼──────────────────────┤
│ aws-credentials │ varset-abc123XYZ456 │ false  │ false    │ organization:firefly │
╰─────────────────┴─────────────────────┴────────┴──────────┴──────────────────────╯
```

**All Organizations Example**

```sh
$ tfx varset list --all
Using config file: /Users/tstraub/.tfx.hcl
Listing variable sets across all organizations
╭──────────────┬─────────────────┬─────────────────────┬────────┬──────────┬──────────────────────╮
│ ORGANIZATION │ NAME            │ ID                  │ GLOBAL │ PRIORITY │ PARENT               │
├──────────────┼─────────────────┼─────────────────────┼────────┼──────────┼──────────────────────┤
│ firefly      │ aws-credentials │ varset-abc123XYZ456 │ false  │ false    │ organization:firefly │
│ staging      │ common-tags     │ varset-def789GHI012 │ true   │ false    │ organization:staging │
╰──────────────┴─────────────────┴─────────────────────┴────────┴──────────┴──────────────────────╯
```

## `tfx variable-set show`

Show details of a variable set by ID or name, including assigned workspaces, projects, and variables. When using `--name`, scope flags help resolve the correct set when names are ambiguous.

| Flag | Description | Required |
|---|---|---|
| `--id` / `-i` | ID of the variable set | One of `--id` or `--name` |
| `--name` / `-n` | Name of the variable set | One of `--id` or `--name` |
| `--organization-name` | Organization scope (optional). | No |
| `--project-name` | Project scope for name resolution (optional). | No |
| `--workspace-name` | Workspace scope for name resolution (optional). | No |

**Example (by ID)**

```sh
$ tfx variable-set show -i varset-abc123XYZ456
Using config file: /Users/tstraub/.tfx.hcl
Showing variable set 'varset-abc123XYZ456'
ID:          varset-abc123XYZ456
Name:        aws-credentials
Description: AWS credentials for production workloads
Global:      false
Priority:    false
Parent:      organization:firefly
Workspaces:
  prod-us-east-1: ws-hLFv8c9bjgXC3mdK
  prod-us-west-2: ws-yN6DnhYxB39qqAre
Projects:
Variables:
  AWS_ACCESS_KEY_ID:     var-GHI789jklMNO345
  AWS_SECRET_ACCESS_KEY: var-JKL012mnoPQR678
```

**Example (by name)**

```sh
$ tfx varset show --name aws-credentials
Using config file: /Users/tstraub/.tfx.hcl
Showing variable set 'aws-credentials'
ID:          varset-abc123XYZ456
Name:        aws-credentials
...
```

**JSON Example**

```sh
$ tfx variable-set show -i varset-abc123XYZ456 --json | jq .
{
  "id": "varset-abc123XYZ456",
  "name": "aws-credentials",
  "description": "AWS credentials for production workloads",
  "global": false,
  "priority": false,
  "parent": {
    "type": "organization",
    "id": "firefly"
  },
  "workspaces": [
    { "id": "ws-hLFv8c9bjgXC3mdK", "name": "prod-us-east-1" },
    { "id": "ws-yN6DnhYxB39qqAre", "name": "prod-us-west-2" }
  ],
  "projects": [],
  "variables": [
    { "id": "var-GHI789jklMNO345", "key": "AWS_ACCESS_KEY_ID",     "value": "" },
    { "id": "var-JKL012mnoPQR678", "key": "AWS_SECRET_ACCESS_KEY", "value": "" }
  ]
}
```

## `tfx variable-set create`

Create a new variable set. Default ownership is the configured organization.

| Flag | Description | Required |
|---|---|---|
| `--name` / `-n` | Name of the variable set | Yes |
| `--description` / `-d` | Description of the variable set | No |
| `--global` | Apply this variable set to all workspaces in the organization | No |
| `--priority` | Variable values in this set override workspace-level values | No |
| `--organization-name` | Organization to create in (defaults to configured org) | No |
| `--project-name` | Create as a project-owned variable set (mutually exclusive with `--global`) | No |
| `--workspace-name` | Apply the variable set to this workspace after creation | No |

**Basic Example (org-owned)**

```sh
$ tfx variable-set create --name aws-credentials --description "AWS credentials for production workloads"
Using config file: /Users/tstraub/.tfx.hcl
Creating variable set 'aws-credentials' in organization 'firefly'
ID:          varset-abc123XYZ456
Name:        aws-credentials
Description: AWS credentials for production workloads
Global:      false
Priority:    false
```

**Global Variable Set Example**

```sh
$ tfx variable-set create --name common-tags --description "Tags applied to all workspaces" --global
Using config file: /Users/tstraub/.tfx.hcl
Creating variable set 'common-tags' in organization 'firefly'
ID:          varset-def789GHI012
Name:        common-tags
Description: Tags applied to all workspaces
Global:      true
Priority:    false
```

**Project-owned Example**

```sh
$ tfx varset create --name project-vars --project-name my-project --description "Owned by my-project"
Using config file: /Users/tstraub/.tfx.hcl
Creating variable set 'project-vars' in organization 'firefly'
ID:          varset-mno345PQR678
Name:        project-vars
Description: Owned by my-project
Global:      false
Priority:    false
Parent:      project:prj-ABC123defGHI789
```

**Workspace Assignment Example**

```sh
$ tfx varset create --name ws-overrides --workspace-name prod-us-east-1 --description "Applied to prod-us-east-1"
Using config file: /Users/tstraub/.tfx.hcl
Creating variable set 'ws-overrides' in organization 'firefly'
ID:          varset-pqr678STU901
Name:        ws-overrides
Description: Applied to prod-us-east-1
Global:      false
Priority:    false
```

## `tfx variable-set delete`

Delete a variable set by ID or name. Scope flags apply when using `--name`.

| Flag | Description | Required |
|---|---|---|
| `--id` / `-i` | ID of the variable set to delete | One of `--id` or `--name` |
| `--name` / `-n` | Name of the variable set to delete | One of `--id` or `--name` |
| `--organization-name` | Organization scope (optional). | No |
| `--project-name` | Project scope for name resolution (optional). | No |
| `--workspace-name` | Workspace scope for name resolution (optional). | No |

:::caution
This permanently deletes the variable set and removes it from all workspaces and projects it is assigned to.
:::

**Example (by ID)**

```sh
$ tfx variable-set delete -i varset-abc123XYZ456
Using config file: /Users/tstraub/.tfx.hcl
Deleting variable set 'varset-abc123XYZ456'
Status: Success
ID:     varset-abc123XYZ456
```

**Example (by name)**

```sh
$ tfx varset delete --name aws-credentials
Using config file: /Users/tstraub/.tfx.hcl
Deleting variable set 'aws-credentials'
Status: Success
ID:     varset-abc123XYZ456
```

## `tfx variable-set variable`

Manage variables within a variable set. All subcommands require `--varset-id` or `--varset-name` (mutually exclusive). When using `--varset-name`, scope flags narrow lookup the same way as show/delete.

Common flags for every variable subcommand:

| Flag | Description |
|---|---|
| `--varset-id` | ID of the variable set |
| `--varset-name` | Name of the variable set |
| `--organization-name` | Organization scope (optional). |
| `--project-name` | Project scope for varset name resolution (optional). |
| `--workspace-name` | Workspace scope for varset name resolution (optional). |

### `tfx varset variable list`

List variables in a variable set.

```sh
$ tfx varset variable list --varset-name aws-credentials
Using config file: /Users/tstraub/.tfx.hcl
Listing variables for variable set 'aws-credentials'
╭───────────────────────┬─────────────────────┬───────────┬───────────┬───────╮
│ KEY                   │ ID                  │ CATEGORY  │ SENSITIVE │ HCL   │
├───────────────────────┼─────────────────────┼───────────┼───────────┼───────┤
│ AWS_ACCESS_KEY_ID     │ var-GHI789jklMNO345 │ terraform │ false     │ false │
│ AWS_SECRET_ACCESS_KEY │ var-JKL012mnoPQR678 │ terraform │ true      │ false │
╰───────────────────────┴─────────────────────┴───────────┴───────────┴───────╯
```

### `tfx varset variable create`

Create a variable in a variable set.

| Flag | Description | Required |
|---|---|---|
| `--key` / `-k` | Variable key | Yes |
| `--value` / `-v` | Variable value | One of `--value` or `--value-file` |
| `--value-file` / `-f` | Read value from file | One of `--value` or `--value-file` |
| `--description` / `-d` | Description | No |
| `--env` / `-e` | Environment variable (default: Terraform variable) | No |
| `--hcl` | Value is HCL | No |
| `--sensitive` / `-s` | Sensitive variable | No |

Environment variable keys (`--env`) must begin with a letter or underscore and contain only letters, numbers, and underscores (no hyphens).

```sh
$ tfx varset variable create --varset-name aws-credentials -k AWS_REGION -v us-east-1 -d "Default AWS region"
Using config file: /Users/tstraub/.tfx.hcl
Creating variable 'AWS_REGION' for variable set 'aws-credentials'
ID:          var-XYZ789abcDEF012
Key:         AWS_REGION
Value:       us-east-1
Sensitive:   false
HCL:         false
Category:    terraform
Description: Default AWS region
```

### `tfx varset variable show`

Show a single variable by key.

```sh
$ tfx varset variable show --varset-name aws-credentials -k AWS_REGION
Using config file: /Users/tstraub/.tfx.hcl
Showing variable 'AWS_REGION' for variable set 'aws-credentials'
ID:          var-XYZ789abcDEF012
Key:         AWS_REGION
Value:       us-east-1
Sensitive:   false
HCL:         false
Category:    terraform
Description: Default AWS region
```

### `tfx varset variable update`

Update a variable by key. Same flags as create.

```sh
$ tfx varset variable update --varset-name aws-credentials -k AWS_REGION -v us-west-2
Using config file: /Users/tstraub/.tfx.hcl
Updating variable 'AWS_REGION' for variable set 'aws-credentials'
ID:          var-XYZ789abcDEF012
Key:         AWS_REGION
Value:       us-west-2
...
```

### `tfx varset variable delete`

Delete a variable by key.

```sh
$ tfx varset variable delete --varset-name aws-credentials -k AWS_REGION
Using config file: /Users/tstraub/.tfx.hcl
Deleting variable 'AWS_REGION' for variable set 'aws-credentials'
Status: Success
Key:    AWS_REGION
```
