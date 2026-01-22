// internal/handler/product_handler.go
package product

import (
	"encoding/json"
	"net/http"

	"sch.dev/my-kasir-gw/internal/http/api"
	"sch.dev/my-kasir-gw/internal/validator"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req upsertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.JSON(w, http.StatusBadRequest, api.Fail("Invalid body", nil))
		return
	}

	if err:= validator.Instance().Struct(req); err != nil {
		api.JSON(w, http.StatusBadRequest, api.Fail("Validation error", api.ValidationError(err)))
		return
	}

	product := &Product{
		Name:  req.Name,
		Price: *req.Price,
		Stock: *req.Stock,
	}

	if err := h.service.Add(r.Context(), product); err != nil {
		api.JSON(w, http.StatusInternalServerError, api.Fail(err.Error(), nil))
		return
	}

	api.JSON(w, http.StatusCreated, api.Success(toResponse(product)))
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAll(r.Context())
	if err != nil {
		api.JSON(w, http.StatusInternalServerError, api.Fail(err.Error(), nil))
		return
	}
	var resp []productResponse
	for _, p := range products {
		resp = append(resp, toResponse(&p))
	}
	api.JSON(w, http.StatusOK, api.Success(resp))
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request, id string) {
	product, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		switch err { 
		case ErrNotFound:
			api.JSON(w, http.StatusNotFound, api.Fail("Product not found", nil))
		default:
			api.JSON(w, http.StatusInternalServerError, api.Fail(err.Error(), nil))
		}
		return
	}
	api.JSON(w, http.StatusOK, api.Success(toResponse(product)))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request, id string) {
	var req upsertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.JSON(w, http.StatusBadRequest, api.Fail("Invalid body", nil))
		return
	}
	if err:= validator.Instance().Struct(req); err != nil {
		api.JSON(w, http.StatusBadRequest, api.Fail("Validation error", api.ValidationError(err)))
		return
	}
	product := &Product{
		ID:    id,
		Name:  req.Name,
		Price: *req.Price,
		Stock: *req.Stock,
	}
	if err := h.service.Update(r.Context(), product); err != nil {
		switch err {
		case ErrNotFound:
			api.JSON(w, http.StatusNotFound, api.Fail("Product not found", nil))
		default:
			api.JSON(w, http.StatusInternalServerError, api.Fail(err.Error(), nil))
		}
		return
	}
	api.JSON(w, http.StatusOK, api.Success(toResponse(product)))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.Delete(r.Context(), id); err != nil {
		switch err {
		case ErrNotFound:
			api.JSON(w, http.StatusNotFound, api.Fail("Product not found", nil))
		default:
			api.JSON(w, http.StatusInternalServerError, api.Fail(err.Error(), nil))
		}
		return
	}
	api.JSON(w, http.StatusOK, api.Success(nil))
}


func toResponse(p *Product) productResponse {
	return productResponse{
		ID:    p.ID,
		Name:  p.Name,
		Price: p.Price,
		Stock: p.Stock,
	}
}