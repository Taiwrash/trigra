# Quick Start Guide

This is a quick reference for common tasks. See [README.md](README.md) for full documentation.

## Local Development

```bash
# 1. Setup
make dev-setup
# Edit .env with your values

# 2. Run locally
make run
```

## Deploy to Kubernetes

```bash
# 1. Build image
make docker-build

# 2. Create secret
cp deployments/kubernetes/secret.yaml.example deployments/kubernetes/secret.yaml
# Edit with your GitHub token and webhook secret

# 3. Deploy
kubectl apply -f deployments/kubernetes/secret.yaml
make deploy

# 4. Get webhook URL
kubectl get svc trigra
```

## Configure GitHub Webhook

1. Go to your repo → Settings → Webhooks → Add webhook
2. Payload URL: `http://<EXTERNAL-IP>/webhook`
3. Content type: `application/json`
4. Secret: Your `WEBHOOK_SECRET` value
5. Events: Just the push event

## Test GitOps Flow

```bash
# Create a deployment
cat > nginx.yaml <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-nginx
spec:
  replicas: 2
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
        image: nginx:alpine
EOF

# Push to Git
git add nginx.yaml
git commit -m "Deploy nginx"
git push

# Watch it deploy!
kubectl get deployments -w
```

## Useful Commands

```bash
make help          # Show all commands
make logs          # View controller logs
make status        # Check deployment status
make example-deploy # Deploy example resources
```

## Troubleshooting

**Webhook not working?**
- Check GitHub webhook deliveries
- View logs: `make logs`
- Verify secret matches

**Resources not applying?**
- Check RBAC permissions
- Validate YAML: `kubectl apply --dry-run=client -f file.yaml`
- Check controller logs

For more help, see [README.md](README.md#-troubleshooting)
