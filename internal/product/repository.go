// internal\product\repository.go
package product

import (
	"context"

	"sch.dev/my-kasir-gw/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, product *domain.Product) error
	FindAll(ctx context.Context) ([]domain.Product, error)
	FindByID(ctx context.Context, id string) (*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, product *domain.Product) error
}
