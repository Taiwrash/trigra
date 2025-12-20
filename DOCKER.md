# Trigra Docker Image

Official multi-platform Docker images for Trigra (Kubernetes GitOps Controller) are available on Docker Hub.

## Quick Start

```bash
docker pull taiwrash/trigra:latest
```

## Supported Architectures

We provide official images for:
- 💻 **AMD64**: Standard x86_64 servers (Cloud, VPS, etc.)
- 🍓 **ARM64**: Raspberry Pi, Apple Silicon (M1/M2/M3), and ARM-based cloud instances.

Trigra uses Docker Manifests, so simply pulling `taiwrash/trigra:latest` will automatically provide the correct architecture for your system.

## Running Locally

### With Kubeconfig (Development)

To test Trigra locally against your cluster:

```bash
docker run -p 8082:8082 \
  -v $HOME/.kube/config:/app/.kube/config:ro \
  -e KUBECONFIG=/app/.kube/config \
  -e WEBHOOK_SECRET=myhooksecret \
  -e GIT_TOKEN=ghp_xxxxxxxx \
  taiwrash/trigra:latest
```

## Image Features

✨ **Production Ready**
- **Non-Root Execution**: Runs as UID 1000 for maximum security.
- **Minimal Footprint**: Built on Alpine Linux, resulting in a ~50MB image.
- **Self-Healing**: Native health checks ensure the controller is always responsive.

🔒 **Identity & Governance**
- Full OCI metadata labels.
- Stripped and optimized binaries.
- Bundled CA certificates for secure API communication.

## Multi-Platform Build

To build multi-platform images yourself, use the included `fast-build.sh` script (requires Docker Buildx):

```bash
# Builds and pushes AMD64 & ARM64 images + Manifest
IMAGE_NAME="your/repo" ./fast-build.sh
```

## Usage in Deployment

```yaml
spec:
  containers:
  - name: trigra
    image: taiwrash/trigra:latest
    ports:
    - containerPort: 8082
    env:
    - name: GIT_PROVIDER
      value: "github"
    - name: GIT_TOKEN
      valueFrom:
        secretKeyRef:
          name: trigra-secret
          key: GIT_TOKEN
```

---
View on [Docker Hub](https://hub.docker.com/r/taiwrash/trigra).
