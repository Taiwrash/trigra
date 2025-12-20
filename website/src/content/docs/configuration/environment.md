---
title: Environment Variables
description: Configure Trigra using environment variables
---

# Environment Variables

Trigra is configured via environment variables, giving you full flexibility to deploy it in any environment—from local development to production Kubernetes clusters.

## Core Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `WEBHOOK_SECRET` | Secret key for webhook validation | **Required** |
| `GIT_PROVIDER` | `github`, `gitlab`, `gitea`, `bitbucket`, `git` | `github` |
| `GIT_TOKEN` | Token for the Git provider API | - |
| `PUBLIC_URL` | Your controller's public URL (for auto-webhooks) | - |

## Advanced Configuration

| Variable | Description |
|----------|-------------|
| `GIT_SSH_KEY_FILE`| Path to SSH private key for the `git` provider |
| `GIT_REPO_URL` | Full URL of the Git repository (required for `git` provider) |
| `GIT_BASE_URL` | API base URL for self-hosted providers (GitLab/Gitea) |
| `GIT_OWNER` | Repository owner/org (overrides discovery) |
| `GIT_REPO` | Repository name (overrides discovery) |
| `SERVER_PORT` | HTTP server port (default: `8082`) |
| `NAMESPACE` | Target namespace for deployments (default: `default`) |

## Detailed Variable Info

### `PUBLIC_URL`
Set this to the public address where Trigra is reachable (e.g., `https://trigra.example.com`). When provided, Trigra will attempt to **automatically register** its webhook endpoint with GitHub, GitLab, or Gitea on startup.

### `GIT_SSH_KEY_FILE`
Path to an SSH private key on the controller's filesystem. Used by the `git` provider to clone private repositories over SSH. 
> **Tip:** In Kubernetes, mount this file from a Secret using a volume.

### `GIT_PROVIDER`
- `github`: Standard GitHub.com integration.
- `gitlab`: Supports both GitLab.com and self-managed instances.
- `gitea`: Works with Gitea and Forgejo. Correct `GIT_BASE_URL` required for self-hosted.
- `bitbucket`: Bitbucket Cloud support.
- `git`: Generic implementation that clones the repo locally. Ideal for custom git servers or local testing.

### `WEBHOOK_SECRET`
A shared secret between your Git provider and Trigra.
- **GitHub**: Paste this in the "Secret" field of the webhook settings.
- **GitLab**: Paste this in the "Secret token" field.
- **Generic Git**: Trigra validates this against the `X-Trigra-Secret` header.

## Kubernetes Secret Example

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: trigra-secret
type: Opaque
stringData:
  GIT_TOKEN: "ghp_xxxxxxxxxxxx"
  WEBHOOK_SECRET: "my-secure-webhook-secret"
  PUBLIC_URL: "https://trigra.my-homelab.com"
  # Optional: for SSH
  # GIT_SSH_KEY_FILE: "/etc/trigra/id_rsa"
```
