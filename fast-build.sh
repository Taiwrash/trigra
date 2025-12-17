#!/bin/bash
# Fast multi-platform build using local Go compilation
# Usage: ./fast-build.sh

set -e

IMAGE_NAME="taiwrash/trigra"
# Generate a timestamp-based tag
TAG="v$(date +%Y%m%d-%H%M%S)"

echo "🚀 Fast Multi-Platform Build"
echo "=============================="
echo "🏷️  Tag: $TAG"
echo ""

# Build for AMD64 (Linux servers)
echo "📦 Building AMD64 binary..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
    -ldflags="-w -s -X main.Version=$TAG -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o bin/trigra-amd64 ./cmd/trigra

# Build for ARM64 (Mac, Raspberry Pi)
echo "📦 Building ARM64 binary..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build \
    -ldflags="-w -s -X main.Version=$TAG -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o bin/trigra-arm64 ./cmd/trigra

echo ""
echo "✅ Binaries built successfully!"
ls -lh bin/

echo ""
echo "🐳 Building and pushing multi-platform images..."

# Build and push AMD64 image
echo "  → AMD64 image..."
docker buildx build \
    --platform linux/amd64 \
    --build-arg TARGETARCH=amd64 \
    -t ${IMAGE_NAME}:latest-amd64 \
    -t ${IMAGE_NAME}:${TAG}-amd64 \
    -f Dockerfile.fast \
    --push .

# Build and push ARM64 image
echo "  → ARM64 image..."
docker buildx build \
    --platform linux/arm64 \
    --build-arg TARGETARCH=arm64 \
    -t ${IMAGE_NAME}:latest-arm64 \
    -t ${IMAGE_NAME}:${TAG}-arm64 \
    -f Dockerfile.fast \
    --push .

echo ""
echo "📋 Creating multi-platform manifest..."
docker buildx imagetools create \
    -t ${IMAGE_NAME}:latest \
    -t ${IMAGE_NAME}:${TAG} \
    ${IMAGE_NAME}:latest-amd64 \
    ${IMAGE_NAME}:latest-arm64

echo ""
echo "✅ Multi-platform image pushed to Docker Hub!"
echo "   Tags: latest, $TAG"

echo ""
echo "📝 Updating deployment manifest..."
DEPLOYMENT_FILE="deployments/kubernetes/deployment.yaml"
if [ -f "$DEPLOYMENT_FILE" ]; then
    # Use sed to replace the image tag
    # This regex looks for 'image: taiwrash/trigra:.*' and replaces it with the new tag
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS sed requires an empty string for -i
        sed -i '' "s|image: ${IMAGE_NAME}:.*|image: ${IMAGE_NAME}:${TAG}|" "$DEPLOYMENT_FILE"
    else
        sed -i "s|image: ${IMAGE_NAME}:.*|image: ${IMAGE_NAME}:${TAG}|" "$DEPLOYMENT_FILE"
    fi
    echo "   Updated $DEPLOYMENT_FILE with image: ${IMAGE_NAME}:${TAG}"
else
    echo "   ⚠️  Deployment file not found: $DEPLOYMENT_FILE"
fi

echo ""
echo "🧪 Test on AMD64 (NixOS):"
echo "   kubectl apply -f $DEPLOYMENT_FILE"

