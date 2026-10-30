package application

import (
	"context"
	"io"
	"strings"
	"testing"
)

type testPhotoStorage struct{}

func (testPhotoStorage) Upload(
	context.Context,
	io.Reader,
	string,
	string,
) (string, error) {
	return "https://example.com/photo.jpg", nil
}

func TestUploadPhotoRejectsMissingExtension(t *testing.T) {
	_, err := (UploadPhoto{Storage: testPhotoStorage{}}).Execute(
		context.Background(),
		strings.NewReader("photo"),
		"photo",
		"image/jpeg",
	)
	if err == nil {
		t.Fatal("expected invalid photo error")
	}
}
