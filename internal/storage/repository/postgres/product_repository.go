// internal/storage/repository/postgres/product_repository.go
package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"sch.dev/my-kasir-gw/internal/domain"
	"sch.dev/my-kasir-gw/internal/product"
	"sch.dev/my-kasir-gw/internal/storage/repository"
)

type productRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) product.Repository {
	return &productRepository{pool: pool}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	query := `
		INSERT INTO products (
			id, 
			name, 
			price, 
			stock, 
			created_at, 
			updated_at,
			category_id
		 )
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.pool.Exec(
		ctx,
		query,
		product.ID,
		product.Name,
		product.Price,
		product.Stock,
		product.CreatedAt,
		product.UpdatedAt,
		product.CategoryID,
	)
	if err != nil {
		return mapPgxError(err)
	}

	return nil
}

func (r *productRepository) FindAll(ctx context.Context, query string, skip, limit int) ([]domain.Product, error) {
	q := `
		SELECT id, name, price, stock, created_at, updated_at, category_id
		FROM products
		WHERE is_active = true
		AND name ILIKE $1
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, q, "%"+query+"%", limit, skip)
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
			&product.CategoryID,
		); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, rows.Err()
}

func (r *productRepository) Count(ctx context.Context, query string) (int, error) {
	q := `
		SELECT COUNT(id)
		FROM products
		WHERE is_active = true
		AND name ILIKE $1
	`
	var count int
	err := r.pool.QueryRow(ctx, q, "%"+query+"%").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *productRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	query := `
		SELECT id, name, price, stock, created_at, updated_at, category_id
		FROM products
		WHERE id = $1 AND is_active = true
	`
	row := r.pool.QueryRow(ctx, query, id)

	var product domain.Product
	if err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.CategoryID,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) FindByIDs(ctx context.Context, ids []string) ([]domain.Product, error) {
	q := `
		SELECT id, name, price, stock, created_at, updated_at, category_id
		FROM products
		WHERE id = ANY($1) AND is_active = true
	`
	rows, err := r.pool.Query(ctx, q, ids)
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
			&product.CategoryID,
		); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, rows.Err()
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	query := `
		UPDATE products
		SET 
			name = $1, 
			price = $2, 
			stock = $3, 
			updated_at = $4, 
			category_id = $5
		WHERE id = $6 AND is_active = true
	`
	ct, err := r.pool.Exec(
		ctx,
		query,
		product.Name,
		product.Price,
		product.Stock,
		product.UpdatedAt,
		product.CategoryID,
		product.ID,
	)
	if err != nil {
		return mapPgxError(err)
	}
	if ct.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *productRepository) Delete(ctx context.Context, product *domain.Product) error {
	query := `
		UPDATE products
		SET 
			is_active = false, 
			deleted_at = $1
		WHERE id = $2
	`
	ct, err := r.pool.Exec(
		ctx,
		query,
		product.DeletedAt,
		product.ID,
	)
	if err != nil {
		return mapPgxError(err)
	}
	if ct.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}
