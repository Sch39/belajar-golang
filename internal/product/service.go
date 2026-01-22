// internal\product\service.go
package product

import (
	"context"
	"time"

	"github.com/google/uuid"
	"sch.dev/my-kasir-gw/internal/domain"
)

type Service interface {
	Add(ctx context.Context, product *domain.Product) error
	GetAll(ctx context.Context) ([]domain.Product, error)
	GetByID(ctx context.Context, id string) (*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Add(ctx context.Context, product *domain.Product) error {
	uuidValue, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	product.CreatedAt = now
	product.UpdatedAt = now
	product.IsActive = true
	product.ID = uuidValue.String()
	return s.repo.Create(ctx, product)
}

func (s *service) GetAll(ctx context.Context) ([]domain.Product, error) {
	return s.repo.FindAll(ctx)
}

func (s *service) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrNotFound
	}

	return product, nil
}

func (s *service) Update(ctx context.Context, product *domain.Product) error {
	_, err := s.GetByID(ctx, product.ID)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	product.UpdatedAt = now
	return s.repo.Update(ctx, product)
}

func (s *service) Delete(ctx context.Context, id string) error {
	p, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	p.IsActive = false
	p.DeletedAt = &now

	return s.repo.Delete(ctx, p)
}
