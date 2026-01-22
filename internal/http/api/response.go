// internal\http\api\response.go
package api

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool              `json:"success"`
	Message string            `json:"message,omitempty"`
	Data    interface{}       `json:"data,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func JSON(w http.ResponseWriter, status int, resp Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

func Success(data interface{}) Response {
	return Response{Success: true, Data: data}
}

func Fail(message string, errors map[string]string) Response {
	return Response{Success: false, Message: message, Errors: errors}
}