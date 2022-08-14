# TFx Docs

_tfx_ documentation.

## Command Output

Most commands support a `--json` (or `-j` for short) flag that will return a proper JSON response.

> Note: Not all commands today support this flag and will ignore it.

**Default Output:**

```sh
$ tfx variable list -w tfx-test               
Using config file: /Users/tstraub/.tfx.hcl
List Variables for Workspace: tfx-test
╭──────────────────────┬───────────┬────────────────┬───────────┬───────┬───────────┬──────────────────────╮
│ ID                   │ KEY       │ VALUE          │ SENSITIVE │ HCL   │ CATEGORY  │ DESCRIPTION          │
├──────────────────────┼───────────┼────────────────┼───────────┼───────┼───────────┼──────────────────────┤
│ var-7XYNuuo4tMjXeXG4 │ variable7 │ {              │ false     │ true  │ terraform │ I am a map in a file │
│                      │           │   "a" = "1",   │           │       │           │                      │
│                      │           │   "b" = "zoo", │           │       │           │                      │
│                      │           │   "c" = "42"   │           │       │           │                      │
│                      │           │ }              │           │       │           │                      │
│ var-MJaLJ7czxKuU48eu │ variable3 │ It is friday   │ false     │ false │ env       │ I am environmental   │
╰──────────────────────┴───────────┴────────────────┴───────────┴───────┴───────────┴──────────────────────╯
```

**JSON Output:**

```sh
$ tfx variable list -w tfx-test --json | jq .
[
  {
    "Category": "terraform",
    "Description": "I am a map in a file",
    "HCL": true,
    "Id": "var-7XYNuuo4tMjXeXG4",
    "Key": "variable7",
    "Sensitive": false,
    "Value": "{\n  \"a\" = \"1\",\n  \"b\" = \"zoo\",\n  \"c\" = \"42\"\n}"
  },
  {
    "Category": "env",
    "Description": "I am environmental",
    "HCL": false,
    "Id": "var-MJaLJ7czxKuU48eu",
    "Key": "variable3",
    "Sensitive": false,
    "Value": "It is friday"
  }
]
```

## API Workflow Commands

Commands:

- [tfx plan](plan.md)
- [tfx apply](apply.md)
- [tfx run](run.md)

### Workspace Run Workflow Example

**Create a Plan**

```sh
# Create a speculative plan that can not be applied
tfx plan -w tfx-test -s

# Create a plan that can be applied
tfx plan -w tfx-test

# Create a Configuration Version based on terraform in the current directory
tfx cv create -w tfx-test

# Create a Configuration Version based on terraform in a supplied directory
tfx cv create -w tfx-test -d ./myterraformfolder/

# Create a plan based on a configuration version
tfx plan -w tfx-test -i cv-HKE8gevVtGBXapcq
```

**Create an Apply**

```sh
tfx apply -r <run-id>
```

## Workspace Commands

Commands:

- [tfx cv](cv.md)
- [tfx state](state.md)
- [tfx variable](variable.md)
- [tfx workspace](workspace.md)

## Private Registry Commands

Commands:

- [tfx registry provider](registry_provider.md)
- [tfx pmr](registry_pmr.md)

## Terraform Enterprise Management Commands

Commands:

- [tfx tfv](tfv.md)
- [tfx release](release.md)


