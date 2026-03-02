---
title: Variable Set Commands
---

General commands to manage Variable Sets.

:::note
All commands below can be used with a `varset` alias.
:::

## `tfx variable-set list`

List Variable Sets available for a given Organization.

Using the `--search` flag allows filtering by variable set name with a given string.

**Basic Example**

```sh
$ tfx variable-set list
Using config file: /Users/tstraub/.tfx.hcl
List Variable Sets for Organization: firefly
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
$ tfx variable-set list --search aws
Using config file: /Users/tstraub/.tfx.hcl
List Variable Sets for Organization: firefly
╭─────────────────┬─────────────────────┬────────┬──────────┬──────────────────────╮
│ NAME            │ ID                  │ GLOBAL │ PRIORITY │ PARENT               │
├─────────────────┼─────────────────────┼────────┼──────────┼──────────────────────┤
│ aws-credentials │ varset-abc123XYZ456 │ false  │ false    │ organization:firefly │
╰─────────────────┴─────────────────────┴────────┴──────────┴──────────────────────╯
```

**Alias Example**

```sh
$ tfx varset list
Using config file: /Users/tstraub/.tfx.hcl
List Variable Sets for Organization: firefly
╭──────────────────────┬─────────────────────┬────────┬──────────┬───────────────────────────────╮
│ NAME                 │ ID                  │ GLOBAL │ PRIORITY │ PARENT                        │
├──────────────────────┼─────────────────────┼────────┼──────────┼───────────────────────────────┤
│ aws-credentials      │ varset-abc123XYZ456 │ false  │ false    │ organization:firefly          │
│ common-tags          │ varset-def789GHI012 │ true   │ false    │ organization:firefly          │
│ production-overrides │ varset-jkl345MNO678 │ false  │ true     │ project:prj-ABC123defGHI789   │
╰──────────────────────┴─────────────────────┴────────┴──────────┴───────────────────────────────╯
```

## `tfx variable-set show`

Show details of a given Variable Set by ID, including assigned workspaces, projects, and variables.

**Example**

```sh
$ tfx variable-set show -i varset-abc123XYZ456
Using config file: /Users/tstraub/.tfx.hcl
Show Variable Set: varset-abc123XYZ456
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

Create a new Variable Set in the Organization.

| Flag | Description | Required |
|---|---|---|
| `--name` / `-n` | Name of the Variable Set | Yes |
| `--description` / `-d` | Description of the Variable Set | No |
| `--global` | Apply this Variable Set to all workspaces in the organization | No |
| `--priority` | Variable values in this set override workspace-level values | No |

**Basic Example**

```sh
$ tfx variable-set create --name aws-credentials --description "AWS credentials for production workloads"
Using config file: /Users/tstraub/.tfx.hcl
Create Variable Set for Organization: firefly
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
Create Variable Set for Organization: firefly
ID:          varset-def789GHI012
Name:        common-tags
Description: Tags applied to all workspaces
Global:      true
Priority:    false
```

**Priority Variable Set Example**

```sh
$ tfx variable-set create --name production-overrides --description "Overrides for production runs" --priority
Using config file: /Users/tstraub/.tfx.hcl
Create Variable Set for Organization: firefly
ID:          varset-jkl345MNO678
Name:        production-overrides
Description: Overrides for production runs
Global:      false
Priority:    true
```

## `tfx variable-set delete`

Delete a Variable Set by ID.

:::caution
This will permanently delete the Variable Set and remove it from all workspaces and projects it is currently assigned to.
:::

**Example**

```sh
$ tfx variable-set delete -i varset-abc123XYZ456
Using config file: /Users/tstraub/.tfx.hcl
Delete Variable Set: varset-abc123XYZ456
Status: Success
ID:     varset-abc123XYZ456
```
