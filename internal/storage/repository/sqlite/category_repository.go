// internal\storage\repository\sqlite\category_repository.go
package sqlite

import (
	"context"
	"database/sql"

	"sch.dev/my-kasir-gw/internal/category"
	"sch.dev/my-kasir-gw/internal/domain"
)

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) category.Repository {
	return &categoryRepository{db: db}
}

// Implement the methods of CategoryRepository interface

func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
	query := "INSERT INTO categories (id, name, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	_, err := r.db.ExecContext(ctx, query, category.ID, category.Name, category.Description, category.CreatedAt, category.UpdatedAt)

	return err
}

func (r *categoryRepository) FindAll(ctx context.Context) ([]domain.Category, error) {
	query := "SELECT id, name, description, created_at, updated_at FROM categories WHERE is_active = 1"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	categories := []domain.Category{}

	for rows.Next() {
		var category domain.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepository) FindByID(ctx context.Context, id string) (*domain.Category, error) {
	query := "SELECT id, name, description, created_at, updated_at FROM categories WHERE id = ? AND is_active = 1"
	row := r.db.QueryRowContext(ctx, query, id)
	var category domain.Category
	if err := row.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *domain.Category) error {
	query := "UPDATE categories SET name = ?, description = ?, updated_at = ? WHERE id = ? AND is_active = 1"
	_, err := r.db.ExecContext(ctx, query, category.Name, category.Description, category.UpdatedAt, category.ID)
	return err
}

func (r *categoryRepository) Delete(ctx context.Context, category *domain.Category) error {
	query := "UPDATE categories SET is_active = 0, deleted_at = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, category.DeletedAt, category.ID)

	return err
}
