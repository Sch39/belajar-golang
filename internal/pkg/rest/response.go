// internal\pkg\rest\response.go
package rest

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message,omitempty" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
}

type FailResponse struct {
	Success bool              `json:"success" example:"false"`
	Message string            `json:"message,omitempty" example:"There were validation errors"`
	Errors  map[string]string `json:"errors,omitempty"`
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

func Fail(message string, errors map[string]string) FailResponse {
	return FailResponse{Success: false, Message: message, Errors: errors}
}
