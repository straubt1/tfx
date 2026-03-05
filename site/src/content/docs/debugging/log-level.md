---
title: TFX_LOG
---

`TFX_LOG` enables terminal log output for debugging TFx commands. When set, structured log messages are written to stderr alongside the normal command output. The spinner is automatically disabled when logging is active to prevent output interference.

## Log Levels

Set `TFX_LOG` to one of the following values:

| Level   | Description                                                                 |
|---------|-----------------------------------------------------------------------------|
| `TRACE` | Full HTTP request and response dumps (headers and body), with sensitive headers redacted |
| `DEBUG` | HTTP request and response summaries (method, URL, status code)              |
| `INFO`  | Informational messages such as client initialization and log file paths     |
| `WARN`  | Warning conditions that do not stop execution                               |
| `ERROR` | Error conditions                                                            |
| `NONE`  | Logging disabled (default when `TFX_LOG` is unset)                         |

Levels are inclusive — setting `DEBUG` also shows `INFO`, `WARN`, and `ERROR` messages. Setting `TRACE` shows everything.

:::note
Sensitive headers (`Authorization`, `Cookie`, `X-Api-Key`) are always redacted in log output regardless of level.
:::

## Log Format

Log lines are written to stderr in the following format:

```
[HH:MM:SS] [LEVEL] message key=value key=value ...
```

Each level is color-coded for easy scanning in the terminal:

- `TRACE` — faint/grey
- `DEBUG` — cyan
- `INFO` — green
- `WARN` — yellow
- `ERROR` — red

## Examples

**DEBUG — HTTP request and response summaries**

```sh
$ TFX_LOG=DEBUG tfx project list
[10:22:01] [DEBUG] Logger initialized level=DEBUG
[10:22:01] [DEBUG] HTTP Request method=GET url=https://app.terraform.io/api/v2/organizations/firefly/projects?page%5Bnumber%5D=1&page%5Bsize%5D=100
[10:22:01] [DEBUG] HTTP Response method=GET url=https://app.terraform.io/api/v2/organizations/firefly/projects?page%5Bnumber%5D=1&page%5Bsize%5D=100 status=200 OK statusCode=200
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

**TRACE — Full HTTP dumps with redacted credentials**

```sh
$ TFX_LOG=TRACE tfx project list
[10:22:01] [DEBUG] Logger initialized level=TRACE
[10:22:01] [TRACE] HTTP Request (full dump) request=GET /api/v2/organizations/firefly/projects?page[number]=1&page[size]=100 HTTP/1.1
Host: app.terraform.io
User-Agent: go-tfe/1.x
Authorization: [REDACTED]
Accept: application/vnd.api+json

[10:22:01] [TRACE] HTTP Response (full dump) response=HTTP/2.0 200 OK
Content-Type: application/vnd.api+json
...
{"data":[{"id":"prj-ABC123defGHI789","type":"projects",...}]}
Using config file: /Users/tstraub/.tfx.hcl
List Projects for Organization: firefly
...
```

**INFO — Lifecycle messages only**

```sh
$ TFX_LOG=INFO tfx project list
[10:22:01] [INFO] HTTP logging to file enabled path=/tmp/tfx-logs/tfx_http_20251007_102201_12345.log
Using config file: /Users/tstraub/.tfx.hcl
List Projects for Organization: firefly
...
```

:::note
`TFX_LOG` can be combined with [`TFX_LOG_PATH`](log-path.md) to write full HTTP dumps to a file while also viewing log summaries in the terminal.
:::
