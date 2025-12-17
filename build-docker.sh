#!/bin/bash
# Build and push TRIGRA Docker image to Docker Hub
# Usage: 
#   ./build-docker.sh              # Build for local platform only
#   ./build-docker.sh --push       # Build multi-platform and push to Docker Hub
#   ./build-docker.sh v1.0.0       # Build specific version for local platform
#   ./build-docker.sh v1.0.0 --push # Build and push specific version

set -e

# Parse arguments
PUSH=false
VERSION="latest"

for arg in "$@"; do
    if [ "$arg" = "--push" ]; then
        PUSH=true
    elif [[ "$arg" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        VERSION="$arg"
    fi
done

IMAGE_NAME="taiwrash/trigra"
FULL_IMAGE="${IMAGE_NAME}:${VERSION}"

echo "🐳 Building Trigra Docker Image"
echo "================================"
echo "Image: $FULL_IMAGE"
echo "Push: $PUSH"
echo ""

if [ "$PUSH" = true ]; then
    # Multi-platform build for Docker Hub
    echo "📦 Building multi-platform image (amd64, arm64) and pushing..."
    docker buildx build --platform linux/amd64,linux/arm64 \
        -t "${FULL_IMAGE}" \
        --push .
    
    # Also tag as latest if building a version
    if [ "$VERSION" != "latest" ]; then
        echo "🏷️  Tagging and pushing as latest..."
        docker buildx build --platform linux/amd64,linux/arm64 \
            -t "${IMAGE_NAME}:latest" \
            --push .
    fi
    
    echo ""
    echo "✅ Multi-platform build pushed to Docker Hub!"
else
    # Local build for current platform only
    echo "📦 Building for local platform..."
    docker build -t "${FULL_IMAGE}" .
    
    # Also tag as latest if building a version
    if [ "$VERSION" != "latest" ]; then
        echo "🏷️  Tagging as latest..."
        docker tag "${FULL_IMAGE}" "${IMAGE_NAME}:latest"
    fi
    
    echo ""
    echo "✅ Build complete!"
    echo ""
    echo "📋 Image details:"
    docker images | grep "${IMAGE_NAME}" | head -2
    
    echo ""
    echo "🚀 To push to Docker Hub:"
    echo "   ./build-docker.sh --push"
    if [ "$VERSION" != "latest" ]; then
        echo "   ./build-docker.sh ${VERSION} --push"
    fi
fi

echo ""
echo "🧪 To test locally:"
echo "   # With kubeconfig:"
echo "   docker run -p 8082:8082 \\"
echo "     -v \$HOME/.kube/config:/app/.kube/config:ro \\"
echo "     -e KUBECONFIG=/app/.kube/config \\"
echo "     -e WEBHOOK_SECRET=test \\"
echo "     ${FULL_IMAGE}"
echo ""
echo "   # Or in-cluster mode (requires K8s deployment):"
echo "   kubectl apply -f deployments/kubernetes/"
