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

// Create godoc
// @Summary      Create a new product
// @Description  Save a new product to database
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      upsertRequest  true  "Product Request Body"
// @Success      201      {object}  api.SuccessResponse{data=productResponse}
// @Failure      400      {object}  api.FailResponse{errors=map[string]string}
// @Router       /api/products [post]
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

// GetAll godoc
// @Summary      Get all products
// @Description  Retrieve a list of all products
// @Tags         products
// @Produce      json
// @Success      200      {object}  api.SuccessResponse{data=[]productResponse}
// @Failure      500      {object}  api.FailResponse
// @Router       /api/products [get]
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

// GetByID godoc
// @Summary      Get a product by ID
// @Description  Get detailed information about a specific product
// @Tags         products
// @Produce      json
// @Param        id   path      string  true  "Product ID" example(550e8400-e29b-41d4-a716-446655440000)
// @Success      200  {object}  api.SuccessResponse{data=productResponse}
// @Failure      404  {object}  api.FailResponse "Product Not Found"
// @Router       /api/products/{id} [get]
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

// Update godoc
// @Summary      Update an existing product
// @Description  Update product details by its ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id       path      string         true  "Product ID" example(550e8400-e29b-41d4-a716-446655440000)
// @Param        product  body      upsertRequest  true  "Updated Product Data"
// @Success      200      {object}  api.SuccessResponse{data=productResponse}
// @Failure      400      {object}  api.FailResponse{errors=map[string]string} "Validation Error"
// @Failure      404      {object}  api.FailResponse "Product Not Found"
// @Router       /api/products/{id} [put]
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

// Delete godoc
// @Summary      Delete a product
// @Description  Remove a product from the database by its ID
// @Tags         products
// @Produce      json
// @Param        id   path      string  true  "Product ID" example(550e8400-e29b-41d4-a716-446655440000)
// @Success      200  {object}  api.SuccessResponse "Success deleted"
// @Failure      404  {object}  api.FailResponse "Product Not Found"
// @Router       /api/products/{id} [delete]
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