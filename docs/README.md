# TFx Docs

_tfx_ documentation.



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


