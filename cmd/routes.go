package main

import (
	"github.com/YasmeenAbdelghany19/Go-s3-file-service/internal/handlers"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func setupRoutes(router *gin.Engine, s3Client *s3.Client, bucketName, s3UploadDir string) {
	// Initialize file handler
	fileHandler := handlers.NewFileHandler(s3Client, bucketName, s3UploadDir)

	// File routes
	router.GET("/download/:filename", fileHandler.DownloadHandler)
	router.POST("/upload", fileHandler.UploadHandler)

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "S3 File Service",
		})
	})
}
