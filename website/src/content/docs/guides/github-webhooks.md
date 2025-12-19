---
title: GitHub Webhooks
description: Configure GitHub webhooks for automated deployments
---

# GitHub Webhooks

This guide explains how to configure GitHub webhooks to trigger automated deployments with TRIGRA.

## How It Works

When you push to your repository:

1. GitHub sends a webhook POST request to your TRIGRA endpoint
2. TRIGRA validates the webhook signature using your secret
3. If valid, TRIGRA fetches the changed files from Git
4. Kubernetes manifests are applied to your cluster

## Setting Up the Webhook

### Step 1: Get Your Webhook URL

After installing TRIGRA, get your webhook URL:

```bash
# For LoadBalancer
kubectl get svc trigra -o jsonpath='{.status.loadBalancer.ingress[0].ip}'

# For Cloudflare Tunnel (shown after install)
# https://random-words.trycloudflare.com/webhook
```

### Step 2: Get Your Webhook Secret

Your webhook secret was either:
- Generated during installation (shown in output)
- Provided by you during installation
- Stored in the Kubernetes secret:

```bash
kubectl get secret trigra-secret -o jsonpath='{.data.WEBHOOK_SECRET}' | base64 -d
```

### Step 3: Configure GitHub

1. Go to your **repository** → **Settings** → **Webhooks**
2. Click **Add webhook**
3. Fill in the form:

| Field | Value |
|-------|-------|
| **Payload URL** | `http://YOUR-IP/webhook` or `https://random.trycloudflare.com/webhook` |
| **Content type** | `application/json` |
| **Secret** | Your webhook secret |
| **SSL verification** | Enable if using HTTPS |
| **Events** | `Just the push event` |

4. Click **Add webhook**

### Step 4: Test the Webhook

GitHub will send a ping event. Check it:

1. Go to **Webhooks** → Your webhook → **Recent Deliveries**
2. You should see a successful ping (green checkmark)

## Webhook Security

TRIGRA uses HMAC-SHA256 to verify webhook authenticity.

**Never expose your webhook secret!** Store it securely.

## Troubleshooting

### Webhook Returns 401 Unauthorized

- **Cause**: Signature mismatch
- **Fix**: Ensure webhook secret matches in GitHub and TRIGRA

### Webhook Returns 404 Not Found

- **Cause**: Wrong URL path
- **Fix**: Ensure URL ends with `/webhook`

### Webhook Times Out

- **Cause**: TRIGRA not accessible externally
- **Fix**: Use [Cloudflare Tunnel](/trigra/guides/cloudflare-tunnel/) or verify network access
