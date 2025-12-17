#!/bin/bash
# KGH - Quick Install Script for Homelab
# This script makes it easy to install the GitOps controller on your Kubernetes cluster

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}╔═══════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║   KGH - Homelab Installation   ║${NC}"
echo -e "${GREEN}╚═══════════════════════════════════════════════╝${NC}"
echo ""

# Check prerequisites
echo -e "${YELLOW}Checking prerequisites...${NC}"

if ! command -v kubectl &> /dev/null; then
    echo -e "${RED}✗ kubectl not found. Please install kubectl first.${NC}"
    exit 1
fi
echo -e "${GREEN}✓ kubectl found${NC}"

if ! command -v helm &> /dev/null; then
    echo -e "${YELLOW}! Helm not found. Installing without Helm...${NC}"
    USE_HELM=false
else
    echo -e "${GREEN}✓ Helm found${NC}"
    USE_HELM=true
fi

# Check cluster connectivity
if ! kubectl cluster-info &> /dev/null; then
    echo -e "${RED}✗ Cannot connect to Kubernetes cluster${NC}"
    echo "Please ensure kubectl is configured correctly"
    exit 1
fi
echo -e "${GREEN}✓ Connected to Kubernetes cluster${NC}"
echo ""

# Get configuration from user
echo -e "${YELLOW}Configuration:${NC}"
echo ""

read -p "Enter your GitHub repository (e.g., username/repo): " GITHUB_REPO
read -p "Enter namespace [default]: " NAMESPACE
NAMESPACE=${NAMESPACE:-default}

# Generate webhook secret if not provided
read -p "Enter webhook secret (or press Enter to generate): " WEBHOOK_SECRET
if [ -z "$WEBHOOK_SECRET" ]; then
    # Try multiple methods to generate random secret
    if command -v openssl &> /dev/null; then
        WEBHOOK_SECRET=$(openssl rand -hex 32)
    elif [ -f /dev/urandom ]; then
        WEBHOOK_SECRET=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 64 | head -n 1)
    else
        # Fallback to date-based random
        WEBHOOK_SECRET=$(date +%s%N | sha256sum | head -c 64)
    fi
    echo -e "${GREEN}Generated webhook secret: ${WEBHOOK_SECRET}${NC}"
fi

read -p "Enter GitHub Personal Access Token (optional, press Enter to skip): " GITHUB_TOKEN

echo ""
echo -e "${YELLOW}Installation Summary:${NC}"
echo "  Repository: $GITHUB_REPO"
echo "  Namespace: $NAMESPACE"
echo "  Webhook Secret: ${WEBHOOK_SECRET:0:10}..."
echo ""

read -p "Proceed with installation? (y/n): " CONFIRM
if [ "$CONFIRM" != "y" ]; then
    echo "Installation cancelled"
    exit 0
fi

echo ""
echo -e "${YELLOW}Installing KGH...${NC}"

if [ "$USE_HELM" = true ]; then
    # Install using Helm
    echo "Installing with Helm..."
    
    helm upgrade --install kgh ./helm/kgh \
        --namespace "$NAMESPACE" \
        --create-namespace \
        --set github.webhookSecret="$WEBHOOK_SECRET" \
        --set github.token="$GITHUB_TOKEN" \
        --set namespace="$NAMESPACE"
    
    echo -e "${GREEN}✓ Installed with Helm${NC}"
else
    # Install using kubectl
    echo "Installing with kubectl..."
    
    # Create namespace if it doesn't exist
    kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -
    
    # Apply RBAC
    kubectl apply -f deployments/kubernetes/rbac.yaml
    
    # Create secret
    kubectl create secret generic kgh-secret \
        --from-literal=GITHUB_TOKEN="$GITHUB_TOKEN" \
        --from-literal=WEBHOOK_SECRET="$WEBHOOK_SECRET" \
        --namespace="$NAMESPACE" \
        --dry-run=client -o yaml | kubectl apply -f -
    
    # Apply deployment and service
    kubectl apply -f deployments/kubernetes/deployment.yaml
    kubectl apply -f deployments/kubernetes/service.yaml
    
    echo -e "${GREEN}✓ Installed with kubectl${NC}"
fi

echo ""
echo -e "${GREEN}╔═══════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║         Installation Complete! 🎉             ║${NC}"
echo -e "${GREEN}╚═══════════════════════════════════════════════╝${NC}"
echo ""

# Wait for deployment
echo -e "${YELLOW}Waiting for deployment to be ready...${NC}"
kubectl wait --for=condition=available --timeout=60s deployment/kgh -n "$NAMESPACE" || true

# Get service info
echo ""
echo -e "${YELLOW}Getting webhook URL...${NC}"
EXTERNAL_IP=""
for i in {1..30}; do
    EXTERNAL_IP=$(kubectl get svc kgh -n "$NAMESPACE" -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "")
    if [ -z "$EXTERNAL_IP" ]; then
        EXTERNAL_IP=$(kubectl get svc kgh -n "$NAMESPACE" -o jsonpath='{.status.loadBalancer.ingress[0].hostname}' 2>/dev/null || echo "")
    fi
    
    if [ ! -z "$EXTERNAL_IP" ]; then
        break
    fi
    sleep 2
done

echo ""
echo -e "${GREEN}Next Steps:${NC}"
echo ""
echo "1. Configure GitHub Webhook:"
echo "   - Go to: https://github.com/$GITHUB_REPO/settings/hooks"
echo "   - Click 'Add webhook'"
if [ ! -z "$EXTERNAL_IP" ]; then
    echo "   - Payload URL: http://$EXTERNAL_IP/webhook"
else
    echo "   - Payload URL: http://<EXTERNAL-IP>/webhook"
    echo "     (Get EXTERNAL-IP with: kubectl get svc kgh -n $NAMESPACE)"
fi
echo "   - Content type: application/json"
echo "   - Secret: $WEBHOOK_SECRET"
echo "   - Events: Just the push event"
echo ""
echo "2. Test the installation:"
echo "   kubectl logs -f deployment/kgh -n $NAMESPACE"
echo ""
echo "3. Push a YAML file to your repo and watch it deploy!"
echo ""
echo -e "${GREEN}Happy GitOps-ing! 🚀${NC}"