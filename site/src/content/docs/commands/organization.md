---
title: Organization Commands
---

General commands to manage Organizations.

:::note
All commands below can be used with an `org` alias.
:::

## `tfx organization list`

List all Organizations available to the authenticated user.

Using the `--search` flag allows filtering by organization name with a given string.

**Basic Example**

```sh
$ tfx organization list
Using config file: /Users/tstraub/.tfx.hcl
Listing all organizations
Found 2 Organizations
╭──────────────┬──────────────────────────────────────────┬──────────────────────────╮
│ NAME         │ ID                                       │ EMAIL                    │
├──────────────┼──────────────────────────────────────────┼──────────────────────────┤
│ firefly      │ org-ABC123defGHI789jklM                  │ admin@firefly.example    │
│ acme-corp    │ org-DEF456ghiJKL012mnoPQ                 │ admin@acme.example       │
╰──────────────┴──────────────────────────────────────────┴──────────────────────────╯
```

**Search Example**

```sh
$ tfx organization list --search firefly
Using config file: /Users/tstraub/.tfx.hcl
Listing organizations matching 'firefly'
Found 1 Organizations
╭──────────────┬──────────────────────────────────────────┬──────────────────────────╮
│ NAME         │ ID                                       │ EMAIL                    │
├──────────────┼──────────────────────────────────────────┼──────────────────────────┤
│ firefly      │ org-ABC123defGHI789jklM                  │ admin@firefly.example    │
╰──────────────┴──────────────────────────────────────────┴──────────────────────────╯
```

## `tfx organization show`

Show details of a given Organization.

**Required Flags**

| Flag     | Short | Description                  |
|----------|-------|------------------------------|
| `--name` | `-n`  | Name of the organization.    |

**Example**

```sh
$ tfx organization show --name firefly
Using config file: /Users/tstraub/.tfx.hcl
Showing organization 'firefly'
Name:                     firefly
ID:                       org-ABC123defGHI789jklM
Email:                    admin@firefly.example
Created At:               2022-01-15 09:23:41 +0000 UTC
Collaborator Auth Policy: password
Cost Estimation Enabled:  true
Owners Team SAML Role ID:
SAML Enabled:             false
Session Remember Minutes: 20160
Session Timeout Minutes:  20160
Two Factor Conformant:    false
Trial Expires At:         0001-01-01 00:00:00 +0000 UTC
Default Execution Mode:   remote
Is Unified:               false
Permissions:
  Can Create Team:                true
  Can Create Workspace:           true
  Can Create Workspace Migration: true
  Can Destroy:                    true
  Can Manage Run Tasks:           true
  Can Traverse:                   true
  Can Update:                     true
```
