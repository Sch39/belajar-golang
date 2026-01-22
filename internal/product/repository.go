// internal\product\repository.go
package product

import (
	"context"
	"database/sql"
)

type Repository interface {
	Create(ctx context.Context, product *Product) (error)
	FindAll(ctx context.Context) ([]Product, error)
	FindByID(ctx context.Context, id string) (*Product, error)
	Update(ctx context.Context, product *Product) (error)
	Delete(ctx context.Context, id string) (error)
}

type sqliteRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &sqliteRepository{db: db}
}

// Implement the methods of ProductRepository interface
func (r *sqliteRepository) Create(ctx context.Context, product *Product) error {
	query := "INSERT INTO products (id, name, price, stock) VALUES (?, ?, ?, ?)"
	_, err := r.db.ExecContext(ctx, query, product.ID, product.Name, product.Price, product.Stock)

	return err
}

func (r *sqliteRepository) FindAll(ctx context.Context) ([]Product, error) {
	query := "SELECT id, name, price, stock FROM products"
	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	products := []Product{}

	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *sqliteRepository) FindByID(ctx context.Context, id string) (*Product, error) {
	query := "SELECT id, name, price, stock FROM products WHERE id = ?"
	row := r.db.QueryRowContext(ctx, query, id)
	var product Product
	if err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *sqliteRepository) Update(ctx context.Context, product *Product) error {
	query := "UPDATE products SET name = ?, price = ?, stock = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, product.Name, product.Price, product.Stock, product.ID)
	return err
}

func (r *sqliteRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM products WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}