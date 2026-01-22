// internal\product\request.go
package product

type upsertRequest struct {
	Name  string `json:"name" validate:"required"`
	Price *int    `json:"price" validate:"required,min=0"`
	Stock *int    `json:"stock" validate:"required,min=0"`
}