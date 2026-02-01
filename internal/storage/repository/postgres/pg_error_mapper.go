// internal/storage/repository/postgres/pg_error_mapper.go
package postgres

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"sch.dev/my-kasir-gw/internal/storage/repository"
)

func mapPgxError(err error) error {
	var pgErr *pgconn.PgError

	if !errors.As(err, &pgErr) {
		return err
	}

	switch pgErr.Code {
	case pgerrcode.UniqueViolation:
		return repository.ErrAlreadyExists
	case pgerrcode.ForeignKeyViolation:
		return repository.ErrForeignKeyViolation
	case pgerrcode.NotNullViolation,
		pgerrcode.CheckViolation:
		return repository.ErrConstraintViolation
	default:
		return err
	}
}
