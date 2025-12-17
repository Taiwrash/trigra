.PHONY: build run clean test docker-build docker-push deploy help

# Variables
BINARY_NAME=trigra
DOCKER_IMAGE=trigra
DOCKER_TAG=latest
NAMESPACE=default

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the GitOps controller binary
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) ./cmd/trigra
	@echo "Build complete: $(BINARY_NAME)"

run: build ## Build and run the controller locally
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

clean: ## Remove build artifacts
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -f main
	@echo "Clean complete"

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

lint: fmt vet ## Run linters

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)"

docker-push: docker-build ## Push Docker image to registry
	@echo "Pushing Docker image..."
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)

deploy-secret: ## Deploy Kubernetes secret (requires secret.yaml)
	@if [ ! -f deployments/kubernetes/secret.yaml ]; then \
		echo "Error: deployments/kubernetes/secret.yaml not found"; \
		echo "Copy secret.yaml.example and fill in your values"; \
		exit 1; \
	fi
	kubectl apply -f deployments/kubernetes/secret.yaml

deploy: ## Deploy controller to Kubernetes
	@echo "Deploying to Kubernetes..."
	kubectl apply -f deployments/kubernetes/rbac.yaml
	kubectl apply -f deployments/kubernetes/deployment.yaml
	kubectl apply -f deployments/kubernetes/service.yaml
	@echo "Deployment complete"

undeploy: ## Remove controller from Kubernetes
	@echo "Removing from Kubernetes..."
	kubectl delete -f deployments/kubernetes/service.yaml --ignore-not-found
	kubectl delete -f deployments/kubernetes/deployment.yaml --ignore-not-found
	kubectl delete -f deployments/kubernetes/rbac.yaml --ignore-not-found
	@echo "Undeployment complete"

logs: ## View controller logs
	kubectl logs -f -l app=trigra -n $(NAMESPACE)

status: ## Check controller status
	@echo "Controller Status:"
	kubectl get deployment trigra -n $(NAMESPACE)
	@echo ""
	@echo "Pods:"
	kubectl get pods -l app=trigra -n $(NAMESPACE)
	@echo ""
	@echo "Service:"
	kubectl get svc trigra -n $(NAMESPACE)

example-deploy: ## Deploy example resources
	@echo "Deploying example resources..."
	kubectl apply -f deployments/examples/
	@echo "Example deployment complete"

example-undeploy: ## Remove example resources
	@echo "Removing example resources..."
	kubectl delete -f deployments/examples/ --ignore-not-found
	@echo "Example removal complete"

dev-setup: ## Set up development environment
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo ".env file created. Please edit it with your values."; \
	else \
		echo ".env file already exists"; \
	fi
	@echo "Development setup complete"

all: clean lint build ## Clean, lint, and build
