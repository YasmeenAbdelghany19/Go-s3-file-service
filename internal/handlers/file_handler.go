package handlers

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	s3Client    *s3.Client
	bucketName  string
	s3UploadDir string
}

func NewFileHandler(s3Client *s3.Client, bucketName, s3UploadDir string) *FileHandler {
	return &FileHandler{
		s3Client:    s3Client,
		bucketName:  bucketName,
		s3UploadDir: s3UploadDir,
	}
}

func (h *FileHandler) DownloadHandler(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	fmt.Printf("Download request received for file: %s\n", filename)

	// Construct the S3 key
	s3Key := filepath.Join(h.s3UploadDir, filename)
	fmt.Printf("S3 key constructed: %s\n", s3Key)

	// Get file info from S3 first
	headObjectOutput, err := h.s3Client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(h.bucketName),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		fmt.Printf("Error getting file info from S3: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Extract file information
	fileSize := *headObjectOutput.ContentLength
	contentType := ""
	if headObjectOutput.ContentType != nil {
		contentType = *headObjectOutput.ContentType
	}
	fileSizeMB := float64(fileSize) / 1024 / 1024
	fmt.Printf("File info - Size: %.2f MB, Type: %s\n", fileSizeMB, contentType)

	// Create presigned URL
	presignClient := s3.NewPresignClient(h.s3Client)
	presignedURL, err := presignClient.PresignGetObject(context.TODO(),
		&s3.GetObjectInput{
			Bucket: aws.String(h.bucketName),
			Key:    aws.String(s3Key),
		},
	)
	if err != nil {
		fmt.Printf("Error generating presigned URL: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate download URL"})
		return
	}

	fmt.Printf("Presigned URL generated successfully for %s\n", filename)

	// Return detailed response
	c.JSON(http.StatusOK, gin.H{
		"filename":    filename,
		"url":         presignedURL.URL,
		"size":        fileSizeMB,
		"contentType": contentType,
		"message":     "Download URL generated successfully",
	})
}

func (h *FileHandler) UploadHandler(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Printf("Error reading file from request: %v\n", err)
		fmt.Println("Body request should be {file: <file>}")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read file"})
		return
	}
	defer file.Close()

	// Check file size (convert to MB)
	fileSizeMB := float64(header.Size) / 1024 / 1024
	if fileSizeMB > 100 { // 100MB limit
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("File too large. Maximum size: 100MB, Received: %.1fMB", fileSizeMB),
		})
		return
	}
	fmt.Printf("Upload request received - File: %s, Size: %.2f MB, Type: %s\n",
		header.Filename, fileSizeMB, header.Header.Get("Content-Type"))

	// Construct file path in S3
	s3Key := filepath.Join(h.s3UploadDir, header.Filename)
	fmt.Printf("S3 key constructed: %s\n", s3Key)

	// Upload to S3
	_, err = h.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(h.bucketName),
		Key:           aws.String(s3Key),
		Body:          file,
		ContentType:   aws.String(header.Header.Get("Content-Type")),
		ContentLength: &header.Size,
	})
	if err != nil {
		fmt.Printf("Error uploading file to S3: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to S3"})
		return
	}

	fmt.Printf("File uploaded successfully to S3: %s\n", header.Filename)
	c.JSON(http.StatusOK, gin.H{
		"message":     "File uploaded successfully",
		"filename":    header.Filename,
		"size":        fileSizeMB,
		"contentType": header.Header.Get("Content-Type"),
		"s3Key":       s3Key,
	})
}
