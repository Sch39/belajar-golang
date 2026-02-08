// internal\product\service.go
package product

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"sch.dev/my-kasir-gw/internal/category"
	"sch.dev/my-kasir-gw/internal/domain"
	"sch.dev/my-kasir-gw/internal/storage/repository"
)

type UpsertProductInput struct {
	Name       string
	Price      int64
	Stock      int
	CategoryID string
}

type CategoryOutput struct {
	ID          string
	Name        string
	Description string
}

type ProductOutput struct {
	ID         string
	Name       string
	Price      int64
	Stock      int
	CategoryID string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ProductDetailOutput struct {
	ProductOutput
	Category *CategoryOutput
}

type Service interface {
	Add(ctx context.Context, product UpsertProductInput) (*ProductOutput, error)
	GetAll(ctx context.Context, query string, page, limit int) ([]ProductOutput, int, error)
	GetByID(ctx context.Context, id string) (*ProductDetailOutput, error)
	Update(ctx context.Context, id string, product UpsertProductInput) (*ProductOutput, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	productRepo  Repository
	categoryRepo category.Repository
}

func NewService(
	productRepo Repository,
	categoryRepo category.Repository,
) Service {
	return &service{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *service) Add(ctx context.Context, input UpsertProductInput) (*ProductOutput, error) {
	_, err := s.categoryRepo.FindByID(ctx, input.CategoryID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()

	product := &domain.Product{
		ID:         id.String(),
		Name:       input.Name,
		Price:      input.Price,
		Stock:      input.Stock,
		CategoryID: input.CategoryID,
		Timestamp: domain.Timestamp{
			CreatedAt: now,
			UpdatedAt: now,
		},
		SoftDelete: domain.SoftDelete{
			IsActive: true,
		},
	}
	if err := mapRepoError(s.productRepo.Create(ctx, product)); err != nil {
		return nil, err
	}
	return mapToProductOutput(*product), nil
}

func (s *service) GetAll(ctx context.Context, query string, page, limit int) ([]ProductOutput, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	skip := (page - 1) * limit

	products, err := s.productRepo.FindAll(ctx, query, skip, limit)
	if err != nil {
		return nil, 0, mapRepoError(err)
	}

	total, err := s.productRepo.Count(ctx, query)
	if err != nil {
		return nil, 0, mapRepoError(err)
	}

	var productOutputs []ProductOutput
	for _, p := range products {
		productOut := mapToProductOutput(p)
		productOutputs = append(productOutputs, *productOut)
	}
	return productOutputs, total, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*ProductDetailOutput, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, mapRepoError(err)
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	category, err := s.categoryRepo.FindByID(ctx, product.CategoryID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return mapToProductDetailOutput(*product, *category), nil
}

func (s *service) Update(ctx context.Context, id string, product UpsertProductInput) (*ProductOutput, error) {
	existing, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, mapRepoError(err)
	}
	if existing == nil {
		return nil, ErrProductNotFound
	}

	_, err = s.categoryRepo.FindByID(ctx, product.CategoryID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	existing.Name = product.Name
	existing.Price = product.Price
	existing.Stock = product.Stock
	existing.CategoryID = product.CategoryID
	existing.UpdatedAt = time.Now().UTC()

	if err := mapRepoError(s.productRepo.Update(ctx, existing)); err != nil {
		return nil, err
	}

	return mapToProductOutput(*existing), nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return mapRepoError(err)
	}
	if product == nil {
		return ErrProductNotFound
	}

	now := time.Now().UTC()
	product.IsActive = false
	product.DeletedAt = &now

	return mapRepoError(s.productRepo.Delete(ctx, product))
}
