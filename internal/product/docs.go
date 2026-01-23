// internal/handler/product/docs.go

package product

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
	Success bool              `json:"success" example:"false"`
	Message string            `json:"message" example:"Validation error"`
	Errors  map[string]string `json:"errors" example:"name:required,price:gt=0"`
}

type InvalidBodyResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Invalid body"`
}

type InternalServerErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Internal server error"`
}
