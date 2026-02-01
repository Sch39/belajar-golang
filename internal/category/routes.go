// internal\category\routes.go
package category

import "net/http"

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/categories", h.Create)
	mux.HandleFunc("GET /api/categories", h.GetAll)
	mux.HandleFunc("GET /api/categories/{id}", withID(h.GetByID))
	mux.HandleFunc("PUT /api/categories/{id}", withID(h.Update))
	mux.HandleFunc("DELETE /api/categories/{id}", withID(h.Delete))
}

func withID(
	handler func(http.ResponseWriter, *http.Request, string),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		handler(w, r, id)
	}
}
