// internal\pkg\rest\response.go
package rest

import (
	"encoding/json"
	"net/http"

	"sch.dev/my-kasir-gw/internal/pkg/apperror"
)

type Pagination struct {
	Total       int `json:"total" example:"100"`
	TotalPage   int `json:"total_page" example:"10"`
	CurrentPage int `json:"current_page" example:"1"`
	Limit       int `json:"limit" example:"10"`
}

type SuccessResponse struct {
	Success    bool        `json:"success" example:"true"`
	Message    string      `json:"message,omitempty" example:"Operation completed successfully"`
	Data       interface{} `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type FailResponse struct {
	Success   bool               `json:"success" example:"false"`
	Message   string             `json:"message,omitempty" example:"There were validation errors"`
	Errors    map[string]string  `json:"errors,omitempty"`
	ErrorCode apperror.ErrorCode `json:"error_code"`
}

func JSON(w http.ResponseWriter, status int, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

func Success(data interface{}, message ...string) SuccessResponse {
	msg := "Operation completed successfully"
	if len(message) > 0 {
		msg = message[0]
	}
	return SuccessResponse{
		Success: true,
		Data:    data,
		Message: msg,
	}
}

func SuccessWithPagination(data interface{}, pagination Pagination, message ...string) SuccessResponse {
	msg := "Operation completed successfully"
	if len(message) > 0 {
		msg = message[0]
	}
	return SuccessResponse{
		Success:    true,
		Data:       data,
		Message:    msg,
		Pagination: &pagination,
	}
}

func Fail(code apperror.ErrorCode, message string, errors map[string]string) FailResponse {
	return FailResponse{
		Success:   false,
		Message:   message,
		Errors:    errors,
		ErrorCode: code,
	}
}
