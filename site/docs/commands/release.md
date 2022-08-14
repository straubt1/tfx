## Release Commands

Manage <a href="https://www.terraform.io/enterprise" target="_blank">Terraform Enterprise</a> and <a href="https://www.replicated.com/" target="_blank">Replicated</a> releases and binaries needed for airgap installations.

These commands are typically used for those running the Self-Managed Terraform Enterprise platform, and do not apply to Terraform Cloud.

> Note: These commands do not communicate with Terraform Enterprise but leverage the replicated REST endpoints.

Authentication does not use an API Token, rather a "License Id" and "Password", these values are given to Customers after purchasing Terraform Enterprise.

### `tfx release tfe list`

List available Terraform Enterprise releases.

`--maxResults` defaults to 10, setting this to a higher number will retrieve more releases.

> Note: Only releases you are licensed for will show up in this list (typically starting with the latest available version on the date of the purchase).

**Example**

```sh
$ tfx release tfe list
Using config file: /Users/tstraub/.tfx.hcl
List Available Terraform Enterprise Releases 
╭──────────┬───────────┬──────────┬───────────────────────╮
│ SEQUENCE │ LABEL     │ REQUIRED │ RELEASE DATE          │
├──────────┼───────────┼──────────┼───────────────────────┤
│      651 │ v202208-2 │ false    │ Fri Aug 12 16:12 2022 │
│      647 │ v202208-1 │ false    │ Wed Aug 10 20:57 2022 │
│      642 │ v202207-2 │ true     │ Fri Jul 15 19:07 2022 │
│      641 │ v202207-1 │ false    │ Wed Jul 13 18:29 2022 │
│      636 │ v202206-1 │ false    │ Wed Jun 15 17:06 2022 │
│      619 │ v202205-1 │ false    │ Tue May 17 19:58 2022 │
│      610 │ v202204-2 │ true     │ Wed Apr 20 22:47 2022 │
│      609 │ v202204-1 │ false    │ Tue Apr 19 16:10 2022 │
│      607 │ v202203-1 │ false    │ Wed Mar 23 18:22 2022 │
│      599 │ v202202-1 │ false    │ Wed Feb 23 18:18 2022 │
╰──────────┴───────────┴──────────┴───────────────────────╯
```

### `tfx release tfe show`

Show details of a Terraform Enterprise release, including release notes.

**Example:**

```sh
$ tfx release tfe show -r 651
Using config file: /Users/tstraub/.tfx.hcl
Show Release details for Terraform Enterprise: 651
Release Sequence: 651
Label:            v202208-2
Release Date:     Fri Aug 12 16:12 2022
Required:         false
Release Notes:    
# TFE Release v202208-2

CHANGES SINCE v202208-1:

<redacted>
```

### `tfx release tfe download`

Download a Terraform Enterprise airgap binary to a directory of your choice (defaults to local directory).

> Note: This file is at least 1GB in size, this command can take a while, but a status will print progress.

**Example:**

```sh
$ tfx release tfe download -r 651
Using config file: /Users/tstraub/.tfx.hcl
Download Release binary for Terraform Enterprise: 651
Downloading from URL: <redacted>
Download Started: ./tfe-651.release
 Download Status: (1.10%) of 1.3G
 Download Status: (2.55%) of 1.3G
 Download Status: (3.83%) of 1.3G
 Download Status: (5.09%) of 1.3G
 ...
 Download Status: (98.09%) of 1.3G
 Download Status: (99.48%) of 1.3G
Release Downloaded! 
```

### `tfx release replicated list`

List available Replicated releases.

`--maxResults` defaults to 10, setting this to a higher number will retrieve more releases.

> Note: This file is at least 1GB in size, this command can take a while, but a status will print progress.

**Example:**

```sh
$ tfx release replicated list
Using config file: /Users/tstraub/.tfx.hcl
List Available Replicated Releases 
╭─────────┬─────────────────────────────────╮
│ VERSION │ PUBLISHED DATE                  │
├─────────┼─────────────────────────────────┤
│ 2.53.7  │ Fri, 24 Jun 2022 00:00:00 -0800 │
│ 2.53.6  │ Thu, 31 Mar 2022 00:00:00 -0800 │
│ 2.53.5  │ Thu, 17 Mar 2022 00:00:00 -0800 │
│ 2.53.4  │ Sat, 12 Mar 2022 00:00:00 -0800 │
│ 2.53.3  │ Fri, 11 Mar 2022 01:00:00 -0800 │
│ 2.53.2  │ Fri, 19 Nov 2021 12:00:00 -0800 │
│ 2.53.1  │ Thu, 23 Sep 2021 12:00:00 -0800 │
│ 2.53.0  │ Tue, 10 Aug 2021 12:00:00 -0800 │
│ 2.52.0  │ Thu, 20 May 2021 12:00:00 -0800 │
│ 2.51.3  │ Tue, 30 Mar 2021 12:00:00 -0800 │
╰─────────┴─────────────────────────────────╯
```

### `tfx release replicated download`

Download a Replicated release binary to a directory of your choice (defaults to local directory).

**Example:**

```sh
$tfx release replicated download -v 2.53.7
Using config file: /Users/tstraub/.tfx.hcl
Download Release binary for Replicated: 2.53.7
Downloading from URL: https://s3.amazonaws.com/replicated-airgap-work/stable/replicated-2.53.7%2B2.53.7%2B2.53.7.tar.gz
Download Started: ./replicated-2.53.7.targz
 Download Status: (1.65%) of 1G
 Download Status: (3.57%) of 1G
 Download Status: (5.27%) of 1G
 Download Status: (7.09%) of 1G
 ...
 Download Status: (95.65%) of 1G
 Download Status: (97.07%) of 1G
 Download Status: (98.83%) of 1G
Release Downloaded! 
```
