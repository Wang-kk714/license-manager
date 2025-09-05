#!/bin/bash

# License Manager Offline Deployment Script
# Optimized for offline environments

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${PURPLE}[HEADER]${NC} $1"
}

# Function to show usage
show_usage() {
    echo "License Manager Offline Deployment Script"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  install     Install License Manager (default)"
    echo "  uninstall   Uninstall License Manager"
    echo "  status      Show deployment status"
    echo "  logs        Show application logs"
    echo "  help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 install     # Install the application"
    echo "  $0 uninstall   # Remove the application"
    echo "  $0 status      # Check deployment status"
    echo "  $0 logs        # View application logs"
}

# Function to check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    # Check if Docker is running
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker daemon."
        exit 1
    fi
    
    # Check if kubectl is available
    if ! command -v kubectl &> /dev/null; then
        print_error "kubectl is not installed or not in PATH."
        exit 1
    fi
    
    # Check if helm is available
    if ! command -v helm &> /dev/null; then
        print_error "helm is not installed or not in PATH."
        exit 1
    fi
    
    print_success "Prerequisites check passed"
}

# Function to load Docker image
load_docker_image() {
    print_status "Loading existing Docker image..."
    if docker images | grep -q "license-manager.*fp216-offline-fixed"; then
        print_success "Docker image license-manager:fp216-offline-fixed already available"
    elif [ -f "license-manager-fp216.tar" ]; then
        docker load -i license-manager-fp216.tar
        docker tag license-manager:fp216 license-manager:fp216-offline-fixed
        print_success "Docker image loaded from license-manager-fp216.tar and tagged as fp216-offline-fixed"
    elif docker images | grep -q "license-manager.*fp216"; then
        docker tag $(docker images | grep "license-manager.*fp216" | head -1 | awk '{print $1":"$2}') license-manager:fp216-offline-fixed
        print_success "Docker image tagged as license-manager:fp216-offline-fixed"
    else
        print_error "No license-manager Docker image found. Please ensure license-manager-fp216.tar exists or the image is available."
        exit 1
    fi
}

# Function to install License Manager
install_license_manager() {
    print_header "🚀 Installing License Manager..."
    
    check_prerequisites
    load_docker_image
    
    # Deploy with Helm
    print_status "Deploying with Helm..."
    helm install --name license-manager-offline ./helm-charts/license-manager \
      --values values.yaml \
      --set image.repository=license-manager \
      --set image.tag=fp216-offline-fixed \
      --set image.pullPolicy=Never
    
    if [ $? -ne 0 ]; then
        print_error "Helm deployment failed"
        exit 1
    fi
    
    print_success "Helm deployment completed"
    
    # Wait for deployment to be ready
    print_status "Waiting for deployment to be ready..."
    kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=license-manager --timeout=300s
    
    if [ $? -ne 0 ]; then
        print_error "Deployment failed to become ready"
        exit 1
    fi
    
    print_success "Deployment is ready!"
    
    # Show deployment status
    echo ""
    print_status "Deployment Status:"
    kubectl get pods -l app.kubernetes.io/name=license-manager
    kubectl get services -l app.kubernetes.io/name=license-manager
    kubectl get ingress -l app.kubernetes.io/name=license-manager
    
    # Get cluster IP
    CLUSTER_IP=$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}')
    if [ -z "$CLUSTER_IP" ]; then
        CLUSTER_IP=$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="ExternalIP")].address}')
    fi
    
    echo ""
    print_success "License Manager deployed successfully!"
    echo ""
    print_status "Access Information:"
    echo "🌐 Application URL: http://${CLUSTER_IP}/license-manager"
    echo ""
    print_status "Useful Commands:"
    echo "📋 Check status: $0 status"
    echo "📋 View logs: $0 logs"
    echo "📋 Port forward: kubectl port-forward service/license-manager-offline 8080:80"
    echo "📋 Uninstall: $0 uninstall"
    echo ""
    print_status "The application is now ready for use!"
}

# Function to uninstall License Manager
uninstall_license_manager() {
    print_header "🗑️  Uninstalling License Manager..."
    
    # Check if deployment exists
    if ! helm list | grep -q "license-manager-offline"; then
        print_warning "License Manager is not currently installed."
        exit 0
    fi
    
    print_status "Removing License Manager deployment..."
    helm delete license-manager-offline --purge
    
    if [ $? -ne 0 ]; then
        print_error "Failed to uninstall License Manager"
        exit 1
    fi
    
    print_success "License Manager uninstalled successfully!"
    
    # Wait for pods to be terminated
    print_status "Waiting for pods to be terminated..."
    kubectl wait --for=delete pod -l app.kubernetes.io/name=license-manager --timeout=60s 2>/dev/null || true
    
    print_success "Cleanup completed!"
}

# Function to show deployment status
show_status() {
    print_header "📊 License Manager Status"
    
    if ! helm list | grep -q "license-manager-offline"; then
        print_warning "License Manager is not currently installed."
        echo ""
        print_status "To install, run: $0 install"
        exit 0
    fi
    
    echo ""
    print_status "Deployment Status:"
    kubectl get pods -l app.kubernetes.io/name=license-manager
    echo ""
    kubectl get services -l app.kubernetes.io/name=license-manager
    echo ""
    kubectl get ingress -l app.kubernetes.io/name=license-manager
    
    # Get cluster IP
    CLUSTER_IP=$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}')
    if [ -z "$CLUSTER_IP" ]; then
        CLUSTER_IP=$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="ExternalIP")].address}')
    fi
    
    echo ""
    print_status "Access Information:"
    echo "🌐 Application URL: http://${CLUSTER_IP}/license-manager"
}

# Function to show logs
show_logs() {
    print_header "📋 License Manager Logs"
    
    if ! helm list | grep -q "license-manager-offline"; then
        print_warning "License Manager is not currently installed."
        exit 0
    fi
    
    print_status "Showing recent logs (last 50 lines):"
    kubectl logs -l app.kubernetes.io/name=license-manager --tail=50
}

# Main script logic
case "${1:-install}" in
    "install")
        install_license_manager
        ;;
    "uninstall")
        uninstall_license_manager
        ;;
    "status")
        show_status
        ;;
    "logs")
        show_logs
        ;;
    "help"|"-h"|"--help")
        show_usage
        ;;
    *)
        print_error "Unknown command: $1"
        echo ""
        show_usage
        exit 1
        ;;
esac