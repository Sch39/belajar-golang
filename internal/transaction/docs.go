// internal/transaction/docs.go
package transaction

type CheckoutSuccessResponse struct {
	Success bool             `json:"success" example:"true"`
	Message string           `json:"message" example:"Operation completed successfully"`
	Data    checkoutResponse `json:"data"`
}

type ValidationErrorResponse struct {
	Success   bool              `json:"success" example:"false"`
	Message   string            `json:"message" example:"Validation error"`
	Errors    map[string]string `json:"errors" example:"items[0].product_id:required"`
	ErrorCode string            `json:"error_code" example:"VALIDATION_ERROR"`
}

type InvalidBodyResponse struct {
	Success   bool   `json:"success" example:"false"`
	Message   string `json:"message" example:"Invalid body"`
	ErrorCode string `json:"error_code" example:"INVALID_PAYLOAD"`
}

type InternalServerErrorResponse struct {
	Success   bool   `json:"success" example:"false"`
	Message   string `json:"message" example:"Internal server error"`
	ErrorCode string `json:"error_code" example:"INTERNAL_ERROR"`
}

type ProductNotFoundResponse struct {
	Success   bool   `json:"success" example:"false"`
	Message   string `json:"message" example:"Product not found or inactive"`
	ErrorCode string `json:"error_code" example:"PRODUCT_NOT_FOUND"`
}

type InsufficientStockResponse struct {
	Success   bool   `json:"success" example:"false"`
	Message   string `json:"message" example:"insufficient stock for product ID: ..."`
	ErrorCode string `json:"error_code" example:"VALIDATION_ERROR"`
}