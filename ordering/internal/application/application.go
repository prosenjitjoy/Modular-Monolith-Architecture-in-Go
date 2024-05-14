package application

import (
	"context"
	"fmt"

	"mall/ordering/internal/domain"
)

type CreateOrder struct {
	ID         string
	CustomerID string
	PaymentID  string
	Items      []*domain.Item
}

type CancelOrder struct {
	ID string
}

type ReadyOrder struct {
	ID string
}

type CompleteOrder struct {
	ID        string
	InvoiceID string
}

type GetOrder struct {
	ID string
}

type App interface {
	CreateOrder(ctx context.Context, cmd CreateOrder) error
	CancelOrder(ctx context.Context, cmd CancelOrder) error
	ReadyOrder(ctx context.Context, cmd ReadyOrder) error
	CompleteOrder(ctx context.Context, cmd CompleteOrder) error
	GetOrder(ctx context.Context, query GetOrder) (*domain.Order, error)
}

type Application struct {
	orders        domain.OrderRepository
	customers     domain.CustomerRepository
	payments      domain.PaymentRepository
	invoices      domain.InvoiceRepository
	shopping      domain.ShoppingRepository
	notifications domain.NotificationRepository
}

var _ App = (*Application)(nil)

func New(orders domain.OrderRepository, customers domain.CustomerRepository, payments domain.PaymentRepository, invoices domain.InvoiceRepository, shopping domain.ShoppingRepository, notifications domain.NotificationRepository) *Application {
	return &Application{
		orders:        orders,
		customers:     customers,
		payments:      payments,
		invoices:      invoices,
		shopping:      shopping,
		notifications: notifications,
	}
}

func (a Application) CreateOrder(ctx context.Context, cmd CreateOrder) error {
	order, err := domain.CreateOrder(cmd.ID, cmd.CustomerID, cmd.PaymentID, cmd.Items)
	if err != nil {
		return fmt.Errorf("create order command: %w", err)
	}

	// authorizeCustomer
	if err = a.customers.Authorize(ctx, order.CustomerID); err != nil {
		return fmt.Errorf("order customer authorization: %w", err)
	}

	// validatePayment
	if err = a.payments.Confirm(ctx, order.PaymentID); err != nil {
		return fmt.Errorf("order payment confirmation: %w", err)
	}

	// scheduleShopping
	if order.ShoppingID, err = a.shopping.Create(ctx, order); err != nil {
		return fmt.Errorf("order shopping scheduling: %w", err)
	}

	// notifyOrderCreated
	if err = a.notifications.NotifyOrderCreated(ctx, order.ID, order.CustomerID); err != nil {
		return err
	}

	err = a.orders.Save(ctx, order)
	if err != nil {
		return fmt.Errorf("create order command: %w", err)
	}

	return nil
}

func (a Application) CancelOrder(ctx context.Context, cmd CancelOrder) error {
	order, err := a.orders.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = order.Cancel(); err != nil {
		return err
	}

	if err = a.shopping.Cancel(ctx, order.ShoppingID); err != nil {
		return err
	}

	if err = a.notifications.NotifyOrderCanceled(ctx, order.ID, order.CustomerID); err != nil {
		return err
	}

	return a.orders.Update(ctx, order)
}

func (a Application) ReadyOrder(ctx context.Context, cmd ReadyOrder) error {
	order, err := a.orders.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = order.Ready(); err != nil {
		return nil
	}

	if err = a.orders.Update(ctx, order); err != nil {
		return err
	}

	if err = a.notifications.NotifyOrderReady(ctx, order.ID, order.CustomerID); err != nil {
		return err
	}

	return a.orders.Update(ctx, order)
}

func (a Application) CompleteOrder(ctx context.Context, cmd CompleteOrder) error {
	order, err := a.orders.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = order.Complete(cmd.InvoiceID)
	if err != nil {
		return nil
	}

	return a.orders.Update(ctx, order)
}

func (a Application) GetOrder(ctx context.Context, query GetOrder) (*domain.Order, error) {
	order, err := a.orders.Find(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("get order query: %w", err)
	}

	return order, nil
}
