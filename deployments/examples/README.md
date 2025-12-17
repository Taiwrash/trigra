# TRIGRA Examples Directory

This directory contains ready-to-use Kubernetes manifests that you can copy and paste to deploy applications using TRIGRA (Kubernetes GitOps Homelab).

## 📋 Available Examples

### Basic Examples

#### 1. **deployment.yaml**
Simple Nginx deployment demonstrating:
- Basic deployment configuration
- Resource limits and requests
- Container port mapping

```bash
kubectl apply -f deployment.yaml
```

#### 2. **service.yaml**
ClusterIP service example showing:
- Service selector configuration
- Port mapping
- Internal cluster networking

```bash
kubectl apply -f service.yaml
```

#### 3. **configmap.yaml**
ConfigMap with multiple data files:
- Application properties
- Database configuration
- Multi-file ConfigMap structure

```bash
kubectl apply -f configmap.yaml
```

---

### Complete Application Examples

#### 4. **homepage-dashboard.yaml** 📊
**Homepage - Application Dashboard for your Homelab**

[Homepage](https://gethomepage.dev/) is a highly customizable dashboard with 100+ service integrations.

Includes:
- ✅ ServiceAccount with RBAC permissions
- ✅ ConfigMap with customizable dashboard settings
- ✅ Kubernetes cluster monitoring widgets
- ✅ Service bookmarks and custom services
- ✅ Resource usage and metrics widgets
- ✅ Deployment with health checks
- ✅ ClusterIP service and Ingress

**Deploy:**
```bash
kubectl apply -f homepage-dashboard.yaml
```

**Access:**
- Update `host: dashboard.local` to your domain in the Ingress section
- Configure your DNS or `/etc/hosts` file
- Visit `http://dashboard.local`

**Customize:**
Edit the ConfigMap to:
- Add your own services and bookmarks
- Change theme and colors
- Configure widgets and layout
- Add service integrations (Plex, Sonarr, Radarr, etc.)

**Features:**
- 📊 Real-time Kubernetes cluster monitoring
- 🔖 Organize all your homelab services
- 🎨 Multiple themes (dark/light)
- 📈 CPU, memory, and network metrics
- 🔍 Integrated search (DuckDuckGo, Google, etc.)
- 🔗 100+ service integrations with status widgets

---

#### 5. **awesome-homepage.yaml** ✨
**Awesome Modern Homepage with Animations**

A visually stunning, animated homepage with modern design elements.

Includes:
- ✅ Animated particle background
- ✅ Glassmorphism design elements
- ✅ Service cards with hover effects
- ✅ Stats dashboard
- ✅ Smooth CSS animations
- ✅ Fully responsive layout
- ✅ Customizable service links

**Deploy:**
```bash
kubectl apply -f awesome-homepage.yaml
```

**Access:**
- Update `host: awesome.local` to your domain in the Ingress section
- Configure your DNS or `/etc/hosts` file
- Visit `http://awesome.local`

**Features:**
- 🎨 Modern gradient backgrounds
- ✨ 50 animated floating particles
- 💎 Glassmorphism UI elements
- 📊 Stats cards with pulse animations
- 🎯 Service grid with 6 customizable cards
- 🔗 Quick links section
- 📱 Fully responsive design

**Customize:**
Edit the ConfigMap HTML to:
- Change service names and descriptions
- Update service URLs in the JavaScript section
- Modify colors and gradients in CSS
- Add or remove service cards
- Customize stats and quick links

---

#### 6. **ollama-ai-model.yaml** 🤖
**Run AI models locally on your cluster**

Includes:
- ✅ 50GB PersistentVolumeClaim for model storage
- ✅ Ollama deployment with resource limits
- ✅ ClusterIP and NodePort services
- ✅ Ingress for API access
- ✅ GPU support (commented out, uncomment if available)
- ✅ Usage instructions ConfigMap

**Deploy:**
```bash
kubectl apply -f ollama-ai-model.yaml
```

**Quick Start:**
```bash
# Access the Ollama pod
kubectl exec -it deployment/ollama -- /bin/bash

# Pull a model (choose one)
ollama pull gemma:2b      # 1.4GB - Fast, good for testing
ollama pull llama2        # 3.8GB - General purpose
ollama pull gemma:7b      # 4.8GB - Better quality
ollama pull mistral       # 4.1GB - High performance
ollama pull codellama     # 3.8GB - Code generation

# Run the model interactively
ollama run gemma:2b
```

**Access Methods:**
1. **From within cluster:** `http://ollama:11434`
2. **NodePort:** `http://<node-ip>:30434`
3. **Ingress:** `http://ollama.local` (configure DNS)

**API Usage:**
```bash
curl http://ollama:11434/api/generate -d '{
  "model": "gemma:2b",
  "prompt": "Explain Kubernetes in simple terms"
}'
```

**Resource Requirements:**
- 2B models: 4GB RAM minimum
- 7B models: 8GB RAM minimum
- 13B models: 16GB RAM minimum
- 70B models: 64GB RAM minimum

---

## 🚀 Using with TRIGRA

All these examples work seamlessly with TRIGRA's GitOps workflow:

1. **Copy the example** to your Git repository
2. **Customize** as needed (change names, namespaces, resources)
3. **Commit and push** to your repository
4. **Watch TRIGRA** automatically apply the changes to your cluster

```bash
# Example workflow
cp deployments/examples/homepage-complete.yaml my-homepage.yaml
vim my-homepage.yaml  # Customize
git add my-homepage.yaml
git commit -m "Deploy my homepage"
git push origin main
# TRIGRA automatically applies the changes! 🎉
```

## 💡 Tips

- **Namespaces:** All examples use `namespace: default`. Change this to organize your applications.
- **Resource Limits:** Adjust CPU and memory based on your cluster capacity.
- **Storage:** Update `storageClassName` in PVCs to match your cluster's storage classes.
- **Ingress:** Configure `ingressClassName` and `host` based on your ingress controller.
- **Combine Examples:** Use `---` separator to combine multiple resources in one file.

## 🔧 Troubleshooting

**Pods not starting?**
```bash
kubectl describe pod <pod-name>
kubectl logs <pod-name>
```

**Service not accessible?**
```bash
kubectl get svc
kubectl get endpoints
```

**Ingress not working?**
```bash
kubectl get ingress
kubectl describe ingress <ingress-name>
```

## 📚 Learn More

- [Main README](../../README.md)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Ollama Documentation](https://ollama.ai/)

---

**Happy deploying! 🚀**
