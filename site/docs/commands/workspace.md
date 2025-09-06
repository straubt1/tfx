# Workspace Commands

General commands to manage Workspaces.

!!! note ""
    All commands below can be used with a `ws` alias.

## `tfx workspace list`

List workspaces available for a given Organization.

Using the `--all` flag allows returning to all Organizations available to the provided API Token.

Additionally there are several filters available, these can be used independently or in combination.

Using the `--search` flag allows filtering by workspace name with a given string.

Using the `--repository` flag allows filtering by workspaces connected to a given repository Identifier (example: straubt1/tfe-workspace-seed).

Using the `--run-status` flag allows filtering by workspaces with a current run with a given status (full list of available run statuses can be found [here](https://www.terraform.io/docs/cloud/api/run.html#run-states)).

Using the `--project-id` flags allows filtering by Project.

**Basic Example**

```sh
$ tfx workspace list
Using config file: /Users/tstraub/.tfx.hcl
List Workspaces for Organization: firefly
Found 6 Workspaces
╭─────────────────────────────┬─────────────────────┬────────────────┬───────────────────────┬──────────────────────┬───────────────┬────────╮
│ NAME                        │ ID                  │ RESOURCE COUNT │ CURRENT RUN CREATED   │ STATUS               │ REPOSITORY    │ LOCKED │
├─────────────────────────────┼─────────────────────┼────────────────┼───────────────────────┼──────────────────────┼───────────────┼────────┤
│ tfx-test-workspace-01       │ ws-hLFv8c9bjgXC3mdK │              5 │ Fri Oct 30 13:30 2020 │ planned_and_finished │ tstraub/demo1 │ false  │
│ tfx-test-workspace-02       │ ws-yN6DnhYxB39qqAre │              0 │                       │                      │               │ false  │
│ tfx-test-workspace-03       │ ws-3NhtTDoX6pqZguuB │              3 │                       │                      │               │ false  │
│ tfx-test-workspace-04       │ ws-uhDDVjE6Q1WxwU5C │              7 │                       │                      │               │ false  │
│ tfx-test-workspace-05       │ ws-yra8oTuc16pgYedk │              2 │                       │                      │               │ false  │
│ tfx-test-workspace-06       │ ws-qsLatjFsibCPAKWr │              1 │                       │                      │               │ false  │
╰─────────────────────────────┴─────────────────────┴────────────────┴───────────────────────┴──────────────────────┴───────────────┴────────╯
```

**Search Example**

```sh
$ tfx workspace list --search aws 
Using config file: /Users/tstraub/.tfx.hcl
List Workspaces for Organization: firefly
Found 4 Workspaces
╭───────────────────┬─────────────────────┬────────────────┬───────────────────────┬────────────────┬───────────────┬────────╮
│ NAME              │ ID                  │ RESOURCE COUNT │ CURRENT RUN CREATED   │ STATUS         │ REPOSITORY    │ LOCKED │
├───────────────────┼─────────────────────┼────────────────┼───────────────────────┼────────────────┼───────────────┼────────┤
│ aws-dev-uswest-1  │ ws-XXn8hDRGA56Wyzxe │             12 │ Fri Oct 30 13:39 2022 │ policy_checked │ tstraub/demo1 │ true   │
│ aws-prod-uswest-1 │ ws-Trm11JYZz9dj46wT │              8 │ Fri Oct 30 13:30 2022 │ errored        │ tstraub/demo2 │ true   │
│ aws-dev-uswest-1  │ ws-BUBSQysttH1FGLqr │              4 │ Fri Oct 30 13:34 2022 │ policy_checked │ tstraub/demo1 │ true   │
│ aws-prod-uswest-2 │ ws-ZWNdqJLrWzHEeevS │             15 │ Fri Oct 30 13:22 2022 │ errored        │ tstraub/demo2 │ true   │
╰───────────────────┴─────────────────────┴────────────────┴───────────────────────┴────────────────┴───────────────┴────────╯
```

**List All Example**

```sh
$ tfx workspace list --all    
Using config file: /Users/tstraub/.tfx.hcl
List Workspaces for all available Organizations 
Found 141 Workspaces
╭──────────────┬───────────────────────┬─────────────────────┬────────────────┬───────────────────────┬──────────────────────┬───────────────┬────────╮
│ ORGANIZATION │ NAME                  │ ID                  │ RESOURCE COUNT │ CURRENT RUN CREATED   │ STATUS               │ REPOSITORY    │ LOCKED │
├──────────────┼───────────────────────┼─────────────────────┼────────────────┼───────────────────────┼──────────────────────┼───────────────┼────────┤
│ firefly      │ tfx-test-workspace-01 │ ws-hLFv8c9bjgXC3mdK │              5 │ Fri Oct 30 13:30 2020 │ planned_and_finished │ tstraub/demo1 │ false  │
│ firefly      │ tfx-test-workspace-02 │ ws-yN6DnhYxB39qqAre │              0 │                       │                      │               │ false  │
│ firefly      │ tfx-test-workspace-03 │ ws-3NhtTDoX6pqZguuB │              3 │                       │                      │               │ false  │
│ firefly      │ tfx-test-workspace-04 │ ws-uhDDVjE6Q1WxwU5C │              7 │                       │                      │               │ false  │
│ firefly      │ tfx-test-workspace-05 │ ws-yra8oTuc16pgYedk │              2 │                       │                      │               │ false  │
│ firefly      │ tfx-test-workspace-06 │ ws-qsLatjFsibCPAKWr │              1 │                       │                      │               │ false  │
╰──────────────┴───────────────────────┴─────────────────────┴────────────────┴───────────────────────┴──────────────────────┴───────────────┴────────╯
```

**List Workspaces with a repository Example**

```sh
$ tfx workspace list --repository straubt1/demo  
Using config file: /Users/tstraub/.tfx.hcl
List Workspaces for Organization: firefly
Found 6 Workspaces
╭───────────────────┬─────────────────────┬────────────────┬───────────────────────┬────────────────┬───────────────┬────────╮
│ NAME              │ ID                  │ RESOURCE COUNT │ CURRENT RUN CREATED   │ STATUS         │ REPOSITORY    │ LOCKED │
├───────────────────┼─────────────────────┼────────────────┼───────────────────────┼────────────────┼───────────────┼────────┤
│ tf-dev-us-west-1  │ ws-XXn8hDRGA56Wyzxe │             10 │ Fri Oct 30 13:39 2022 │ errored        │ straubt1/demo │ true   │
│ tf-stg-us-west-1  │ ws-Trm11JYZz9dj46wT │              8 │ Fri Oct 30 13:30 2022 │ policy_checked │ straubt1/demo │ true   │
│ tf-prd-us-west-1  │ ws-BUBSQysttH1FGLqr │             15 │ Fri Oct 30 13:30 2022 │ policy_checked │ straubt1/demo │ true   │
│ tf-dev-us-west-2  │ ws-ZWNdqJLrWzHEeevS │              6 │ Fri Oct 30 13:30 2022 │ policy_checked │ straubt1/demo │ true   │
│ tf-stg-us-west-2  │ ws-ZPgFWuQ79qFeX7MD │             12 │ Fri Oct 30 13:30 2022 │ policy_checked │ straubt1/demo │ true   │
│ tf-prd-us-west-2  │ ws-axCg6rXwKByChA6f │             20 │ Fri Oct 30 13:30 2022 │ policy_checked │ straubt1/demo │ true   │
╰───────────────────┴─────────────────────┴────────────────┴───────────────────────┴────────────────┴───────────────┴────────╯
```

**List Workspaces with a repository and a run status Example**

```sh
$ tfx workspace list --repository straubt1/demo --run-status errored
Using config file: /Users/tstraub/.tfx.hcl
List Workspaces for Organization: firefly
Found 1 Workspaces
╭───────────────────┬─────────────────────┬────────────────┬───────────────────────┬────────────────┬───────────────┬────────╮
│ NAME              │ ID                  │ RESOURCE COUNT │ CURRENT RUN CREATED   │ STATUS         │ REPOSITORY    │ LOCKED │
├───────────────────┼─────────────────────┼────────────────┼───────────────────────┼────────────────┼───────────────┼────────┤
│ tf-dev-us-west-1  │ ws-XXn8hDRGA56Wyzxe │             10 │ Fri Oct 30 13:39 2022 │ errored        │ straubt1/demo │ true   │
╰───────────────────┴─────────────────────┴────────────────┴───────────────────────┴────────────────┴───────────────┴────────╯
```

## `tfx workspace show`

Show details of a given Workspace, include Team Access and State sharing.

**Example**

```sh
$ tfx workspace show -n tfx-test          
Using config file: /Users/tstraub/.tfx.hcl
Show Workspace: tfx-test
ID:                   ws-VxepewkunumUbR9V
Resource Count:       8
Terraform Version:    1.0.0
Execution Mode:       remote
Auto Apply:           false
Working Directory:    
Locked:               false
Global State Sharing: false
Current Run Id:       run-tNGxao7zMos5YrY1
Current Run Status:   errored
Current Run Created:  Sun Aug 21 16:40 2022
Team Access:         
  appteam-read
  ws-outputs
Remote State Sharing:
  tfx-test-workspace-16
  tfx-test-workspace-17
```
