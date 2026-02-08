// internal/transaction/repository.go
package transaction

import (
	"context"

	"sch.dev/my-kasir-gw/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, transaction *domain.Transaction, details []domain.TransactionDetail) error
}
