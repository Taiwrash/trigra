# Trigra - Kubernetes Quick Start Card

A cheat-sheet for deploying Trigra quickly in any environment.

## 1️⃣ One-Step Installation (Installer)

```bash
export GIT_PROVIDER="github"        # or gitlab, gitea
export GIT_TOKEN="ghp_xxx"          # Personal Access Token
export PUBLIC_URL="https://xxx.com" # Required for Auto-Webhooks
export WEBHOOK_SECRET="secure-key"  # Used for webhook validation

curl -fsSL https://raw.githubusercontent.com/Taiwrash/trigra/main/quick-install.sh | bash
```

## 2️⃣ Manual Manifest Deployment

```bash
# Get the manifests
mkdir trigra && cd trigra
curl -O https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/kubernetes/trigra-config.yaml
curl -O https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/kubernetes/example-secret.yaml
curl -O https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/kubernetes/deployment.yaml
curl -O https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/kubernetes/rbac.yaml
curl -O https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/kubernetes/service.yaml

# 1. Update config and secrets
vim trigra-config.yaml
vim example-secret.yaml # rename to secret.yaml

# 2. Apply everything
kubectl apply -f .
```

## 3️⃣ Verification

| Command | Purpose |
|---------|---------|
| `kubectl get pods -l app=trigra` | Check if Trigra is running |
| `kubectl logs -f deployment/trigra` | Watch real-time sync activities |
| `kubectl describe deployment trigra` | Check configuration and events |

## 💡 Troubleshooting

- **No sync?**: Check if `PUBLIC_URL` is set or webhooks are configured in your Git provider.
- **Auth error?**: Ensure `GIT_TOKEN` has `repo` (private) or `public_repo` scope.
- **Traverse error?**: Check your manifest paths; Trigra sanitizes all file paths for security.

---
📖 **Full Guide**: [Main README](../../README.md)
