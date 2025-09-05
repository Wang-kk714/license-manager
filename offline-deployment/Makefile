# License Manager Makefile

# Variables
APP_NAME = license-manager
DOCKER_IMAGE = $(APP_NAME):latest
NAMESPACE = default
HELM_RELEASE = license-manager

# Main targets
.PHONY: help
help:
	@echo "License Manager - Available Commands:"
	@echo ""
	@echo "Development:"
	@echo "  dev-deploy     - Build and deploy to Minikube (recommended)"
	@echo "  dev-cleanup    - Clean up development environment"
	@echo "  get-url        - Show application URL and setup instructions"
	@echo ""
	@echo "Docker:"
	@echo "  build          - Build Docker image"
	@echo "  clean          - Clean up Docker resources"
	@echo ""
	@echo "Kubernetes:"
	@echo "  helm-upgrade   - Deploy/upgrade with Helm"
	@echo "  helm-uninstall - Uninstall from Kubernetes"
	@echo "  port-forward   - Port forward to local machine"
	@echo ""
	@echo "Testing:"
	@echo "  test           - Run all tests"
	@echo "  test-unit      - Run unit tests only"
	@echo "  test-integration - Run integration tests only"

# Development workflow
.PHONY: dev-deploy
dev-deploy: build helm-upgrade
	@echo "🚀 Deploying License Manager to Minikube..."
	@kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=license-manager --timeout=300s
	@echo "✅ Application deployed successfully!"
	@echo ""
	@echo "📋 Next steps:"
	@echo "1. Add to /etc/hosts: $(shell minikube ip) license-manager.local"
	@echo "2. Access: http://license-manager.local"
	@echo ""
	@echo "💡 Run 'make get-url' for detailed setup instructions"

.PHONY: dev-cleanup
dev-cleanup: helm-uninstall
	@echo "🧹 Development environment cleaned up"

.PHONY: get-url
get-url:
	@echo "🌐 License Manager Access Information:"
	@echo ""
	@echo "Minikube IP: $(shell minikube ip)"
	@echo "Application URL: http://license-manager.local"
	@echo ""
	@echo "📝 Setup Instructions:"
	@echo "1. Add this line to your /etc/hosts file:"
	@echo "   $(shell minikube ip) license-manager.local"
	@echo ""
	@echo "2. Then access: http://license-manager.local"
	@echo ""
	@echo "🔧 Quick setup:"
	@echo "   echo '$(shell minikube ip) license-manager.local' | sudo tee -a /etc/hosts"

# Docker commands
.PHONY: build
build:
	@echo "🔨 Building Docker image..."
	@eval $$(minikube docker-env) && docker build -t $(DOCKER_IMAGE) .
	@echo "✅ Docker image built successfully"

.PHONY: clean
clean:
	@echo "🧹 Cleaning up Docker resources..."
	@docker rmi $(DOCKER_IMAGE) 2>/dev/null || true
	@echo "✅ Cleanup completed"

# Kubernetes/Helm commands
.PHONY: helm-upgrade
helm-upgrade:
	@echo "📦 Deploying with Helm..."
	@helm upgrade --install $(HELM_RELEASE) ./helm-charts/license-manager --set image.tag=latest
	@echo "✅ Helm deployment completed"

.PHONY: helm-uninstall
helm-uninstall:
	@echo "🗑️  Uninstalling from Kubernetes..."
	@helm uninstall $(HELM_RELEASE) 2>/dev/null || true
	@echo "✅ Uninstall completed"

.PHONY: port-forward
port-forward:
	@echo "🔗 Port forwarding to localhost:8080..."
	@echo "Access at: http://localhost:8080"
	@kubectl port-forward service/$(HELM_RELEASE) 8080:80

# Testing
.PHONY: test
test:
	@echo "🧪 Running all tests..."
	@docker run --rm -v $(PWD):/app -w /app golang:1.21-alpine go test ./tests/...

.PHONY: test-unit
test-unit:
	@echo "🧪 Running unit tests..."
	@docker run --rm -v $(PWD):/app -w /app golang:1.21-alpine go test ./tests/unit/...

.PHONY: test-integration
test-integration:
	@echo "🧪 Running integration tests..."
	@docker run --rm -v $(PWD):/app -w /app golang:1.21-alpine go test ./tests/integration/...

# Utility commands
.PHONY: minikube-start
minikube-start:
	@echo "🚀 Starting Minikube..."
	@minikube start
	@minikube addons enable ingress
	@echo "✅ Minikube started with ingress enabled"

.PHONY: minikube-status
minikube-status:
	@minikube status

.PHONY: logs
logs:
	@kubectl logs -l app.kubernetes.io/name=license-manager --tail=50 -f