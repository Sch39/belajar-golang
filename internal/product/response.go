// internal\product\response.go
package product

import "time"

type productResponse struct {
	ID         string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name       string    `json:"name" example:"Kopi Susu Gula Aren"`
	Price      int64     `json:"price" example:"15000"`
	Stock      int       `json:"stock" example:"10"`
	CategoryID string    `json:"category_id" example:"9825b44a-101f-4c6e-8d8a-675204481359"`
	CreatedAt  time.Time `json:"created_at" example:"2023-01-01T12:00:00Z"`
	UpdatedAt  time.Time `json:"updated_at" example:"2023-01-01T12:00:00Z"`
}

type categoryResponse struct {
	ID          string `json:"id" example:"9825b44a-101f-4c6e-8d8a-675204481359"`
	Name        string `json:"name" example:"Minuman"`
	Description string `json:"description" example:"Aneka minuman segar"`
}

type productDetailResponse struct {
	productResponse
	Category *categoryResponse `json:"category"`
}
