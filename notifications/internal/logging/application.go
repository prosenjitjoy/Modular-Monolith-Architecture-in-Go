package logging

import (
	"context"
	"log/slog"

	"mall/notifications/internal/application"
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

func (a Application) NotifyOrderCreated(ctx context.Context, notify application.OrderCreated) error {
	a.logger.Info("--> Notifications.NotifyOrderCreated")
	defer func() {
		a.logger.Info("<-- Notifications.NotifyOrderCreated")
	}()

	return a.App.NotifyOrderCreated(ctx, notify)
}

func (a Application) NotifyOrderCanceled(ctx context.Context, notify application.OrderCanceled) error {
	a.logger.Info("--> Notifications.NotifyOrderCanceled")
	defer func() {
		a.logger.Info("<-- Notifications.NotifyOrderCanceled")
	}()

	return a.App.NotifyOrderCanceled(ctx, notify)
}

func (a Application) NotifyOrderReady(ctx context.Context, notify application.OrderReady) error {
	a.logger.Info("--> Notifications.NotifyOrderReady")
	defer func() {
		a.logger.Info("<-- Notifications.NotifyOrderReady")
	}()

	return a.App.NotifyOrderReady(ctx, notify)
}
