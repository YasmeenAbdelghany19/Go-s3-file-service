package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/YasmeenAbdelghany19/Go-s3-file-service/internal/config"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	// Load AWS config
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(cfg.S3Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		)),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to load AWS config: %v", err))
	}

	s3Client := s3.NewFromConfig(awsCfg)
	fmt.Println("S3 client initialized")

	// Initialize Gin router
	router := gin.Default()

	// Setup routes
	setupRoutes(router, s3Client, cfg.BucketName, cfg.S3UploadDir)

	// Print configuration
	cfg.PrintConfig()
	fmt.Printf("Server started at :%s\n", cfg.ServerPort)

	// Start server
	router.Run(":" + cfg.ServerPort)
}
