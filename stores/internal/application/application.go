package application

import (
	"context"
	"fmt"

	"mall/stores/internal/domain"
)

type CreateStore struct {
	ID       string
	Name     string
	Location string
}

type EnableParticipation struct {
	ID string
}

type DisableParticipation struct {
	ID string
}

type AddProduct struct {
	ID          string
	StoreID     string
	Name        string
	Description string
	SKU         string
	Price       float64
}

type RemoveProduct struct {
	ID string
}

type GetStore struct {
	ID string
}

type GetStores struct {
}

type GetParticipatingStores struct {
}

type GetCatalog struct {
	StoreID string
}

type GetProduct struct {
	ID string
}

type App interface {
	CreateStore(ctx context.Context, cmd CreateStore) error
	EnableParticipation(ctx context.Context, cmd EnableParticipation) error
	DisableParticipation(ctx context.Context, cmd DisableParticipation) error
	AddProduct(ctx context.Context, cmd AddProduct) error
	RemoveProduct(ctx context.Context, cmd RemoveProduct) error
	GetStore(ctx context.Context, query GetStore) (*domain.Store, error)
	GetStores(ctx context.Context, query GetStores) ([]*domain.Store, error)
	GetParticipatingStores(ctx context.Context, query GetParticipatingStores) ([]*domain.Store, error)
	GetCatalog(ctx context.Context, query GetCatalog) ([]*domain.Product, error)
	GetProduct(ctx context.Context, query GetProduct) (*domain.Product, error)
}

type Application struct {
	stores              domain.StoreRepository
	participatingStores domain.ParticipatingStoreRepository
	products            domain.ProductRepository
}

var _ App = (*Application)(nil)

func New(stores domain.StoreRepository, participatingStores domain.ParticipatingStoreRepository, products domain.ProductRepository) *Application {
	return &Application{
		stores:              stores,
		participatingStores: participatingStores,
		products:            products,
	}
}

func (a Application) CreateStore(ctx context.Context, cmd CreateStore) error {
	store, err := domain.CreateStore(cmd.ID, cmd.Name, cmd.Location)
	if err != nil {
		return err
	}

	return a.stores.Save(ctx, store)
}

func (a Application) EnableParticipation(ctx context.Context, cmd EnableParticipation) error {
	store, err := a.stores.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = store.EnableParticipation()
	if err != nil {
		return err
	}

	return a.stores.Update(ctx, store)
}

func (a Application) DisableParticipation(ctx context.Context, cmd DisableParticipation) error {
	store, err := a.stores.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = store.DisableParticipation()
	if err != nil {
		return err
	}

	return a.stores.Update(ctx, store)
}

func (a Application) AddProduct(ctx context.Context, cmd AddProduct) error {
	_, err := a.stores.Find(ctx, cmd.StoreID)
	if err != nil {
		return fmt.Errorf("error adding product: %w", err)
	}

	product, err := domain.CreateProduct(cmd.ID, cmd.StoreID, cmd.Name, cmd.Description, cmd.SKU, cmd.Price)
	if err != nil {
		return fmt.Errorf("error adding product: %w", err)
	}

	err = a.products.AddProduct(ctx, product)
	if err != nil {
		return fmt.Errorf("error adding product: %w", err)
	}

	return nil
}

func (a Application) RemoveProduct(ctx context.Context, cmd RemoveProduct) error {
	return a.products.RemoveProduct(ctx, cmd.ID)
}

func (a Application) GetStore(ctx context.Context, query GetStore) (*domain.Store, error) {
	return a.stores.Find(ctx, query.ID)
}

func (a Application) GetStores(ctx context.Context, _ GetStores) ([]*domain.Store, error) {
	return a.stores.FindAll(ctx)
}

func (a Application) GetParticipatingStores(ctx context.Context, _ GetParticipatingStores) ([]*domain.Store, error) {
	return a.participatingStores.FindAll(ctx)
}

func (a Application) GetCatalog(ctx context.Context, query GetCatalog) ([]*domain.Product, error) {
	return a.products.GetCatalog(ctx, query.StoreID)
}

func (a Application) GetProduct(ctx context.Context, query GetProduct) (*domain.Product, error) {
	return a.products.FindProduct(ctx, query.ID)
}
