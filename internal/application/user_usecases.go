package application

import (
	"context"

	"github.com/yohannes/sint-new-backend/internal/domain"
)

type SyncUser struct {
	Users UserRepository
}

func (useCase SyncUser) Execute(ctx context.Context, user domain.User) (domain.User, error) {
	return useCase.Users.Upsert(ctx, user)
}
