// internal\domain\product.go
package domain

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int64    `json:"price"`
	Stock int    `json:"stock"`

	Timestamp
	SoftDelete
}