package application

import (
	"context"
	"errors"
	"strings"

	"github.com/yohannes/sint-new-backend/internal/domain"
)

var ErrInvalidPriceUpdate = errors.New("product, store, and a positive price are required")

type SubmitPrice struct {
	Updates PriceUpdateRepository
}

type ListPriceSubmissions struct {
	Updates PriceUpdateRepository
}

func (useCase ListPriceSubmissions) Execute(
	ctx context.Context,
	user domain.User,
) ([]domain.PriceSubmission, error) {
	return useCase.Updates.ListPriceSubmissions(ctx, user)
}

func (useCase SubmitPrice) Execute(
	ctx context.Context,
	user domain.User,
	update domain.PriceUpdate,
) error {
	update.ProductName = strings.TrimSpace(update.ProductName)
	update.StoreName = strings.TrimSpace(update.StoreName)
	if update.ProductName == "" || update.StoreName == "" || update.Price <= 0 {
		return ErrInvalidPriceUpdate
	}
	return useCase.Updates.CreatePriceUpdate(ctx, user, update)
}
