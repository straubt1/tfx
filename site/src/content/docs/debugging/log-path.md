---
title: TFX_LOG_PATH
---

`TFX_LOG_PATH` enables HTTP request and response logging to a file. When set to a directory path, TFx writes a timestamped log file containing the full details of every API call made during the command.

This is useful for auditing, sharing debug output with others, or capturing traces for issue reports without needing to read log output in the terminal.

## Behavior

- The directory specified by `TFX_LOG_PATH` is created automatically if it does not exist.
- A new log file is created for each TFx invocation, named using the current timestamp and process ID:
  ```
  tfx_http_YYYYMMDD_HHMMSS_<pid>.log
  ```
- All HTTP requests and responses are written to the file, including headers and body.
- Sensitive headers (`Authorization`, `Cookie`, `X-Api-Key`) are redacted and replaced with `[REDACTED]`.

:::note
`TFX_LOG_PATH` works independently of [`TFX_LOG`](log-level.md). You can use either or both at the same time.
:::

## Log File Format

Each log file begins with a header, followed by request/response pairs separated by dividers:

```
################################################################################
# TFX HTTP LOG - Started at 2025-10-07T10:22:01Z
################################################################################

================================================================================
REQUEST @ 2025-10-07T10:22:01Z
================================================================================
GET /api/v2/organizations/firefly/projects?page[number]=1&page[size]=100 HTTP/1.1
Host: app.terraform.io
User-Agent: go-tfe/1.x
Authorization: [REDACTED]
Accept: application/vnd.api+json

--------------------------------------------------------------------------------
RESPONSE @ 2025-10-07T10:22:02Z
--------------------------------------------------------------------------------
HTTP/2.0 200 OK
Content-Type: application/vnd.api+json

{"data":[{"id":"prj-ABC123defGHI789","type":"projects","attributes":{"name":"infrastructure-core",...}},...]}
```

If an HTTP transport error occurs (e.g. connection refused), an error entry is written instead:

```
*** ERROR @ 2025-10-07T10:22:01Z ***
dial tcp: connection refused
```

## Examples

**Write HTTP logs to a directory**

```sh
$ TFX_LOG_PATH=/tmp/tfx-logs tfx project list
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

$ ls /tmp/tfx-logs/
tfx_http_20251007_102201_84321.log
```

**Combine with TFX_LOG for terminal summaries and full file logs**

```sh
$ TFX_LOG=DEBUG TFX_LOG_PATH=/tmp/tfx-logs tfx workspace list
[10:22:01] [DEBUG] Logger initialized level=DEBUG
[10:22:01] [INFO]  HTTP logging to file enabled path=/tmp/tfx-logs/tfx_http_20251007_102201_84322.log
[10:22:01] [DEBUG] HTTP Request method=GET url=https://app.terraform.io/api/v2/organizations/firefly/workspaces?page%5Bnumber%5D=1&page%5Bsize%5D=100
[10:22:01] [DEBUG] HTTP Response method=GET url=https://app.terraform.io/api/v2/organizations/firefly/workspaces?page%5Bnumber%5D=1&page%5Bsize%5D=100 status=200 OK statusCode=200
Using config file: /Users/tstraub/.tfx.hcl
List Workspaces for Organization: firefly
...
```
