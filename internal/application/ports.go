package application

import (
	"context"
	"io"

	"github.com/yohannes/sint-new-backend/internal/domain"
)

type Authenticator interface {
	VerifyIDToken(ctx context.Context, idToken string) (domain.User, error)
}

type PhotoStorage interface {
	Upload(ctx context.Context, content io.Reader, filename string, contentType string) (string, error)
}

type DatabaseHealth interface {
	Ping(ctx context.Context) error
}
