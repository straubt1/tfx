# Gitlab CI

## Module Release
The following snippet serves as a starting point to release Terraform Modules using `tfx`:

```yaml
# .gitlab-ci.yml
variables:
  # tfx auth
  TFE_HOSTNAME: ""
  TFE_ORGANIZATION: ""
  TFE_TOKEN:  ""

  # module vars
  MODULE_NAME: aws-s3
  MODULE_PROVIDER: aws

stages:
  - release

terraform_module_release:
  stage: release
  image:
    name: ghcr.io/straubt1/tfx:latest
    entrypoint: [""]
  script:
    - tfx registry module show --name "${MODULE_NAME}" --provider "${MODULE_PROVIDER}" || tfx registry module create --name "${MODULE_NAME}" --provider "${MODULE_PROVIDER}"  
    - |
      tfx registry module version create \
        --name "${MODULE_NAME}" \
        --provider "${MODULE_PROVIDER}" \
        --version "${CI_COMMIT_TAG#v}" \
        --directory "${CI_PROJECT_DIR}"
  rules:
    - if: $CI_COMMIT_TAG
```

## Provider Release
The following snippet serves as a starting point to release Terraform Provider using `tfx`:

```yaml
# .gitlab-ci.yml
variables:
  # tfx auth
  TFE_HOSTNAME: ""
  TFE_ORGANIZATION: ""
  TFE_TOKEN:  ""

  # goreleaser vars
  GITLAB_TOKEN: ""
  GPG_FINGERPRINT: ""

  # provider vars
  PROVIDER_NAME: custom-provider

stages:
  - release
  - publish

# most likely you  will call goreleaser before publishing
goreleaser:
  stage: release
  image:
    name: goreleaser/goreleaser:latest
    entrypoint: [""]
  script:
    - goreleaser release
  artifacts:
    paths:
      - ${CI_PROJECT_DIR}/dist
  rules:
    - if: $CI_COMMIT_TAG

version:
  stage: publish
  image:
    name: ghcr.io/straubt1/tfx:latest
    entrypoint: [""]
  needs:
    - job: goreleaser_release
      artifacts: true
  script:
    - tfx registry provider version create \
        --name="${PROVIDER_NAME}" \
        --version="${CI_COMMIT_TAG#v}" \
        --key-id="${GPG_FINGERPRINT}" \
        --shasums="${CI_PROJECT_DIR}/dist/terraform-provider-${PROVIDER_NAME}_${CI_COMMIT_TAG#v}_SHA256SUMS" \
        --shasums-sig="${CI_PROJECT_DIR}/dist/terraform-provider-${PROVIDER_NAME}_${CI_COMMIT_TAG#v}_SHA256SUMS.sig"
  rules:
    - if: $CI_COMMIT_TAG
   
platforms:
  stage: publish
  image:
    name: ghcr.io/straubt1/tfx:latest
    entrypoint: [""]
  needs:
    - version
  parallel:
    matrix:
      PLATFORMS:
        - OS: linux
          ARCH: amd64
        - OS: darwin
          ARCH: arm64
        - OS: windows
          ARCH: amd64
  script:
    - tfx registry provider version platform create \
        --name="${PROVIDER_NAME}" \
        --version="${CI_COMMIT_TAG#v}" \
        --os="${OS}" \
        --arch="${ARCH}" \
        -f="${CI_PROJECT_DIR}/dist/terraform-provider-${PROVIDER_NAME}_${CI_COMMIT_TAG#v}_${OS}_${ARCH}.zip";
  rules:
    - if: $CI_COMMIT_TAG
```
