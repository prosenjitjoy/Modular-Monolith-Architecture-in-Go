package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"mall/stores/internal/domain"
)

type StoreRepository struct {
	tableName string
	db        *sql.DB
}

var _ domain.StoreRepository = (*StoreRepository)(nil)

func NewStoreRepository(tableName string, db *sql.DB) StoreRepository {
	return StoreRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r StoreRepository) Find(ctx context.Context, storeID string) (*domain.Store, error) {
	const query = "SELECT name, location, participating FROM %s WHERE id = $1 LIMIT 1"

	store := &domain.Store{
		ID: storeID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), storeID).Scan(&store.Name, &store.Location, &store.Participating)
	if err != nil {
		return nil, fmt.Errorf("scanning store: %w", err)
	}

	return store, nil
}

func (r StoreRepository) FindAll(ctx context.Context) ([]*domain.Store, error) {
	const query = "SELECT id, name, location, participating FROM %s"

	rows, err := r.db.QueryContext(ctx, r.table(query))
	if err != nil {
		return nil, fmt.Errorf("querying stores: %w", err)
	}
	defer rows.Close()

	var stores []*domain.Store
	for rows.Next() {
		store := &domain.Store{}
		err := rows.Scan(&store.ID, &store.Name, &store.Location, &store.Participating)
		if err != nil {
			return nil, fmt.Errorf("scanning store: %w", err)
		}

		stores = append(stores, store)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("finishing store rows: %w", err)
	}

	return stores, nil
}

func (r StoreRepository) FindParticipatingStores(ctx context.Context) ([]*domain.Store, error) {
	const query = "SELECT id, name, location, participating FROM %s WHERE participating is true"

	rows, err := r.db.QueryContext(ctx, r.table(query))
	if err != nil {
		return nil, fmt.Errorf("querying participating stores: %w", err)
	}
	defer rows.Close()

	var stores []*domain.Store
	for rows.Next() {
		store := &domain.Store{}
		err := rows.Scan(&store.ID, &store.Name, &store.Location, &store.Participating)
		if err != nil {
			return nil, fmt.Errorf("scanning participating store: %w", err)
		}

		stores = append(stores, store)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("finishing participating store rows: %w", err)
	}

	return stores, nil
}

func (r StoreRepository) Save(ctx context.Context, store *domain.Store) error {
	const query = "INSERT INTO %s (id, name, location, participating) VALUES ($1, $2, $3, $4)"

	_, err := r.db.ExecContext(ctx, r.table(query), store.ID, store.Name, store.Location, store.Participating)
	if err != nil {
		return fmt.Errorf("inserting store: %w", err)
	}

	return nil
}

func (r StoreRepository) Update(ctx context.Context, store *domain.Store) error {
	const query = "UPDATE %s SET name = $2, location = $3, participating = $4 WHERE id = $1"

	_, err := r.db.ExecContext(ctx, r.table(query), store.ID, store.Name, store.Location, store.Participating)
	if err != nil {
		return fmt.Errorf("updating store: %w", err)
	}

	return nil
}

func (r StoreRepository) Delete(ctx context.Context, storeID string) error {
	const query = "DELETE FROM %s WHERE id = $1"

	_, err := r.db.ExecContext(ctx, r.table(query), storeID)
	if err != nil {
		return fmt.Errorf("deleting store: %w", err)
	}

	return nil
}

func (r StoreRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
