// internal\product\model.go
package product

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}