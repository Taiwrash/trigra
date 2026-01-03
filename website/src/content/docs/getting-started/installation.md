---
title: Installation Guide
description: Complete guide to installing TRIGRA on your Kubernetes cluster
---


This guide covers all installation methods for TRIGRA.

## Installation Methods

### Method 1: Quick Install Script (Recommended)

The easiest way to install TRIGRA:

```bash
# Install to default namespace
curl -fsSL https://raw.githubusercontent.com/Taiwrash/trigra/main/quick-install.sh | bash -s -- default

# Install to custom namespace
curl -fsSL https://raw.githubusercontent.com/Taiwrash/trigra/main/quick-install.sh | bash -s -- my-namespace

# With custom webhook secret
curl -fsSL https://raw.githubusercontent.com/Taiwrash/trigra/main/quick-install.sh | bash -s -- default my-secret-key
```

The script will:
- ✅ Check prerequisites (kubectl, cluster connection)
- ✅ Create namespace if needed
- ✅ Generate webhook secret (if not provided)
- ✅ Deploy RBAC, controller, and service
- ✅ Optionally start Cloudflare Tunnel

### Method 2: Helm Chart

For more control over installation:

```bash
# Clone the repository
git clone https://github.com/Taiwrash/trigra.git
cd trigra

# Generate webhook secret
WEBHOOK_SECRET=$(openssl rand -hex 32)

# Install with Helm
helm install trigra ./helm/trigra \
  --namespace default \
  --set github.webhookSecret="$WEBHOOK_SECRET" \
  --set github.token="YOUR_GITHUB_TOKEN"  # Optional for private repos

# Get webhook URL
kubectl get svc trigra
```

### Method 3: kubectl (Manual)

For full manual control:

```bash
# 1. Apply RBAC
kubectl apply -f https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/kubernetes/rbac.yaml

# 2. Create secret
kubectl create secret generic trigra-secret \
  --from-literal=GITHUB_TOKEN="YOUR_TOKEN" \
  --from-literal=WEBHOOK_SECRET="$(openssl rand -hex 32)"

# 3. Deploy controller
kubectl apply -f https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/kubernetes/deployment.yaml
kubectl apply -f https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/kubernetes/service.yaml
```

## Service Types

### LoadBalancer (Default)

Best for clusters with LoadBalancer support (cloud providers, MetalLB):

```yaml
service:
  type: LoadBalancer
```

### NodePort

For clusters without LoadBalancer:

```bash
helm install trigra ./helm/trigra \
  --set service.type=NodePort \
  --set service.nodePort=30080
```

Access via: `http://<NODE-IP>:30080/webhook`

### ClusterIP with Ingress

For production with Ingress controller:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: trigra
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  tls:
  - hosts:
    - gitops.yourdomain.com
    secretName: trigra-tls
  rules:
  - host: gitops.yourdomain.com
    http:
      paths:
      - path: /webhook
        pathType: Prefix
        backend:
          service:
            name: trigra
            port:
              number: 80
```

## Verify Installation

```bash
# Check deployment status
kubectl get deployment trigra
kubectl get pods -l app=trigra

# Check service
kubectl get svc trigra

# View logs
kubectl logs -f deployment/trigra

# Test health endpoint
curl http://<SERVICE-IP>/health
```

## Upgrading

### Helm Upgrade

```bash
helm upgrade trigra ./helm/trigra \
  --reuse-values \
  --set image.tag=new-version
```

### kubectl Upgrade

```bash
kubectl apply -f deployments/kubernetes/deployment.yaml
kubectl rollout restart deployment/trigra
```

## Uninstalling

### Helm

```bash
helm uninstall trigra
```

### kubectl

```bash
kubectl delete -f deployments/kubernetes/
kubectl delete clusterrolebinding trigra
kubectl delete clusterrole trigra
```

## Troubleshooting

### Pods Not Starting

```bash
kubectl describe pod -l app=trigra
kubectl logs -l app=trigra
```

### Webhook Not Working

```bash
# Check service is accessible
kubectl get svc trigra

# Test health endpoint
curl http://<SERVICE-IP>/health

# Check GitHub webhook deliveries in your repo settings
```

### RBAC Issues

```bash
kubectl get sa trigra
kubectl get clusterrolebinding trigra
```
