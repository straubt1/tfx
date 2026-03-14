# TFx TUI Layout Reference

ASCII diagrams for all views. Edit these diagrams to spec a layout change, then hand to the AI to implement.

---

## Chrome Structure (all views)

```
fixedLines = 9
  Row 1:     header
  Rows 2–5:  profile bar (4 rows, always present — "…" while loading)
  Row 6:     content box top border  ┌─ breadcrumb ──────────────────┐
  Rows 7–N:  content rows            │ ...                           │
  Row N+1:   content box bot border  └───────────────────────────────┘
  Row N+2:   status bar
  Row N+3:   CLI hint bar

  content height = terminal_height - fixedLines
```

---

## Header

```
 TFx   app.terraform.io   ⬥  HCP Terraform                        v0.2.2-local
│─────│─────────────────────────────────────────────────────────│─────────────│
 blue  dim (hostname + remote app/version, omitted on HCP TF)    purple
```

- `⬥  RemoteApp  TFEVersion` segment only appears for TFE instances (not HCP Terraform)
- Gap between remote info and version is filled with `colorHeaderBg`

---

## Profile Bar (4 rows)

```
  profile:   default
  username:  tstraub
  email:     tstraub@hashi.com
  expires:   2027-01-15 (312 days)
```

- Background: `colorHeaderBg`
- Labels (`profile:` etc.) padded to 9 chars, `colorDim`
- Values: `colorAccent` (blue) when loaded, `colorDim` italic `…` while loading
- `expires`: `"never"` when zero, `"n/a"` on error, `"YYYY-MM-DD (N days/hours/minutes)"` otherwise
- Full terminal width, gap filled with `colorHeaderBg`

---

## Content Box Top Border (breadcrumb embedded as title)

```
┌─ organizations ──────────────────────────────────────────────────────────────────────────┐
┌─ org: my-org  /   projects ──────────────────────────────────────────────────────────────┐
┌─ org: my-org  /   project: my-proj  /   workspaces ──────────────────────────────────────┐
┌─ org: my-org  /   project: my-proj  /   workspace: my-ws  /   runs ──────────────────────┐
┌─ org: my-org  /   project: my-proj  /   workspace: my-ws  /   detail ────────────────────┐
┌─ org: my-org  /   project: my-proj  /   workspace: my-ws  /   runs  /   run: run-abc123 ─┐
```

- Active segment: `colorAccent` bold
- Inactive segments: `colorDim`
- Separator `  /  `: `colorDim`
- Border chars `┌ ─ ┐`: `colorDim`

---

## Standard List View (Organizations / Projects / Workspaces)

```
┌─ org: my-org  /   workspaces ────────────────────────────────────────────────┐
│  / filter-text                                                               │  ← filter bar (only when filtering active)
│  NAME              PROJECT         TERRAFORM   LOCKED   UPDATED              │  ← table header (colorHeaderBg bg, colorAccent fg, bold)
│ ──────────────────────────────────────────────────────────────────────────── │  ← accent divider
│    my-workspace    default         1.9.0       false    2025-06-01           │  ← unselected row
│  ▶ prod-infra      networking      1.8.2       true     2025-05-30           │  ← selected row (colorSelected bg, colorAccent fg, bold)
│    staging-api     default         1.7.4       false    2025-05-28           │
│    ...                                                                       │
└──────────────────────────────────────────────────────────────────────────────┘
  42 workspaces                                                                   ← status bar: count (or "N / M  •  filter: text")
  cmd: tfx workspace list   •   enter runs   v vars   d detail   •   q quit      ← CLI hint bar
```

---

## Workspace Sub-Views (Runs / Variables / Config Versions / State Versions)

Tab strip is the first line inside the content box:

