// internal\product\errors.go
package product

import "errors"

var (
	ErrProductNotFound  = errors.New("product not found or inactive")
	ErrCategoryNotFound = errors.New("category not found")
	ErrInvalidInput     = errors.New("invalid product input")
	ErrConflict         = errors.New("product conflict")
)
