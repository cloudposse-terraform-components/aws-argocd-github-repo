## `argocd-github-repo` Component PR [#17](https://github.com/cloudposse-terraform-components/aws-argocd-github-repo/pull/17)

Corrected the spelling of "succeded" to "succeeded" in the `on-deploy-succeded` notification. As a result, both components (`argocd-github-repo` and `eks/argocd`) need to be updated to make this change.

See the [PR for eks/argocd](https://github.com/cloudposse-terraform-components/aws-eks-argocd/pull/16)

## Components PR [#851](https://github.com/cloudposse/terraform-aws-components/pull/851)

This is a bug fix and feature enhancement update. There are few actions necessary to upgrade.

## Upgrade actions

1. Enable `github_default_notifications_enabled` (set `true`)

```yaml
components:
  terraform:
    argocd-repo-defaults:
      metadata:
        type: abstract
      vars:
        enabled: true
        github_default_notifications_enabled: true
```

2. Apply changes with Atmos

## Features

- Support predefined GitHub commit status notifications for CD sync mode:
  - `on-deploy-started`
    - `app-repo-github-commit-status`
    - `argocd-repo-github-commit-status`
  - `on-deploy-succeded`
    - `app-repo-github-commit-status`
    - `argocd-repo-github-commit-status`
  - `on-deploy-failed`
    - `app-repo-github-commit-status`
    - `argocd-repo-github-commit-status`

### Bug Fixes

- Remove legacy unnecessary helm values used in old ArgoCD versions (ex. `workflow auth` configs) and dropped
  notifications services
