package application

import (
	"context"

	"github.com/yohannes/sint-new-backend/internal/domain"
)

type UserRepository interface {
	Upsert(ctx context.Context, user domain.User) (domain.User, error)
}

type ProductRepository interface {
	List(ctx context.Context) ([]domain.Product, error)
}

type PriceUpdateRepository interface {
	CreatePriceUpdate(ctx context.Context, user domain.User, update domain.PriceUpdate) error
}
