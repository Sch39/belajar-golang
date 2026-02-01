// internal/product/routes.go
package product

import "net/http"

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/products", h.Create)
	mux.HandleFunc("GET /api/products", h.GetAll)
	mux.HandleFunc("GET /api/products/{id}", withID(h.GetByID))
	mux.HandleFunc("PUT /api/products/{id}", withID(h.Update))
	mux.HandleFunc("DELETE /api/products/{id}", withID(h.Delete))
}

func withID(
	handler func(http.ResponseWriter, *http.Request, string),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		handler(w, r, id)
	}
}
