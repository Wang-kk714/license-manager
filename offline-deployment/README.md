# License Manager - Offline Deployment

This directory contains the optimized offline deployment configuration for the License Manager application.

## Features

- **Offline Ready**: No internet connection required
- **Portable**: Works in any Kubernetes environment
- **Clean Configuration**: Optimized Helm charts and values
- **Easy Deployment**: Single script deployment
- **Proper Ingress**: Configured for path-based routing

## Quick Start

### Prerequisites

- Kubernetes cluster with Helm installed
- Docker daemon running
- kubectl configured to access your cluster

### Deploy

```bash
# Navigate to offline deployment directory
cd offline-deployment

# Run the deployment script
./deploy.sh
```

### Access

After deployment, access the application at:
```
http://<cluster-ip>/license-manager
```

## Configuration

### Ingress Configuration

The application is configured to be accessible at `/license-manager` path:

- **Main Application**: `http://<cluster-ip>/license-manager/`
- **API Endpoints**: `http://<cluster-ip>/license-manager/api/*`
- **Static Files**: `http://<cluster-ip>/license-manager/static/*`

### Resource Limits

Default resource configuration:
- **CPU**: 100m request, 200m limit
- **Memory**: 128Mi request, 256Mi limit

### Storage

- **Uploads**: 1Gi temporary storage
- **Downloads**: 1Gi temporary storage
- Uses `emptyDir` volumes (temporary)

## Management Commands

### Check Status
```bash
kubectl get pods -l app.kubernetes.io/name=license-manager
kubectl get services -l app.kubernetes.io/name=license-manager
kubectl get ingress -l app.kubernetes.io/name=license-manager
```

### View Logs
```bash
kubectl logs -l app.kubernetes.io/name=license-manager
```

### Port Forward (for testing)
```bash
kubectl port-forward service/license-manager-offline 8080:80
# Access at: http://localhost:8080
```

### Uninstall
```bash
helm delete license-manager-offline --purge
```

## Customization

### Modify Values

Edit `values.yaml` to customize:
- Resource limits
- Ingress configuration
- Service settings
- Storage configuration

### Custom Ingress

To change the ingress path, modify `values.yaml`:
```yaml
ingress:
  hosts:
    - host: ""
      paths:
        - path: /your-custom-path(/|$)(.*)
          pathType: Prefix
```

## Troubleshooting

### Common Issues

1. **Pod not starting**: Check resource limits and node capacity
2. **Ingress not working**: Verify ingress controller is running
3. **API calls failing**: Check ingress path configuration

### Debug Commands

```bash
# Check pod status
kubectl describe pod -l app.kubernetes.io/name=license-manager

# Check ingress configuration
kubectl describe ingress -l app.kubernetes.io/name=license-manager

# Check service endpoints
kubectl get endpoints -l app.kubernetes.io/name=license-manager
```

## Architecture

```
Internet → Ingress Controller → Service → Pod
    ↓
/license-manager/* → license-manager:80 → license-manager:8080
```

The application uses:
- **Nginx Ingress**: Path-based routing
- **ClusterIP Service**: Internal service discovery
- **Deployment**: Single replica with resource limits
- **ConfigMap**: Application configuration
- **ServiceAccount**: Pod security context