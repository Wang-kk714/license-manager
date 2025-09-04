package fixtures

import "license-manager/internal/services"

// TestSSHConfigs provides test SSH configurations
var TestSSHConfigs = struct {
	Valid   *services.SSHConfig
	Invalid *services.SSHConfig
	Empty   *services.SSHConfig
}{
	Valid: &services.SSHConfig{
		Host:     "localhost",
		Port:     "22",
		Username: "testuser",
		Password: "testpass",
	},
	Invalid: &services.SSHConfig{
		Host:     "invalid-host",
		Port:     "9999",
		Username: "invaliduser",
		Password: "invalidpass",
	},
	Empty: &services.SSHConfig{
		Host:     "",
		Port:     "",
		Username: "",
		Password: "",
	},
}

// TestCommands provides test SSH commands
var TestCommands = struct {
	Valid   string
	Invalid string
}{
	Valid:   "ls -la",
	Invalid: "nonexistentcommand",
}

// TestFiles provides test file names
var TestFiles = struct {
	Valid   string
	Invalid string
}{
	Valid:   "test.txt",
	Invalid: "/nonexistent/path/file.txt",
}
