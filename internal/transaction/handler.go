package transaction

import (
	"encoding/json"
	"net/http"

	"sch.dev/my-kasir-gw/internal/pkg/apperror"
	"sch.dev/my-kasir-gw/internal/pkg/rest"
	"sch.dev/my-kasir-gw/internal/pkg/validator"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Checkout(w http.ResponseWriter, r *http.Request) {
	var req checkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		code := apperror.ErrInvalidPayload
		rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, "Invalid body", nil))
		return
	}

	if err := validator.Instance().Struct(req); err != nil {
		code := apperror.ErrValidation
		rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, "Validation error", rest.ValidationError(err)))
		return
	}

	// Map HTTP Request to Service DTO
	input := CheckoutInput{
		Items: make([]CheckoutItemInput, len(req.Items)),
	}
	for i, item := range req.Items {
		input.Items[i] = CheckoutItemInput{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	output, err := h.service.Checkout(r.Context(), input)
	if err != nil {
		code := mapServiceError(err)
		rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, err.Error(), nil))
		return
	}

	// Map Service Output to HTTP Response
	res := mapToCheckoutResponse(output)

	rest.JSON(w, http.StatusCreated, rest.Success(res))
}
