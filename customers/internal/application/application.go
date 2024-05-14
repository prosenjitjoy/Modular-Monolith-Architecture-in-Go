package application

import (
	"context"
	"fmt"

	"mall/customers/internal/domain"
)

type RegisterCustomer struct {
	ID        string
	Name      string
	SmsNumber string
}

type AuthorizeCustomer struct {
	ID string
}

type GetCustomer struct {
	ID string
}

type EnableCustomer struct {
	ID string
}

type DisableCustomer struct {
	ID string
}

type App interface {
	RegisterCustomer(ctx context.Context, cmd RegisterCustomer) error
	EnableCustomer(ctx context.Context, cmd EnableCustomer) error
	DisableCustomer(ctx context.Context, cmd DisableCustomer) error
	AuthorizeCustomer(ctx context.Context, query AuthorizeCustomer) error
	GetCustomer(ctx context.Context, query GetCustomer) (*domain.Customer, error)
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

func (a Application) RegisterCustomer(ctx context.Context, cmd RegisterCustomer) error {
	customer, err := domain.RegisterCustomer(cmd.ID, cmd.Name, cmd.SmsNumber)
	if err != nil {
		return err
	}

	return a.customers.Save(ctx, customer)
}

func (a Application) EnableCustomer(ctx context.Context, cmd EnableCustomer) error {
	customer, err := a.customers.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = customer.Enable()
	if err != nil {
		return err
	}

	return a.customers.Update(ctx, customer)
}

func (a Application) DisableCustomer(ctx context.Context, cmd DisableCustomer) error {
	customer, err := a.customers.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = customer.Disable()
	if err != nil {
		return err
	}

	return a.customers.Update(ctx, customer)
}

func (a Application) AuthorizeCustomer(ctx context.Context, query AuthorizeCustomer) error {
	customer, err := a.customers.Find(ctx, query.ID)
	if err != nil {
		return err
	}

	if !customer.Enabled {
		return fmt.Errorf("customer is not authorized")
	}

	return nil
}

func (a Application) GetCustomer(ctx context.Context, query GetCustomer) (*domain.Customer, error) {
	return a.customers.Find(ctx, query.ID)
}
