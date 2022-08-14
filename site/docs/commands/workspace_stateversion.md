## State Version Commands

Managing Workspace State Files (State Versions) in a Workspace.

!!! note ""
    All commands below can be used with a `sv` alias.

## `tfx workspace state-version list`

List all State Versions for a supplied Workspace.

## `tfx workspace state-version show`

Show state details for a supplied State Version.

## `tfx workspace state-version download`

Download a specific State Version.

## `tfx workspace state-version create`

Create a new State Version with a supplied state file.
- There is no way to delete State Versions
- The LAST State Version to be created is the "current" state file that will be used by the Workspace
- A Workspace must be locked to create new State Version
- The "serial" attribute must be incremented
- The "lineage" attribute must be the same for any newly created State Version
- The API does not return a state versions lineage, you must download the file and parse to get the lineage
