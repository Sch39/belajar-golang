// internal/category/docs.go

package category

// success responses
type BaseSuccessResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Category deleted successfully"`
}

type CategorySuccessResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Operation completed successfully"`
	Data    categoryResponse `json:"data"`
}

type CategoryListSuccessResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Operation completed successfully"`
	Data    []categoryResponse `json:"data"`
}

// error responses
type ValidateErrorResponse struct {
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

type CategoryNotFoundResponse struct {
	Success   bool   `json:"success" example:"false"`
	Message   string `json:"message" example:"Category not found"`
	ErrorCode string `json:"error_code" example:"CATEGORY_NOT_FOUND"`
}