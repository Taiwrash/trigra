#!/bin/bash
# Super simple build - build locally on Mac, then user builds on NixOS
set -e

echo "🚀 Building TRIGRA Docker Image (Local Platform)"
echo "=============================================="
echo ""

# Build for current platform
docker build -t taiwrash/trigra:latest .

echo ""
echo "✅ Image built for $(uname -m) architecture!"
echo ""
echo "📤 Pushing to Docker Hub..."
docker push taiwrash/trigra:latest

echo ""
echo "✅ Done!"
echo ""
echo "⚠️  IMPORTANT: This image is for ARM64 (Mac)"
echo ""
echo "📋 On your NixOS server (AMD64), run:"
echo "   git clone https://github.com/Taiwrash/trigra.git"
echo "   cd trigra"
echo "   docker build -t trigra:local ."
echo "   docker run -p 8082:8082 \\"
echo "     -v \$HOME/.kube/config:/app/.kube/config:ro \\"
echo "     -e KUBECONFIG=/app/.kube/config \\"
echo "     -e WEBHOOK_SECRET=\$WEBHOOK_SECRET \\"
echo "     trigra:local"
