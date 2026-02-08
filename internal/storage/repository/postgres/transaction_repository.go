// internal/storage/repository/postgres/transaction_repository.go
package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"sch.dev/my-kasir-gw/internal/domain"
	"sch.dev/my-kasir-gw/internal/transaction"
)

type transactionRepository struct {
	pool *pgxpool.Pool
}

func NewTransactionRepository(pool *pgxpool.Pool) transaction.Repository {
	return &transactionRepository{pool: pool}
}

func (r *transactionRepository) Create(ctx context.Context, tx *domain.Transaction, details []domain.TransactionDetail) error {
	dbTx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer dbTx.Rollback(ctx)

	queryTx := `
		INSERT INTO transactions (id, total_price, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err = dbTx.Exec(ctx, queryTx, tx.ID, tx.TotalPrice, tx.CreatedAt, tx.UpdatedAt)
	if err != nil {
		return err
	}

	queryDetail := `
		INSERT INTO transaction_details (id, transaction_id, product_id, quantity, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	queryStock := `
		UPDATE products 
		SET stock = stock - $1, updated_at = $2 
		WHERE id = $3 AND stock >= $1
	`

	for _, d := range details {
		_, err = dbTx.Exec(ctx, queryDetail, d.ID, d.TransactionID, d.ProductID, d.Quantity, d.Price, d.CreatedAt, d.UpdatedAt)
		if err != nil {
			return err
		}

		cmdTag, err := dbTx.Exec(ctx, queryStock, d.Quantity, d.UpdatedAt, d.ProductID)
		if err != nil {
			return err
		}

		if cmdTag.RowsAffected() == 0 {
			return fmt.Errorf("insufficient stock for product ID: %s", d.ProductID)
		}
	}

	return dbTx.Commit(ctx)
}
