---
title: Quick Start
description: Get TRIGRA running on your cluster in 5 minutes
---

# Quick Start

Get TRIGRA up and running on your Kubernetes cluster in under 5 minutes.

## Prerequisites

Before you begin, ensure you have:

- ✅ A running Kubernetes cluster
- ✅ `kubectl` installed and configured
- ✅ Access to create namespaces and deployments

## One-Command Install

The fastest way to install TRIGRA:

```bash
curl -fsSL https://raw.githubusercontent.com/Taiwrash/trigra/main/quick-install.sh | bash -s -- default
```

This will:

1. ✅ Create the `default` namespace (or use existing)
2. ✅ Generate a secure webhook secret
3. ✅ Deploy TRIGRA controller
4. ✅ Install Cloudflare Tunnel (optional)
5. ✅ Display your webhook URL

## What Happens Next

After installation, you'll see output like:

```
╔══════════════════════════════════════════════════════════════════╗
║  🔗 YOUR WEBHOOK URL (copy this to GitHub):                     ║
╠══════════════════════════════════════════════════════════════════╣
║                                                                  ║
║  https://random-words.trycloudflare.com/webhook
║                                                                  ║
╚══════════════════════════════════════════════════════════════════╝

✓ Tunnel running in background (PID: 12345)
```

## Configure GitHub Webhook

1. Go to your **GitHub repository** → **Settings** → **Webhooks**
2. Click **Add webhook**
3. Configure:
   - **Payload URL**: The webhook URL from installation
   - **Content type**: `application/json`
   - **Secret**: Your webhook secret (shown during install)
   - **Events**: Just the push event

## Test Your Setup

Create a simple deployment in your repository:

```yaml
# test-app.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-nginx
spec:
  replicas: 2
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:alpine
```

Commit and push:

```bash
git add test-app.yaml
git commit -m "Test GitOps deployment"
git push
```

Watch it deploy:

```bash
kubectl get deployments -w
```

## Verify Installation

Check that everything is running:

```bash
# Check deployment
kubectl get deployment trigra

# Check pods
kubectl get pods -l app=trigra

# View logs
kubectl logs -f deployment/trigra
```

## Next Steps

- 📖 Read the full [Installation Guide](/getting-started/installation/) for more options
- 🔗 Learn about [GitHub Webhooks](/guides/github-webhooks/) configuration
- 🌐 Set up [Cloudflare Tunnel](/guides/cloudflare-tunnel/) for external access
- 📦 Deploy [pre-built examples](/guides/deploy-examples/)
