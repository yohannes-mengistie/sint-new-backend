package storage

import (
	"context"
	"io"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryStorage struct {
	client *cloudinary.Cloudinary
	folder string
}

func NewCloudinaryStorage(client *cloudinary.Cloudinary, folder string) *CloudinaryStorage {
	return &CloudinaryStorage{client: client, folder: folder}
}

func (storage *CloudinaryStorage) Upload(
	ctx context.Context,
	content io.Reader,
	filename string,
	contentType string,
) (string, error) {
	result, err := storage.client.Upload.Upload(ctx, content, uploader.UploadParams{
		Folder:       storage.folder,
		PublicID:     filename,
		ResourceType: "image",
	})
	if err != nil {
		return "", err
	}
	return result.SecureURL, nil
}
