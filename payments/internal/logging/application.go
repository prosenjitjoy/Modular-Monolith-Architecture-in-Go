package logging

import (
	"context"
	"log/slog"

	"mall/payments/internal/application"
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

func (a Application) AuthorizePayment(ctx context.Context, authorize application.AuthorizePayment) error {
	a.logger.Info("--> Payments.AuthorizePayment")
	defer func() {
		a.logger.Info("<-- Payments.AuthorizePayment")
	}()

	return a.App.AuthorizePayment(ctx, authorize)
}

func (a Application) ConfirmPayment(ctx context.Context, confirm application.ConfirmPayment) error {
	a.logger.Info("--> Payments.ConfirmPayment")
	defer func() {
		a.logger.Info("<-- Payments.ConfirmPayment")
	}()

	return a.App.ConfirmPayment(ctx, confirm)
}

func (a Application) CreateInvoice(ctx context.Context, create application.CreateInvoice) error {
	a.logger.Info("--> Payments.CreateInvoice")
	defer func() {
		a.logger.Info("<-- Payments.CreateInvoice")
	}()

	return a.App.CreateInvoice(ctx, create)
}

func (a Application) AdjustInvoice(ctx context.Context, adjust application.AdjustInvoice) error {
	a.logger.Info("--> Payments.AdjustInvoice")
	defer func() {
		a.logger.Info("<-- Payments.AdjustInvoice")
	}()

	return a.App.AdjustInvoice(ctx, adjust)
}

func (a Application) PayInvoice(ctx context.Context, pay application.PayInvoice) error {
	a.logger.Info("--> Payments.PayInvoice")
	defer func() {
		a.logger.Info("<-- Payments.PayInvoice")
	}()

	return a.App.PayInvoice(ctx, pay)
}

func (a Application) CancelInvoice(ctx context.Context, cancel application.CancelInvoice) error {
	a.logger.Info("--> Payments.CancelInvoice")
	defer func() {
		a.logger.Info("<-- Payments.CancelInvoice")
	}()

	return a.App.CancelInvoice(ctx, cancel)
}
