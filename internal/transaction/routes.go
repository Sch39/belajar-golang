package transaction

import "net/http"

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/transactions/checkout", h.Checkout)
	mux.HandleFunc("GET /api/report/hari-ini", h.GetReportToday)
	mux.HandleFunc("GET /api/report", h.GetReport)
}
