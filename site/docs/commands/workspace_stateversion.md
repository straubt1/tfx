# Workspace State Version Commands

Managing Workspace State Files (State Versions) in a Workspace.

!!! note ""
    All commands below can be used with a `sv` alias.

## `tfx workspace state-version list`

List all State Versions for a supplied Workspace.

`--maxItems` defaults to 10, setting this to a higher number will retrieve more items.

**Example**

```sh
$ tfx workspace state-version list -w tfx-test
Using config file: /Users/tstraub/.tfx.hcl
List State Versions for Workspace: tfx-test
╭─────────────────────┬───────────────────┬────────┬──────────────────────┬───────────────────────╮
│ ID                  │ TERRAFORM VERSION │ SERIAL │ RUN ID               │ CREATED               │
├─────────────────────┼───────────────────┼────────┼──────────────────────┼───────────────────────┤
│ sv-eoYznk6PbJY1o9XY │ 0.15.3            │     21 │ run-HmJRanghKXFCoRpe │ Thu Jul 15 18:03 2021 │
│ sv-VfpmiWSw5NUVWe5W │ 0.15.3            │     20 │ run-UjgDJwAeinyzzxAX │ Thu Jul 15 18:02 2021 │
│ sv-PQrMeSHA5DApQWyD │ 0.15.3            │     19 │ run-JMVCbHt6QYGKSpjS │ Thu Jul 15 18:01 2021 │
│ sv-eeX8tPfiEiUCgRsj │ 0.15.3            │     18 │ run-31xdejiW9JyLjkKz │ Thu Jul 15 17:58 2021 │
│ sv-HLWLncRuKwXkXHvz │ 0.15.3            │     17 │ run-tiC3MEGYbuueUg2X │ Thu Jun 24 22:21 2021 │
│ sv-zdZruvurj7K2GYpL │ 0.15.3            │     16 │ run-8VTakdYndfsBEfdY │ Wed Jun 23 12:33 2021 │
│ sv-Jdx81MWz2NVLCQNY │ 0.15.3            │     15 │ run-AbJ8yAgbgdBVhXkA │ Sun May 23 18:42 2021 │
│ sv-RVbc8e8qQkhn1s6e │ 0.15.3            │     14 │                      │ Sun May 23 18:40 2021 │
│ sv-EUYx9TPGi8BkSySL │ 0.15.3            │     13 │ run-ZWEtD3KWuur1rKdu │ Sat May 22 17:54 2021 │
│ sv-NAjiA8UvuFe5oUPb │ 0.15.3            │     12 │                      │ Sat May 22 17:54 2021 │
╰─────────────────────┴───────────────────┴────────┴──────────────────────┴───────────────────────╯
```

## `tfx workspace state-version show`

Show state details for a supplied State Version.

**Example**

```sh
$ tfx workspace state-version show --stateId sv-VfpmiWSw5NUVWe5W
Using config file: /Users/tstraub/.tfx.hcl
Show State Version for Workspace from Id: sv-VfpmiWSw5NUVWe5W
ID:                sv-VfpmiWSw5NUVWe5W
Created:           Thu Jul 15 18:02 2021
Terraform Version: 0.15.3
Serial:            20
State Version:     4
Run Id:            run-UjgDJwAeinyzzxAX
```


## `tfx workspace state-version download`

Download a specific State Version.

**Temp Folder Example**

```sh
$ tfx workspace state-version download --stateId sv-VfpmiWSw5NUVWe5W
Using config file: /Users/tstraub/.tfx.hcl
Directory not supplied, creating a temp directory 
Downloading State Version from Id: sv-VfpmiWSw5NUVWe5W
State Version Found, download started... 
Status: Success
File:   /var/folders/99/srh_6psj6g5520gwyv8v3nbw0000gn/T/slug3100435901/sv-VfpmiWSw5NUVWe5W.state
```


## `tfx workspace state-version create`

Create a new State Version with a supplied state file.
- There is no way to delete State Versions
- The LAST State Version to be created is the "current" state file that will be used by the Workspace
- A Workspace must be locked to create new State Version
- The "serial" attribute must be incremented
- The "lineage" attribute must be the same for any newly created State Version
- The API does not return a state versions lineage, you must download the file and parse to get the lineage

**Example**

```sh
$ 
```

