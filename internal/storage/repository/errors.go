// internal/storage/repository/errors.go
package repository

import "errors"

var (
	ErrNotFound            = errors.New("data not found")
	ErrNoRowsAffected      = errors.New("no rows affected")
	ErrAlreadyExists       = errors.New("data already exists")
	ErrForeignKeyViolation = errors.New("foreign key violation")
	ErrConstraintViolation = errors.New("constraint violation")
	ErrTransactionFailed   = errors.New("transaction failed")
)
