package logging

import (
	"context"
	"log/slog"

	"mall/ordering/internal/application"
	"mall/ordering/internal/domain"
)

type Application struct {
	application.App
	logger *slog.Logger
}

var _ application.App = (*Application)(nil)

func LogApplicationAccess(application application.App, logger *slog.Logger) Application {
	return Application{
		App:    application,
		logger: logger,
	}
}

func (a Application) CreateOrder(ctx context.Context, cmd application.CreateOrder) error {
	a.logger.Info("--> Ordering.CreateOrder")
	defer func() {
		a.logger.Info("<-- Ordering.CreateOrder")
	}()

	return a.App.CreateOrder(ctx, cmd)
}

func (a Application) CancelOrder(ctx context.Context, cmd application.CancelOrder) error {
	a.logger.Info("--> Ordering.CancelOrder")
	defer func() {
		a.logger.Info("<-- Ordering.CancelOrder")
	}()

	return a.App.CancelOrder(ctx, cmd)
}

func (a Application) ReadyOrder(ctx context.Context, cmd application.ReadyOrder) error {
	a.logger.Info("--> Ordering.ReadyOrder")
	defer func() {
		a.logger.Info("<-- Ordering.ReadyOrder")
	}()

	return a.App.ReadyOrder(ctx, cmd)
}

func (a Application) CompleteOrder(ctx context.Context, cmd application.CompleteOrder) error {
	a.logger.Info("--> Ordering.CompleteOrder")
	defer func() {
		a.logger.Info("<-- Ordering.CompleteOrder")
	}()

	return a.App.CompleteOrder(ctx, cmd)
}

func (a Application) GetOrder(ctx context.Context, query application.GetOrder) (*domain.Order, error) {
	a.logger.Info("--> Ordering.GetOrder")
	defer func() {
		a.logger.Info("<-- Ordering.GetOrder")
	}()

	return a.App.GetOrder(ctx, query)
}
