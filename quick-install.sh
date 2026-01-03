#!/bin/bash
# TRIGRA - Non-Interactive Quick Install
REPO_URL="https://github.com/Taiwrash/trigra"
BINARY_NAME="trigra"
# Usage: curl -fsSL https://raw.githubusercontent.com/Taiwrash/trigra/main/quick-install.sh | bash -s -- <namespace> <webhook-secret>

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}╔═══════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║   TRIGRA - Quick Install          ║${NC}"
echo -e "${GREEN}╚═══════════════════════════════════════════════╝${NC}"
echo ""

# Parse arguments
NAMESPACE="$1"
WEBHOOK_SECRET="$2"

# Generate webhook secret if not provided
if [ -z "$WEBHOOK_SECRET" ]; then
    if command -v openssl &> /dev/null; then
        WEBHOOK_SECRET=$(openssl rand -hex 32)
    elif [ -f /dev/urandom ]; then
        WEBHOOK_SECRET=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 64 | head -n 1)
    else
        WEBHOOK_SECRET=$(date +%s%N | sha256sum 2>/dev/null | head -c 64 || date +%s | md5sum | head -c 64)
    fi
    echo -e "${GREEN}Generated webhook secret ${NC}"
fi

GITHUB_TOKEN="${3:-}"

# Check kubectl
if ! command -v kubectl &> /dev/null; then
    echo -e "${RED}✗ kubectl not found${NC}"
    exit 1
fi
echo -e "${GREEN}✓ kubectl found${NC}"

# Check cluster
if ! kubectl cluster-info &> /dev/null; then
    echo -e "${RED}✗ Cannot connect to cluster${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Connected to cluster${NC}"

echo ""
echo -e "${YELLOW}Installing to namespace: ${NAMESPACE}${NC}"

# Create namespace
kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f - 2>/dev/null || true

# Create temporary directory for manifests
TEMP_DIR=$(mktemp -d)
echo -e "${YELLOW}Downloading manifests to ${TEMP_DIR}...${NC}"

# Function to download file
download_file() {
    local url=$1
    local filename=$2
    if command -v curl &> /dev/null; then
        curl -fsSL "$url" -o "${TEMP_DIR}/${filename}"
    elif command -v wget &> /dev/null; then
        wget -q "$url" -O "${TEMP_DIR}/${filename}"
    else
        echo -e "${RED}✗ Neither curl nor wget found${NC}"
        rm -rf "$TEMP_DIR"
        exit 1
    fi
}

# Download manifests
download_file "https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/kubernetes/rbac.yaml" "rbac.yaml"
download_file "https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/kubernetes/deployment.yaml" "deployment.yaml"
download_file "https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/kubernetes/service.yaml" "service.yaml"

# Apply RBAC
echo "Applying RBAC..."
kubectl apply -f "${TEMP_DIR}/rbac.yaml"

# Create secret
echo "Creating secret..."
kubectl create secret generic trigra-secret \
    --from-literal=GITHUB_TOKEN="$GITHUB_TOKEN" \
    --from-literal=WEBHOOK_SECRET="$WEBHOOK_SECRET" \
    --namespace="$NAMESPACE" \
    --dry-run=client -o yaml | kubectl apply -f -

# Apply deployment
echo "Deploying controller..."
kubectl apply -f "${TEMP_DIR}/deployment.yaml" -n "$NAMESPACE"
kubectl apply -f "${TEMP_DIR}/service.yaml" -n "$NAMESPACE"

# Cleanup
rm -rf "$TEMP_DIR"

echo ""
echo -e "${GREEN}╔═══════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║         Installation Complete! 🎉             ║${NC}"
echo -e "${GREEN}╚═══════════════════════════════════════════════╝${NC}"
echo ""

# Wait for deployment
echo "Waiting for deployment..."
kubectl wait --for=condition=available --timeout=60s deployment/trigra -n "$NAMESPACE" 2>/dev/null || true

