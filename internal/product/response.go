// internal\product\response.go
package product

import "time"

type productResponse struct {
	ID         string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name       string    `json:"name" example:"Kopi Susu Gula Aren"`
	Price      int64     `json:"price" example:"15000"`
	Stock      int       `json:"stock" example:"10"`
	CategoryID string    `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type categoryResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type productDetailResponse struct {
	productResponse
	Category *categoryResponse `json:"category"`
}
