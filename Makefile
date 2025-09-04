# License Manager Makefile

# Variables
APP_NAME = license-manager
DOCKER_IMAGE = $(APP_NAME):latest
NAMESPACE = default
HELM_RELEASE = license-manager

# Docker commands
.PHONY: build
build:
	docker build -t $(DOCKER_IMAGE) .

.PHONY: clean
clean:
	docker rmi $(DOCKER_IMAGE) 2>/dev/null || true

# Minikube commands
.PHONY: minikube-start
minikube-start:
	minikube start
	minikube addons enable ingress
	minikube addons enable ingress-dns

.PHONY: minikube-stop
minikube-stop:
	minikube stop

.PHONY: minikube-status
minikube-status:
	minikube status

.PHONY: minikube-dashboard
minikube-dashboard:
	minikube dashboard

# Helm commands
.PHONY: helm-install
helm-install:
	helm install $(HELM_RELEASE) ./helm-charts/license-manager -n $(NAMESPACE)

.PHONY: helm-upgrade
helm-upgrade:
	helm upgrade $(HELM_RELEASE) ./helm-charts/license-manager -n $(NAMESPACE)

.PHONY: helm-uninstall
helm-uninstall:
	helm uninstall $(HELM_RELEASE) -n $(NAMESPACE)

.PHONY: helm-status
helm-status:
	helm status $(HELM_RELEASE) -n $(NAMESPACE)

.PHONY: helm-list
helm-list:
	helm list -n $(NAMESPACE)

# Development commands
.PHONY: dev-setup
dev-setup: minikube-start
	@echo "Setting up development environment..."
	@echo "Minikube started. You can now run 'make helm-install' to deploy the application."

.PHONY: dev-deploy
dev-deploy: build helm-upgrade
	@echo "Application deployed to minikube"
	@echo "Waiting for ingress to be ready..."
	@kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=license-manager --timeout=300s
	@echo "Adding license-manager.local to /etc/hosts..."
	@echo "$(shell minikube ip) license-manager.local" | sudo tee -a /etc/hosts > /dev/null || echo "Please add '$(shell minikube ip) license-manager.local' to your /etc/hosts file"
	@echo "Access the application at: http://license-manager.local"

.PHONY: dev-cleanup
dev-cleanup: helm-uninstall minikube-stop
	@echo "Development environment cleaned up"

# Test commands
.PHONY: test
test:
	@echo "Running all tests with Docker..."
	@docker run --rm -v $(PWD):/app -w /app golang:1.21-alpine go test ./tests/...

.PHONY: test-unit
test-unit:
	@echo "Running unit tests..."
	@docker run --rm -v $(PWD):/app -w /app golang:1.21-alpine go test ./tests/unit/...

.PHONY: test-integration
test-integration:
	@echo "Running integration tests..."
	@docker run --rm -v $(PWD):/app -w /app golang:1.21-alpine go test ./tests/integration/...

.PHONY: test-verbose
test-verbose:
	@echo "Running tests with verbose output..."
	@docker run --rm -v $(PWD):/app -w /app golang:1.21-alpine go test -v ./tests/...

.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@docker run --rm -v $(PWD):/app -w /app golang:1.21-alpine go test -cover ./tests/...

.PHONY: test-benchmark
test-benchmark:
	@echo "Running benchmarks..."
	@docker run --rm -v $(PWD):/app -w /app golang:1.21-alpine go test -bench=. ./tests/...

# Utility commands
.PHONY: port-forward
port-forward:
	kubectl port-forward service/$(HELM_RELEASE) 8080:80

.PHONY: get-url
get-url:
	@echo "Application URL: http://license-manager.local"
	@echo "Minikube IP: $(shell minikube ip)"
	@echo "Make sure '$(shell minikube ip) license-manager.local' is in your /etc/hosts file"

.PHONY: check-ingress
check-ingress:
	@echo "Checking ingress status..."
	@kubectl get ingress
	@kubectl describe ingress $(HELM_RELEASE)

.PHONY: add-hosts
add-hosts:
	@echo "Adding license-manager.local to /etc/hosts..."
	@echo "$(shell minikube ip) license-manager.local" | sudo tee -a /etc/hosts > /dev/null || echo "Please add '$(shell minikube ip) license-manager.local' to your /etc/hosts file"

.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build          - Build Docker image"
	@echo "  clean          - Clean up Docker resources"
	@echo ""
	@echo "Minikube commands:"
	@echo "  minikube-start    - Start minikube cluster"
	@echo "  minikube-stop     - Stop minikube cluster"
	@echo "  minikube-status   - Show minikube status"
	@echo "  minikube-dashboard - Open minikube dashboard"
	@echo ""
	@echo "Helm commands:"
	@echo "  helm-install   - Install with Helm"
	@echo "  helm-upgrade   - Upgrade with Helm"
	@echo "  helm-uninstall - Uninstall with Helm"
	@echo "  helm-status    - Show Helm release status"
	@echo "  helm-list      - List Helm releases"
	@echo ""
	@echo "Testing:"
	@echo "  test           - Run all tests with Docker"
	@echo "  test-unit      - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  test-verbose   - Run tests with verbose output"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-benchmark - Run benchmark tests"
	@echo ""
	@echo "Development:"
	@echo "  dev-setup      - Setup development environment"
	@echo "  dev-deploy     - Deploy to minikube"
	@echo "  dev-cleanup    - Cleanup development environment"
	@echo "  port-forward   - Port forward to local machine"
	@echo "  get-url        - Get application URL"
	@echo "  check-ingress  - Check ingress status"
	@echo "  add-hosts      - Add license-manager.local to /etc/hosts"
