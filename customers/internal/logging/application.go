package logging

import (
	"context"
	"log/slog"

	"mall/customers/internal/application"
	"mall/customers/internal/domain"
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

func (a Application) RegisterCustomer(ctx context.Context, register application.RegisterCustomer) error {
	a.logger.Info("--> Customers.RegisterCustomer")
	defer func() {
		a.logger.Info("<-- Customers.RegisterCustomer")
	}()

	return a.App.RegisterCustomer(ctx, register)
}

func (a Application) AuthorizeCustomer(ctx context.Context, authorize application.AuthorizeCustomer) error {
	a.logger.Info("--> Customers.AuthorizeCustomer")
	defer func() {
		a.logger.Info("<-- Customers.AuthorizeCustomer")
	}()

	return a.App.AuthorizeCustomer(ctx, authorize)
}

func (a Application) GetCustomer(ctx context.Context, get application.GetCustomer) (*domain.Customer, error) {
	a.logger.Info("--> Customers.GetCustomer")
	defer func() {
		a.logger.Info("<-- Customers.GetCustomer")
	}()

	return a.App.GetCustomer(ctx, get)
}

func (a Application) EnableCustomer(ctx context.Context, enable application.EnableCustomer) error {
	a.logger.Info("--> Customers.EnableCustomer")
	defer func() {
		a.logger.Info("<-- Customers.EnableCustomer")
	}()

	return a.App.EnableCustomer(ctx, enable)
}

func (a Application) DisableCustomer(ctx context.Context, disable application.DisableCustomer) error {
	a.logger.Info("--> Customers.DisableCustomer")
	defer func() {
		a.logger.Info("<-- Customers.DisableCustomer")
	}()

	return a.App.DisableCustomer(ctx, disable)
}
