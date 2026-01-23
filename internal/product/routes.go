// internal/product/routes.go
package product

import "net/http"

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/products", h.handles)
	mux.HandleFunc("/api/products/", h.handleByID)
}

func (h *Handler) handles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Create(w, r)
	case http.MethodGet:
		h.GetAll(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/products/"):]
	if id == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetByID(w, r, id)
	case http.MethodPut:
		h.Update(w, r, id)
	case http.MethodDelete:
		h.Delete(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
