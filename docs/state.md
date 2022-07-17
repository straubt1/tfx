# `tfx state` Commands

Managing Workspace State Files (State Versions).

## `tfx state list`

List all State Versions for a supplied Workspace

## `tfx state show`

Show state details for a supplied State Version

## `tfx state download`

Download a specific State Version

## `tfx state create`

Create a new State Version with a supplied state file
- There is no way to delete State Versions
- The LAST State Version to be created is the "current" state file that will be used by the Workspace
- A Workspace must be locked to create new State Versions
- The "serial" attribute must be incremented
- The "lineage" attribute must be the same for any newly created State Version
- The API does not return a state versions lineage, you must download the file and parse to get the lineage

