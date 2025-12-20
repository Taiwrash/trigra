# Trigra - Kubernetes GitOps Controller

[![Security & Tests](https://github.com/Taiwrash/trigra/actions/workflows/security.yml/badge.svg)](https://github.com/Taiwrash/trigra/actions/workflows/security.yml)

A lightweight GitOps controller for Kubernetes clusters that automatically applies changes from your Git repository. Edit YAML files, commit to Git, and watch your cluster update automatically!

## 🚀 Features

- **Universal Resource Support**: Deploy any Kubernetes resource type (Deployments, Services, ConfigMaps, Secrets, StatefulSets, DaemonSets, Jobs, CronJobs, Ingress, PVCs, and more)
- **GitHub Integration**: Webhook-driven updates triggered by Git pushes
- **Smart Deployment**: Automatically creates or updates resources based on current cluster state
- **Multi-Document Support**: Process multiple YAML resources in a single file
- **In-Cluster & Local**: Auto-detects running environment (in-cluster vs local development)
- **Health Checks**: Built-in liveness and readiness probes
- **Graceful Shutdown**: Proper signal handling for clean shutdowns

## 📋 Prerequisites

- Kubernetes cluster (homelab, Minikube, Kind, K3s, etc.)
- Go 1.25+ (for building from source)
- GitHub repository for your manifests
- `kubectl` configured to access your cluster

## 🏗️ Architecture

```
┌─────────────┐
│   GitHub    │
│  Repository │
└──────┬──────┘
       │ Push Event
       │ (Webhook)
       ▼
┌─────────────────────┐
│ Trigra - GitOps Controller      │
│                     │
│ ┌─────────────────┐ │
│ │ Webhook Handler │ │
│ └────────┬────────┘ │
│          │          │
│ ┌────────▼────────┐ │
│ │ Resource Applier│ │
│ └────────┬────────┘ │
└──────────┼──────────┘
           │
           ▼
    ┌─────────────┐
    │ Kubernetes  │
    │   Cluster   │
    └─────────────┘
```

## 🛠️ Installation

### Option 1: Local Development Mode

1. **Clone the repository**:
   ```bash
   git clone https://github.com/Taiwrash/trigra.git
   cd trigra
   ```

2. **Set up environment variables**:
   ```bash
   cp .env.example .env
   # Edit .env with your values
   ```

3. **Generate a webhook secret**:
   ```bash
   openssl rand -hex 32
   ```

4. **Configure `.env`**:
   ```env
   GITHUB_TOKEN=ghp_your_token_here
   WEBHOOK_SECRET=your_generated_secret_here
   SERVER_PORT=8082
   NAMESPACE=default
   ```

5. **Build and run**:
   ```bash
   go build -o trigra ./cmd/trigra
   ./trigra
   ```

### Option 2: In-Cluster Deployment (Recommended for Production)

> **📖 For detailed deployment instructions**, see the [Kubernetes Deployment Guide](deployments/kubernetes/README.md)

**Quick Start:**

1. **Using Docker Hub image** (no build required):
   ```bash
   # Download manifests
   curl -O https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/kubernetes/secret.yaml.example
   curl -O https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/kubernetes/deployment.yaml
   curl -O https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/kubernetes/rbac.yaml
   curl -O https://raw.githubusercontent.com/Taiwrash/trigra/main/deployments/kubernetes/service.yaml
   
   # Create your secret
   cp secret.yaml.example secret.yaml
   vim secret.yaml  # Add your GitHub token and repo
   
   # Deploy
   kubectl apply -f secret.yaml
   kubectl apply -f rbac.yaml
   kubectl apply -f deployment.yaml
   kubectl apply -f service.yaml
   ```

2. **Or build from source**:
   ```bash
   # Build Docker image
   ./build-docker.sh
   
   # Create secret
   cp deployments/kubernetes/secret.yaml.example deployments/kubernetes/secret.yaml
   vim deployments/kubernetes/secret.yaml
   
   # Deploy
   kubectl apply -f deployments/kubernetes/
   ```

3. **Verify deployment**:
   ```bash
   kubectl get pods -l app=trigra
   kubectl logs -l app=trigra -f
   ```

## 🔧 GitHub Webhook Configuration

1. Go to your GitHub repository → **Settings** → **Webhooks** → **Add webhook**

2. Configure:
   - **Payload URL**: `http://your-controller-ip/webhook`
   - **Content type**: `application/json`
   - **Secret**: Your `WEBHOOK_SECRET` value
   - **Events**: Select "Just the push event"
   - **Active**: ✓ Checked

3. Click **Add webhook**

4. Test by pushing a YAML file to your repository!

## 🎯 Quick Example Deployment

Want to try TRIGRA with ready-made examples? Deploy any example with one command:

```bash
curl -fsSL https://raw.githubusercontent.com/Taiwrash/trigra/main/deploy-example.sh | bash
```

This interactive script lets you choose from:
- 📊 **Homepage Dashboard** - Full homelab dashboard
- 🤖 **Ollama AI** - Run LLMs locally
- 🏠 **Simple Homepage** - Clean landing page
- ✨ **Awesome Homepage** - Stunning animated design
- 🚀 **All Examples** - Deploy everything!

## 📝 Usage

### Basic Workflow

1. **Create a Kubernetes manifest** in your Git repository:
   ```yaml
   # nginx-deployment.yaml
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: nginx
     namespace: default
   spec:
     replicas: 3
     selector:
       matchLabels:
         app: nginx
     template:
       metadata:
         labels:
           app: nginx
       spec:
         containers:
         - name: nginx
           image: nginx:1.25-alpine
           ports:
           - containerPort: 80
   ```

2. **Commit and push**:
   ```bash
   git add nginx-deployment.yaml
   git commit -m "Deploy nginx"
   git push origin main
   ```

3. **Watch the magic happen**:
   ```bash
   kubectl get deployments -w
   ```

### Multi-Resource Files

You can include multiple resources in a single YAML file using `---` separator:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  environment: production
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  # ... deployment spec
---
apiVersion: v1
kind: Service
metadata:
  name: myapp
spec:
  # ... service spec
```

### Supported Resource Types

The controller supports **all Kubernetes resource types**, including:

- **Workloads**: Deployments, StatefulSets, DaemonSets, Jobs, CronJobs, Pods
- **Services**: Service, Ingress, NetworkPolicy
- **Config**: ConfigMap, Secret
- **Storage**: PersistentVolumeClaim, PersistentVolume, StorageClass
- **RBAC**: Role, RoleBinding, ClusterRole, ClusterRoleBinding, ServiceAccount
- **Custom Resources**: Any CRD installed in your cluster

## 🔍 Monitoring

### Health Checks

- **Liveness**: `http://controller:8082/health`
- **Readiness**: `http://controller:8082/ready`

### Logs

```bash
# View controller logs
kubectl logs -f deployment/trigra

# Follow logs in real-time
kubectl logs -f -l app=trigra
```

## 🔒 Security Best Practices

1. **Never commit secrets**: Use `.gitignore` to exclude `.env` and `secret.yaml` files
2. **Use Kubernetes Secrets**: Store sensitive data in Kubernetes secrets, not in Git
3. **Limit RBAC permissions**: Adjust `rbac.yaml` to grant only necessary permissions
4. **Use private repositories**: Keep your infrastructure manifests in private repos
5. **Rotate tokens**: Regularly rotate your GitHub PAT and webhook secrets

## 🐛 Troubleshooting

### Controller not receiving webhooks

1. Check webhook delivery in GitHub (Settings → Webhooks → Recent Deliveries)
2. Verify the service is accessible: `kubectl get svc trigra`
3. Check controller logs: `kubectl logs -f deployment/trigra`
4. Ensure webhook secret matches in both GitHub and Kubernetes secret

### Resources not applying

1. Check controller logs for errors
2. Verify RBAC permissions: `kubectl auth can-i create deployments --as=system:serviceaccount:default:trigra`
3. Validate YAML syntax: `kubectl apply --dry-run=client -f your-file.yaml`

### Permission denied errors

Update RBAC permissions in `deployments/kubernetes/rbac.yaml` to include the required resources.

## 📚 Examples

See the `deployments/examples/` directory for ready-to-use manifests:

### Basic Examples
- **`deployment.yaml`** - Simple Nginx deployment with resource limits
- **`service.yaml`** - ClusterIP service example
- **`configmap.yaml`** - ConfigMap with multiple data files

### Complete Application Examples

#### 📊 Homepage Dashboard (`homepage-dashboard.yaml`)
Deploy [Homepage](https://gethomepage.dev/) - a highly customizable application dashboard for your homelab:
- ServiceAccount with proper RBAC permissions
- ConfigMap with customizable dashboard settings
- Kubernetes integration to monitor your cluster
- Widgets for cluster metrics, resources, and search
- Service bookmarks and custom services
- Ingress for external access

**Quick Deploy:**
```bash
kubectl apply -f deployments/examples/homepage-dashboard.yaml
# Access via http://dashboard.local (configure your ingress/DNS)
```

**Features:**
- 📊 Real-time Kubernetes cluster monitoring
- 🔖 Customizable service bookmarks
- 🎨 Multiple themes and layouts
- 📈 Resource usage widgets
- 🔍 Integrated search
- 🔗 100+ service integrations

#### 🤖 AI Model with Ollama (`ollama-ai-model.yaml`)
Run local LLMs (Llama, Gemma, Mistral, etc.) on your cluster:
- PersistentVolumeClaim for model storage (50GB)
- Deployment with resource limits (configurable for GPU)
- Both ClusterIP and NodePort services
- Ingress for API access
- Usage instructions ConfigMap

**Quick Deploy:**
```bash
kubectl apply -f deployments/examples/ollama-ai-model.yaml

# Access the pod and pull a model
kubectl exec -it deployment/ollama -- ollama pull gemma:2b
kubectl exec -it deployment/ollama -- ollama run gemma:2b

# Or access via NodePort: http://<node-ip>:30434
```

**Supported Models:** llama2, gemma (2b/7b), mistral, codellama, and more!

## 🤝 Contributing

We love your input! We want to makeBy contributing to Trigra, you agree that your contributions will be licensed under its MIT License.porting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features

Please read our [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## 🌈 Code of Conduct

We are committed to making our community inclusive and welcoming. Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md).

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

Built with:
- [client-go](https://github.com/kubernetes/client-go) - Kubernetes Go client
- [go-github](https://github.com/google/go-github) - GitHub API client
- [oauth2](https://golang.org/x/oauth2) - OAuth2 authentication

---

**Happy GitOps-ing! 🚀**
