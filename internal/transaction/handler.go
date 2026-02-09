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

// Checkout godoc
// @Summary      Create a new transaction
// @Description  Create a new transaction (checkout)
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        request  body      checkoutRequest  true  "Checkout Request Body"
// @Success      201      {object}  transaction.CheckoutSuccessResponse
// @Failure      400      {object}  transaction.InvalidBodyResponse "Invalid JSON body"
// @Failure      404      {object}  transaction.ProductNotFoundResponse "Product not found"
// @Failure      422      {object}  transaction.ValidationErrorResponse "Validation error"
// @Failure      422      {object}  transaction.InsufficientStockResponse "Insufficient stock"
// @Failure      500      {object}  transaction.InternalServerErrorResponse "Internal server error"
// @Router       /api/transactions/checkout [post]
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

// GetReportToday godoc
// @Summary      Get transaction report for today
// @Description  Get transaction report (revenue, total transactions, best selling product) for today
// @Tags         report
// @Produce      json
// @Success      200      {object}  transaction.ReportSuccessResponse
// @Failure      500      {object}  transaction.InternalServerErrorResponse "Internal server error"
// @Router       /api/report/hari-ini [get]
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

// GetReport godoc
// @Summary      Get transaction report by date range
// @Description  Get transaction report (revenue, total transactions, best selling product) for a specific date range
// @Tags         report
// @Produce      json
// @Param        start_date  query     string  false  "Start Date (YYYY-MM-DD)"
// @Param        end_date    query     string  false  "End Date (YYYY-MM-DD)"
// @Success      200         {object}  transaction.ReportSuccessResponse
// @Failure      400         {object}  transaction.InvalidBodyResponse "Invalid date format"
// @Failure      500         {object}  transaction.InternalServerErrorResponse "Internal server error"
// @Router       /api/report [get]
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
