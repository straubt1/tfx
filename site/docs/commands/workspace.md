# Workspace Commands

General commands to manage Workspaces.

!!! note ""
    All commands below can be used with a `ws` alias.

## `tfx workspace list`

Using the `--search` flag allows filtering by workspaces with a given string.

Using the `--run-status` flag allows filtering by workspaces with a current run with a given status (full list of available run statuses can be found [here](https://www.terraform.io/docs/cloud/api/run.html#run-states)).

**Basic Example**

```sh
$ tfx workspace list
Using config file: /Users/tstraub/.tfx.hcl
List Workspaces for Organization: firefly
Found 6 Workspaces
╭─────────────────────────────┬─────────────────────┬───────────────────────┬──────────────────────┬───────────────┬────────╮
│ NAME                        │ ID                  │ CURRENT RUN CREATED   │ STATUS               │ REPOSITORY    │ LOCKED │
├─────────────────────────────┼─────────────────────┼───────────────────────┼──────────────────────┼───────────────┼────────┤
│ tfx-test-workspace-01       │ ws-hLFv8c9bjgXC3mdK │ Fri Oct 30 13:30 2020 │ planned_and_finished │ tstraub/demo1 │ false  │
│ tfx-test-workspace-02       │ ws-yN6DnhYxB39qqAre │                       │                      │               │ false  │
│ tfx-test-workspace-03       │ ws-3NhtTDoX6pqZguuB │                       │                      │               │ false  │
│ tfx-test-workspace-04       │ ws-uhDDVjE6Q1WxwU5C │                       │                      │               │ false  │
│ tfx-test-workspace-05       │ ws-yra8oTuc16pgYedk │                       │                      │               │ false  │
│ tfx-test-workspace-06       │ ws-qsLatjFsibCPAKWr │                       │                      │               │ false  │
╰─────────────────────────────┴─────────────────────┴───────────────────────┴──────────────────────┴───────────────┴────────╯
```

**Search Example**

```sh
$ tfx workspace list --search aws 
Using config file: /Users/tstraub/.tfx.hcl
List Workspaces for Organization: firefly
Found 4 Workspaces
╭───────────────────┬─────────────────────┬───────────────────────┬────────────────┬───────────────┬────────╮
│ NAME              │ ID                  │ CURRENT RUN CREATED   │ STATUS         │ REPOSITORY    │ LOCKED │
├───────────────────┼─────────────────────┼───────────────────────┼────────────────┼───────────────┼────────┤
│ aws-dev-uswest-1  │ ws-XXn8hDRGA56Wyzxe │ Fri Oct 30 13:39 2022 │ policy_checked │ tstraub/demo1 │ true   │
│ aws-prod-uswest-1 │ ws-Trm11JYZz9dj46wT │ Fri Oct 30 13:30 2022 │ errored        │ tstraub/demo2 │ true   │
│ aws-dev-uswest-1  │ ws-BUBSQysttH1FGLqr │ Fri Oct 30 13:34 2022 │ policy_checked │ tstraub/demo1 │ true   │
│ aws-prod-uswest-2 │ ws-ZWNdqJLrWzHEeevS │ Fri Oct 30 13:22 2022 │ errored        │ tstraub/demo2 │ true   │
╰───────────────────┴─────────────────────┴───────────────────────┴────────────────┴───────────────┴────────╯
```

**List All Example**

```sh
$ tfx workspace list --all    
Using config file: /Users/tstraub/.tfx.hcl
List Workspaces for all available Organizations 
Found 141 Workspaces
╭──────────────┬───────────────────────┬─────────────────────┬───────────────────────┬──────────────────────┬───────────────┬────────╮
│ ORGANIZATION │ NAME                  │ ID                  │ CURRENT RUN CREATED   │ STATUS               │ REPOSITORY    │ LOCKED │
├──────────────┼───────────────────────┼─────────────────────┼───────────────────────┼──────────────────────┼───────────────┼────────┤
│ firefly      │ tfx-test-workspace-01 │ ws-hLFv8c9bjgXC3mdK │ Fri Oct 30 13:30 2020 │ planned_and_finished │ tstraub/demo1 │ false  │
│ firefly      │ tfx-test-workspace-02 │ ws-yN6DnhYxB39qqAre │                       │                      │               │ false  │
│ firefly      │ tfx-test-workspace-03 │ ws-3NhtTDoX6pqZguuB │                       │                      │               │ false  │
│ firefly      │ tfx-test-workspace-04 │ ws-uhDDVjE6Q1WxwU5C │                       │                      │               │ false  │
│ firefly      │ tfx-test-workspace-05 │ ws-yra8oTuc16pgYedk │                       │                      │               │ false  │
│ firefly      │ tfx-test-workspace-06 │ ws-qsLatjFsibCPAKWr │                       │                      │               │ false  │
╰──────────────┴───────────────────────┴─────────────────────┴───────────────────────┴──────────────────────┴───────────────┴────────╯
```

## `tfx workspace show`

Show details of a given Workspace, include Team Access and State sharing.

**Example**

```sh
$ tfx workspace show -n tfx-test
Using config file: /Users/tstraub/.tfx.hcl
Show Workspace: tfx-test
ID:                   ws-VxepewkunumUbR9V
Terraform Version:    1.0.0
Execution Mode:       remote
Auto Apply:           false
Working Directory:    
Locked:               false
Global State Sharing: false
Current Run Id:       run-muJzD4EXcYXeb6aY
Current Run Status:   planned_and_finished
Current Run Created:  Sat Aug 20 14:45 2022
Team Access:          appteam-read,appteam-custom
Remote State Sharing: tfx-test-workspace-16,tfx-test-workspace-17
```
