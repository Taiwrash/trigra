# Trigra Examples Directory

This directory contains ready-to-use Kubernetes manifests that demonstrate the power of Trigra's GitOps flow. Copy these to your repository to see Trigra in action!

## 📋 Available Examples

### 1. **Core Resources**

- **`deployment.yaml`**: Simple Nginx deployment with resource limits.
- **`service.yaml`**: Standard ClusterIP service for internal networking.
- **`configmap.yaml`**: Configuration management with multi-file data.

### 2. **Complete Applications**

#### 📊 Homepage Dashboard (`homepage-dashboard.yaml`)
A stunning, high-performance dashboard for your cluster:
- **Features**: Real-time K8s stats, sleek dark mode, and integrated service bookmarks.
- **Components**: Deployment, Service, ConfigMap (with full HTML/JS/CSS), and Ingress.

#### ✨ Awesome Homepage (`awesome-homepage.yaml`)
A modern, animated landing page with glassmorphism and particles:
- **Features**: 60fps animations, responsive design, and easy customization via ConfigMap.
- **Components**: Full-stack application packaged as a single K8s manifest.

#### 🤖 AI with Ollama (`ollama-ai-model.yaml`)
Run large language models locally on your cluster:
- **Features**: Persistent storage for models, API access via Ingress, and GPU support.
- **Components**: PVC (50GB), Deployment, Service, and Ingress.

## 🚀 How to use with Trigra

1. **Pick an example**:
   ```bash
   cp deployments/examples/homepage-dashboard.yaml my-dashboard.yaml
   ```

2. **Customize if needed**:
   Open the file and update namespaces, image tags, or configuration values.

3. **Commit and Push**:
   ```bash
   git add my-dashboard.yaml
   git commit -m "feat: deploy homepage dashboard via Trigra"
   git push origin main
   ```

4. **Watch Trigra Sync**:
   Monitor the Trigra logs:
   ```bash
   kubectl logs -f deployment/trigra
   ```

## 🔧 Tips for Best Results

- **Namespaces**: Always check the `metadata.namespace` field. By default, most examples use `default`.
- **Ingress**: If you use an Ingress controller, update the `hosts` and `ingressClassName` to match your environment.
- **Storage**: For the AI model (Ollama), ensure your cluster has a default `StorageClass` or specify one in the `PersistentVolumeClaim`.

---
Return to [Main Documentation](../../README.md).
