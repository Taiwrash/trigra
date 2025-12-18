---
title: Helm Values
description: Complete reference for TRIGRA Helm chart values
---

# Helm Values Reference

This page documents all available Helm values for customizing your TRIGRA installation.

## Quick Reference

```bash
helm install trigra ./helm/trigra \
  --set replicaCount=2 \
  --set github.webhookSecret="your-secret" \
  --set service.type=NodePort
```

## All Values

### Replica Count

```yaml
# Number of pod replicas
replicaCount: 1
```

For high availability, set to 2 or more.

---

### Image Configuration

```yaml
image:
  repository: taiwrash/trigra
  pullPolicy: IfNotPresent
  tag: "latest"  # Overrides chart appVersion
```

| Value | Default | Description |
|-------|---------|-------------|
| `repository` | `taiwrash/trigra` | Docker image repository |
| `pullPolicy` | `IfNotPresent` | When to pull the image |
| `tag` | `latest` | Image tag to deploy |

---

### Service Configuration

```yaml
service:
  type: LoadBalancer
  port: 80
  nodePort: ""  # Only used if type is NodePort
```

| Value | Default | Description |
|-------|---------|-------------|
| `type` | `LoadBalancer` | Service type: LoadBalancer, NodePort, ClusterIP |
| `port` | `80` | Service port |
| `nodePort` | `""` | Specific NodePort (30000-32767) |

**Examples:**

```bash
# NodePort with specific port
helm install trigra ./helm/trigra \
  --set service.type=NodePort \
  --set service.nodePort=30080

# ClusterIP (internal only)
helm install trigra ./helm/trigra \
  --set service.type=ClusterIP
```

---

### GitHub Configuration

```yaml
github:
  token: ""
  webhookSecret: ""
```

| Value | Default | Description |
|-------|---------|-------------|
| `token` | `""` | GitHub token for private repos (optional) |
| `webhookSecret` | `""` | **Required.** Secret for webhook validation |

:::caution
Never commit your webhook secret to Git! Use `--set` or external secrets management.
:::

---

### Resource Limits

```yaml
resources:
  limits:
    cpu: 500m
    memory: 256Mi
  requests:
    cpu: 100m
    memory: 128Mi
```

Adjust based on your cluster capacity.

---

### Node Selection

```yaml
nodeSelector: {}
# Example:
# nodeSelector:
#   kubernetes.io/arch: amd64
#   node-type: worker
```

---

### Tolerations

```yaml
tolerations: []
# Example:
# tolerations:
# - key: "dedicated"
#   operator: "Equal"
#   value: "trigra"
#   effect: "NoSchedule"
```

---

### Affinity Rules

```yaml
affinity: {}
# Example for anti-affinity:
# affinity:
#   podAntiAffinity:
#     requiredDuringSchedulingIgnoredDuringExecution:
#     - labelSelector:
#         matchLabels:
#           app: trigra
#       topologyKey: kubernetes.io/hostname
```

---

### Target Namespace

```yaml
namespace: default
```

The namespace where TRIGRA deploys resources. Different from the namespace where TRIGRA itself runs.

---

### Service Account

```yaml
serviceAccount:
  create: true
  name: ""  # If empty, uses release name
  annotations: {}
```

---

### Pod Annotations

```yaml
podAnnotations: {}
# Example:
# podAnnotations:
#   prometheus.io/scrape: "true"
#   prometheus.io/port: "8080"
```

---

### Security Context

```yaml
securityContext:
  runAsNonRoot: true
  runAsUser: 1000
  fsGroup: 1000
```

---

## Complete Example

```yaml
# values-production.yaml
replicaCount: 3

image:
  repository: taiwrash/trigra
  tag: "v1.0.0"
  pullPolicy: Always

service:
  type: LoadBalancer
  port: 80

github:
  webhookSecret: ""  # Set via --set

namespace: production

resources:
  limits:
    cpu: 1000m
    memory: 512Mi
  requests:
    cpu: 200m
    memory: 256Mi

affinity:
  podAntiAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
    - labelSelector:
        matchLabels:
          app: trigra
      topologyKey: kubernetes.io/hostname

podAnnotations:
  prometheus.io/scrape: "true"
```

Deploy with:

```bash
helm install trigra ./helm/trigra \
  -f values-production.yaml \
  --set github.webhookSecret="$WEBHOOK_SECRET"
```
