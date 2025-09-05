# Offline Deployment Instructions

## Prerequisites
- Kubernetes cluster with Helm installed
- Docker daemon running
- No internet access required

## Deployment Steps

### 1. Load Docker Image
```bash
# Load the pre-built Docker image
docker load -i license-manager-fp216.tar

# Verify image is loaded
docker images | grep license-manager
```

### 2. Deploy with Helm
```bash
# Create namespace
kubectl create namespace license-manager

# Deploy the application
helm upgrade --install license-manager ./helm-charts/license-manager \
  --namespace license-manager \
  --create-namespace \
  --set image.repository=license-manager \
  --set image.tag=fp216 \
  --set image.pullPolicy=Never \
  --set namespace=license-manager

# Check deployment status
kubectl get pods -l app.kubernetes.io/name=license-manager -n license-manager
kubectl get services -n license-manager
kubectl get ingress -n license-manager
```

### 3. Access the Application
```bash
# Get the service URL
kubectl get ingress license-manager -n license-manager

# Or port-forward for testing
kubectl port-forward service/license-manager 8080:80 -n license-manager
# Access at: http://localhost:8080

# Ingress URL (if ingress controller is available)
# Access at: http://192.168.5.152/license-manager
```

### 4. Verify Deployment
```bash
# Check pod logs
kubectl logs -l app.kubernetes.io/name=license-manager -n license-manager

# Check pod status
kubectl describe pods -l app.kubernetes.io/name=license-manager -n license-manager
```

## Troubleshooting

### If pods are not starting:
```bash
# Check events
kubectl get events --sort-by=.metadata.creationTimestamp

# Check pod details
kubectl describe pod <pod-name>
```

### If image pull fails:
- Ensure Docker image is loaded: `docker images | grep license-manager`
- Check image pull policy is set to `Never` in Helm values

### If service is not accessible:
```bash
# Check service endpoints
kubectl get endpoints

# Check ingress controller
kubectl get pods -n ingress-nginx
```

## Cleanup
```bash
# Uninstall the application
helm uninstall license-manager -n license-manager

# Remove namespace (optional)
kubectl delete namespace license-manager

# Remove Docker image (optional)
docker rmi license-manager:fp216
```
