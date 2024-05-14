package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"mall/stores/internal/domain"
)

type ParticipatingStoreRepository struct {
	tableName string
	db        *sql.DB
}

var _ domain.ParticipatingStoreRepository = (*ParticipatingStoreRepository)(nil)

func NewParticipatingStoreRepository(tableName string, db *sql.DB) ParticipatingStoreRepository {
	return ParticipatingStoreRepository{tableName: tableName, db: db}
}

func (r ParticipatingStoreRepository) FindAll(ctx context.Context) ([]*domain.Store, error) {
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

func (r ParticipatingStoreRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
