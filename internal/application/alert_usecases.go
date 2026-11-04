package application

import (
	"context"

	"github.com/yohannes/sint-new-backend/internal/domain"
)

type ListAlerts struct {
	Alerts AlertRepository
}

func (useCase ListAlerts) Execute(
	ctx context.Context,
	user domain.User,
) ([]domain.Alert, error) {
	return useCase.Alerts.ListAlerts(ctx, user)
}

type DeleteAlert struct {
	Alerts AlertRepository
}

func (useCase DeleteAlert) Execute(
	ctx context.Context,
	user domain.User,
	alertID int64,
) error {
	return useCase.Alerts.DeleteAlert(ctx, user, alertID)
}
