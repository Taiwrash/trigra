---
title: Environment Variables
description: Configure TRIGRA using environment variables
---

# Environment Variables

TRIGRA can be configured using environment variables, which are set via Kubernetes secrets or ConfigMaps.

## Required Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `WEBHOOK_SECRET` | Secret key for webhook validation | **Required** |

## Optional Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `GIT_PROVIDER` | Git provider to use (`github`, `gitlab`, `bitbucket`, `git`) | `github` |
| `GIT_TOKEN` | Token for the Git provider API | `""` |
| `GIT_REPO_URL` | Full URL of the Git repository (required for `git` provider) | `""` |
| `BITBUCKET_USER` | Username for Bitbucket Basic Auth (required for `bitbucket` provider) | `""` |
| `SERVER_PORT` | HTTP server port | `8082` |
| `NAMESPACE` | Target namespace for deployments | `default` |

## Environment Variable Details

### GIT_PROVIDER

The Git provider where your manifests are stored.
- `github` (default): Native GitHub API support.
- `gitlab`: Native GitLab API support.
- `bitbucket`: Bitbucket Cloud support. Requires `BITBUCKET_USER`.
- `git`: Generic Git support using local cloning. Requires `GIT_REPO_URL`.

### GIT_TOKEN

**Recommended.** API token for the provider.
- **GitHub**: Personal Access Token (classic or fine-grained).
- **GitLab**: Personal Access Token or Project Access Token.
- **Bitbucket**: App Password.

### GIT_REPO_URL

Only used with `GIT_PROVIDER=git`. This can be any URL accessible by the controller (e.g., `https://gitea.local/owner/repo.git`).

### WEBHOOK_SECRET

**Required.** Used to validate incoming webhooks.
- **GitHub**: Set in the "Secret" field of the webhook.
- **GitLab**: Set in the "Secret token" field.
- **Bitbucket**: Not natively supported by Bitbucket Cloud webhooks, but can be used for secondary validation if implemented.
- **Generic Git**: Checked against the `X-Trigra-Secret` header.

## Setting via Kubernetes Secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: trigra-secret
type: Opaque
stringData:
  WEBHOOK_SECRET: "your-webhook-secret-here"
  GIT_TOKEN: "your-provider-token"
  GIT_REPO_URL: "https://your-git-server.com/repo.git" # Optional
```
