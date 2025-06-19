# S3 File Upload Service

A simple Go-based file upload service that allows users to upload and download files using AWS S3 storage.

## Project Structure

```
file-upload/
├── cmd/
│   ├── main.go          # Application entry point
│   └── routes.go        # Route definitions
├── internal/
│   ├── config/
│   │   └── config.go    # Configuration management
│   └── handlers/
│       └── file_handler.go  # File upload/download handlers
├── env.example          # Environment variables template
├── go.mod              # Go module dependencies
├── go.sum              # Go module checksums
└── README.md           # This file
```

## Features

- **File Upload**: Upload files to AWS S3 with size validation (max 100MB)
- **File Download**: Generate presigned URLs for secure file downloads

## Prerequisites

- Go 1.23 or higher
- AWS S3 bucket
- AWS credentials with S3 permissions

## Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/YasmeenAbdelghany19/Go-s3-file-service.git
   cd Go-s3-file-service
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment variables**
   ```bash
   cp env.example .env
   ```
   
   Edit `.env` file with your AWS credentials (you can use the example file as a template):
   ```env
   # AWS Configuration
   AWS_REGION=us-east-1
   AWS_S3_BUCKET=your-s3-bucket-name
   AWS_ACCESS_KEY_ID=your-access-key-id
   AWS_SECRET_ACCESS_KEY=your-secret-access-key
   
   # Server Configuration
   SERVER_PORT=8080
   ```

4. **Run the application**
   ```bash
   go run ./cmd
   ```

   The server will start on `http://localhost:8080`

## API Endpoints

### Health Check
- **GET** `/health` - Check if service is running

### Upload File
- **POST** `/upload` - Upload a file to S3
  - Send file with field name `file`
  - Max file size: 100MB

### Download File
- **GET** `/download/:filename` - Get download URL for a file
  - Returns a presigned S3 URL for secure download
