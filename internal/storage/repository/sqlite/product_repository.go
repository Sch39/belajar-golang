// internal/storage/repository/sqlite/product_repository.go
package sqlite

import (
	"context"
	"database/sql"

	"sch.dev/my-kasir-gw/internal/domain"
	"sch.dev/my-kasir-gw/internal/product"
)

type repository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) product.Repository {
	return &repository{db: db}
}

// Implement the methods of ProductRepository interface
func (r *repository) Create(ctx context.Context, product *domain.Product) error {
	query := "INSERT INTO products (id, name, price, stock, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := r.db.ExecContext(ctx, query, product.ID, product.Name, product.Price, product.Stock, product.CreatedAt, product.UpdatedAt)

	return err
}

func (r *repository) FindAll(ctx context.Context) ([]domain.Product, error) {
	query := "SELECT id, name, price, stock, created_at, updated_at FROM products WHERE is_active = 1"
	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	products := []domain.Product{}

	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	query := "SELECT id, name, price, stock, created_at, updated_at FROM products WHERE id = ? AND is_active = 1"
	row := r.db.QueryRowContext(ctx, query, id)
	var product domain.Product
	if err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *repository) Update(ctx context.Context, product *domain.Product) error {
	query := "UPDATE products SET name = ?, price = ?, stock = ?, updated_at = ? WHERE id = ? AND is_active = 1"
	_, err := r.db.ExecContext(ctx, query, product.Name, product.Price, product.Stock, product.UpdatedAt, product.ID)
	return err
}

func (r *repository) Delete(ctx context.Context, product *domain.Product) error {
	query := "UPDATE products SET is_active = 0, deleted_at = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, product.DeletedAt, product.ID)
	return err
}
