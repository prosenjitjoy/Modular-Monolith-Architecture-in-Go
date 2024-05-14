package application

import (
	"context"
	"mall/notifications/internal/domain"
)

type OrderCreated struct {
	OrderID    string
	CustomerID string
}

type OrderCanceled struct {
	OrderID    string
	CustomerID string
}

type OrderReady struct {
	OrderID    string
	CustomerID string
}

type App interface {
	NotifyOrderCreated(ctx context.Context, notify OrderCreated) error
	NotifyOrderCanceled(ctx context.Context, notify OrderCanceled) error
	NotifyOrderReady(ctx context.Context, notify OrderReady) error
}

type Application struct {
	customers domain.CustomerRepository
}

var _ App = (*Application)(nil)

func New(customers domain.CustomerRepository) *Application {
	return &Application{
		customers: customers,
	}
}

func (a Application) NotifyOrderCreated(ctx context.Context, notify OrderCreated) error {
	// not implemented

	return nil
}

func (a Application) NotifyOrderCanceled(ctx context.Context, notify OrderCanceled) error {
	// not implemented

	return nil
}

func (a Application) NotifyOrderReady(ctx context.Context, notify OrderReady) error {
	// not implemented

	return nil
}
