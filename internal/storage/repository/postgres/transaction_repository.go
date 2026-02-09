// internal/storage/repository/postgres/transaction_repository.go
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
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
	// 1. Aggregate quantities per product to handle duplicates and prepare for locking
	productQuantities := make(map[string]int32)
	productIDs := make([]string, 0, len(details))
	seenIDs := make(map[string]bool)

	for _, d := range details {
		productQuantities[d.ProductID] += d.Quantity
		if !seenIDs[d.ProductID] {
			productIDs = append(productIDs, d.ProductID)
			seenIDs[d.ProductID] = true
		}
	}

	dbTx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer dbTx.Rollback(ctx)

	// 2. Lock & Validate (Batch)
	// Lock products in ID order to prevent deadlocks and validate stock/status
	queryLock := `
		SELECT id, stock, is_active 
		FROM products 
		WHERE id = ANY($1) 
		ORDER BY id 
		FOR UPDATE
	`
	rowsLock, err := dbTx.Query(ctx, queryLock, productIDs)
	if err != nil {
		return err
	}

	foundCount := 0
	for rowsLock.Next() {
		var id string
		var stock int32
		var isActive bool
		if err := rowsLock.Scan(&id, &stock, &isActive); err != nil {
			rowsLock.Close()
			return err
		}

		foundCount++

		if !isActive {
			rowsLock.Close()
			return fmt.Errorf("product %s is inactive", id)
		}

		needed := productQuantities[id]
		if stock < needed {
			rowsLock.Close()
			return fmt.Errorf("insufficient stock for product %s: have %d, need %d", id, stock, needed)
		}
	}
	rowsLock.Close()

	if foundCount != len(productIDs) {
		return fmt.Errorf("some products not found or invalid")
	}

	// 3. Insert Transaction
	queryTx := `
		INSERT INTO transactions (id, total_price, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err = dbTx.Exec(ctx, queryTx, tx.ID, tx.TotalPrice, tx.CreatedAt, tx.UpdatedAt)
	if err != nil {
		return err
	}

	// 4. Bulk Insert Transaction Details using CopyFrom
	var rows [][]interface{}
	for _, d := range details {
		rows = append(rows, []interface{}{
			d.ID, d.TransactionID, d.ProductID, d.Quantity, d.Price, d.CreatedAt, d.UpdatedAt,
		})
	}

	_, err = dbTx.CopyFrom(
		ctx,
		pgx.Identifier{"transaction_details"},
		[]string{"id", "transaction_id", "product_id", "quantity", "price", "created_at", "updated_at"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return err
	}

	// 5. Batch Update Stock
	// Update stock based on aggregated quantities
	batch := &pgx.Batch{}
	queryStock := `
		UPDATE products 
		SET stock = stock - $1, updated_at = $2 
		WHERE id = $3
	`
	for id, qty := range productQuantities {
		batch.Queue(queryStock, qty, tx.UpdatedAt, id)
	}

	br := dbTx.SendBatch(ctx, batch)
	defer br.Close()

	for range productQuantities {
		_, err := br.Exec()
		if err != nil {
			return err
		}
	}

	if err := br.Close(); err != nil {
		return err
	}

	return dbTx.Commit(ctx)
}

func (r *transactionRepository) GetReport(ctx context.Context, startDate, endDate time.Time) (int64, int64, string, int64, error) {
	var totalRevenue int64
	var totalTransaction int64
	var bestProductName string
	var bestProductQty int64

	// 1. Get Total Revenue and Total Transaction
	queryStats := `
		SELECT COALESCE(SUM(total_price), 0), COUNT(id)
		FROM transactions
		WHERE created_at >= $1 AND created_at <= $2
	`
	err := r.pool.QueryRow(ctx, queryStats, startDate, endDate).Scan(&totalRevenue, &totalTransaction)
	if err != nil {
		return 0, 0, "", 0, err
	}

	// 2. Get Best Selling Product
	queryBestSeller := `
		SELECT p.name, COALESCE(SUM(td.quantity), 0)
		FROM transaction_details td
		JOIN products p ON p.id = td.product_id
		JOIN transactions t ON t.id = td.transaction_id
		WHERE t.created_at >= $1 AND t.created_at <= $2
		GROUP BY p.name
		ORDER BY SUM(td.quantity) DESC
		LIMIT 1
	`
	err = r.pool.QueryRow(ctx, queryBestSeller, startDate, endDate).Scan(&bestProductName, &bestProductQty)
	if err != nil {
		if err == pgx.ErrNoRows {
			// No transactions yet
			return totalRevenue, totalTransaction, "", 0, nil
		}
		return 0, 0, "", 0, err
	}

	return totalRevenue, totalTransaction, bestProductName, bestProductQty, nil
}
