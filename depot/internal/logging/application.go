package logging

import (
	"context"
	"log/slog"
	"mall/depot/internal/application"
	"mall/depot/internal/domain"
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

func (a Application) CreateShoppingList(ctx context.Context, cmd application.CreateShoppingList) error {
	a.logger.Info("--> Depot.CreateShoppingList")
	defer func() {
		a.logger.Info("<-- Depot.CreateShoppingList")
	}()

	return a.App.CreateShoppingList(ctx, cmd)
}

func (a Application) CancelShoppingList(ctx context.Context, cmd application.CancelShoppingList) error {
	a.logger.Info("--> Depot.CancelShoppingList")
	defer func() {
		a.logger.Info("<-- Depot.CancelShoppingList")
	}()

	return a.App.CancelShoppingList(ctx, cmd)
}

func (a Application) AssignShoppingList(ctx context.Context, cmd application.AssignShoppingList) error {
	a.logger.Info("--> Depot.AssignShoppingList")
	defer func() {
		a.logger.Info("<-- Depot.AssignShoppingList")
	}()

	return a.App.AssignShoppingList(ctx, cmd)
}

func (a Application) CompleteShoppingList(ctx context.Context, cmd application.CompleteShoppingList) error {
	a.logger.Info("--> Depot.CompleteShoppingList")
	defer func() {
		a.logger.Info("<-- Depot.CompleteShoppingList")
	}()

	return a.App.CompleteShoppingList(ctx, cmd)
}

func (a Application) GetShoppingList(ctx context.Context, query application.GetShoppingList) (*domain.ShoppingList, error) {
	a.logger.Info("--> Depot.GetShoppingList")
	defer func() {
		a.logger.Info("<-- Depot.GetShoppingList")
	}()

	return a.App.GetShoppingList(ctx, query)
}
