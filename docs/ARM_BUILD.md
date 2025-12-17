# Building TRIGRA for ARM Architecture

## Problem

The Docker image `taiwrash/trigra:latest` on Docker Hub was built for AMD64 architecture, but your Kubernetes cluster nodes are running ARM64 (likely Raspberry Pi or similar ARM-based hardware).

**Error**: `exec ./trigra: exec format error`

## Solution: Build Locally on Your Cluster

Since your cluster nodes are ARM-based, the easiest solution is to build the image directly on one of your nodes.

### Option 1: Build on Cluster Node (Recommended)

SSH into one of your Kubernetes nodes and run:

```bash
# Clone the repository
git clone https://github.com/Taiwrash/trigra.git
cd trigra

# Build the image locally (will automatically use ARM64)
docker build -t taiwrash/trigra:latest .

# The image is now available locally on this node
# Kubernetes will use it when deploying
```

### Option 2: Build and Push to Local Registry

If you have a local container registry:

```bash
# Build for ARM64
docker build -t your-registry.local/trigra:latest .

# Push to your registry
docker push your-registry.local/trigra:latest

# Update deployment to use your registry
kubectl set image deployment/trigra trigra=your-registry.local/trigra:latest
```

### Option 3: Use ImagePullPolicy Never

Build on each node that might run the pod:

```bash
# On each node, build the image
# On each node, build the image
docker build -t taiwrash/trigra:latest .
```

Then update the deployment:

```yaml
spec:
  template:
    spec:
      containers:
      - name: controller
        image: taiwrash/trigra:latest
        imagePullPolicy: Never  # Use local image only
```

Apply the change:

```bash
kubectl patch deployment trigra -p '{"spec":{"template":{"spec":{"containers":[{"name":"trigra","imagePullPolicy":"Never"}]}}}}'
```

### Option 4: Quick Fix Script

Save this as `build-and-deploy-arm.sh`:

```bash
#!/bin/bash
set -e

echo "Building TRIGRA for ARM64..."

# Clone if not exists
if [ ! -d "trigra" ]; then
    git clone https://github.com/Taiwrash/trigra.git
fi

cd trigra

# Build the image
docker build -t taiwrash/trigra:latest .

echo "✅ Image built successfully!"
echo ""
echo "Now restart your deployment:"
echo "  kubectl rollout restart deployment trigra"
echo ""
echo "Or if using imagePullPolicy: Never, just delete the pod:"
echo "  kubectl delete pod -l app=trigra"
```

Make it executable and run:

```bash
chmod +x build-and-deploy-arm.sh
./build-and-deploy-arm.sh
```

## Verify Architecture

After building, verify the image is ARM64:

```bash
docker inspect taiwrash/trigra:latest | grep Architecture
```

Should show: `"Architecture": "arm64"`

## Deploy

Once the image is built locally:

```bash
# Restart the deployment to use the new local image
kubectl rollout restart deployment trigra

# Watch it come up
kubectl get pods -l app=trigra -w
```

## Long-term Solution

I'll work on pushing a proper multi-platform image to Docker Hub that supports both AMD64 and ARM64. In the meantime, building locally is the fastest solution.
