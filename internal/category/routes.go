// internal\category\routes.go
package category

import "net/http"

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/categories", h.Create)
	mux.HandleFunc("GET /api/categories", h.GetAll)
	mux.HandleFunc("GET /api/categories/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		h.GetByID(w, r, id)
	})
	mux.HandleFunc("PUT /api/categories/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		h.Update(w, r, id)
	})
	mux.HandleFunc("DELETE /api/categories/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		h.Delete(w, r, id)
	})
}
