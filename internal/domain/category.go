// internal/domain/category.go
package domain

type Category struct {
	ID          string
	Name        string
	Description string

	Timestamp
	SoftDelete
}
