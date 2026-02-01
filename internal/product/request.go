// internal\product\request.go
package product

type upsertRequest struct {
	Name       string `json:"name" validate:"required" example:"Kopi Susu Gula Aren"`
	Price      *int64 `json:"price" validate:"required,min=0" example:"15000"`
	Stock      *int   `json:"stock" validate:"required,min=0" example:"10"`
	CategoryID string `json:"category_id" validate:"required" example:"pohon-sawit-0001"`
}
