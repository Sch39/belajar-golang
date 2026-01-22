// internal\product\response.go
package product

type productResponse struct {
	ID    string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name  string `json:"name" example:"Kopi Susu Gula Aren"`
	Price int64  `json:"price" example:"15000"`
	Stock int    `json:"stock" example:"10"`
}
