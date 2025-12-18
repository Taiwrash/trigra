---
title: RBAC & Security
description: Security best practices and RBAC configuration for TRIGRA
---

# RBAC & Security

This guide covers security best practices and RBAC configuration for TRIGRA.

## RBAC Overview

TRIGRA requires cluster-level permissions to deploy resources across namespaces.

### Default RBAC Configuration

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: trigra
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: trigra
rules:
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["*"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: trigra
subjects:
- kind: ServiceAccount
  name: trigra
  namespace: default
roleRef:
  kind: ClusterRole
  name: trigra
  apiGroup: rbac.authorization.k8s.io
```

:::caution
The default RBAC grants full cluster access. For production, use [least privilege](#least-privilege-rbac).
:::

## Least Privilege RBAC

For production environments, limit permissions to only what's needed:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: trigra
rules:
# Core resources
- apiGroups: [""]
  resources:
    - configmaps
    - secrets
    - services
    - pods
    - persistentvolumeclaims
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]

# Deployments, StatefulSets, DaemonSets
- apiGroups: ["apps"]
  resources:
    - deployments
    - statefulsets
    - daemonsets
    - replicasets
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]

# Ingress
- apiGroups: ["networking.k8s.io"]
  resources:
    - ingresses
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]

# Jobs and CronJobs
- apiGroups: ["batch"]
  resources:
    - jobs
    - cronjobs
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
```

## Namespace-Scoped Access

To restrict TRIGRA to specific namespaces:

```yaml
# Use Role instead of ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: trigra
  namespace: production  # Only this namespace
rules:
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["*"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: trigra
  namespace: production
subjects:
- kind: ServiceAccount
  name: trigra
  namespace: default
roleRef:
  kind: Role
  name: trigra
  apiGroup: rbac.authorization.k8s.io
```

## Webhook Security

### Signature Validation

TRIGRA validates all incoming webhooks using HMAC-SHA256:

1. GitHub signs each webhook with your secret
2. TRIGRA computes expected signature
3. Request is rejected if signatures don't match

```go
// Validation logic
mac := hmac.New(sha256.New, []byte(webhookSecret))
mac.Write(body)
expectedSig := "sha256=" + hex.EncodeToString(mac.Sum(nil))

if !hmac.Equal([]byte(expectedSig), []byte(receivedSig)) {
    // Reject request
}
```

### Strong Secrets

Always use cryptographically secure secrets:

```bash
# Good: 256-bit random key
openssl rand -hex 32

# Bad: Weak password
# webhook-secret-123
```

## Network Security

### Network Policies

Restrict network access to TRIGRA:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: trigra-network-policy
spec:
  podSelector:
    matchLabels:
      app: trigra
  policyTypes:
  - Ingress
  - Egress
  ingress:
  # Allow webhook traffic from anywhere (GitHub)
  - from: []
    ports:
    - protocol: TCP
      port: 8080
  egress:
  # Allow DNS
  - to: []
    ports:
    - protocol: UDP
      port: 53
  # Allow Kubernetes API
  - to:
    - namespaceSelector: {}
    ports:
    - protocol: TCP
      port: 443
  # Allow GitHub API
  - to: []
    ports:
    - protocol: TCP
      port: 443
```

### TLS/HTTPS

For production, use HTTPS:

1. **Cloudflare Tunnel** - Automatic TLS
2. **Ingress with cert-manager** - Let's Encrypt certificates
3. **Service Mesh** - mTLS between services

```yaml
# Ingress with TLS
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: trigra
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  tls:
  - hosts:
    - webhook.yourdomain.com
    secretName: trigra-tls
  rules:
  - host: webhook.yourdomain.com
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

## Pod Security

### Security Context

Run TRIGRA with restricted privileges:

```yaml
spec:
  template:
    spec:
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
        fsGroup: 1000
      containers:
      - name: trigra
        securityContext:
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
          capabilities:
            drop:
            - ALL
```

### Pod Security Standards

Apply Pod Security Standards:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: trigra-system
  labels:
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
```

## Secret Management

### External Secrets

Use external secret management:

```yaml
# With External Secrets Operator
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: trigra-secret
spec:
  refreshInterval: 1h
  secretStoreRef:
    kind: ClusterSecretStore
    name: vault
  target:
    name: trigra-secret
  data:
  - secretKey: WEBHOOK_SECRET
    remoteRef:
      key: trigra/webhook-secret
```

### Secret Rotation

Rotate secrets regularly:

1. Generate new secret
2. Update GitHub webhook with new secret
3. Update Kubernetes secret
4. Restart TRIGRA pods

```bash
# Rotation script
NEW_SECRET=$(openssl rand -hex 32)
kubectl create secret generic trigra-secret \
  --from-literal=WEBHOOK_SECRET="$NEW_SECRET" \
  --dry-run=client -o yaml | kubectl apply -f -
kubectl rollout restart deployment/trigra
echo "Update GitHub webhook with: $NEW_SECRET"
```

## Audit Logging

Enable Kubernetes audit logging for TRIGRA actions:

```yaml
# Audit policy
apiVersion: audit.k8s.io/v1
kind: Policy
rules:
- level: Metadata
  users: ["system:serviceaccount:default:trigra"]
  verbs: ["create", "update", "patch", "delete"]
```

## Security Checklist

- [ ] Use strong, randomly generated webhook secret
- [ ] Enable HTTPS (Cloudflare Tunnel or Ingress)
- [ ] Apply least privilege RBAC
- [ ] Enable network policies
- [ ] Run as non-root user
- [ ] Use external secret management for production
- [ ] Rotate secrets periodically
- [ ] Enable audit logging
- [ ] Monitor for suspicious activity
