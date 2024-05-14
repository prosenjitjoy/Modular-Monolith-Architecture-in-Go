package domain

import "errors"

var (
	ErrShoppingCannotBeCancelled = errors.New("shopping list cannot be cancelled")
	ErrShoppingCannotBeAssigned  = errors.New("shopping list cannot be assigned")
	ErrShoppingCannotBeCompleted = errors.New("shopping list cannot be completed")
)

type ShoppingListStatus string

const (
	ShoppingListUnknown   ShoppingListStatus = ""
	ShoppingListAvailable ShoppingListStatus = "available"
	ShoppingListAssigned  ShoppingListStatus = "assigned"
	ShoppingListActive    ShoppingListStatus = "active"
	ShoppingListCompleted ShoppingListStatus = "completed"
	ShoppingListCancelled ShoppingListStatus = "cancelled"
)

func (s ShoppingListStatus) String() string {
	switch s {
	case ShoppingListAvailable, ShoppingListAssigned, ShoppingListActive, ShoppingListCompleted, ShoppingListCancelled:
		return string(s)
	default:
		return ""
	}
}

func ToShoppingListStatus(status string) ShoppingListStatus {
	switch status {
	case ShoppingListAvailable.String():
		return ShoppingListAvailable
	case ShoppingListAssigned.String():
		return ShoppingListAssigned
	case ShoppingListActive.String():
		return ShoppingListActive
	case ShoppingListCompleted.String():
		return ShoppingListCompleted
	case ShoppingListCancelled.String():
		return ShoppingListCancelled
	default:
		return ShoppingListUnknown
	}
}

type ShoppingList struct {
	ID            string
	OrderID       string
	Stops         Stops
	AssignedBotID string
	Status        ShoppingListStatus
}

func CreateShopping(id, orderID string) *ShoppingList {
	return &ShoppingList{
		ID:      id,
		OrderID: orderID,
		Status:  ShoppingListAvailable,
		Stops:   make(Stops),
	}
}

func (sl *ShoppingList) AddItem(store *Store, product *Product, quantity int) error {
	if _, exists := sl.Stops[store.ID]; !exists {
		sl.Stops[store.ID] = &Stop{
			StoreName:     store.Name,
			StoreLocation: store.Location,
			Items:         make(Items),
		}
	}

	return sl.Stops[store.ID].AddItem(product, quantity)
}

func (sl *ShoppingList) Cancel() error {
	if sl.Status == ShoppingListUnknown || sl.Status == ShoppingListCompleted {
		return ErrShoppingCannotBeCancelled
	}

	sl.Status = ShoppingListCancelled

	return nil
}

func (sl *ShoppingList) Assign(id string) error {
	if sl.Status != ShoppingListAvailable {
		return ErrShoppingCannotBeAssigned
	}

	sl.AssignedBotID = id
	sl.Status = ShoppingListAssigned

	return nil
}

func (sl *ShoppingList) Complete() error {
	if sl.Status != ShoppingListAssigned {
		return ErrShoppingCannotBeCompleted
	}

	sl.Status = ShoppingListCompleted

	return nil
}
