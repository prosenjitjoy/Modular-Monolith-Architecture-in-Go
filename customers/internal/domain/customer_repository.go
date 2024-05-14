package domain

import "context"

type CustomerRepository interface {
	Find(ctx context.Context, customerID string) (*Customer, error)
	Save(ctx context.Context, customer *Customer) error
	Update(ctx context.Context, customer *Customer) error
}
