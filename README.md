# License Manager

A web-based license management tool for remote license2_cli operations.

## Features

- **Server Connection**: Connect to remote servers via SSH
- **License CLI Check**: Verify if license2_cli exists on the target server
- **Sysinfo Download**: Execute `license2_cli getsysinfo -f 10` and download the generated file
- **License Upload**: Upload renewed license files and execute `license2_cli import -l {file}`

## Prerequisites

- Go 1.21 or later
- SSH access to target servers
- license2_cli installed on target servers

## Installation

### Minikube/Kubernetes (Recommended)

1. Start minikube:
```bash
make minikube-start
# or
minikube start
minikube addons enable ingress
```

2. Deploy to Kubernetes:
```bash
make dev-deploy
# or
make build
make helm-install
```

3. Get the application URL:
```bash
make get-url
# Access at: http://license-manager.local
```

**Note**: The application uses Ingress for external access. Make sure `license-manager.local` is added to your `/etc/hosts` file pointing to the Minikube IP.

### Local Development (requires Go)

1. Install Go 1.21 or later
2. Install dependencies:
```bash
go mod tidy
```
3. Run the application:
```bash
go run main.go
```

## Usage

1. **Connect to Server**: Enter the server IP, SSH port, username, and password
2. **Check License CLI**: Click "Check License2_CLI" to verify the tool exists on the server
3. **Download Sysinfo**: If license2_cli is found, you can download system information files
4. **Upload License**: Upload renewed license files to import them to the server

## API Endpoints

- `GET /` - Web interface
- `POST /api/check-license-cli` - Check if license2_cli exists on server
- `POST /api/download-sysinfo` - Download sysinfo file from server
- `POST /api/upload-license` - Upload and import license file

## Security Notes

- SSH connections use password authentication (consider using SSH keys for production)
- Host key verification is disabled for convenience (enable for production)
- Temporary files are cleaned up automatically
- All operations are logged

## Project Structure

```
license-manager/
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── internal/
│   ├── handlers/          # HTTP request handlers
│   ├── services/          # Business logic services
│   └── middleware/        # HTTP middleware
├── templates/             # HTML templates
├── static/               # Static assets (CSS, JS)
├── downloads/            # Downloaded sysinfo files
└── uploads/              # Temporary upload directory
```

## Development

### Docker Commands

```bash
# Build Docker image
make build

# Clean up Docker resources
make clean
```

### Kubernetes Commands

```bash
# Setup development environment
make dev-setup

# Deploy to minikube
make dev-deploy

# Port forward to local machine
make port-forward

# Get application URL
make get-url

# Cleanup
make dev-cleanup
```

### Local Development

```bash
# Build the application
go build -o license-manager main.go

# Run tests
go test ./...

# Run locally
go run main.go
```

### Available Make Commands

Run `make help` to see all available commands.

## License

This project is licensed under the MIT License.
