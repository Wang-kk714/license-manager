package handlers

import (
	"license-manager/internal/services"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	Host     string `json:"host" binding:"required"`
	Port     string `json:"port" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CheckLicenseCLIResponse struct {
	Exists bool   `json:"exists"`
	Error  string `json:"error,omitempty"`
}

type DownloadSysinfoResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type UploadLicenseResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "License Manager",
	})
}

func CheckLicenseCLIHandler(c *gin.Context) {
	var config ServerConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, CheckLicenseCLIResponse{
			Exists: false,
			Error:  "Invalid request data: " + err.Error(),
		})
		return
	}

	sshConfig := &services.SSHConfig{
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
	}

	sshService := services.NewSSHService(sshConfig)
	defer sshService.Close()

	// Connect to server
	if err := sshService.Connect(); err != nil {
		c.JSON(http.StatusInternalServerError, CheckLicenseCLIResponse{
			Exists: false,
			Error:  "Failed to connect to server: " + err.Error(),
		})
		return
	}

	// Check if license2_cli exists
	exists, err := sshService.CheckLicenseCLI()
	if err != nil {
		c.JSON(http.StatusInternalServerError, CheckLicenseCLIResponse{
			Exists: false,
			Error:  "Failed to check license2_cli: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, CheckLicenseCLIResponse{
		Exists: exists,
	})
}

func DownloadSysinfoHandler(c *gin.Context) {
	var config ServerConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, DownloadSysinfoResponse{
			Success: false,
			Error:   "Invalid request data: " + err.Error(),
		})
		return
	}

	sshConfig := &services.SSHConfig{
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
	}

	sshService := services.NewSSHService(sshConfig)
	defer sshService.Close()

	// Connect to server
	if err := sshService.Connect(); err != nil {
		c.JSON(http.StatusInternalServerError, DownloadSysinfoResponse{
			Success: false,
			Error:   "Failed to connect to server: " + err.Error(),
		})
		return
	}

	// Check if license2_cli exists first
	exists, err := sshService.CheckLicenseCLI()
	if err != nil {
		c.JSON(http.StatusInternalServerError, DownloadSysinfoResponse{
			Success: false,
			Error:   "Failed to check license2_cli: " + err.Error(),
		})
		return
	}

	if !exists {
		c.JSON(http.StatusBadRequest, DownloadSysinfoResponse{
			Success: false,
			Error:   "license2_cli not found on server",
		})
		return
	}

	// Generate sysinfo file
	sysinfoFile, err := sshService.GenerateSysinfoFile()
	if err != nil {
		log.Printf("Error generating sysinfo file: %v", err)
		c.JSON(http.StatusInternalServerError, DownloadSysinfoResponse{
			Success: false,
			Error:   "Failed to generate sysinfo file: " + err.Error(),
		})
		return
	}

	// Extract IP address from host (remove port if present)
	hostIP := config.Host
	if strings.Contains(hostIP, ":") {
		hostIP = strings.Split(hostIP, ":")[0]
	}
	
	// Extract last octet of IP address (e.g., "192.168.5.152" -> "152")
	ipParts := strings.Split(hostIP, ".")
	lastOctet := "unknown"
	if len(ipParts) > 0 {
		lastOctet = ipParts[len(ipParts)-1]
	}

	// Create filename with last octet appended (preserve original extension)
	originalExt := filepath.Ext(sysinfoFile)
	baseName := strings.TrimSuffix(sysinfoFile, originalExt)
	downloadFilename := baseName + originalExt + "_" + lastOctet

	// Stream file directly to browser
	if err := sshService.StreamFileToResponse(c, sysinfoFile, downloadFilename); err != nil {
		log.Printf("Error streaming file: %v", err)
		c.JSON(http.StatusInternalServerError, DownloadSysinfoResponse{
			Success: false,
			Error:   "Failed to download sysinfo file: " + err.Error(),
		})
		return
	}
}

func UploadLicenseHandler(c *gin.Context) {
	// Get server config from form data
	config := ServerConfig{
		Host:     c.PostForm("host"),
		Port:     c.PostForm("port"),
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
	}

	if config.Host == "" || config.Port == "" || config.Username == "" || config.Password == "" {
		c.JSON(http.StatusBadRequest, UploadLicenseResponse{
			Success: false,
			Error:   "Missing server configuration",
		})
		return
	}

	// Get uploaded file
	file, err := c.FormFile("license_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, UploadLicenseResponse{
			Success: false,
			Error:   "No license file uploaded: " + err.Error(),
		})
		return
	}

	// Create uploads directory if it doesn't exist
	uploadsDir := "uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, UploadLicenseResponse{
			Success: false,
			Error:   "Failed to create uploads directory: " + err.Error(),
		})
		return
	}

	// Save uploaded file temporarily
	tempFile := filepath.Join(uploadsDir, file.Filename)
	if err := c.SaveUploadedFile(file, tempFile); err != nil {
		c.JSON(http.StatusInternalServerError, UploadLicenseResponse{
			Success: false,
			Error:   "Failed to save uploaded file: " + err.Error(),
		})
		return
	}
	defer os.Remove(tempFile) // Clean up temp file

	sshConfig := &services.SSHConfig{
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
	}

	sshService := services.NewSSHService(sshConfig)
	defer sshService.Close()

	// Connect to server
	if err := sshService.Connect(); err != nil {
		c.JSON(http.StatusInternalServerError, UploadLicenseResponse{
			Success: false,
			Error:   "Failed to connect to server: " + err.Error(),
		})
		return
	}

	// Check if license2_cli exists first
	exists, err := sshService.CheckLicenseCLI()
	if err != nil {
		c.JSON(http.StatusInternalServerError, UploadLicenseResponse{
			Success: false,
			Error:   "Failed to check license2_cli: " + err.Error(),
		})
		return
	}

	if !exists {
		c.JSON(http.StatusBadRequest, UploadLicenseResponse{
			Success: false,
			Error:   "license2_cli not found on server",
		})
		return
	}

	// Upload license file to server
	remoteFile := "/tmp/" + file.Filename
	if err := sshService.UploadFile(tempFile, remoteFile); err != nil {
		c.JSON(http.StatusInternalServerError, UploadLicenseResponse{
			Success: false,
			Error:   "Failed to upload license file: " + err.Error(),
		})
		return
	}

	// Execute license2_cli import command
	importCmd := "license2_cli import -l " + remoteFile
	output, err := sshService.ExecuteCommand(importCmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UploadLicenseResponse{
			Success: false,
			Error:   "Failed to import license: " + err.Error(),
		})
		return
	}

	// Clean up remote file
	sshService.ExecuteCommand("rm -f " + remoteFile)

	// Execute license2_cli check to verify the license
	checkCmd := "license2_cli check"
	checkOutput, err := sshService.ExecuteCommand(checkCmd)
	if err != nil {
		c.JSON(http.StatusOK, UploadLicenseResponse{
			Success: true,
			Message: "License imported successfully, but check command failed: " + err.Error() + "\n\nImport Output:\n```\n" + output + "\n```",
		})
		return
	}

	c.JSON(http.StatusOK, UploadLicenseResponse{
		Success: true,
		Message: "License imported successfully!\n\nImport Output:\n```\n" + output + "\n```\n\nLicense Check Output:\n```\n" + checkOutput + "\n```",
	})
}
