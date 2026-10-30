package application

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/yohannes/sint-new-backend/internal/domain"
)

var (
	ErrMissingToken = errors.New("missing bearer token")
	ErrInvalidPhoto = errors.New("invalid photo")
)

type GetCurrentUser struct {
	Authenticator Authenticator
}

func (useCase GetCurrentUser) Execute(ctx context.Context, token string) (domain.User, error) {
	if token == "" {
		return domain.User{}, ErrMissingToken
	}
	return useCase.Authenticator.VerifyIDToken(ctx, token)
}

type UploadPhoto struct {
	Storage PhotoStorage
}

func (useCase UploadPhoto) Execute(
	ctx context.Context,
	content io.Reader,
	filename string,
	contentType string,
) (string, error) {
	if content == nil || filename == "" || contentType == "" {
		return "", ErrInvalidPhoto
	}
	extension := filepath.Ext(filename)
	if extension == "" {
		return "", fmt.Errorf("%w: filename must include an extension", ErrInvalidPhoto)
	}
	return useCase.Storage.Upload(ctx, content, filepath.Base(filename), contentType)
}
