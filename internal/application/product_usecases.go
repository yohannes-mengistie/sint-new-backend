package application

import (
	"context"

	"github.com/yohannes/sint-new-backend/internal/domain"
)

type ListProducts struct {
	Products ProductRepository
}

func (useCase ListProducts) Execute(ctx context.Context) ([]domain.Product, error) {
	return useCase.Products.List(ctx)
}
