package main

import (
	"license-manager/internal/handlers"
	"license-manager/internal/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create Gin router
	r := gin.Default()

	// Add middleware
	r.Use(middleware.CORS())

	// Serve static files
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// Routes
	r.GET("/", handlers.IndexHandler)
	r.POST("/api/check-license-cli", handlers.CheckLicenseCLIHandler)
	r.POST("/api/download-sysinfo", handlers.DownloadSysinfoHandler)
	r.POST("/api/upload-license", handlers.UploadLicenseHandler)

	// Start server
	log.Println("License Manager starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
