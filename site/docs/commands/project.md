# Project Commands

General commands to manage Projects.

!!! note ""
    All commands below can be used with a `prj` alias.

## `tfx project list`

List projects available for a given Organization.

Using the `--all` flag allows returning projects from all Organizations available to the provided API Token.

Using the `--search` flag allows filtering by project name with a given string.

**Basic Example**

```sh
$ tfx project list
Using config file: /Users/tstraub/.tfx.hcl
List Projects for Organization: firefly
Found 4 Projects
╭─────────────────────────────┬─────────────────────┬──────────────────────────────────────────────╮
│ NAME                        │ ID                  │ DESCRIPTION                                  │
├─────────────────────────────┼─────────────────────┼──────────────────────────────────────────────┤
│ infrastructure-core         │ prj-ABC123defGHI789 │ Core infrastructure components               │
│ application-platform        │ prj-DEF456ghiJKL012 │ Application platform and services            │
│ security-compliance         │ prj-GHI789jklMNO345 │ Security and compliance resources            │
│ development-environments    │ prj-JKL012mnoPQR678 │ Development and testing environments         │
╰─────────────────────────────┴─────────────────────┴──────────────────────────────────────────────╯
```

**Search Example**

```sh
$ tfx project list --search infrastructure
Using config file: /Users/tstraub/.tfx.hcl
List Projects for Organization: firefly
Found 2 Projects
╭─────────────────────────────┬─────────────────────┬──────────────────────────────────────────────╮
│ NAME                        │ ID                  │ DESCRIPTION                                  │
├─────────────────────────────┼─────────────────────┼──────────────────────────────────────────────┤
│ infrastructure-core         │ prj-ABC123defGHI789 │ Core infrastructure components               │
│ infrastructure-monitoring   │ prj-STU901vwxYZA234 │ Infrastructure monitoring and alerting       │
╰─────────────────────────────┴─────────────────────┴──────────────────────────────────────────────╯
```

**List All Example**

```sh
$ tfx project list --all    
Using config file: /Users/tstraub/.tfx.hcl
List Projects for all available Organizations 
Found 12 Projects
╭──────────────┬─────────────────────────────┬─────────────────────┬──────────────────────────────────────────────╮
│ ORGANIZATION │ NAME                        │ ID                  │ DESCRIPTION                                  │
├──────────────┼─────────────────────────────┼─────────────────────┼──────────────────────────────────────────────┤
│ firefly      │ infrastructure-core         │ prj-ABC123defGHI789 │ Core infrastructure components               │
│ firefly      │ application-platform        │ prj-DEF456ghiJKL012 │ Application platform and services            │
│ firefly      │ security-compliance         │ prj-GHI789jklMNO345 │ Security and compliance resources            │
│ firefly      │ development-environments    │ prj-JKL012mnoPQR678 │ Development and testing environments         │
│ acme-corp    │ web-services                │ prj-MNO345pqrSTU789 │ Web application services                     │
│ acme-corp    │ data-platform               │ prj-PQR678stuvWX012 │ Data processing and analytics platform       │
│ acme-corp    │ mobile-backend              │ prj-STU901vwxYZA234 │ Mobile application backend services          │
│ acme-corp    │ ci-cd-pipeline              │ prj-VWX234yzaBCD567 │ Continuous integration and deployment        │
╰──────────────┴─────────────────────────────┴─────────────────────┴──────────────────────────────────────────────╯
```

**List Projects with search across all organizations Example**

```sh
$ tfx project list --all --search platform
Using config file: /Users/tstraub/.tfx.hcl
List Projects for all available Organizations 
Found 2 Projects
╭──────────────┬─────────────────────────────┬─────────────────────┬──────────────────────────────────────────────╮
│ ORGANIZATION │ NAME                        │ ID                  │ DESCRIPTION                                  │
├──────────────┼─────────────────────────────┼─────────────────────┼──────────────────────────────────────────────┤
│ firefly      │ application-platform        │ prj-DEF456ghiJKL012 │ Application platform and services            │
│ acme-corp    │ data-platform               │ prj-PQR678stuvWX012 │ Data processing and analytics platform       │
╰──────────────┴─────────────────────────────┴─────────────────────┴──────────────────────────────────────────────╯
```

## `tfx project show`

Show details of a given Project, including configuration and tags.

**Example**

```sh
$ tfx project show -i prj-ABC123defGHI789          
Using config file: /Users/tstraub/.tfx.hcl
Show Project: prj-ABC123defGHI789
Name: infrastructure-core
ID: prj-ABC123defGHI789
Description: Core infrastructure components
DefaultExecutionMode: remote
Auto Destroy Activity Duration: 7d
Tags:                         
  environment: production
  team: platform
  cost-center: engineering
  compliance: required
```
