# TRIGRA Kubernetes Deployment Guide

This guide will help you deploy Trigra (Kubernetes GitOps Controller) to your Kubernetes cluster.

![](https://app.eraser.io/workspace/QdXTK61OUqZqUE2ocAE7/preview?elements=gJntNtLMFIXvvvYV2TYMNw&type=embed)

## Prerequisites

- A running Kubernetes cluster
- `kubectl` configured to access your cluster
- A GitHub Personal Access Token
- Your GitHub repository URL

## Quick Start

### 1. Create Your Secret

First, copy the example secret file and configure it with your credentials:

```bash
# Download the example secret
curl -O https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/kubernetes/secret.yaml.example

# Copy it to create your actual secret file
cp secret.yaml.example secret.yaml
```

Edit `secret.yaml` and replace the placeholder values:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: trigra-secret
  namespace: default
type: Opaque
stringData:
  GITHUB_TOKEN: "ghp_your_actual_token_here"
  WEBHOOK_SECRET: "your_webhook_secret_generated_with_openssl"
```

**Important:** Never commit `secret.yaml` to version control!

### 2. Get Your GitHub Token

1. Go to [GitHub Settings > Personal Access Tokens](https://github.com/settings/tokens)
2. Click "Generate new token (classic)"
3. Give it a descriptive name (e.g., "Trigra GitOps Controller")
4. Select scopes:
   - For **private repositories**: Select `repo` (Full control of private repositories)
   - For **public repositories**: Select `public_repo` (Access public repositories)
5. Click "Generate token"
6. **Copy the token immediately** (you won't be able to see it again!)

### 3. Deploy to Kubernetes

Deploy all components in the correct order:

```bash
# 1. Create the secret (must be first!)
kubectl apply -f secret.yaml

# 2. Create RBAC (service account, role, role binding)
kubectl apply -f rbac.yaml

# 3. Deploy the controller
kubectl apply -f deployment.yaml

# 4. Create the service (optional, for webhooks)
kubectl apply -f service.yaml
```

Or deploy everything at once:

```bash
kubectl apply -f secret.yaml
kubectl apply -f .
```

### 4. Verify Deployment

Check that everything is running:

```bash
# Check if the pod is running
kubectl get pods -l app=trigra

# Check the logs
kubectl logs -l app=trigra -f

# Check the secret was created
kubectl get secret trigra-secret
```

You should see output like:
```
NAME                   READY   STATUS    RESTARTS   AGE
trigra-7d8f9b5c4d-x7k2m   1/1     Running   0          30s
```

## Configuration Options

### Environment Variables (via Secret)

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `GITHUB_TOKEN` | Yes | - | GitHub Personal Access Token |
| `WEBHOOK_SECRET` | Yes | - | Secret for validating GitHub webhooks |
| `SERVER_PORT` | No | `8082` | Port for the webhook server |
| `NAMESPACE` | No | `default` | Kubernetes namespace to deploy resources |

### Deployment Configuration

Edit `deployment.yaml` to customize:

- **Replicas**: Change `replicas: 1` to scale horizontally
- **Resources**: Adjust CPU/memory limits and requests
- **Image**: Use a specific version tag instead of `latest`
- **Namespace**: Deploy to a different namespace

Example:
```yaml
spec:
  replicas: 2  # Run 2 instances for high availability
  template:
    spec:
      containers:
      - name: controller
        image: taiwrash/trigra:v1.0.0  # Use specific version
        resources:
          limits:
            memory: "256Mi"
            cpu: "500m"
          requests:
            memory: "128Mi"
            cpu: "100m"
```

## Troubleshooting

### Pod is not starting

```bash
# Check pod status
kubectl describe pod -l app=trigra

# Common issues:
# - Secret not created: Create secret.yaml first
# - Image pull error: Check Docker Hub or use a different image tag
# - RBAC issues: Ensure rbac.yaml is applied
```

### Controller is not detecting changes

```bash
# Check the logs
kubectl logs -l app=trigra -f

# Verify your secret values
kubectl get secret trigra-secret -o yaml

# Common issues:
# - Invalid GitHub token: Generate a new one
# - Webhook secret mismatch: Ensure GitHub webhook secret matches WEBHOOK_SECRET
# - Webhook not configured: Set up webhook in GitHub repository settings
```

### Authentication errors

```bash
# Check if token has correct permissions
# Token needs 'repo' scope for private repos or 'public_repo' for public repos

# Recreate the secret with correct token
kubectl delete secret trigra-secret
kubectl apply -f secret.yaml
kubectl rollout restart deployment trigra
```

## Updating the Deployment

### Update Configuration

```bash
# Edit your secret
vim secret.yaml

# Apply changes
kubectl apply -f secret.yaml

# Restart the deployment to pick up new values
kubectl rollout restart deployment trigra
```

### Update to New Version

```bash
# Update the image tag in deployment.yaml
kubectl set image deployment/trigra trigra=taiwrash/trigra:v1.1.0

# Or edit deployment.yaml and apply
kubectl apply -f deployment.yaml
```

## Uninstalling

To remove Trigra from your cluster:

```bash
# Delete all resources
kubectl delete -f .

# Or delete individually
kubectl delete deployment trigra
kubectl delete service trigra
kubectl delete secret trigra-secret
kubectl delete serviceaccount trigra
kubectl delete role trigra-role
kubectl delete rolebinding trigra-rolebinding
```

## Security Best Practices

1. **Never commit secrets**: Always keep `secret.yaml` in `.gitignore`
2. **Use least privilege**: Only grant necessary GitHub token scopes
3. **Rotate tokens regularly**: Generate new tokens periodically
4. **Use namespaces**: Deploy to a dedicated namespace for isolation
5. **Limit RBAC**: Review `rbac.yaml` and restrict permissions as needed
6. **Use specific image tags**: Avoid `latest` in production

## Getting the Manifests

If you don't have access to the codebase, download the manifests directly:

```bash
# Create a directory
mkdir trigra-deployment
cd trigra-deployment

# Download all manifests
curl -O https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/kubernetes/secret.yaml.example
curl -O https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/kubernetes/deployment.yaml
curl -O https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/kubernetes/rbac.yaml
curl -O https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/kubernetes/service.yaml

# Create your secret
cp secret.yaml.example secret.yaml
vim secret.yaml  # Edit with your values
```

## Support

For issues or questions:
- Check the logs: `kubectl logs -l app=trigra -f`
- Review this guide's troubleshooting section
- Open an issue on GitHub
