// internal/http/api/validation.go
package api

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

func ValidationError(err error) map[string]string {
	out := make(map[string]string)

	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	for _, e := range ve {
		out[e.Field()] = translateError(e)
	}

	return out
}


func translateError(e validator.FieldError) string {
	messages := map[string]string{
		"required": "bidang ini wajib diisi",
		"email":    "format email tidak valid",
		"numeric":  "harus berupa angka",
	}

	if msg, ok := messages[e.Tag()]; ok {
		return msg
	}

	switch e.Tag() {
	case "min":
		if e.Kind() == reflect.String {
			return fmt.Sprintf("minimal harus %s karakter", e.Param())
		}
		return fmt.Sprintf("nilai minimal adalah %s", e.Param())
	case "max":
		if e.Kind() == reflect.String {
			return fmt.Sprintf("maksimal harus %s karakter", e.Param())
		}
		return fmt.Sprintf("nilai maksimal adalah %s", e.Param())
	}

	return fmt.Sprintf("format tidak valid pada tag: %s", e.Tag())
}