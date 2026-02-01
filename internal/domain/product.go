// internal\domain\product.go
package domain

type Product struct {
	ID    string
	Name  string
	Price int64
	Stock int

	CategoryID string
	Category   *Category

	Timestamp
	SoftDelete
}
