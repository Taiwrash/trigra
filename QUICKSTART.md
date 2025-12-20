# Quick Start Guide

This guide will help you get Trigra up and running in minutes.

## 🚀 Instant Deployment

The fastest way to deploy Trigra to your cluster is using our interactive installer:

```bash
# Set your Git provider configuration
export GIT_PROVIDER="github"
export GIT_TOKEN="ghp_your_personal_access_token"
export WEBHOOK_SECRET=$(openssl rand -hex 32)
export PUBLIC_URL="https://trigra.yourdomain.com" # Required for auto-webhooks

# Run the installer
curl -fsSL https://raw.githubusercontent.com/Taiwrash/trigra/main/quick-install.sh | bash
```

## 🛠 Manual Installation

If you prefer to maintain full control, follow these steps:

1. **Clone the Repo**:
   ```bash
   git clone https://github.com/Taiwrash/trigra.git
   cd trigra
   ```

2. **Configure Secrets**:
   Edit `deployments/kubernetes/example-secret.yaml` and apply it:
   ```bash
   kubectl apply -f deployments/kubernetes/example-secret.yaml
   ```

3. **Deploy Core Components**:
   ```bash
   kubectl apply -f deployments/kubernetes/rbac.yaml
   # Note: Ensure you update the image in deployment.yaml if using a custom build
   kubectl apply -f deployments/kubernetes/deployment.yaml
   kubectl apply -f deployments/kubernetes/service.yaml
   ```

4. **Automated Webhook**:
   If you provided `PUBLIC_URL`, Trigra will automatically register its webhook with your Git provider (GitHub, GitLab, or Gitea) on startup.

## 🧪 Testing the Flow

1. Create a simple manifest in your Git repo:
   ```yaml
   # test-app.yaml
   apiVersion: v1
   kind: ConfigMap
   metadata:
     name: trigra-test
   data:
     message: "Trigra is working!"
   ```

2. Commit and Push:
   ```bash
   git add test-app.yaml
   git commit -m "Test Trigra sync"
   git push
   ```

3. Verify:
   ```bash
   kubectl get configmap trigra-test -o yaml
   ```

## 🔍 Useful Commands

| Task | Command |
|------|---------|
| View Logs | `kubectl logs -f deployment/trigra` |
| Check Status | `kubectl get pods -l app=trigra` |
| Health Check | `curl http://<trigra-ip>/health` |
| Ready Check | `curl http://<trigra-ip>/ready` |

---
For detailed information, visit the [Full Documentation](README.md).
