// internal/transaction/repository.go
package transaction

import (
	"context"
	"time"

	"sch.dev/my-kasir-gw/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, transaction *domain.Transaction, details []domain.TransactionDetail) error
	GetReport(ctx context.Context, startDate, endDate time.Time) (int64, int64, string, int64, error)
}
