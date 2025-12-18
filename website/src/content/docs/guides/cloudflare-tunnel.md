---
title: Cloudflare Tunnel
description: Securely expose your TRIGRA webhook to the internet
---

# Cloudflare Tunnel

Cloudflare Tunnel allows you to expose your TRIGRA webhook endpoint to the internet without opening ports or configuring firewalls.

## Benefits

- ✅ No port forwarding required
- ✅ Automatic HTTPS with valid certificate
- ✅ DDoS protection included
- ✅ No public IP needed

## Automatic Setup

The TRIGRA quick-install script automatically sets up Cloudflare Tunnel:

```bash
curl -fsSL https://raw.githubusercontent.com/Taiwrash/trigra/main/quick-install.sh | bash -s -- default
```

After installation, you'll see your webhook URL displayed in a styled box.

## Manual Setup

### Step 1: Install cloudflared

**macOS:**
```bash
brew install cloudflare/cloudflare/cloudflared
```

**Debian/Ubuntu:**
```bash
curl -L --output cloudflared.deb https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
sudo dpkg -i cloudflared.deb
```

### Step 2: Get Service IP/Port

```bash
SERVICE_IP=$(kubectl get svc trigra -o jsonpath='{.spec.clusterIP}')
PORT=$(kubectl get svc trigra -o jsonpath='{.spec.ports[0].port}')
echo "Service: http://${SERVICE_IP}:${PORT}"
```

### Step 3: Start Tunnel

```bash
cloudflared tunnel --url http://${SERVICE_IP}:${PORT}
```

## Running as a Service

For persistent tunnels, run cloudflared as a systemd service. See our systemd configuration in the full documentation.

## Named Tunnels (Production)

For production, use named tunnels with persistent URLs:

```bash
cloudflared tunnel login
cloudflared tunnel create trigra-webhook
cloudflared tunnel route dns trigra-webhook webhook.yourdomain.com
cloudflared tunnel run trigra-webhook
```

Your permanent webhook URL: `https://webhook.yourdomain.com/webhook`
