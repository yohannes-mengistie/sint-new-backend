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
	GetDetail(ctx context.Context, id int64) (domain.ProductDetail, error)
}

type StoreRepository interface {
	ListStores(ctx context.Context) ([]domain.Store, error)
}

type FavoriteRepository interface {
	ListProductIDs(ctx context.Context, user domain.User) ([]int64, error)
	Add(ctx context.Context, user domain.User, productID int64) error
	Remove(ctx context.Context, user domain.User, productID int64) error
}

type PriceUpdateRepository interface {
	CreatePriceUpdate(ctx context.Context, user domain.User, update domain.PriceUpdate) error
	ListPriceSubmissions(ctx context.Context, user domain.User) ([]domain.PriceSubmission, error)
}

type AlertRepository interface {
	ListAlerts(ctx context.Context, user domain.User) ([]domain.Alert, error)
	DeleteAlert(ctx context.Context, user domain.User, alertID int64) error
}
