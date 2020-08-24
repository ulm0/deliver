# Deliver


> A dead simple tool for creating releases *(plus artifacts)* on Github using GitLab CI
>
> Deliver is based on [drone-github-release](https://github.com/drone-plugins/drone-github-release)

This tool is a companion for [GitLab CI/CD for external repositories](https://docs.gitlab.com/ce/ci/ci_cd_for_external_repos/)

# Usage

An API token is required, save it as `GITHUB_TOKEN` in GitLab CI/CD variables

```yaml
stages:
  - build
  - publish

build:
  image: alpine
  stage: build
  before_script:
    - apk add --no-cache make
  script:
    - make
  artifacts:
    paths:
      - build/*
  only:
    - master
    - tags
    - external_pull_requests

github:release:
  image: ulm0/deliver
  stage: publish
  variables:
    OVERRIDE_RELEASE: "true"
    RELEASE_CHECKSUM: sha512
    RELEASE_DESCRIPTION: "My awesome release"
    RELEASE_FILES: $CI_PROJECT_DIR/build/*
    RELEASE_NAME: $CI_COMMIT_TAG
  script:
    - deliver
  only:
    - tags
  dependencies:
    - build
```

## Options

```sh
--api-key value, -k value     API token to access github [$GITHUB_TOKEN, $GITHUB_API_TOKEN]
--api-url value, -a value     API endpoint (default: "https://api.github.com/") [$GITHUB_API_URL]
--checksum value, -c value    Methods for generating files checksums [$RELEASE_CHECKSUM]
--checksum-file value         Name for checksum file. "CHECKSUM" is replaced with chosen method (default: "CHECKSUMsum.txt") [$CHECKSUM_FILE]
--checksum-flatten, -l        Include only the basename of the file in the checksum file (default: true) [$CHECKSUM_FLATTEN]
--description value           File or string containing release description [$RELEASE_DESCRIPTION]
--draft, -d                   This is a draft release (default: false) [$DRAFT_RELEASE]
--file-exists value           Behavior in case a file previously exists (default: "overwrite") [$RELEASE_FILE_EXISTS]
--files value, -f value       List of files to release [$RELEASE_FILES]
--name value, -n value        Name for this release [$RELEASE_NAME]
--only-tags                   Release only on tag refs. If set to false --tag must be set (default: true) [$RELEASE_ONLY_TAGS]
--override, -o                Override existing release information (default: false) [$OVERRIDE_RELEASE]
--pre-release, -p             This is a pre-release (default: false) [$PRE_RELEASE]
--tag value, -t value         Tag ref for this release. If --only-tags is false, this value must be set [$CI_COMMIT_TAG, $RELEASE_TAG]
--upload-url value, -u value  API endpoint for uploading assets (default: "https://uploads.github.com/") [$GITHUB_UPLOAD_URL]
--help, -h                    show help (default: false)
--version, -v                 print the version (default: false)
```
