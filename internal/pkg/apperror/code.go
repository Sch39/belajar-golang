// internal/pkg/apperror/code.go
package apperror

type ErrorCode string

const (
	// general errors
	ErrInternal       ErrorCode = "INTERNAL_ERROR"
	ErrCodeNotFound   ErrorCode = "NOT_FOUND"
	ErrValidation     ErrorCode = "VALIDATION_ERROR"
	ErrInvalidPayload ErrorCode = "INVALID_PAYLOAD"

	// product errors
	ErrProductNotFound ErrorCode = "PRODUCT_NOT_FOUND"

	// category errors
	ErrCategoryNotFound ErrorCode = "CATEGORY_NOT_FOUND"
)

func (e ErrorCode) ToHttpStatus() int {
	switch e {
	case ErrInvalidPayload:
		return 400
	case ErrCodeNotFound, ErrProductNotFound, ErrCategoryNotFound:
		return 404
	case ErrValidation:
		return 422
	default:
		return 500
	}
}
