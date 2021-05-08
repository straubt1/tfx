# TFx

## References

https://github.com/hashicorp/go-tfe

https://github.com/spf13/cobra#installing

## Development

VS Code Plugin -> https://marketplace.visualstudio.com/items?itemName=golang.Go

go install github.com/go-delve/delve/cmd/dlv@latest

https://github.com/golang/vscode-go/blob/master/docs/debugging.md#launch-configuration


```
go get -u github.com/spf13/cobra
go get github.com/spf13/cobra/cobra
mkdir tfx && cd tfx
~/go/bin/cobra init --pkg-name github.com/straubt1/tfx
~/go/bin/cobra add plan


export GOPATH=/usr/local/bin/go
```

## Commands

- `tfx run list`
- `tfx cv list`
- `tfx plan`
- `tfx apply`



## Notes

`tfx init` could still be valuable, maybe pull state file locally
