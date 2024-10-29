# Workspace Commands

General commands to manage Workspace Locks.

!!! note ""
    All commands below can be used with a `ws` alias.

## `tfx workspace lock`

Lock a given workspace by name, in a given organization.

```sh
$ tfx workspace lock -n tfx-test-workspace-1 
Using config file: /Users/tstraub/.tfx.hcl
Lock Workspace in Organization: firefly
tfx-test-workspace-1: Locked
```

## `tfx workspace lock all`

Lock all workspaces in a given organization (sequentially).

This command will ignore individual errors and attempt to execute on all Workspaces.

```sh
$ tfx workspace lock all  
Using config file: /Users/tstraub/.tfx.hcl
Lock All Workspace in Organization: firefly
Locking 6 Workspaces, please wait...
tfx-test-workspace-01:        Locked
tfx-test-workspace-02:        Locked
tfx-test-workspace-03:        Locked
tfx-test-workspace-04:        Locked
tfx-test-workspace-05:        Locked
tfx-test-workspace-06:        Locked
```

## `tfx workspace unlock`

Unlock a given workspace by name, in a given organization.

```sh
$ tfx git:(tt-additional-refactor) ✗ tfx workspace unlock -n tfx-test-workspace-1
Using config file: /Users/tstraub/.tfx.hcl
Unlock Workspace in Organization: firefly
tfx-test-workspace-1: Unlocked
```

## `tfx workspace unlock all`

Unlock all workspaces in a given organization (sequentially).

This command will ignore individual errors and attempt to execute on all Workspaces.

```sh
$ tfx workspace unlock all  
Using config file: /Users/tstraub/.tfx.hcl
Unlock All Workspace in Organization: firefly
Unlocking 6 Workspaces, please wait...
tfx-test-workspace-01:        Unlocked
tfx-test-workspace-02:        Unlocked
tfx-test-workspace-03:        Unlocked
tfx-test-workspace-04:        Unlocked
tfx-test-workspace-05:        Unlocked
tfx-test-workspace-06:        Unlocked
```