```
┌─ org: my-org  /   project: my-proj  /   workspace: my-ws  /   runs ──────────┐
│   Runs   Variables   Config Versions   State Versions                        │  ← tab strip (active: colorAccent bold underline; inactive: colorDim)
│  ID               STATUS     MESSAGE              CREATED          DURATION  │
│ ──────────────────────────────────────────────────────────────────────────── │
│    run-xyz789     applied    terraform apply       2025-06-01       1m 23s   │
│  ▶ run-abc123     planned    speculative plan       2025-05-30       0m 44s  │
│    ...                                                                       │
└──────────────────────────────────────────────────────────────────────────────┘
  8 runs                                                                          ← status bar
  cmd: tfx run list   •   enter detail   •   ? help   •   q quit                 ← CLI hint bar
```

---

## Detail View (Workspace / Org / Project / Run / Variable / State Version / Config Version)

```
┌─ org: my-org  /   project: my-proj  /   workspace: my-ws  /   detail ────────┐
│  Name              prod-infra                                                │
│  ID                ws-ABcd1234EFgh5678                                       │
│  Description       Production networking infrastructure                      │
│  ─────────────────────────────────────────────────────────────────────────── │  ← section divider
│  Terraform Version  1.8.2                                                    │
│  Locked             true                                                     │
│  ...                                                                         │
└──────────────────────────────────────────────────────────────────────────────┘
  workspace: prod-infra  •  ↑ ↓ to scroll                                        ← status bar
  cmd: tfx workspace show   •   ↑ ↓ scroll   u url   U browser   •   q quit      ← CLI hint bar
```

---

## State Version JSON Viewer

```
┌─ org: my-org  /  ...  /   state versions  /   sv: 42  /   json ──────────────┐
│  {                                                                           │
│    "version": 4,                                                             │  ← syntax-highlighted JSON
│    "terraform_version": "1.8.2",                                             │     keys: colorAccent
│    "resources": [                                                            │     strings: colorSuccess
│      ...                                                                     │     numbers: colorPurple
│    ]                                                                         │     true/false/null: colorLoading
│  }                                                                           │     punctuation: colorDim
└──────────────────────────────────────────────────────────────────────────────┘
  state JSON  •  line 1 of 312  (48 KB)                                          ← status bar
  cmd: tfx sv show   •   ↑ ↓ scroll   g top   G bottom   •   q quit              ← CLI hint bar
```

---

## Config Version File Browser

```
┌─ org: my-org  /  ...  /   config versions  /   cv: cv-abc123  /   files ─────┐
│  ├── main.tf                                                                 │
│  ├── variables.tf                                                            │
│  ├── outputs.tf                                                              │
│  └── modules/                                                                │
│      ├── networking/                                                         │
│      │   ├── main.tf                                                         │
│      │   └── variables.tf                                                    │
│      └── compute/                                                            │
│          └── main.tf                                                         │
└──────────────────────────────────────────────────────────────────────────────┘
  config version files  •  ~/.../.tfx-cv-cache/cv-abc123  •  9 files             ← status bar (path is OSC 8 hyperlink)
  cmd: tfx cv show   •   enter view file   •   q quit                            ← CLI hint bar
```

---

## Config Version File Content Viewer

```
┌─ org: my-org  /  ...  /   cv: cv-abc123  /   files  /   main.tf ─────────────┐
│  terraform {                                                                 │
│    required_providers {                                                      │  ← raw file content
│      aws = {                                                                 │     HCL block keywords: colorPurple
│        source  = "hashicorp/aws"                                             │
│        version = "~> 5.0"                                                    │
│      }                                                                       │
│    }                                                                         │
│  }                                                                           │
└──────────────────────────────────────────────────────────────────────────────┘
  main.tf  •  line 1 of 47  (1 KB)                                               ← status bar
  cmd:    •   ↑ ↓ scroll   g top   G bottom   •   q quit                         ← CLI hint bar
```

---

## Loading State

