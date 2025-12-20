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
| `GIT_PROVIDER` | Git provider to use (`github`, `gitlab`) | `github` |
| `GIT_TOKEN` | Token for the Git provider API | `""` |
| `GITHUB_TOKEN` | Alias for `GIT_TOKEN` (backward compatibility) | `""` |
| `SERVER_PORT` | HTTP server port | `8082` |
| `NAMESPACE` | Target namespace for deployments | `default` |
| `LOG_LEVEL` | Logging verbosity (debug, info, warn, error) | `info` |

## Setting via Kubernetes Secret

The recommended way to set sensitive values:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: trigra-secret
type: Opaque
stringData:
  WEBHOOK_SECRET: "your-webhook-secret-here"
  GIT_TOKEN: "your-provider-token"
```

Apply:

```bash
kubectl apply -f secret.yaml
```

## Setting via ConfigMap

For non-sensitive configuration:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: trigra-config
data:
  GIT_PROVIDER: "github"
  NAMESPACE: "production"
  LOG_LEVEL: "debug"
```

## Environment Variable Details

### GIT_PROVIDER

The Git provider where your manifests are stored. Currently supported:
- `github` (default)
- `gitlab` (experimental)

### GIT_TOKEN / GITHUB_TOKEN

**Optional.** Required only for private repositories or to increase API rate limits.

### WEBHOOK_SECRET

**Required.** Used to validate incoming webhooks. This must match the secret configured in your provider's webhook settings.

Generate a secure secret:

```bash
openssl rand -hex 32
```

### NAMESPACE

The default namespace where resources are deployed when not specified in the YAML.

### LOG_LEVEL

Control logging verbosity: `debug`, `info`, `warn`, `error`.
