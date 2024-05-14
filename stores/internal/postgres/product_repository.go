package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"mall/stores/internal/domain"
)

type ProductRepository struct {
	tableName string
	db        *sql.DB
}

var _ domain.ProductRepository = (*ProductRepository)(nil)

func NewProductRepository(tableName string, db *sql.DB) ProductRepository {
	return ProductRepository{tableName: tableName, db: db}
}

func (r ProductRepository) FindProduct(ctx context.Context, id string) (*domain.Product, error) {
	const query = "SELECT store_id, name, description, sku, price FROM %s WHERE id = $1 LIMIT 1"

	product := &domain.Product{ID: id}

	err := r.db.QueryRowContext(ctx, r.table(query), id).Scan(&product.StoreID, &product.Name, &product.Description, &product.SKU, &product.Price)
	if err != nil {
		return nil, fmt.Errorf("scanning product: %w", err)
	}

	return product, nil
}

func (r ProductRepository) AddProduct(ctx context.Context, product *domain.Product) error {
	const query = "INSERT INTO %s (id, store_id, name, description, sku, price) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := r.db.ExecContext(ctx, r.table(query), product.ID, product.StoreID, product.Name, product.Description, product.SKU, product.Price)
	if err != nil {
		return fmt.Errorf("inserting product: %w", err)
	}

	return nil
}

func (r ProductRepository) RemoveProduct(ctx context.Context, id string) error {
	const query = "DELETE FROM %s WHERE id = $1"

	_, err := r.db.ExecContext(ctx, r.table(query), id)
	if err != nil {
		return fmt.Errorf("deleting product: %w", err)
	}

	return nil
}

func (r ProductRepository) GetCatalog(ctx context.Context, storeID string) ([]*domain.Product, error) {
	const query = "SELECT id, name, description, sku, price FROM %s WHERE store_id = $1"

	rows, err := r.db.QueryContext(ctx, r.table(query), storeID)
	if err != nil {
		return nil, fmt.Errorf("querying products: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		product := &domain.Product{
			StoreID: storeID,
		}
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.SKU, &product.Price)
		if err != nil {
			return nil, fmt.Errorf("scanning product: %w", err)
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("finishing product rows: %w", err)
	}

	return products, nil
}

func (r ProductRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
