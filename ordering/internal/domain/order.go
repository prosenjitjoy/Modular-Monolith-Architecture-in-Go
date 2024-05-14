package domain

import (
	"errors"
)

var (
	ErrOrderHasNoItems         = errors.New("order has no items")
	ErrOrderCannotBeCancelled  = errors.New("order cannot be cancelled")
	ErrCustomerIDCannotBeBlank = errors.New("customer id cannot be blank")
	ErrPaymentIDCannotBeBlank  = errors.New("payment id cannot be blank")
	ErrOrderCannotBeReady      = errors.New("order cannot be ready")
	ErrOrderCannotBeCompleted  = errors.New("order cannot be completed")
)

type OrderStatus string

const (
	OrderUnknown   OrderStatus = ""
	OrderPending   OrderStatus = "pending"
	OrderInProcess OrderStatus = "in-progress"
	OrderReady     OrderStatus = "ready"
	OrderCompleted OrderStatus = "completed"
	OrderCancelled OrderStatus = "cancelled"
)

func (s OrderStatus) String() string {
	switch s {
	case OrderPending, OrderInProcess, OrderReady, OrderCompleted, OrderCancelled:
		return string(s)
	default:
		return ""
	}
}

func ToOrderStatus(status string) OrderStatus {
	switch status {
	case OrderPending.String():
		return OrderPending
	case OrderInProcess.String():
		return OrderInProcess
	case OrderReady.String():
		return OrderReady
	case OrderCancelled.String():
		return OrderCancelled
	case OrderCompleted.String():
		return OrderCompleted
	default:
		return OrderUnknown
	}
}

type Order struct {
	ID         string
	CustomerID string
	PaymentID  string
	InvoiceID  string
	ShoppingID string
	Items      []*Item
	Status     OrderStatus
}

func CreateOrder(id, customerID, paymentID string, items []*Item) (*Order, error) {
	if len(items) == 0 {
		return nil, ErrOrderHasNoItems
	}

	if customerID == "" {
		return nil, ErrCustomerIDCannotBeBlank
	}

	if paymentID == "" {
		return nil, ErrPaymentIDCannotBeBlank
	}

	order := &Order{
		ID:         id,
		CustomerID: customerID,
		PaymentID:  paymentID,
		Items:      items,
		Status:     OrderPending,
	}

	return order, nil
}

func (o *Order) Cancel() error {
	if o.Status != OrderPending {
		return ErrOrderCannotBeCancelled
	}

	o.Status = OrderCancelled

	return nil
}

func (o *Order) Ready() error {
	if o.Status != OrderPending {
		return ErrOrderCannotBeReady
	}

	o.Status = OrderReady

	return nil
}

func (o *Order) Complete(invoiceID string) error {
	// validate invoice exists

	if o.Status != OrderReady {
		return ErrOrderCannotBeCompleted
	}

	o.InvoiceID = invoiceID
	o.Status = OrderCompleted

	return nil
}

func (o Order) GetTotal() float64 {
	var total float64

	for _, item := range o.Items {
		total += item.ProductPrice * float64(item.Quantity)
	}

	return total
}
