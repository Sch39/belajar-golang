// internal\product\repository.go
package product

import (
	"context"

	"sch.dev/my-kasir-gw/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, product *domain.Product) error
	FindAll(ctx context.Context, query string, skip, limit int) ([]domain.Product, error)
	Count(ctx context.Context, query string) (int, error)
	FindByID(ctx context.Context, id string) (*domain.Product, error)
	FindByIDs(ctx context.Context, ids []string) ([]domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, product *domain.Product) error
}
