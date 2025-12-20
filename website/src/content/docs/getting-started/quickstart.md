---
title: Quick Start
description: Get Trigra running on your cluster in 5 minutes
---

# Quick Start

Get Trigra up and running on your Kubernetes cluster in under 5 minutes.

## Prerequisites

Before you begin, ensure you have:

- ✅ A running Kubernetes cluster
- ✅ `kubectl` installed and configured
- ✅ Access to create namespaces and deployments

## One-Command Install

The fastest way to install Trigra is using our interactive installer. By providing your Git provider details, Trigra can even **automatically register** its webhook for you.

```bash
# Optional: Pre-set configuration for zero-touch install
export GIT_PROVIDER="github"        # or gitlab, gitea
export GIT_TOKEN="ghp_xxx"          # your PAT
export PUBLIC_URL="https://xxx.trycloudflare.com"

# Run the installer
curl -fsSL https://raw.githubusercontent.com/Taiwrash/trigra/main/quick-install.sh | bash
```

This will:

1. ✅ Setup the target namespace
2. ✅ Configure your Git credentials securely
3. ✅ Deploy the Trigra controller
4. ✅ Setup a Cloudflare Tunnel (if requested)
5. ✅ **Automate** webhook registration (if `PUBLIC_URL` provided)

## What Happens Next

After installation, if you didn't provide a `PUBLIC_URL` for auto-setup, you'll see your webhook details:

```
╔══════════════════════════════════════════════════════════════════╗
║  🔗 YOUR WEBHOOK URL (copy this to GitHub):                     ║
╠══════════════════════════════════════════════════════════════════╣
║                                                                  ║
║  https://random-words.trycloudflare.com/webhook
║                                                                  ║
╚══════════════════════════════════════════════════════════════════╝
```

## Configure Your Git Provider

If not already automated, follow these steps:

1. Go to your **repository settings** → **Webhooks**
2. Click **Add webhook**
3. Configure:
   - **Payload URL**: The webhook URL from installation
   - **Content type**: `application/json`
   - **Secret**: Your webhook secret (shown during install)
   - **Events**: Just the push event

## Test Your Flow

Create a simple manifest in your repository:

```yaml
# test-trigra.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: trigra-test
data:
  status: "Synchronized by Trigra"
```

Commit and push:

```bash
git add test-trigra.yaml
git commit -m "Test Trigra deployment"
git push
```

Watch it deploy instantly:

```bash
kubectl get configmap trigra-test -o yaml
```

## Verify Installation

Check that everything is running:

```bash
# Check status
kubectl get pods -l app=trigra

# View logs
kubectl logs -f deployment/trigra
```

## Next Steps

- 📖 Read the full [Environment Config](/trigra/configuration/environment/)
- 📦 Deploy [ready-to-use examples](/trigra/guides/deploy-examples/)
- 🔑 Set up [Private Repos via SSH](/trigra/configuration/environment#git_ssh_key_file)
- 🌐 Setup [Cloudflare Tunnel](/trigra/guides/cloudflare-tunnel/) manually
