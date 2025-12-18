---
title: Deploy Examples
description: Ready-to-deploy example applications for your cluster
---

# Deploy Examples

TRIGRA comes with ready-to-deploy examples for common homelab use cases.

## Quick Deploy Script

```bash
curl -fsSL https://raw.githubusercontent.com/Taiwrash/trigra/main/deploy-example.sh | bash
```

## Homepage Dashboard

A beautiful dashboard to manage all your homelab services.

```bash
kubectl apply -f https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/examples/homepage-dashboard.yaml
```

**Features:**
- Real-time cluster monitoring
- 100+ service integrations
- Customizable themes

## Ollama AI

Run local LLMs on your cluster.

```bash
kubectl apply -f https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/examples/ollama-ai-model.yaml
```

**Features:**
- Run Llama 2, Gemma, Mistral
- 50GB persistent storage
- GPU support ready

## Awesome Homepage

Stunning homepage with animations.

```bash
kubectl apply -f https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/examples/awesome-homepage.yaml
```

**Features:**
- Animated particle background
- Glassmorphism UI
- Smooth animations
