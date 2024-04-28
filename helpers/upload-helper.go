package helpers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

func UploadToS3(ctx context.Context, file io.Reader, eventID string) (string, error) {
	// Load default configuration
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		return "", fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create an S3 client
	svc := s3.NewFromConfig(cfg)

	// Read the entire file into a buffer
	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, file); err != nil {
		return "", fmt.Errorf("failed to read file contents: %w", err)
	}

	// Use the first 512 bytes of the file to detect the content type
	contentType := http.DetectContentType(buffer.Bytes()[:512])
	fileExtension, err := getFileExtension(contentType)
	if err != nil {
		return "", fmt.Errorf("failed to get file extension: %w", err)
	}

	// Generate a unique key for the S3 object
	objectKey := fmt.Sprintf("%s/%s%s", eventID, generateUniqueKey(), fileExtension)

	// Upload the entire file to S3
	_, err = svc.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String("thaparevents"), // Change to your S3 bucket name
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(buffer.Bytes()),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Return the URL of the uploaded file
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", "thaparevents", objectKey), nil
}

func generateUniqueKey() string {
	// Generate a new UUID
	uuid := uuid.New()

	// Convert UUID to string
	return uuid.String()
}

func getFileExtension(mimeType string) (string, error) {
	// Custom mapping of MIME types to file extensions
	mimeExtensions := map[string]string{
		"image/jpeg":         ".jpg",
		"image/png":          ".png",
		"image/gif":          ".gif",
		"image/bmp":          ".bmp",
		"image/svg+xml":      ".svg",
		"image/tiff":         ".tiff",
		"image/webp":         ".webp",
		"image/x-icon":       ".ico",
		"application/pdf":    ".pdf",
		"application/msword": ".doc",
		"video/mp4":          ".mp4",
		"video/mpeg":         ".mpeg",
	}

	// Attempt to get the file extension from the MIME type
	ext, ok := mimeExtensions[mimeType]
	if !ok {
		// If the MIME type is not in the map, attempt to get a default extension
		exts, err := mime.ExtensionsByType(mimeType)
		if err != nil || len(exts) == 0 {
			return "", fmt.Errorf("no file extension found for MIME type: %s", mimeType)
		}
		ext = exts[0] // Use the first extension from the list
	}
	return ext, nil
}
