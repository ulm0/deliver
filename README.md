# Deliver


> A dead simple tool for creating releases *(plus artifacts)* on Github using GitLab CI
>
> Deliver is based on [drone-github-release](https://github.com/drone-plugins/drone-github-release)

This tool is a companion for [GitLab CI/CD for external repositories](https://docs.gitlab.com/ce/ci/ci_cd_for_external_repos/) and is supposed to be used in tags.

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
  image: ulm0/deliver:1.0
  stage: publish
  variables:
    OVERRIDE_RELEASE: "true"
    RELEASE_CHECKSUMS: sha512
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
