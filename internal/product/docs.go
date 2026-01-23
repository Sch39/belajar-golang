// internal/handler/product/docs.go

package product

type BaseSuccessResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Product deleted successfully"`
}

type ProductSuccessResponse struct {
	Success bool            `json:"success" example:"true"`
	Message string          `json:"message" example:"Operation completed successfully"`
	Data    productResponse `json:"data"`
}

type ProductsSuccessResponse struct {
	Success bool              `json:"success" example:"true"`
	Message string            `json:"message" example:"Operation completed successfully"`
	Data    []productResponse `json:"data"`
}

type ValidationErrorResponse struct {
	Success   bool              `json:"success" example:"false"`
	Message   string            `json:"message" example:"Validation error"`
	Errors    map[string]string `json:"errors" example:"name:required,price:gt=0"`
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
	Message   string `json:"message" example:"Product not found"`
	ErrorCode string `json:"error_code" example:"PRODUCT_NOT_FOUND"`
}
