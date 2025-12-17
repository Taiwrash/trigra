# TRIGRA Quick Start Card

## For Users Without Codebase Access

### 1️⃣ Get the Files

```bash
mkdir trigra-deploy && cd trigra-deploy

# Download manifests
curl -O https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/kubernetes/secret.yaml.example
curl -O https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/kubernetes/deployment.yaml
curl -O https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/kubernetes/rbac.yaml
curl -O https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/kubernetes/service.yaml
```

### 2️⃣ Get GitHub Token

1. Go to: https://github.com/settings/tokens
2. Click "Generate new token (classic)"
3. Select scope: `repo` (for private) or `public_repo` (for public)
4. Copy the token

### 3️⃣ Configure Secret

```bash
# Create your secret file
cp secret.yaml.example secret.yaml

# Edit it
vim secret.yaml
```

Replace these values:
```yaml
GITHUB_TOKEN: "ghp_your_actual_token_here"
WEBHOOK_SECRET: "your_webhook_secret_here"  # Generate with: openssl rand -hex 32
```

### 4️⃣ Deploy

```bash
kubectl apply -f secret.yaml
kubectl apply -f rbac.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```

### 5️⃣ Verify

```bash
# Check if running
kubectl get pods -l app=trigra

# View logs
kubectl logs -l app=trigra -f
```

## Common Issues

| Problem | Solution |
|---------|----------|
| Pod not starting | Check: `kubectl describe pod -l app=trigra` |
| Auth errors | Verify token has `repo` or `public_repo` scope |
| Not detecting changes | Check `GITHUB_REPO` URL is correct |
| Image pull error | Image is public, check internet connection |

## Update Configuration

```bash
# Edit secret
vim secret.yaml

# Apply changes
kubectl apply -f secret.yaml

# Restart
kubectl rollout restart deployment trigra
```

## Uninstall

```bash
kubectl delete -f .
```

---

📖 **Full Guide**: See [deployments/kubernetes/README.md](README.md) for detailed instructions
