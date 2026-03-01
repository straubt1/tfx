# Workspace Plan Commands

Managing Workspace Plans.

## `tfx workspace plan show`

Show Plan details for a supplied Plan ID.

**Example**

```sh
$ tfx workspace plan show --id plan-CKNawhfgSGdJoGPx
Using config file: /Users/tstraub/.tfx.hcl
Show Plan
ID:     plan-CKNawhfgSGdJoGPx
Status: finished
```

## `tfx workspace plan logs`

Stream logs for a supplied Plan ID.

**Example**

```sh
$ tfx workspace plan show --id plan-CKNawhfgSGdJoGPx
Using config file: /Users/tstraub/.tfx.hcl
Plan Logs for: plan-CKNawhfgSGdJoGPx
...
```

## `tfx workspace plan jsonoutput`

Show JSON execution output for a supplied Plan ID.

**Example**

```sh
$ tfx workspace plan jsonoutput --id plan-CKNawhfgSGdJoGPx
Using config file: /Users/tstraub/.tfx.hcl
Plan JSON Output for: plan-CKNawhfgSGdJoGPx
...
```

## `tfx workspace plan create`

Create a Plan for a supplied Workspace.

Optionally supply a directory of Terraform configuration, a specific Configuration Version ID, a run message, or flags to perform a speculative or destroy plan.

**Basic Example**

```sh
$ tfx workspace plan create --name tt-workspace
Using config file: /Users/tstraub/.tfx.hcl
Create Plan for Workspace: tt-workspace
...
```

**Speculative Plan Example**

```sh
$ tfx workspace plan create --name tt-workspace --speculative
Using config file: /Users/tstraub/.tfx.hcl
Create Plan for Workspace: tt-workspace
...
```

**Destroy Plan Example**

```sh
$ tfx workspace plan create --name tt-workspace --destroy
Using config file: /Users/tstraub/.tfx.hcl
Create Plan for Workspace: tt-workspace
...
```
