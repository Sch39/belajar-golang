// internal/storage/repository/postgres/category_repository.go
package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"sch.dev/my-kasir-gw/internal/category"
	"sch.dev/my-kasir-gw/internal/domain"
	"sch.dev/my-kasir-gw/internal/storage/repository"
)

type categoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) category.Repository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
	query := `
		INSERT INTO categories (
		id, 
		name, 
		description, 
		created_at, 
		updated_at
		)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(
		ctx,
		query,
		category.ID,
		category.Name,
		category.Description,
		category.CreatedAt,
		category.UpdatedAt,
	)

	if err != nil {
		return mapPgxError(err)
	}

	return nil
}

func (r *categoryRepository) FindAll(ctx context.Context) ([]domain.Category, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM categories
		WHERE is_active = true
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var category domain.Category
		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.CreatedAt,
			&category.UpdatedAt,
		); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, rows.Err()
}

func (r *categoryRepository) FindByID(ctx context.Context, id string) (*domain.Category, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM categories
		WHERE id = $1 AND is_active = true
	`
	row := r.db.QueryRow(ctx, query, id)

	var category domain.Category
	if err := row.Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *domain.Category) error {
	query := `
		UPDATE categories
		SET 
			name = $1, 
			description = $2, 
			updated_at = $3
		WHERE id = $4 AND is_active = true
	`
	ct, err := r.db.Exec(
		ctx,
		query,
		category.Name,
		category.Description,
		category.UpdatedAt,
		category.ID,
	)
	if err != nil {
		return mapPgxError(err)
	}
	if ct.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, category *domain.Category) error {
	query := `
		UPDATE categories
		SET 
			is_active = false, 
			deleted_at = $1
		WHERE id = $2
	`
	ct, err := r.db.Exec(ctx, query, category.DeletedAt, category.ID)
	if err != nil {
		return mapPgxError(err)
	}
	if ct.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}
