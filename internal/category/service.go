// internal\category\service.go
package category

import (
	"context"
	"time"

	"github.com/google/uuid"
	"sch.dev/my-kasir-gw/internal/domain"
)

type UpsertCategoryInput struct {
	Name        string
	Description string
}

type CategoryOutput struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Service interface {
	Add(ctx context.Context, input UpsertCategoryInput) (*CategoryOutput, error)
	GetAll(ctx context.Context) ([]CategoryOutput, error)
	GetByID(ctx context.Context, id string) (*CategoryOutput, error)
	Update(ctx context.Context, id string, input UpsertCategoryInput) (*CategoryOutput, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Add(ctx context.Context, input UpsertCategoryInput) (*CategoryOutput, error) {
	uuidValue, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()

	category := &domain.Category{
		ID:          uuidValue.String(),
		Name:        input.Name,
		Description: input.Description,
		Timestamp: domain.Timestamp{
			CreatedAt: now,
			UpdatedAt: now,
		},
		SoftDelete: domain.SoftDelete{
			IsActive: true,
		},
	}

	if err := mapRepoError(s.repo.Create(ctx, category)); err != nil {
		return nil, err
	}

	return mapToCategoryOutput(*category), nil
}

func (s *service) GetAll(ctx context.Context) ([]CategoryOutput, error) {
	categories, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, mapRepoError(err)
	}

	var outputs []CategoryOutput
	for _, c := range categories {
		outputs = append(outputs, *mapToCategoryOutput(c))
	}
	return outputs, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*CategoryOutput, error) {
	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, mapRepoError(err)
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}

	return mapToCategoryOutput(*category), nil
}

func (s *service) Update(ctx context.Context, id string, input UpsertCategoryInput) (*CategoryOutput, error) {
	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, mapRepoError(err)
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}

	now := time.Now().UTC()
	category.Name = input.Name
	category.Description = input.Description
	category.UpdatedAt = now

	if err := mapRepoError(s.repo.Update(ctx, category)); err != nil {
		return nil, err
	}

	return mapToCategoryOutput(*category), nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return mapRepoError(err)
	}
	if c == nil {
		return ErrCategoryNotFound
	}

	now := time.Now().UTC()
	c.IsActive = false
	c.DeletedAt = &now

	return mapRepoError(s.repo.Delete(ctx, c))
}
