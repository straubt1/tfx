# Workspace Variable Commands

Managing Workspace Variables.

!!! note ""
    All commands below can be used with a `var` alias.

## `tfx workspace variable list`

List all Workspace variables.

**Example**

```sh
$ tfx workspace variable list -w tt-workspace                  
Using config file: /Users/tstraub/.tfx.hcl
List Variables for Workspace: tt-workspace
╭──────────────────────┬───────────┬──────────────┬───────────┬───────┬───────────┬───────────────────────────────────╮
│ ID                   │ KEY       │ VALUE        │ SENSITIVE │ HCL   │ CATEGORY  │ DESCRIPTION                       │
├──────────────────────┼───────────┼──────────────┼───────────┼───────┼───────────┼───────────────────────────────────┤
│ var-ALQUrgeMBDPE9wiy │ variable3 │ It is friday │ false     │ false │ env       │ I am environmental                │
│ var-zrX5efBRKdiRQuUN │ variable2 │              │ true      │ false │ terraform │ I am sensitive                   │
│ var-viy2a1iMKp6Hxgmn │ variable5 │ ./list.hcl   │ false     │ true  │ terraform │ I am a list in a file             │
│ var-DmKBqRHJb34uTfmu │ variable4 │ ./string.hcl │ false     │ true  │ terraform │ I am a string in a file           │
│ var-XvP33JGaRQ3m7FP8 │ variable6 │ ./map.hcl    │ false     │ true  │ terraform │ I am a map in a file              │
│ var-bNCzmaMNtUUDaSzN │ variable1 │ It is friday │ false     │ false │ terraform │ some important info about this... │
╰──────────────────────┴───────────┴──────────────┴───────────┴───────┴───────────┴───────────────────────────────────╯
```

## `tfx workspace variable create`

Create a Workspace variable, optionally the value can read from a filename.

**Basic Example**

```sh
$ tfx workspace variable create -w tt-workspace -k variable1 -v "It is friday" -d "some important info about this..."
Using config file: /Users/tstraub/.tfx.hcl
Create Variable for Workspace: tt-workspace
Variable Created: variable1
ID:          var-bNCzmaMNtUUDaSzN
Key:         variable1
Value:       It is friday
Sensitive:   false
HCL:         false
Category:    terraform
Description: some important info about this...
```

**Sensitive Example**

```sh
$ tfx workspace variable create -w tt-workspace -k variable2 -v "It is friday" -d "I am sensitive" --sensitive
Using config file: /Users/tstraub/.tfx.hcl
Create Variable for Workspace: tt-workspace
Variable Created: variable2
ID:          var-zrX5efBRKdiRQuUN
Key:         variable2
Value:       
Sensitive:   true
HCL:         false
Category:    terraform
Description: I am sensitive.
```

**Environment Variable Example**

```sh
$ tfx workspace variable create -w tt-workspace -k variable3 -v "It is friday" -d "I am environmental" --env
Using config file: /Users/tstraub/.tfx.hcl
Create Variable for Workspace: tt-workspace
Variable Created: variable3
ID:          var-ALQUrgeMBDPE9wiy
Key:         variable3
Value:       It is friday
Sensitive:   false
HCL:         false
Category:    env
Description: I am environmental
```

**HCL String Example**

```sh
$ tfx workspace variable create -w tt-workspace -k variable4 -v ./string.hcl -d "I am a string in a file" --hcl
Using config file: /Users/tstraub/.tfx.hcl
Create Variable for Workspace: tt-workspace
Variable Created: variable4
ID:          var-DmKBqRHJb34uTfmu
Key:         variable4
Value:       ./string.hcl
Sensitive:   false
HCL:         true
Category:    terraform
Description: I am a string in a file
```

**HCL List Example**

```sh
$ tfx workspace variable create -w tt-workspace -k variable5 -v ./list.hcl -d "I am a list in a file" --hcl
Using config file: /Users/tstraub/.tfx.hcl
Create Variable for Workspace: tt-workspace
Variable Created: variable5
ID:          var-viy2a1iMKp6Hxgmn
Key:         variable5
Value:       ./list.hcl
Sensitive:   false
HCL:         true
Category:    terraform
Description: I am a list in a file
```

**HCL Map Example**

```sh
$ tfx workspace variable create -w tt-workspace -k variable6 -v ./map.hcl -d "I am a map in a file" --hcl
Using config file: /Users/tstraub/.tfx.hcl
Create Variable for Workspace: tt-workspace
Variable Created: variable6
ID:          var-XvP33JGaRQ3m7FP8
Key:         variable6
Value:       ./map.hcl
Sensitive:   false
HCL:         true
Category:    terraform
Description: I am a map in a file
```

## `tfx workspace variable update`

Update an existing Workspace Variable, optionally the value can read from a filename.

Variables do not have an "upsert" functionality, so you either need to delete then recreate, or update.

**Example**

```sh
$ tfx workspace variable update -w tt-workspace -k variable1 -v "It is July" -d "(update) I made the mistakes"
Using config file: /Users/tstraub/.tfx.hcl
Update Variable for Workspace: tt-workspace
Variable Updated
ID:          var-bNCzmaMNtUUDaSzN
Key:         variable1
Value:       It is July
Sensitive:   false
HCL:         false
Category:    terraform
Description: (update) I made the mistakes
```

## `tfx workspace variable show`

Show details of a Workspace Variable.

**Example**

```sh
$ 
```

## `tfx workspace variable delete`

Delete a Workspace Variable.

**Example**

```sh
$ tfx workspace variable delete -w tt-workspace --key variable7
Using config file: /Users/tstraub/.tfx.hcl
Delete Variable for Workspace: tt-workspace
Variable Deleted: variable7
Status: Success
```