```
┌─ organizations ──────────────────────────────────────────────────────────────┐
│                                                                              │
│                                                                              │
│                          ⣾  Loading…                                         │  ← centered spinner (colorDim italic)
│                                                                              │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
  ⣾  Loading…                                                                    ← status bar (colorLoading)
  cmd: ...
```

---

## API Inspector Panel (`l` key — 33% for the main view/67% for the inspector view horizontal split)

Header, profile bar, status bar, and CLI hint bar remain **full terminal width**.
The content box splits into left (main view) and right (inspector) panels.

```
 TFx   app.terraform.io   ⬥  HCP Terraform                        v0.2.2-local   ← full width
  profile:   default                                                               ← full width
  username:  tstraub
  email:     tstraub@hashi.com
  expires:   2027-01-15 (312 days)
┌─ org: my-org  /   workspaces ─────────────────┬── API Inspector ──────────────────┐  ← split top border
│   NAME              LOCKED   UPDATED          │  GET  /api/v2/organizations  12ms │
│  ───────────────────────────────────────────  │ ▶ POST /api/v2/runs          34ms │  ← selected row in inspector highlighted
│    my-workspace     false    2025-06-01       │  GET  /api/v2/workspaces      8ms │
│  ▶ prod-infra       true     2025-05-30       │  GET  /api/v2/projects        6ms │
│    staging-api      false    2025-05-28       │  ...                              │
└───────────────────────────────────────────────┴───────────────────────────────────┘  ← split bottom border
  42 workspaces                                              [api inspector]           ← status bar (badge right-aligned when focused)
  cmd: tfx workspace list   •   Tab focus inspector   •   l close   •   q quit         ← CLI hint bar
```

### Inspector Detail View (Enter on a call)

```
┌─ org: my-org  /   workspaces ─────────────────┬── API Inspector ──────────────┐
│   NAME              LOCKED   UPDATED          │  POST /api/v2/runs            │  ← title row
│  ...                                          │ ───────────────────────────── │  ← divider
│                                               │  Request Headers              │
│                                               │  Content-Type: application/.. │
│                                               │  Authorization: Bearer [redac]│
│                                               │                               │
│                                               │  Request Body                 │
│                                               │  { "data": { ... } }          │
│                                               │                               │
│                                               │  Response  200  34ms          │
│                                               │  { "data": { ... } }          │
└───────────────────────────────────────────────┴───────────────────────────────┘
  42 workspaces                                        [api inspector › detail]    ← status bar badge changes
```

---

## Status Bar States

| State | Content |
|---|---|
| Loading | `⣾  Loading…` (colorLoading) |
| Error | `✗  <error message>` (colorError) |
| Clipboard feedback | `<message>` (colorSuccess) |
| Normal | `  N organizations / projects / workspaces / runs / etc.` |
| Filtered | `  N / M workspaces  •  filter: text` |
| Detail scroll | `  workspace: name  •  ↑ ↓ to scroll` |
| Inspector focused | left msg + right-aligned `[api inspector]` (colorAccent) |
| Inspector detail | left msg + right-aligned `[api inspector › detail]` |

---

## Color Palette Reference

| Name | Hex | Used for |
|---|---|---|
| `colorBg` | `#0D1117` | content area background |
| `colorFg` | `#E6EDF3` | default text |
| `colorAccent` | `#58A6FF` | selected rows, active breadcrumb, table headers, values |
| `colorDim` | `#8B949E` | labels, inactive items, borders |
| `colorPurple` | `#BC8CFF` | version string, HCL keywords, JSON numbers |
| `colorBorder` | `#30363D` | box-drawing characters |
| `colorHeaderBg` | `#161B22` | header, profile bar, status bar, table header row bg |
| `colorAppBg` | `#1F6FEB` | "TFx" app name badge background |
| `colorSelected` | `#1C2128` | selected table row background |
| `colorError` | `#F85149` | error messages |
| `colorLoading` | `#D29922` | spinner, amber warnings |
| `colorSuccess` | `#3FB950` | clipboard feedback, JSON strings |
