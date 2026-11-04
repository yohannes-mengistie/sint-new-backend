package application

import (
	"context"

	"github.com/yohannes/sint-new-backend/internal/domain"
)

type ListFavorites struct {
	Favorites FavoriteRepository
}

func (useCase ListFavorites) Execute(
	ctx context.Context,
	user domain.User,
) ([]int64, error) {
	return useCase.Favorites.ListProductIDs(ctx, user)
}

type AddFavorite struct {
	Favorites FavoriteRepository
}

func (useCase AddFavorite) Execute(
	ctx context.Context,
	user domain.User,
	productID int64,
) error {
	return useCase.Favorites.Add(ctx, user, productID)
}

type RemoveFavorite struct {
	Favorites FavoriteRepository
}

func (useCase RemoveFavorite) Execute(
	ctx context.Context,
	user domain.User,
	productID int64,
) error {
	return useCase.Favorites.Remove(ctx, user, productID)
}
