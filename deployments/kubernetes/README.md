# Trigra - Kubernetes Deployment Guide

This guide provides detailed instructions for deploying Trigra to your Kubernetes cluster for production use.

## 📋 Prerequisites

- A running Kubernetes cluster (K8s 1.25+)
- `kubectl` configured to access your cluster
- API Token for your Git provider (GitHub PAT, GitLab Access Token, Gitea Token, etc.)

## 🚀 Quick Deployment

### 1. Configure Secrets

Copy the example secret and update it with your credentials:

```bash
cp deployments/kubernetes/example-secret.yaml deployments/kubernetes/secret.yaml
```

Update `secret.yaml`:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: trigra-secret
stringData:
  GIT_TOKEN: "your_api_token"
  WEBHOOK_SECRET: "your_webhook_secret"
  # Required for Private Repos via SSH
  # GIT_SSH_KEY_FILE: "/etc/trigra/id_rsa" 
```

### 2. Configure Settings

Update `deployments/kubernetes/trigra-config.yaml` (or create one):

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: trigra-config
data:
  GIT_PROVIDER: "github" # Or gitlab, gitea, bitbucket, git
  NAMESPACE: "default"   # Target namespace for deployments
  PUBLIC_URL: "https://trigra.example.com" # Required for automated webhooks
```

### 3. Deploy

```bash
# 1. Apply configuration
kubectl apply -f deployments/kubernetes/trigra-config.yaml
kubectl apply -f deployments/kubernetes/secret.yaml

# 2. Setup RBAC
kubectl apply -f deployments/kubernetes/rbac.yaml

# 3. Deploy the Controller
kubectl apply -f deployments/kubernetes/deployment.yaml

# 4. Expose the Webhook Endpoint
kubectl apply -f deployments/kubernetes/service.yaml
# Use an Ingress or LoadBalancer to make this reachable at your PUBLIC_URL
```

## 🔐 Advanced Configuration

### SSH Support for Private Repositories

If you are using the `git` provider with a private repository:

1. Create a secret from your SSH private key:
   ```bash
   kubectl create secret generic trigra-ssh-key --from-file=id_rsa=/path/to/your/key
   ```

2. Mount it in `deployment.yaml`:
   ```yaml
   spec:
     template:
       spec:
         volumes:
         - name: ssh-key
           secret:
             secretName: trigra-ssh-key
             defaultMode: 0400
         containers:
         - name: trigra
           volumeMounts:
           - name: ssh-key
             mountPath: /etc/trigra
             readOnly: true
           env:
           - name: GIT_SSH_KEY_FILE
             value: "/etc/trigra/id_rsa"
   ```

### Automated Webhook Setup

When `PUBLIC_URL` is set, Trigra will attempt to register its webhook endpoint (`$PUBLIC_URL/webhook`) with your Git provider on startup.

Supported providers for auto-registration:
- **GitHub**: Requires `GIT_TOKEN` with `admin:repo_hook` scope.
- **GitLab**: Requires `GIT_TOKEN` with `api` scope.
- **Gitea**: Requires `GIT_TOKEN` permission to manage hooks.

## 🔍 Monitoring & Logs

Verify that Trigra is successfully syncing:

```bash
# Watch logs
kubectl logs -f deployment/trigra

# Verify sync
# If you see "SUCCESS: Applied resources", your GitOps flow is working!
```

---
Return to [Main Documentation](../../README.md).
