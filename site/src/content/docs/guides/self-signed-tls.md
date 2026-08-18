---
title: Self-Signed TLS
description: Skip TLS verification for local Terraform Enterprise development and testing when certificates are not in your trust store.
---

Production **Terraform Enterprise** and **HCP Terraform** deployments typically use certificates your system already trusts — no extra configuration is required.

This page is for a narrower case: **local development and testing**, where TFE may run with a self-signed certificate or a private CA that is not installed on your machine (for example a `local.tfe.rocks` dev instance). In that situation, TFx API calls can fail with:

```
tls: failed to verify certificate: x509: certificate signed by unknown authority
```

TFx can skip TLS certificate verification in those environments. **This is off by default** — you must opt in explicitly.

:::danger[Security implications]
When `ssl_skip_verify` is enabled, TFx **does not validate** the server's identity. A network attacker could intercept API traffic, read your token, and impersonate your TFE instance.

**Only use this for local dev or isolated test environments** where you accept that risk. Do not enable it for production or shared staging systems unless you fully understand the exposure.

**Prefer trusting the certificate instead** — add the signing CA or cert to your OS trust store (see [Alternatives](#alternatives) below). That keeps encryption and identity verification intact.
:::

:::caution[Not the default path]
Most TFE users never need this setting. If your instance presents a trusted certificate, leave `ssl_skip_verify` unset (or `false`) and skip this page entirely.
:::

## Profile setting (recommended)

Add `ssl_skip_verify = true` to the profile block in `~/.tfx.hcl`:

```hcl
profile "local" {
  hostname        = "local.tfe.rocks"
  organization    = "org-alpha"
  token           = "your-token"
  ssl_skip_verify = true
}
```

When the key is **absent**, TFx treats it as `false` and performs normal certificate verification.

Use the profile as usual:

```sh
tfx varset list --profile local
tfx --profile local          # TUI
```

`ssl_skip_verify` is preserved when you re-authenticate with `tfx login` on an existing profile — you do not need to set it again after updating a token.

## One-off CLI flag

For a single command without editing the config file:

```sh
tfx workspace list --profile local --ssl-skip-verify
```

## Environment variable

Set `TFE_SSL_SKIP_VERIFY` when you need skip-verify before a profile exists (for example during the first `tfx login` against a self-signed host):

```sh
export TFE_SSL_SKIP_VERIFY=true
tfx login local.tfe.rocks
```

You can combine it with a profile that already has `ssl_skip_verify = true`:

```sh
TFE_SSL_SKIP_VERIFY=true tfx varset list --profile local
```

## Configuration precedence

When multiple sources are set, the highest-precedence source wins:

1. **CLI flag** — `--ssl-skip-verify`
2. **Environment variable** — `TFE_SSL_SKIP_VERIFY`
3. **Profile value** — `ssl_skip_verify` in `.tfx.hcl`
4. **Default** — `false`

## Alternatives (recommended)

Before enabling `ssl_skip_verify`, consider installing the certificate in your system trust store. TFx uses the same trust store as other HTTPS clients — once the CA or cert is trusted, no skip-verify setting is needed.

For local development, tools like [mkcert](https://github.com/FiloSottile/mkcert) can issue certificates your machine trusts by default. That is the safer long-term approach for repeated dev work against a local TFE instance.
