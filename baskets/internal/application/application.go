package application

import (
	"context"
	"fmt"

	"mall/baskets/internal/domain"
)

type StartBasket struct {
	ID         string
	CustomerID string
}

type CancelBasket struct {
	ID string
}

type CheckoutBasket struct {
	ID        string
	PaymentID string
}

type AddItem struct {
	ID        string
	ProductID string
	Quantity  int
}

type RemoveItem struct {
	ID        string
	ProductID string
	Quantity  int
}

type GetBasket struct {
	ID string
}

type App interface {
	StartBasket(ctx context.Context, cmd StartBasket) error
	CancelBasket(ctx context.Context, cmd CancelBasket) error
	CheckoutBasket(ctx context.Context, cmd CheckoutBasket) error
	AddItem(ctx context.Context, cmd AddItem) error
	RemoveItem(ctx context.Context, cmd RemoveItem) error
	GetBasket(ctx context.Context, query GetBasket) (*domain.Basket, error)
}

type Application struct {
	baskets  domain.BasketRepository
	stores   domain.StoreRepository
	products domain.ProductRepository
	orders   domain.OrderRepository
}

var _ App = (*Application)(nil)

func New(baskets domain.BasketRepository, stores domain.StoreRepository, products domain.ProductRepository, orders domain.OrderRepository) *Application {
	return &Application{
		baskets:  baskets,
		stores:   stores,
		products: products,
		orders:   orders,
	}
}

func (a Application) StartBasket(ctx context.Context, cmd StartBasket) error {
	basket, err := domain.StartBasket(cmd.ID, cmd.CustomerID)
	if err != nil {
		return err
	}

	return a.baskets.Save(ctx, basket)
}

func (a Application) CancelBasket(ctx context.Context, cmd CancelBasket) error {
	basket, err := a.baskets.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = basket.Cancel()
	if err != nil {
		return err
	}

	return a.baskets.Update(ctx, basket)
}

func (a Application) CheckoutBasket(ctx context.Context, cmd CheckoutBasket) error {
	basket, err := a.baskets.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = basket.Checkout(cmd.PaymentID)
	if err != nil {
		return fmt.Errorf("baskets checkout: %w", err)
	}

	// submit the basket to the order module
	_, err = a.orders.Save(ctx, basket)
	if err != nil {
		return fmt.Errorf("baskets checkout: %w", err)
	}

	return a.baskets.Update(ctx, basket)
}

func (a Application) AddItem(ctx context.Context, cmd AddItem) error {
	basket, err := a.baskets.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	product, err := a.products.Find(ctx, cmd.ProductID)
	if err != nil {
		return err
	}

	store, err := a.stores.Find(ctx, product.StoreID)
	if err != nil {
		return err
	}
	err = basket.AddItem(store, product, cmd.Quantity)
	if err != nil {
		return err
	}

	return a.baskets.Update(ctx, basket)
}

func (a Application) RemoveItem(ctx context.Context, cmd RemoveItem) error {
	product, err := a.products.Find(ctx, cmd.ProductID)
	if err != nil {
		return err
	}

	basket, err := a.baskets.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = basket.RemoveItem(product, cmd.Quantity)
	if err != nil {
		return err
	}

	return a.baskets.Update(ctx, basket)
}

func (a Application) GetBasket(ctx context.Context, query GetBasket) (*domain.Basket, error) {
	return a.baskets.Find(ctx, query.ID)
}
