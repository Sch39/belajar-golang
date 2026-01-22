// internal/service/product_service.go
package product

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	Add(ctx context.Context, product *Product) error
	GetAll(ctx context.Context) ([]Product, error)
	GetByID(ctx context.Context, id string) (*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Add(ctx context.Context, product *Product) error {
	uuidValue, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	product.ID = uuidValue.String()
	return s.repo.Create(ctx, product)
}


func (s *service) GetAll(ctx context.Context) ([]Product, error) {
	return s.repo.FindAll(ctx)
}

func (s *service) GetByID(ctx context.Context, id string) (*Product, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Update(ctx context.Context, product *Product) error {
	return s.repo.Update(ctx, product)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}