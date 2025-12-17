# TRIGRA Docker Image

Production-ready Docker image for TRIGRA (Kubernetes GitOps Homelab) available on Docker Hub.

## Quick Start

```bash
docker pull taiwrash/trigra:latest
```

## Available Tags

- `latest` - Latest stable release
- `1.0.0` - Specific version

## Running Locally

### With Kubeconfig (Local Development)

```bash
docker run -p 8082:8082 \
  -v $HOME/.kube/config:/app/.kube/config:ro \
  -e KUBECONFIG=/app/.kube/config \
  -e WEBHOOK_SECRET=your-secret \
  -e GITHUB_TOKEN=your-token \
  taiwrash/trigra:latest
```

### In-Cluster (Production)

Deploy to Kubernetes where it will automatically use in-cluster config:

```bash
kubectl apply -f deployments/kubernetes/
```

## Image Features

✨ **Production Optimized**
- Multi-stage build for minimal image size
- Non-root user for security
- Built-in health checks
- Optimized binary with stripped symbols

🔒 **Security**
- Runs as non-root user (UID 1000)
- Minimal attack surface
- CA certificates included
- Timezone data for accurate logging

📊 **Metadata**
- Full OCI image labels
- Version information
- Build timestamp
- Source repository link

## Building from Source

```bash
# Build image
./build-docker.sh

# Build specific version
./build-docker.sh 1.0.0

# Push to Docker Hub
docker login
docker push taiwrash/trigra:latest
```

## Image Details

- **Base**: Alpine Linux (minimal)
- **Size**: ~50MB
- **Architecture**: amd64
- **Health Check**: HTTP GET /health every 30s
- **Exposed Port**: 8082

## Usage in Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: trigra
spec:
  replicas: 1
  selector:
    matchLabels:
      app: trigra
  template:
    metadata:
      labels:
        app: trigra
    spec:
      containers:
      - name: trigra
        image: taiwrash/trigra:latest
        imagePullPolicy: Always
        ports:
        - containerPort: 8082
        env:
        - name: WEBHOOK_SECRET
          valueFrom:
            secretKeyRef:
              name: trigra-secret
              key: webhook-secret
```

## Docker Hub

View on Docker Hub: https://hub.docker.com/r/taiwrash/trigra
