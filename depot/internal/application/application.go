package application

import (
	"context"
	"fmt"

	"mall/depot/internal/domain"
)

type OrderItem struct {
	StoreID   string
	ProductID string
	Quantity  int
}

type CreateShoppingList struct {
	ID      string
	OrderID string
	Items   []OrderItem
}

type CancelShoppingList struct {
	ID string
}

type AssignShoppingList struct {
	ID    string
	BotID string
}

type CompleteShoppingList struct {
	ID string
}

type GetShoppingList struct {
	ID string
}

type App interface {
	CreateShoppingList(ctx context.Context, cmd CreateShoppingList) error
	CancelShoppingList(ctx context.Context, cmd CancelShoppingList) error
	AssignShoppingList(ctx context.Context, cmd AssignShoppingList) error
	CompleteShoppingList(ctx context.Context, cmd CompleteShoppingList) error
	GetShoppingList(ctx context.Context, query GetShoppingList) (*domain.ShoppingList, error)
}

type Application struct {
	shoppingLists domain.ShoppingListRepository
	stores        domain.StoreRepository
	products      domain.ProductRepository
	orders        domain.OrderRepository
}

var _ App = (*Application)(nil)

func New(shoppingLists domain.ShoppingListRepository, stores domain.StoreRepository, products domain.ProductRepository, orders domain.OrderRepository) *Application {
	return &Application{
		shoppingLists: shoppingLists,
		stores:        stores,
		products:      products,
		orders:        orders,
	}
}

func (a Application) CreateShoppingList(ctx context.Context, cmd CreateShoppingList) error {
	list := domain.CreateShopping(cmd.ID, cmd.OrderID)

	for _, item := range cmd.Items {
		// horribly inefficient
		store, err := a.stores.Find(ctx, item.StoreID)
		if err != nil {
			return fmt.Errorf("building shopping list: %w", err)
		}

		product, err := a.products.Find(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("building shopping list: %w", err)
		}

		err = list.AddItem(store, product, item.Quantity)
		if err != nil {
			return fmt.Errorf("building shopping list: %w", err)
		}
	}

	err := a.shoppingLists.Save(ctx, list)
	if err != nil {
		return fmt.Errorf("scheduling shopping: %w", err)
	}

	return nil
}

func (a Application) CancelShoppingList(ctx context.Context, cmd CancelShoppingList) error {
	list, err := a.shoppingLists.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = list.Cancel()
	if err != nil {
		return err
	}

	return a.shoppingLists.Update(ctx, list)
}

func (a Application) AssignShoppingList(ctx context.Context, cmd AssignShoppingList) error {
	list, err := a.shoppingLists.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = list.Assign(cmd.BotID)
	if err != nil {
		return err
	}

	return a.shoppingLists.Update(ctx, list)
}

func (a Application) CompleteShoppingList(ctx context.Context, cmd CompleteShoppingList) error {
	list, err := a.shoppingLists.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = list.Complete()
	if err != nil {
		return err
	}

	err = a.orders.Ready(ctx, list.OrderID)
	if err != nil {
		return err
	}

	return a.shoppingLists.Update(ctx, list)
}

func (a Application) GetShoppingList(ctx context.Context, query GetShoppingList) (*domain.ShoppingList, error) {
	return a.shoppingLists.Find(ctx, query.ID)
}
