package unit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"license-manager/internal/handlers"

	"github.com/gin-gonic/gin"
)

func TestIndexHandler(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new router
	router := gin.New()
	router.LoadHTMLGlob("../../templates/*")
	router.GET("/", handlers.IndexHandler)

	// Create a test request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check the status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check if the response contains expected content
	body := w.Body.String()
	if !bytes.Contains([]byte(body), []byte("License Manager")) {
		t.Error("Expected response to contain 'License Manager'")
	}

	if !bytes.Contains([]byte(body), []byte("Server Connection")) {
		t.Error("Expected response to contain 'Server Connection'")
	}
}

func TestCheckLicenseCLIHandler_InvalidJSON(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new router
	router := gin.New()
	router.POST("/api/check-license-cli", handlers.CheckLicenseCLIHandler)

	// Create a test request with invalid JSON
	req, err := http.NewRequest("POST", "/api/check-license-cli", bytes.NewBufferString("invalid json"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check the status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Check the response body
	var response handlers.CheckLicenseCLIResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Exists {
		t.Error("Expected Exists to be false for invalid request")
	}

	if response.Error == "" {
		t.Error("Expected Error to be set for invalid request")
	}
}

func TestCheckLicenseCLIHandler_MissingFields(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new router
	router := gin.New()
	router.POST("/api/check-license-cli", handlers.CheckLicenseCLIHandler)

	// Create a test request with missing fields
	config := handlers.ServerConfig{
		Host:     "localhost",
		Port:     "", // Missing port
		Username: "testuser",
		Password: "testpass",
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

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check the status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Check the response body
	var response handlers.CheckLicenseCLIResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Exists {
		t.Error("Expected Exists to be false for missing fields")
	}

	if response.Error == "" {
		t.Error("Expected Error to be set for missing fields")
	}
}

func TestDownloadSysinfoHandler_InvalidJSON(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new router
	router := gin.New()
	router.POST("/api/download-sysinfo", handlers.DownloadSysinfoHandler)

	// Create a test request with invalid JSON
	req, err := http.NewRequest("POST", "/api/download-sysinfo", bytes.NewBufferString("invalid json"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check the status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Check the response body
	var response handlers.DownloadSysinfoResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Success {
		t.Error("Expected Success to be false for invalid request")
	}

	if response.Error == "" {
		t.Error("Expected Error to be set for invalid request")
	}
}

func TestUploadLicenseHandler_NoFile(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new router
	router := gin.New()
	router.POST("/api/upload-license", handlers.UploadLicenseHandler)

	// Create a test request without file
	req, err := http.NewRequest("POST", "/api/upload-license", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check the status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Check the response body
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
}