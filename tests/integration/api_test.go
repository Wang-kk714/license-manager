package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"license-manager/internal/handlers"
	"license-manager/tests/fixtures"

	"github.com/gin-gonic/gin"
)

func TestAPIEndpointsIntegration(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new router with all endpoints
	router := gin.New()
	router.LoadHTMLGlob("../../templates/*")
	
	// Register all routes
	router.GET("/", handlers.IndexHandler)
	router.POST("/api/check-license-cli", handlers.CheckLicenseCLIHandler)
	router.POST("/api/download-sysinfo", handlers.DownloadSysinfoHandler)
	router.POST("/api/upload-license", handlers.UploadLicenseHandler)

	t.Run("Index Page", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		body := w.Body.String()
		if !bytes.Contains([]byte(body), []byte("License Manager")) {
			t.Error("Expected response to contain 'License Manager'")
		}
	})

	t.Run("Check License CLI - Valid Request", func(t *testing.T) {
		config := handlers.ServerConfig{
			Host:     fixtures.TestSSHConfigs.Valid.Host,
			Port:     fixtures.TestSSHConfigs.Valid.Port,
			Username: fixtures.TestSSHConfigs.Valid.Username,
			Password: fixtures.TestSSHConfigs.Valid.Password,
		}

		jsonData, err := json.Marshal(config)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("POST", "/api/check-license-cli", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return an error since we can't actually connect to localhost:22
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
		}

		var response handlers.CheckLicenseCLIResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatal(err)
		}

		if response.Exists {
			t.Error("Expected Exists to be false for connection failure")
		}

		if response.Error == "" {
			t.Error("Expected Error to be set for connection failure")
		}
	})

	t.Run("Download Sysinfo - Valid Request", func(t *testing.T) {
		config := handlers.ServerConfig{
			Host:     fixtures.TestSSHConfigs.Valid.Host,
			Port:     fixtures.TestSSHConfigs.Valid.Port,
			Username: fixtures.TestSSHConfigs.Valid.Username,
			Password: fixtures.TestSSHConfigs.Valid.Password,
		}

		jsonData, err := json.Marshal(config)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("POST", "/api/download-sysinfo", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return an error since we can't actually connect
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
		}

		var response handlers.DownloadSysinfoResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatal(err)
		}

		if response.Success {
			t.Error("Expected Success to be false for connection failure")
		}

		if response.Error == "" {
			t.Error("Expected Error to be set for connection failure")
		}
	})

	t.Run("Upload License - No File", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/api/upload-license", nil)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}

		var response handlers.UploadLicenseResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatal(err)
		}

		if response.Success {
			t.Error("Expected Success to be false for no file request")
		}

		if response.Error == "" {
			t.Error("Expected Error to be set for no file request")
		}
	})
}
