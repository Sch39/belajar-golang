// internal/storage/repository/postgres/product_repository.go
package postgres

import (
	"context"
	"database/sql"

	"sch.dev/my-kasir-gw/internal/domain"
	"sch.dev/my-kasir-gw/internal/product"
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) product.Repository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	query := `
		INSERT INTO products (id, name, price, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		product.ID,
		product.Name,
		product.Price,
		product.Stock,
		product.CreatedAt,
		product.UpdatedAt,
	)
	return err
}

func (r *productRepository) FindAll(ctx context.Context) ([]domain.Product, error) {
	query := `
		SELECT id, name, price, stock, created_at, updated_at
		FROM products
		WHERE is_active = true
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Stock,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, rows.Err()
}

func (r *productRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	query := `
		SELECT id, name, price, stock, created_at, updated_at
		FROM products
		WHERE id = $1 AND is_active = true
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var product domain.Product
	if err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	query := `
		UPDATE products
		SET name = $1, price = $2, stock = $3, updated_at = $4
		WHERE id = $5 AND is_active = true
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		product.Name,
		product.Price,
		product.Stock,
		product.UpdatedAt,
		product.ID,
	)
	return err
}

func (r *productRepository) Delete(ctx context.Context, product *domain.Product) error {
	query := `
		UPDATE products
		SET is_active = false, deleted_at = $1
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, product.DeletedAt, product.ID)
	return err
}
