package transaction

import (
	"errors"

	"sch.dev/my-kasir-gw/internal/pkg/apperror"
	"sch.dev/my-kasir-gw/internal/storage/repository"
)

func mapRepoError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, repository.ErrNotFound):
		return ErrTransactionNotFound

	default:
		return err

	}
}

func mapServiceError(err error) apperror.ErrorCode {
	switch err {
	case ErrTransactionNotFound:
		return apperror.ErrCodeNotFound
	case ErrProductNotFound:
		return apperror.ErrProductNotFound
	case ErrInsufficientStock:
		return apperror.ErrValidation // Atau code khusus seperti INSUFFICIENT_STOCK
	case ErrOptimisticLock:
		return apperror.ErrConflict
	default:
		return apperror.ErrInternal
	}
}

func mapToCheckoutResponse(output CheckoutOutput) checkoutResponse {
	return checkoutResponse{
		ID:         output.ID,
		TotalPrice: output.TotalPrice,
	}
}