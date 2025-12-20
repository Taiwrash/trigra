# Trigra 
> Deploy to Kubernetes with Git. Simple, Fast, Reliable.

[![Security & Tests](https://github.com/Taiwrash/trigra/actions/workflows/security.yml/badge.svg)](https://github.com/Taiwrash/trigra/actions/workflows/security.yml)
[![Build & Push](https://github.com/Taiwrash/trigra/actions/workflows/build-push.yml/badge.svg)](https://github.com/Taiwrash/trigra/actions/workflows/build-push.yml)

Trigra is a lightweight, high-performance GitOps controller for Kubernetes. It synchronizes your cluster state with your Git repository, ensuring that your infrastructure is always up-to-date, version-controlled, and reproducible.

## 🚀 Features

- **Multi-Provider Support**: Native support for **GitHub**, **GitLab**, **Gitea**, **Bitbucket**, and **Generic Git**.
- **Automated Webhook Setup**: Zero-touch configuration for GitHub, GitLab, and Gitea.
- **SSH Authentication**: Securely clone from private repositories using SSH keys.
- **Universal Resource Support**: Deploy any Kubernetes resource (Deployments, CRDs, RBAC, etc.).
- **Smart Sync**: Automatically detects changes and applies them using optimistic concurrency.
- **Multi-Platform**: Official images for `amd64` and `arm64` (Raspberry Pi, Apple Silicon).
- **Security-First**: Built-in SAST, SCA, and secret scanning. Non-root execution by default.

## 🏗️ Architecture

```
┌─────────────┐
│    Git      │
│  Provider   │
└──────┬──────┘
       │ Push Event
       │ (Webhook)
       ▼
┌─────────────────────┐
│       Trigra        │
│ ┌─────────────────┐ │
│ │ Webhook Handler │ │
│ └────────┬────────┘ │
│          │          │
│ ┌────────▼────────┐ │
│ │ Resource Applier│ │
│ └────────┬────────┘ │
└──────────┼──────────┘
           │
           ▼
    ┌─────────────┐
    │ Kubernetes  │
    │   Cluster   │
    └─────────────┘
```

## 🛠️ Quick Start

### 1. In-Cluster (Production)

The easiest way to get started is using our **Quick Install** script:

```bash
curl -fsSL https://raw.githubusercontent.com/Taiwrash/trigra/main/quick-install.sh | bash
```

### 2. Manual Deployment

See the [Kubernetes Guide](deployments/kubernetes/README.md) for detailed manual setup instructions.

### 3. Local Development

```bash
git clone https://github.com/Taiwrash/trigra.git
cd trigra
go run ./cmd/trigra
```

## 📝 Configuration

Trigra is configured primarily via environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `GIT_PROVIDER` | `github`, `gitlab`, `gitea`, `bitbucket`, `git` | `github` |
| `GIT_TOKEN` | API Token for your provider | - |
| `WEBHOOK_SECRET`| Secret to validate webhooks | (Required) |
| `PUBLIC_URL` | Your controller's public URL | - |
| `GIT_SSH_KEY_FILE`| Path to SSH private key | - |

See [Configuration Guide](website/src/content/docs/configuration/environment.md) for all options.

## 🎯 Example Usage

Create a file named `app.yaml` in your Git repo:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: helloworld
spec:
  replicas: 2
  template:
    spec:
      containers:
      - name: app
        image: nginx:alpine
```

Commit and push. Trigra will detect the change and deploy it instantly.

## 🔒 Security

Trigra is designed with security in mind:
- **Least Privilege**: Runs as a non-root user.
- **Secret Integration**: Supports Kubernetes Secrets and SecretManagers.
- **Validation**: Strict webhook signature validation for all providers.

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## 📄 License

Trigra is licensed under the [MIT License](LICENSE).

---
**Happy GitOps-ing! 🚀**
