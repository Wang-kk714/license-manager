# License Manager

A modern web-based license management tool for batch operations on remote servers using `license2_cli`.

## Features

- **Multi-Server Management**: Connect to multiple servers simultaneously with server cards
- **Batch Operations**: Check, download, and upload operations across multiple servers
- **Multi-File Upload**: Assign different license files to specific servers
- **SSH Connectivity**: Secure remote server operations
- **Modern Web UI**: Responsive interface with real-time status updates
- **Kubernetes Ready**: Complete Helm charts for production deployment
- **Docker Containerized**: Easy deployment and scaling

## Quick Start

### Prerequisites
- Docker and Minikube (recommended)
- SSH access to target servers
- `license2_cli` installed on target servers

### Deploy with Minikube

```bash
# Start Minikube and deploy
make dev-deploy

# Access the application
make get-url
# Open: http://license-manager.local
```

### Manual Setup

```bash
# Start Minikube
minikube start
minikube addons enable ingress

# Build and deploy
make build
make helm-upgrade

# Add to hosts file
echo "$(minikube ip) license-manager.local" | sudo tee -a /etc/hosts
```

## Usage

### 1. Server Connection
- Enter server details (IP:Port, Username, Password)
- Click "Connect" to add servers to your session
- Manage connected servers with compact server cards

### 2. Batch Operations
- **Check**: Verify `license2_cli` exists on all connected servers
- **Download**: Download system info files from all servers
- **Upload**: Assign license files to specific servers

### 3. Multi-File Upload
- Click "+ Add File" to create upload assignments
- Select target server from dropdown
- Choose license file for each server
- Upload all files simultaneously

## API Endpoints

- `GET /` - Web interface
- `POST /api/check-license-cli` - Check license2_cli availability
- `POST /api/download-sysinfo` - Download system info files
- `POST /api/upload-license` - Upload and import license files

## Development

### Local Development (requires Go 1.21+)

```bash
go mod tidy
go run main.go
```

### Docker Development

```bash
# Build image
make build

# Run tests
make test

# Deploy to Minikube
make dev-deploy
```

### Available Commands

```bash
make help  # Show all available commands
```

## Project Structure

```
license-manager/
├── main.go                    # Application entry point
├── internal/
│   ├── handlers/             # HTTP request handlers
│   ├── services/             # SSH and business logic
│   └── middleware/           # CORS middleware
├── templates/                # HTML templates
├── static/                   # JavaScript and CSS
├── helm-charts/              # Kubernetes deployment
└── tests/                    # Unit and integration tests
```

## Security Notes

- SSH password authentication (consider SSH keys for production)
- Temporary files are automatically cleaned up
- All operations are logged
- Uses `emptyDir` volumes for temporary storage

## License

MIT License