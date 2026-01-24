// internal/category/repository.go
package category

import (
	"context"

	"sch.dev/my-kasir-gw/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, category *domain.Category) error
	FindAll(ctx context.Context) ([]domain.Category, error)
	FindByID(ctx context.Context, id string) (*domain.Category, error)
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, category *domain.Category) error
}
