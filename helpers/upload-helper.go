package helpers

import (
	"context"
	"io"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadToCloudinary(ctx context.Context, file io.Reader) (*uploader.UploadResult, error) {
	// Initialize Cloudinary uploader
	cloudinaryClient, _ := cloudinary.NewFromParams(os.Getenv("Cloud"), os.Getenv("Key"), os.Getenv("Secret"))

	// Upload file to Cloudinary
	uploadResult, err := cloudinaryClient.Upload.Upload(ctx, file, uploader.UploadParams{})
	if err != nil {
		return nil, err
	}

	return uploadResult, nil
}
