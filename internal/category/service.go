// internal\category\service.go
package category

import (
	"context"
	"time"

	"github.com/google/uuid"
	"sch.dev/my-kasir-gw/internal/domain"
)

type Service interface {
	Add(ctx context.Context, category *domain.Category) error
	GetAll(ctx context.Context) ([]domain.Category, error)
	GetByID(ctx context.Context, id string) (*domain.Category, error)
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Add(ctx context.Context, category *domain.Category) error {
	uuidValue, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	category.CreatedAt = now
	category.UpdatedAt = now
	category.IsActive = true
	category.ID = uuidValue.String()
	return s.repo.Create(ctx, category)
}

func (s *service) GetAll(ctx context.Context) ([]domain.Category, error) {
	return s.repo.FindAll(ctx)
}

func (s *service) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrNotFound
	}

	return category, nil
}

func (s *service) Update(ctx context.Context, category *domain.Category) error {
	c, err := s.GetByID(ctx, category.ID)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	category.UpdatedAt = now
	category.CreatedAt = c.CreatedAt
	return s.repo.Update(ctx, category)
}

func (s *service) Delete(ctx context.Context, id string) error {
	c, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	c.IsActive = false
	c.DeletedAt = &now

	return s.repo.Delete(ctx, c)
}
