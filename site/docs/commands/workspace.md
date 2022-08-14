# Workspace Commands

General commands to manage Workspaces.

## `tfx workspace list`

Using the `--search` flag allows filtering by workspaces with a given string

Using the `--run-status` flag allows filtering by workspaces with a current run with a given status (full list of available run statuses can be found [here](https://www.terraform.io/docs/cloud/api/run.html#run-states))

## `tfx workspace list all`

Using the "--search" flag allows filtering by workspaces with a given string

## `tfx workspace lock`

Lock a given workspace by name, in a given organization

## `tfx workspace lock all`

Lock all workspaces in a given organization

## `tfx workspace unlock`

Unlock a given workspace by name, in a given organization

## `tfx workspace unlock all`

Unlock all workspaces in a given organization
