// internal/transaction/request.go
package transaction

type checkoutItemRequest struct {
	ProductID string `json:"product_id" validate:"required" example:"9825b44a-101f-4c6e-8d8a-675204481359"`
	Quantity  int32  `json:"quantity" validate:"required,min=1" example:"2"`
}

type checkoutRequest struct {
	Items []checkoutItemRequest `json:"items" validate:"required,min=1,dive" example:"[]"`
}
