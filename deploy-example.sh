#!/bin/bash
# TRIGRA Example Deployment Script
# One-command deployment for TRIGRA examples

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Banner
echo -e "${BLUE}"
cat << "EOF"
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║   TRIGRA - Kubernetes GitOps Homelab                        ║
║   Example Deployment Script                              ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo -e "${RED}❌ kubectl is not installed. Please install kubectl first.${NC}"
    exit 1
fi

# Check if cluster is accessible
if ! kubectl cluster-info &> /dev/null; then
    echo -e "${RED}❌ Cannot connect to Kubernetes cluster. Please check your kubeconfig.${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Connected to Kubernetes cluster${NC}\n"

# Display menu
echo -e "${BLUE}📦 Available Examples:${NC}\n"
echo "1. 📊 Homepage Dashboard - Full homelab dashboard with cluster monitoring"
echo "2. 🤖 Ollama AI Models - Run LLMs (Llama, Gemma, Mistral) locally"
echo "3. 🏠 Simple Homepage - Clean, minimal landing page"
echo "4. ✨ Awesome Homepage - Stunning animated homepage"
echo "5. 🚀 All Examples - Deploy everything!"
echo ""
read -p "Choose an example (1-5): " choice

# Base URL for raw GitHub files
BASE_URL="https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/examples"

# Function to deploy an example
deploy_example() {
    local name=$1
    local file=$2
    local url="${BASE_URL}/${file}"
    
    echo -e "\n${YELLOW}📥 Deploying ${name}...${NC}"
    
    if kubectl apply -f "$url"; then
        echo -e "${GREEN}✓ ${name} deployed successfully!${NC}"
        return 0
    else
        echo -e "${RED}❌ Failed to deploy ${name}${NC}"
        return 1
    fi
}

# Function to show access instructions
show_access() {
    local service=$1
    local host=$2
    
    echo -e "\n${BLUE}🌐 Access Instructions for ${service}:${NC}"
    echo -e "1. ${YELLOW}Port Forward:${NC}"
    echo -e "   kubectl port-forward svc/${service} 8080:80 --address 0.0.0.0"
    echo -e "   Then visit: http://<your-server-ip>:8080"
    echo ""
    echo -e "2. ${YELLOW}NodePort:${NC}"
    echo -e "   kubectl patch svc ${service} -p '{\"spec\":{\"type\":\"NodePort\"}}'"
    echo -e "   kubectl get svc ${service}"
    echo -e "   Then visit: http://<your-server-ip>:<node-port>"
    echo ""
    echo -e "3. ${YELLOW}Ingress:${NC}"
    echo -e "   Add to /etc/hosts: <ingress-ip> ${host}"
    echo -e "   Then visit: http://${host}"
    echo ""
}

# Deploy based on choice
case $choice in
    1)
        deploy_example "Homepage Dashboard" "homepage-dashboard.yaml"
        show_access "homepage" "dashboard.local"
        ;;
    2)
        deploy_example "Ollama AI Models" "ollama-ai-model.yaml"
        echo -e "\n${BLUE}🤖 Ollama Quick Start:${NC}"
        echo "kubectl exec -it deployment/ollama -- ollama pull gemma:2b"
        echo "kubectl exec -it deployment/ollama -- ollama run gemma:2b"
        show_access "ollama" "ollama.local"
        ;;
    3)
        deploy_example "Simple Homepage" "homepage-complete.yaml"
        show_access "homepage" "home.local"
        ;;
    4)
        deploy_example "Awesome Homepage" "awesome-homepage.yaml"
        show_access "awesome-homepage" "awesome.local"
        ;;
    5)
        echo -e "\n${YELLOW}🚀 Deploying all examples...${NC}\n"
        deploy_example "Homepage Dashboard" "homepage-dashboard.yaml"
        deploy_example "Ollama AI Models" "ollama-ai-model.yaml"
        deploy_example "Simple Homepage" "homepage-complete.yaml"
        deploy_example "Awesome Homepage" "awesome-homepage.yaml"
        
        echo -e "\n${GREEN}✓ All examples deployed!${NC}"
        echo -e "\n${BLUE}📋 Deployed Services:${NC}"
        kubectl get svc | grep -E "homepage|ollama"
        ;;
    *)
        echo -e "${RED}❌ Invalid choice. Please run the script again.${NC}"
        exit 1
        ;;
esac

echo -e "\n${GREEN}════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}✓ Deployment Complete!${NC}"
echo -e "${GREEN}════════════════════════════════════════════════════════${NC}\n"

# Show all deployed resources
echo -e "${BLUE}📊 Deployed Resources:${NC}"
kubectl get pods,svc,ingress 2>/dev/null | grep -E "homepage|ollama" || echo "No resources found"

echo -e "\n${YELLOW}💡 Tip:${NC} View all resources with: kubectl get all -A | grep -E 'homepage|ollama'"
echo -e "${YELLOW}💡 Tip:${NC} Check logs with: kubectl logs -l app=<app-name> -f"
echo ""
