package application

import (
	"context"
	"fmt"

	"mall/payments/internal/domain"
)

type AuthorizePayment struct {
	ID         string
	CustomerID string
	Amount     float64
}

type ConfirmPayment struct {
	ID string
}

type CreateInvoice struct {
	ID        string
	OrderID   string
	PaymentID string
	Amount    float64
}

type AdjustInvoice struct {
	ID     string
	Amount float64
}

type PayInvoice struct {
	ID string
}

type CancelInvoice struct {
	ID string
}

type App interface {
	AuthorizePayment(ctx context.Context, cmd AuthorizePayment) error
	CreateInvoice(ctx context.Context, cmd CreateInvoice) error
	AdjustInvoice(ctx context.Context, cmd AdjustInvoice) error
	PayInvoice(ctx context.Context, cmd PayInvoice) error
	CancelInvoice(ctx context.Context, cmd CancelInvoice) error
	ConfirmPayment(ctx context.Context, query ConfirmPayment) error
}

type Application struct {
	invoices domain.InvoiceRepository
	payments domain.PaymentRepository
	orders   domain.OrderRepository
}

var _ App = (*Application)(nil)

func New(invoices domain.InvoiceRepository, payments domain.PaymentRepository, orders domain.OrderRepository) *Application {
	return &Application{
		invoices: invoices,
		payments: payments,
		orders:   orders,
	}
}

func (a Application) AuthorizePayment(ctx context.Context, cmd AuthorizePayment) error {
	return a.payments.Save(ctx, &domain.Payment{
		ID:         cmd.ID,
		CustomerID: cmd.CustomerID,
		Amount:     cmd.Amount,
	})
}

func (a Application) CreateInvoice(ctx context.Context, cmd CreateInvoice) error {
	return a.invoices.Save(ctx, &domain.Invoice{
		ID:      cmd.ID,
		OrderID: cmd.OrderID,
		Amount:  cmd.Amount,
		Status:  domain.InvoicePending,
	})
}

func (a Application) AdjustInvoice(ctx context.Context, cmd AdjustInvoice) error {
	invoice, err := a.invoices.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	invoice.Amount = cmd.Amount

	return a.invoices.Update(ctx, invoice)
}

func (a Application) PayInvoice(ctx context.Context, cmd PayInvoice) error {
	invoice, err := a.invoices.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if invoice.Status != domain.InvoicePending {
		return fmt.Errorf("invoice cannot be paid for")
	}

	invoice.Status = domain.InvoicePaid

	if err = a.orders.Complete(ctx, invoice.ID, invoice.OrderID); err != nil {
		return err
	}

	return a.invoices.Update(ctx, invoice)
}

func (a Application) CancelInvoice(ctx context.Context, cmd CancelInvoice) error {
	invoice, err := a.invoices.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if invoice.Status != domain.InvoicePending {
		return fmt.Errorf("invoice cannot be canceled")
	}

	invoice.Status = domain.InvoiceCanceled

	return a.invoices.Update(ctx, invoice)
}

func (a Application) ConfirmPayment(ctx context.Context, query ConfirmPayment) error {
	_, err := a.payments.Find(ctx, query.ID)
	if err != nil {
		return err
	}

	return nil
}
