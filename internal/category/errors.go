// internal\category\errors.go
package category

import "errors"

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrInvalidInput     = errors.New("invalid category input")
	ErrConflict         = errors.New("category conflict")
)
