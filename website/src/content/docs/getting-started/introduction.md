---
title: Introduction
description: What is TRIGRA and why should you use it?
---

# Welcome to TRIGRA

**TRIGRA** (Trigger + Infra) is a lightweight GitOps controller for Kubernetes that enables continuous deployment through Git commits. Simply push changes to your repository, and TRIGRA automatically deploys them to your cluster.

## What is GitOps?

GitOps is a modern approach to infrastructure and application deployment where:

- **Git is the single source of truth** - All configuration is stored in Git
- **Changes are declarative** - You describe the desired state, not imperative steps
- **Automation handles deployment** - Push changes, and they're automatically applied
- **Full audit trail** - Every change is tracked in Git history

## Why TRIGRA?

### 🎯 Simplicity First

Unlike complex GitOps tools like ArgoCD or Flux, TRIGRA is designed to be:

- **Single binary** - No complex operator dependencies
- **Minimal configuration** - Works out of the box
- **Easy to understand** - Simple webhook-based architecture

### 🏠 Perfect for Homelabs

TRIGRA is optimized for homelab environments:

- **Low resource usage** - Runs on minimal hardware
- **Multi-node support** - Works across your cluster
- **Built-in Cloudflare Tunnel** - Easy external access

### ⚡ Instant Feedback

See your changes deployed immediately:

1. Edit your Kubernetes manifests
2. Push to Git
3. TRIGRA receives webhook
4. Resources are applied to your cluster

## How It Works

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Your Git  │────▶│   GitHub    │────▶│   TRIGRA    │
│   Commit    │     │   Webhook   │     │  Controller │
└─────────────┘     └─────────────┘     └─────┬───────┘
                                              │
                                              ▼
                                    ┌─────────────────┐
                                    │   Kubernetes    │
                                    │    Cluster      │
                                    └─────────────────┘
```

## Key Features

| Feature | Description |
|---------|-------------|
| **Universal Resource Support** | Deploy any Kubernetes resource type |
| **Webhook Security** | HMAC signature validation for all requests |
| **Auto-Detection** | Works in-cluster or with kubeconfig |
| **Helm Chart** | Easy installation with customizable values |
| **Cloudflare Integration** | Built-in tunnel for external access |
| **Multi-Cluster Ready** | Deploy to multiple namespaces |

## Next Steps

Ready to get started? Head to the [Quick Start](/getting-started/quickstart/) guide to deploy TRIGRA in under 5 minutes.
