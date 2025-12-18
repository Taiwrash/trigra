---
title: Environment Variables
description: Configure TRIGRA using environment variables
---

# Environment Variables

TRIGRA can be configured using environment variables, which are set via Kubernetes secrets or ConfigMaps.

## Required Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `WEBHOOK_SECRET` | Secret key for GitHub webhook validation | **Required** |

## Optional Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `GITHUB_TOKEN` | GitHub token for private repositories | `""` |
| `PORT` | HTTP server port | `8080` |
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
  GITHUB_TOKEN: "ghp_xxxxxxxxxxxx"
```

Apply:

```bash
kubectl apply -f secret.yaml
```

Or create directly:

```bash
kubectl create secret generic trigra-secret \
  --from-literal=WEBHOOK_SECRET="$(openssl rand -hex 32)" \
  --from-literal=GITHUB_TOKEN="your-token"
```

## Setting via ConfigMap

For non-sensitive configuration:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: trigra-config
data:
  NAMESPACE: "production"
  LOG_LEVEL: "debug"
```

## Using in Deployment

Reference in your deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: trigra
spec:
  template:
    spec:
      containers:
      - name: trigra
        image: taiwrash/trigra:latest
        env:
        # From Secret
        - name: WEBHOOK_SECRET
          valueFrom:
            secretKeyRef:
              name: trigra-secret
              key: WEBHOOK_SECRET
        - name: GITHUB_TOKEN
          valueFrom:
            secretKeyRef:
              name: trigra-secret
              key: GITHUB_TOKEN
        # From ConfigMap
        - name: NAMESPACE
          valueFrom:
            configMapKeyRef:
              name: trigra-config
              key: NAMESPACE
        # Direct value
        - name: LOG_LEVEL
          value: "info"
```

## Helm Configuration

When using Helm, values are automatically mapped to environment variables:

```bash
helm install trigra ./helm/trigra \
  --set github.webhookSecret="your-secret" \
  --set github.token="your-token" \
  --set namespace="production"
```

## Environment Variable Details

### WEBHOOK_SECRET

**Required.** Used to validate incoming GitHub webhooks.

Generate a secure secret:

```bash
# Using OpenSSL
openssl rand -hex 32

# Using /dev/urandom
cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 64 | head -n 1
```

This must match the secret configured in your GitHub webhook settings.

### GITHUB_TOKEN

**Optional.** Required only for private repositories.

Create a token at: [GitHub Settings → Developer settings → Personal access tokens](https://github.com/settings/tokens)

Required scopes:
- `repo` - For private repositories
- `read:org` - If using organization repositories

### PORT

The HTTP port the server listens on. Defaults to `8080`.

:::note
The Kubernetes Service typically maps port 80 → 8080, so you rarely need to change this.
:::

### NAMESPACE

The default namespace where resources are deployed when not specified in the YAML.

```bash
# Deploy to staging namespace
--set namespace=staging
```

### LOG_LEVEL

Control logging verbosity:

| Level | Description |
|-------|-------------|
| `debug` | Verbose output, useful for troubleshooting |
| `info` | Standard operational logs |
| `warn` | Warnings and errors only |
| `error` | Errors only |

## Rotating Secrets

To rotate the webhook secret:

1. Generate new secret:
   ```bash
   NEW_SECRET=$(openssl rand -hex 32)
   ```

2. Update Kubernetes secret:
   ```bash
   kubectl create secret generic trigra-secret \
     --from-literal=WEBHOOK_SECRET="$NEW_SECRET" \
     --dry-run=client -o yaml | kubectl apply -f -
   ```

3. Restart deployment:
   ```bash
   kubectl rollout restart deployment/trigra
   ```

4. Update GitHub webhook with new secret
