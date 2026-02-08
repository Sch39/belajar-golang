package transaction

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"sch.dev/my-kasir-gw/internal/domain"
	"sch.dev/my-kasir-gw/internal/product"
	"sch.dev/my-kasir-gw/internal/storage/repository"
)

type CheckoutItemInput struct {
	ProductID string
	Quantity  int32
}

type CheckoutInput struct {
	Items []CheckoutItemInput
}

type CheckoutOutput struct {
	ID         string
	TotalPrice int64
	CreatedAt  time.Time
}

type Service interface {
	Checkout(ctx context.Context, input CheckoutInput) (CheckoutOutput, error)
}

type service struct {
	repo        Repository
	productRepo product.Repository
}

func NewService(repo Repository, productRepo product.Repository) Service {
	return &service{
		repo:        repo,
		productRepo: productRepo,
	}
}

func (s *service) Checkout(ctx context.Context, input CheckoutInput) (CheckoutOutput, error) {
	var totalPrice int64
	var details []domain.TransactionDetail

	txID := uuid.New().String()
	now := time.Now()

	for _, item := range input.Items {
		prod, err := s.productRepo.FindByID(ctx, item.ProductID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return CheckoutOutput{}, ErrProductNotFound
			}
			return CheckoutOutput{}, err
		}

		subTotal := prod.Price * int64(item.Quantity)
		totalPrice += subTotal

		details = append(details, domain.TransactionDetail{
			ID:            uuid.New().String(),
			TransactionID: txID,
			ProductID:     prod.ID,
			Quantity:      item.Quantity,
			Price:         prod.Price,
			Timestamp:     domain.Timestamp{CreatedAt: now, UpdatedAt: now},
		})
	}

	tx := &domain.Transaction{
		ID:         txID,
		TotalPrice: totalPrice,
		Timestamp:  domain.Timestamp{CreatedAt: now, UpdatedAt: now},
	}

	if err := s.repo.Create(ctx, tx, details); err != nil {
		return CheckoutOutput{}, mapRepoError(err)
	}

	return CheckoutOutput{
		ID:         txID,
		TotalPrice: totalPrice,
		CreatedAt:  now,
	}, nil
}
