package application

import (
	"context"
	"errors"

	"github.com/yohannes/sint-new-backend/internal/domain"
)

var ErrProductNotFound = errors.New("product not found")

type ListProducts struct {
	Products ProductRepository
}

func (useCase ListProducts) Execute(ctx context.Context) ([]domain.Product, error) {
	return useCase.Products.List(ctx)
}

type GetProductDetail struct {
	Products ProductRepository
}

func (useCase GetProductDetail) Execute(ctx context.Context, id int64) (domain.ProductDetail, error) {
	detail, err := useCase.Products.GetDetail(ctx, id)
	if errors.Is(err, ErrProductNotFound) {
		return domain.ProductDetail{}, ErrProductNotFound
	}

	type ListStores struct {
		Stores StoreRepository
	}

	func (useCase ListStores) Execute(ctx context.Context) ([]domain.Store, error) {
		return useCase.Stores.List(ctx)
	}
	return detail, err
}
