package services

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
)

type SSHConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

type SSHService struct {
	config *SSHConfig
	client *ssh.Client
}

// Config returns the SSH configuration (for testing)
func (s *SSHService) Config() *SSHConfig {
	return s.config
}

func NewSSHService(config *SSHConfig) *SSHService {
	return &SSHService{
		config: config,
	}
}

func (s *SSHService) Connect() error {
	config := &ssh.ClientConfig{
		User: s.config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.config.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	addr := net.JoinHostPort(s.config.Host, s.config.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %v", addr, err)
	}

	s.client = client
	return nil
}

func (s *SSHService) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}

func (s *SSHService) CheckLicenseCLI() (bool, error) {
	if s.client == nil {
		return false, fmt.Errorf("not connected to server")
	}

	session, err := s.client.NewSession()
	if err != nil {
		return false, fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	// Check if license2_cli exists
	output, err := session.CombinedOutput("which license2_cli")
	if err != nil {
		return false, nil // Command failed, likely means license2_cli doesn't exist
	}

	return strings.TrimSpace(string(output)) != "", nil
}

func (s *SSHService) ExecuteCommand(command string) (string, error) {
	if s.client == nil {
		return "", fmt.Errorf("not connected to server")
	}

	session, err := s.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", fmt.Errorf("command failed: %v", err)
	}

	return string(output), nil
}

func (s *SSHService) DownloadFile(remotePath, localPath string) error {
	if s.client == nil {
		return fmt.Errorf("not connected to server")
	}

	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	// Create local file
	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %v", err)
	}
	defer localFile.Close()

	// Create remote file reader
	remoteFile, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	// Start the command
	if err := session.Start(fmt.Sprintf("cat %s", remotePath)); err != nil {
		return fmt.Errorf("failed to start command: %v", err)
	}

	// Copy data
	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		return fmt.Errorf("failed to copy data: %v", err)
	}

	// Wait for command to complete
	return session.Wait()
}

func (s *SSHService) UploadFile(localPath, remotePath string) error {
	if s.client == nil {
		return fmt.Errorf("not connected to server")
	}

	// Read local file
	localFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %v", err)
	}
	defer localFile.Close()

	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	// Create remote file writer
	remoteFile, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %v", err)
	}

	// Start the command
	if err := session.Start(fmt.Sprintf("cat > %s", remotePath)); err != nil {
		return fmt.Errorf("failed to start command: %v", err)
	}

	// Copy data
	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		return fmt.Errorf("failed to copy data: %v", err)
	}

	// Close the pipe and wait for command to complete
	remoteFile.Close()
	return session.Wait()
}

func (s *SSHService) GenerateSysinfoFile() (string, error) {
	if s.client == nil {
		return "", fmt.Errorf("not connected to server")
	}

	// First, check if license2_cli exists and get its version
	_, err := s.ExecuteCommand("license2_cli --version 2>/dev/null || license2_cli -v 2>/dev/null || echo 'version command failed'")
	if err != nil {
		return "", fmt.Errorf("failed to check license2_cli version: %v", err)
	}

	// Execute license2_cli getsysinfo -f 10
	output, err := s.ExecuteCommand("license2_cli getsysinfo -f 10")
	if err != nil {
		return "", fmt.Errorf("failed to generate sysinfo: %v, output: %s", err, output)
	}

	// The file should be generated in current directory
	// We need to find the generated file
	session, err := s.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	// List all files in current directory to see what was created
	allFiles, err := session.CombinedOutput("ls -la")
	if err != nil {
		return "", fmt.Errorf("failed to list all files: %v", err)
	}

	// Close the current session and create a new one for the next command
	session.Close()
	session, err = s.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create new session: %v", err)
	}
	defer session.Close()

	// List files to find the generated sysinfo file (license2_cli generates sys_info.bin)
	fileList, err := session.CombinedOutput("ls -la sys_info.bin 2>/dev/null || ls -la *.sysinfo 2>/dev/null || ls -la sysinfo* 2>/dev/null || echo 'No sysinfo files found'")
	if err != nil {
		return "", fmt.Errorf("failed to list sysinfo files: %v", err)
	}

	output = strings.TrimSpace(string(fileList))
	if output == "" || strings.Contains(output, "No sysinfo files found") {
		// Close current session and create new one for alternative search
		session.Close()
		session, err = s.client.NewSession()
		if err != nil {
			return "", fmt.Errorf("failed to create session for alternative search: %v", err)
		}
		defer session.Close()
		
		// Try alternative file patterns
		altFiles, err := session.CombinedOutput("find . -name 'sys_info*' -o -name '*sysinfo*' -o -name '*.info' -o -name 'system*' 2>/dev/null || echo 'No alternative files found'")
		if err == nil {
			altOutput := strings.TrimSpace(string(altFiles))
			if altOutput != "" && !strings.Contains(altOutput, "No alternative files found") {
				return "", fmt.Errorf("no sysinfo file generated, but found alternative files: %s. All files: %s", altOutput, string(allFiles))
			}
		}
		return "", fmt.Errorf("no sysinfo file generated. Command output: %s. All files: %s", output, string(allFiles))
	}

	lines := strings.Split(output, "\n")
	latestFile := ""
	
	// Get the most recent sysinfo file
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Check if line contains sysinfo file (sys_info.bin or .sysinfo files)
		if strings.Contains(line, "sys_info.bin") || strings.Contains(line, ".sysinfo") || strings.Contains(line, "sysinfo") {
			parts := strings.Fields(line)
			if len(parts) >= 9 {
				// The filename is the last part of the ls -la output
				latestFile = parts[len(parts)-1]
				break
			}
		}
	}

	if latestFile == "" {
		return "", fmt.Errorf("could not determine sysinfo filename from output: %s. All files: %s", output, string(allFiles))
	}

	return latestFile, nil
}

// StreamFileToResponse streams a remote file directly to the HTTP response
func (s *SSHService) StreamFileToResponse(c *gin.Context, remotePath, downloadFilename string) error {
	if s.client == nil {
		return fmt.Errorf("not connected to server")
	}

	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	// Create remote file reader
	remoteFile, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	// Start the command
	if err := session.Start(fmt.Sprintf("cat %s", remotePath)); err != nil {
		return fmt.Errorf("failed to start command: %v", err)
	}

	// Set response headers for file download
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", downloadFilename))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")

	// Stream the file content directly to the response
	_, err = io.Copy(c.Writer, remoteFile)
	if err != nil {
		return fmt.Errorf("failed to copy data to response: %v", err)
	}

	// Wait for command to complete
	return session.Wait()
}
