// internal/domain/category.go
package domain

type Category struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Timestamp
	SoftDelete
}
