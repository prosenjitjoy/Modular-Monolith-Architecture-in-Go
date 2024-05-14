package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"mall/payments/internal/domain"
)

type PaymentRepository struct {
	tableName string
	db        *sql.DB
}

var _ domain.PaymentRepository = (*PaymentRepository)(nil)

func NewPaymentRepository(tableName string, db *sql.DB) PaymentRepository {
	return PaymentRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r PaymentRepository) Save(ctx context.Context, payment *domain.Payment) error {
	const query = "INSERT INTO %s (id, customer_id, amount) VALUES ($1, $2, $3)"

	_, err := r.db.ExecContext(ctx, r.table(query), payment.ID, payment.CustomerID, payment.Amount)

	return err
}

func (r PaymentRepository) Find(ctx context.Context, paymentID string) (*domain.Payment, error) {
	const query = "SELECT customer_id, amount FROM %s WHERE id = $1 LIMIT 1"

	payment := &domain.Payment{
		ID: paymentID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), paymentID).Scan(&payment.CustomerID, &payment.Amount)

	return payment, err
}

func (r PaymentRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
