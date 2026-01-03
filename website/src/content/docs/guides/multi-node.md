---
title: Multi-Node Clusters
description: Running TRIGRA on multi-node Kubernetes clusters
---


Best practices for running TRIGRA on multi-node clusters.

## High Availability

Run multiple replicas:

```bash
helm upgrade trigra ./helm/trigra --set replicaCount=2
```

## Node Affinity

Run on control plane nodes:

```yaml
spec:
  template:
    spec:
      nodeSelector:
        node-role.kubernetes.io/control-plane: ""
```

## Anti-Affinity

Ensure replicas run on different nodes:

```yaml
affinity:
  podAntiAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
    - labelSelector:
        matchLabels:
          app: trigra
      topologyKey: kubernetes.io/hostname
```

## Resource Management

```yaml
resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 500m
    memory: 256Mi
```

## Scaling Recommendations

| Cluster Size | Replicas |
|--------------|----------|
| 1-3 nodes | 1 |
| 4-10 nodes | 2 |
| 10+ nodes | 3 |
