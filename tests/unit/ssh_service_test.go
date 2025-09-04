package unit

import (
	"testing"

	"license-manager/internal/services"
)

func TestSSHConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   *services.SSHConfig
		expected services.SSHConfig
	}{
		{
			name: "valid config",
			config: &services.SSHConfig{
				Host:     "localhost",
				Port:     "22",
				Username: "testuser",
				Password: "testpass",
			},
			expected: services.SSHConfig{
				Host:     "localhost",
				Port:     "22",
				Username: "testuser",
				Password: "testpass",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.config.Host != tt.expected.Host {
				t.Errorf("Expected Host to be '%s', got '%s'", tt.expected.Host, tt.config.Host)
			}
			if tt.config.Port != tt.expected.Port {
				t.Errorf("Expected Port to be '%s', got '%s'", tt.expected.Port, tt.config.Port)
			}
			if tt.config.Username != tt.expected.Username {
				t.Errorf("Expected Username to be '%s', got '%s'", tt.expected.Username, tt.config.Username)
			}
			if tt.config.Password != tt.expected.Password {
				t.Errorf("Expected Password to be '%s', got '%s'", tt.expected.Password, tt.config.Password)
			}
		})
	}
}

func TestNewSSHService(t *testing.T) {
	config := &services.SSHConfig{
		Host:     "localhost",
		Port:     "22",
		Username: "testuser",
		Password: "testpass",
	}

	service := services.NewSSHService(config)

	if service == nil {
		t.Error("Expected service to be created, got nil")
	}

	if service.Config() != config {
		t.Error("Expected service config to match input config")
	}
}

func TestSSHService_CheckLicenseCLI_NotConnected(t *testing.T) {
	config := &services.SSHConfig{
		Host:     "localhost",
		Port:     "22",
		Username: "testuser",
		Password: "testpass",
	}

	service := services.NewSSHService(config)

	exists, err := service.CheckLicenseCLI()

	if err == nil {
		t.Error("Expected error when not connected, got nil")
	}

	if exists {
		t.Error("Expected exists to be false when not connected")
	}

	expectedError := "not connected to server"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestSSHService_ExecuteCommand_NotConnected(t *testing.T) {
	config := &services.SSHConfig{
		Host:     "localhost",
		Port:     "22",
		Username: "testuser",
		Password: "testpass",
	}

	service := services.NewSSHService(config)

	output, err := service.ExecuteCommand("ls")

	if err == nil {
		t.Error("Expected error when not connected, got nil")
	}

	if output != "" {
		t.Error("Expected empty output when not connected")
	}

	expectedError := "not connected to server"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestSSHService_GenerateSysinfoFile_NotConnected(t *testing.T) {
	config := &services.SSHConfig{
		Host:     "localhost",
		Port:     "22",
		Username: "testuser",
		Password: "testpass",
	}

	service := services.NewSSHService(config)

	filename, err := service.GenerateSysinfoFile()

	if err == nil {
		t.Error("Expected error when not connected, got nil")
	}

	if filename != "" {
		t.Error("Expected empty filename when not connected")
	}

	expectedError := "not connected to server"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}
