package transaction

import "errors"

var (
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrProductNotFound     = errors.New("product not found or inactive")
	ErrInsufficientStock   = errors.New("insufficient stock")
	ErrOptimisticLock      = errors.New("transaction conflict, please retry")
)