# Get service info
echo ""
echo -e "${YELLOW}Webhook Configuration:${NC}"
echo "  Secret: $WEBHOOK_SECRET"
echo ""
# Get Service IP
# Get Service IP (Try LoadBalancer first, then ClusterIP)
SERVICE_IP=$(kubectl get svc trigra -n "$NAMESPACE" -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
if [ -z "$SERVICE_IP" ]; then
    SERVICE_IP=$(kubectl get svc trigra -n "$NAMESPACE" -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
fi
if [ -z "$SERVICE_IP" ]; then
    SERVICE_IP=$(kubectl get svc trigra -n "$NAMESPACE" -o jsonpath='{.spec.clusterIP}')
fi
PORT=$(kubectl get svc trigra -n "$NAMESPACE" -o jsonpath='{.spec.ports[0].port}')

echo "Get webhook URL with:"
echo "  kubectl get svc trigra -n $NAMESPACE"
echo ""

# Cloudflare Tunnel Setup
install_cloudflared() {
    echo ""
    echo -e "${YELLOW}Checking Cloudflare Tunnel (cloudflared)...${NC}"
    
    if command -v cloudflared &> /dev/null; then
        echo -e "${GREEN}✓ cloudflared is already installed${NC}"
        return
    fi
    
    echo "Installing cloudflared..."
    
    if [[ "$OSTYPE" == "darwin"* ]]; then
        if command -v brew &> /dev/null; then
            brew install cloudflare/cloudflare/cloudflared
        else
             echo -e "${YELLOW}! Homebrew not found. skipping cloudflared installation.${NC}"
             return
        fi
    else
        # Detect architecture
        ARCH=$(uname -m)
        case $ARCH in
            x86_64) DEB_ARCH="amd64" ;;
            aarch64) DEB_ARCH="arm64" ;;
            armv7l) DEB_ARCH="armhf" ;;
            i386|i686) DEB_ARCH="386" ;;
            *)
                echo -e "${YELLOW}! Architecture $ARCH not supported for auto-install.${NC}"
                return
                ;;
        esac

        if command -v dpkg &> /dev/null; then
            if [ "$EUID" -ne 0 ]; then
                echo -e "${YELLOW}! Root privileges required for .deb install. Skipping.${NC}"
                return
            fi
            curl -L --output cloudflared.deb "https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-${DEB_ARCH}.deb"
            dpkg -i cloudflared.deb
            rm cloudflared.deb
        elif command -v rpm &> /dev/null; then
             if [ "$EUID" -ne 0 ]; then
                echo -e "${YELLOW}! Root privileges required for .rpm install. Skipping.${NC}"
                return
            fi
             RPM_ARCH=$DEB_ARCH
             if [ "$DEB_ARCH" == "amd64" ]; then RPM_ARCH="x86_64"; fi
             if [ "$DEB_ARCH" == "arm64" ]; then RPM_ARCH="aarch64"; fi
             
             curl -L --output cloudflared.rpm "https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-${RPM_ARCH}.rpm"
             rpm -ivh cloudflared.rpm
             rm cloudflared.rpm
        else
             # Binary install
             curl -L --output cloudflared "https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-${DEB_ARCH}"
             chmod +x cloudflared
             if [ "$EUID" -eq 0 ]; then
                mv cloudflared /usr/local/bin/
             else
                echo -e "${YELLOW}! Root required to move binary. Leaving in current dir.${NC}"
             fi
        fi
    fi
}

# Attempt installation
install_cloudflared

echo ""
echo -e "${GREEN}Next Steps:${NC}"
echo "1. Configure GitHub webhook with the URL below"
echo ""
echo -e "${YELLOW}Starting Cloudflare Tunnel as background service...${NC}"

# Create log file for cloudflared
CLOUDFLARED_LOG="/tmp/cloudflared-trigra.log"
CLOUDFLARED_PID_FILE="/tmp/cloudflared-trigra.pid"

# Kill any existing tunnel
if [ -f "$CLOUDFLARED_PID_FILE" ]; then
    OLD_PID=$(cat "$CLOUDFLARED_PID_FILE")
    kill "$OLD_PID" 2>/dev/null || true
    rm -f "$CLOUDFLARED_PID_FILE"
fi

# Start cloudflared in background
nohup cloudflared tunnel --url http://${SERVICE_IP}:${PORT} > "$CLOUDFLARED_LOG" 2>&1 &
CLOUDFLARED_PID=$!
echo "$CLOUDFLARED_PID" > "$CLOUDFLARED_PID_FILE"

# Wait for URL to be generated (max 30 seconds)
echo -e "${YELLOW}Waiting for tunnel URL...${NC}"
WEBHOOK_URL=""
for i in {1..30}; do
    if [ -f "$CLOUDFLARED_LOG" ]; then
        URL=$(grep -o 'https://[-a-z0-9]*\.trycloudflare\.com' "$CLOUDFLARED_LOG" 2>/dev/null | head -1)
        if [ -n "$URL" ]; then
            WEBHOOK_URL="${URL}/webhook"
            break
        fi
    fi
    sleep 1
done

if [ -n "$WEBHOOK_URL" ]; then
    echo ""
    echo -e "${GREEN}╔══════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║${NC}  ${YELLOW}🔗 YOUR WEBHOOK URL (copy this to GitHub):${NC}                     ${GREEN}║${NC}"
    echo -e "${GREEN}╠══════════════════════════════════════════════════════════════════╣${NC}"
    echo -e "${GREEN}║${NC}                                                                  ${GREEN}║${NC}"
    echo -e "${GREEN}║${NC}  \033[1;36m${WEBHOOK_URL}\033[0m"
    echo -e "${GREEN}║${NC}                                                                  ${GREEN}║${NC}"
    echo -e "${GREEN}╚══════════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${GREEN}✓ Tunnel running in background (PID: ${CLOUDFLARED_PID})${NC}"
    echo ""
    echo -e "${YELLOW}Useful commands:${NC}"
    echo "  View logs:    tail -f $CLOUDFLARED_LOG"
    echo "  Stop tunnel:  kill \$(cat $CLOUDFLARED_PID_FILE)"
    echo ""
else
    echo -e "${RED}✗ Failed to get tunnel URL. Check logs: $CLOUDFLARED_LOG${NC}"
    echo "  You can also run manually: cloudflared tunnel --url http://${SERVICE_IP}:${PORT}"
fi

# Anonymous telemetry to help improve TRIGRA (completely anonymous)
# This helps us understand which environments TRIGRA is used in.
(curl -s -X POST "https://telemetry.trigra.dev/v1/install" \
    -H "Content-Type: application/json" \
    -d "{
        \"os\": \"$(uname -s)\",
        \"arch\": \"$(uname -m)\",
        \"installer\": \"quick-install\"
    }" -o /dev/null 2>&1 &) || true





