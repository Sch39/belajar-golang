package transaction

import (
	"encoding/json"
	"net/http"
	"time"

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

func (h *Handler) GetReportToday(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

	output, err := h.service.GetReport(r.Context(), startDate, endDate)
	if err != nil {
		code := mapServiceError(err)
		rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, err.Error(), nil))
		return
	}

	rest.JSON(w, http.StatusOK, rest.Success(output))
}

func (h *Handler) GetReport(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" && endDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			code := apperror.ErrInvalidPayload
			rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, "Invalid start_date format (YYYY-MM-DD)", nil))
			return
		}

		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			code := apperror.ErrInvalidPayload
			rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, "Invalid end_date format (YYYY-MM-DD)", nil))
			return
		}

		// Set endDate to end of day
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())
	} else {
		// Default to today if not provided? Or return error?
		// User requirement says: Get api/report?start_date=...&end_date=...
		// It implies params are expected. But let's fallback to today or error.
		// Let's fallback to today for convenience.
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endDate = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	}

	output, err := h.service.GetReport(r.Context(), startDate, endDate)
	if err != nil {
		code := mapServiceError(err)
		rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, err.Error(), nil))
		return
	}

	rest.JSON(w, http.StatusOK, rest.Success(output))
}
