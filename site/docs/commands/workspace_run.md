# Workspace Run Commands

Managing Workspace Runs.

## `tfx workspace run list`

List all Runs for a supplied Workspace.

`--max-items` defaults to 10, setting this to a higher number will retrieve more items.

**Example**

```sh
$ tfx workspace run list --workspace-name tt--workspace-nameorkspace
Using config file: /Users/tstraub/.tfx.hcl
List Runs for Workspace: tt--workspace-nameorkspace
╭──────────────────────┬───────────────────────┬──────────────────────┬───────────┬───────────────────┬───────────────────────┬──────────────────────────────────────────────────╮
│ ID                   │ CONFIGURATION VERSION │ STATUS               │ PLAN ONLY │ TERRAFORM VERSION │ CREATED               │ MESSAGE                                          │
├──────────────────────┼───────────────────────┼──────────────────────┼───────────┼───────────────────┼───────────────────────┼──────────────────────────────────────────────────┤
│ run-muJzD4EXcYXeb6aY │ cv-q9yhRwv73u6UFJdq   │ planned_and_finished │ false     │ 1.0.0             │ Sat Aug 20 14:45 2022 │ Run created from the TFx CLI                     │
│ run-xucJV6uRyz1Ntf9u │ cv-q9yhRwv73u6UFJdq   │ planned_and_finished │ false     │ 1.0.0             │ Sat Aug 20 14:44 2022 │ Triggered via API                                │
│ run-GqEv5WwffMRDQew2 │ cv-VYikVwjgfHNnUYfr   │ discarded            │ false     │ 1.0.0             │ Tue Jun 28 17:46 2022 │ Triggered via API                                │
│ run-uUi5cTRXLBeDdHoB │ cv-LcVa2hMVZg1nAX6R   │ discarded            │ false     │ 1.0.0             │ Tue Jun 28 17:44 2022 │ Triggered via API                                │
│ run-8tJTJPZUv24bwxdj │ cv-q9yhRwv73u6UFJdq   │ discarded            │ false     │ 1.0.0             │ Thu Jul 15 18:04 2021 │ Queued manually to destroy infrastructure        │
│ run-HmJRanghKXFCoRpe │ cv-pn7T5L8J58FV5PSZ   │ applied              │ false     │ 0.15.3            │ Thu Jul 15 18:03 2021 │ Queued manually via the Terraform Enterprise API │
│ run-UjgDJwAeinyzzxAX │ cv-p8XXa5rcph3W1MoF   │ applied              │ false     │ 0.15.3            │ Thu Jul 15 18:02 2021 │ Queued manually to destroy infrastructure        │
│ run-yVXxdJ8vav52UwpH │ cv-BGP2Q8WwAM9RfzcN   │ planned_and_finished │ false     │ 0.15.3            │ Fri Jun 25 12:31 2021 │ Queued manually via the Terraform Enterprise API │
╰──────────────────────┴───────────────────────┴──────────────────────┴───────────┴───────────────────┴───────────────────────┴──────────────────────────────────────────────────╯
```

## `tfx workspace run create`

Create a Run for a supplied Workspace.

**Latest Configuration Version Example**

```sh
$ tfx workspace run create --workspace-name tt--workspace-nameorkspace          
Using config file: /Users/tstraub/.tfx.hcl
Create Run for Workspace: tt--workspace-nameorkspace
The run will be created using the workspace's latest configuration version 
Run Created 
ID:                    run-RZntt2QgVmD5w9xa
Configuration Version: cv-e83GeSpjVKXuUGmU
Terraform Version:     1.0.0
Link:                  https://tfe.rocks/app/firefly/workspaces/tt--workspace-nameorkspace/runs/run-RZntt2QgVmD5w9xa
```

**Specific Configuration Version Example**

!!! warning ""
  Executing this command with a specific Configuration Version will result in that Configuration Version to be the **latest**

```sh
$  tfx workspace run create --workspace-name tfx-test --id cv-q9yhRwv73u6UFJdq
Using config file: /Users/tstraub/.tfx.hcl
Create Run for Workspace: tfx-test
Configuration Version Provided: cv-q9yhRwv73u6UFJdq
Run Created 
ID:                    run-Q7cVGhK77dukA41G
Configuration Version: cv-q9yhRwv73u6UFJdq
Terraform Version:     1.0.0
Link:                  https://tfe.rocks/app/firefly/workspaces/tfx-test/runs/run-Q7cVGhK77dukA41G
```

**Message Example**

```sh
$ tfx workspace run create --workspace-name tfx-test --message "Run created from the TFx CLI"
Using config file: /Users/tstraub/.tfx.hcl
Create Run for Workspace: tfx-test
The run will be created using the workspace's latest configuration version 
Run Created 
ID:                    run-muJzD4EXcYXeb6aY
Configuration Version: cv-q9yhRwv73u6UFJdq
Terraform Version:     1.0.0
Link:                  https://tfe.rocks/app/firefly/workspaces/tfx-test/runs/run-muJzD4EXcYXeb6aY
```


## `tfx workspace run show`

Show Run details for a supplied Run.

**Example**

```sh
$ tfx workspace run show -i run-GqEv5WwffMRDQew2 
Using config file: /Users/tstraub/.tfx.hcl
Show Run for Workspace: run-GqEv5WwffMRDQew2
ID:                    run-GqEv5WwffMRDQew2
Configuration Version: cv-VYikVwjgfHNnUYfr
Status:                discarded
Message:               Triggered via API
Terraform Version:     1.0.0
Created:               Tue Jun 28 17:46 2022
```
