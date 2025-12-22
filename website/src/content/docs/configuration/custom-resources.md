---
title: Custom Resource Configuration
description: Configure Trigra repositories using Kubernetes YAML (CRDs)
---

Trigra supports declarative configuration using a **Custom Resource Definition (CRD)**. This allows you to manage multiple repositories from a single Trigra instance using standard Kubernetes YAML manifests.

## 🚀 The GitRepo CRD

The `GitRepo` resource defines a repository to be synced and its configuration.

### Installation

The CRD is included in the default installation, but you can also apply it manually:

```bash
kubectl apply -f https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/kubernetes/crds/gitrepo.yaml
```

### Example Usage

Create a file named `my-repo.yaml`:

```yaml
apiVersion: trigra.io/v1alpha1
kind: GitRepo
metadata:
  name: my-app
  namespace: default
spec:
  provider: github
  url: https://github.com/your-user/your-repo
  targetNamespace: prod
  tokenSecretRef:
    name: git-secrets
    key: GITHUB_TOKEN
  webhookSecretRef:
    name: git-secrets
    key: WEBHOOK_SECRET
```

Apply it to your cluster:

```bash
kubectl apply -f my-repo.yaml
```

## ⚙️ Configuration Fields

| Field | Description | Required |
|-------|-------------|----------|
| `url` | Full Git repository URL | Yes |
| `provider` | Git provider (`github`, `gitlab`, `gitea`, `bitbucket`, `git`) | Yes |
| `targetNamespace` | Namespace where YAMLs should be applied | No (defaults to CRD namespace) |
| `tokenSecretRef` | Reference to a Secret containing the API token | No |
| `webhookSecretRef` | Reference to a Secret containing the webhook secret | No |
| `branch` | The branch to sync | No (defaults to `main`) |

## 🔗 Webhook Setup

When using CRDs, the webhook URL for your repository follows this pattern:

`https://your-trigra-domain.com/webhook/{name}`

Where `{name}` matches the `metadata.name` of your `GitRepo` resource.

**Example:**
For a `GitRepo` named `my-app`, the payload URL in GitHub/GitLab settings should be:
`https://trigra.example.com/webhook/my-app`

## 🛡️ Secret References

Trigra looks for tokens and secrets in the same namespace as the `GitRepo` resource. 

**Example Secret:**

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: git-secrets
  namespace: default
type: Opaque
data:
  GITHUB_TOKEN: <base64-token>
  WEBHOOK_SECRET: <base64-secret>
```
