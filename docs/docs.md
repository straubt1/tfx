# TFx Docs

TFx is designed to make interacting with Terraform Cloud and Terraform Enterprise easier,

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

## Private Registry

Commands:

- [tfx registry provider](registry_provider.md)
- [tfx pmr](registry_pmr.md)

## Terraform Enterprise Management

Commands:

- [tfx tfv](tfv.md)
- [tfx release](release.md)


