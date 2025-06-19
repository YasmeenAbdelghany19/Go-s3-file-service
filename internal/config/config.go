package config

import (
	"fmt"
	"os"
)

type Config struct {
	S3Region    string
	BucketName  string
	S3UploadDir string
	ServerPort  string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		S3Region:    os.Getenv("AWS_REGION"),
		BucketName:  os.Getenv("AWS_S3_BUCKET"),
		S3UploadDir: "uploads/",
		ServerPort:  os.Getenv("SERVER_PORT"),
	}

	// Validate required configuration
	if config.BucketName == "" {
		return nil, fmt.Errorf("AWS_S3_BUCKET environment variable is required")
	}
	if config.S3Region == "" {
		return nil, fmt.Errorf("AWS_REGION environment variable is required")
	}
	if config.ServerPort == "" {
		return nil, fmt.Errorf("SERVER_PORT environment variable is required")
	}

	return config, nil
}

func (c *Config) PrintConfig() {
	fmt.Printf("=== S3 File Service Configuration ===\n")
	fmt.Printf("S3 Region: %s\n", c.S3Region)
	fmt.Printf("S3 Bucket: %s\n", c.BucketName)
	fmt.Printf("Upload Directory: %s\n", c.S3UploadDir)
	fmt.Printf("Server Port: %s\n", c.ServerPort)
}
