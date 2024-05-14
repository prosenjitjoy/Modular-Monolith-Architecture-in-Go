package domain

import (
	"context"
)

type OrderRepository interface {
	Complete(ctx context.Context, invoiceID, orderID string) error
}
