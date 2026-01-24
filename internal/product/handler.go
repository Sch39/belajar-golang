// internal/handler/handler.go
package product

import (
	"encoding/json"
	"net/http"

	"sch.dev/my-kasir-gw/internal/domain"
	"sch.dev/my-kasir-gw/internal/pkg/apperror"
	"sch.dev/my-kasir-gw/internal/pkg/rest"
	"sch.dev/my-kasir-gw/internal/pkg/validator"
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
// @Success      201      {object}  product.ProductSuccessResponse
// @Failure      400      {object}  product.InvalidBodyResponse "Invalid JSON body"
// @Failure      422      {object}  product.ValidationErrorResponse "Validation error"
// @Failure 500 {object} product.InternalServerErrorResponse "Internal server error"
// @Router       /api/products [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req upsertRequest
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

	product := &domain.Product{
		Name:  req.Name,
		Price: *req.Price,
		Stock: *req.Stock,
	}

	if err := h.service.Add(r.Context(), product); err != nil {
		code := apperror.ErrInternal
		rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, err.Error(), nil))
		return
	}

	rest.JSON(w, http.StatusCreated, rest.Success(toResponse(product), "Product created successfully"))
}

// GetAll godoc
// @Summary      Get all products
// @Description  Retrieve a list of all products
// @Tags         products
// @Produce      json
// @Success      200      {object}  product.ProductsSuccessResponse
// @Failure      500      {object}  product.InternalServerErrorResponse "Internal server error"
// @Router       /api/products [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAll(r.Context())
	if err != nil {
		code := apperror.ErrInternal
		rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, err.Error(), nil))
		return
	}
	resp := []productResponse{}
	for _, p := range products {
		resp = append(resp, toResponse(&p))
	}
	rest.JSON(w, http.StatusOK, rest.Success(resp))
}

// GetByID godoc
// @Summary      Get a product by ID
// @Description  Get detailed information about a specific product
// @Tags         products
// @Produce      json
// @Param        id   path      string  true  "Product ID" example(550e8400-e29b-41d4-a716-446655440000)
// @Success      200  {object}  product.ProductSuccessResponse
// @Failure      404  {object}  product.ProductNotFoundResponse "Product Not Found"
// @Failure      500  {object}  product.InternalServerErrorResponse "Internal server error"
// @Router       /api/products/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request, id string) {
	product, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case ErrNotFound:
			code := apperror.ErrProductNotFound
			rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, "Product not found", nil))
			return
		default:
			code := apperror.ErrInternal
			rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, err.Error(), nil))
			return
		}
	}
	rest.JSON(w, http.StatusOK, rest.Success(toResponse(product)))
}

// Update godoc
// @Summary      Update an existing product
// @Description  Update product details by its ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id       path      string         true  "Product ID" example(550e8400-e29b-41d4-a716-446655440000)
// @Param        product  body      upsertRequest  true  "Updated Product Data"
// @Success      200      {object}  product.ProductSuccessResponse
// @Failure      400      {object}  product.InvalidBodyResponse "Invalid JSON body"
// @Failure      404      {object}  product.ProductNotFoundResponse "Product Not Found"
// @Failure      422      {object}  product.ValidationErrorResponse "Validation error"
// @Failure 500 {object} product.InternalServerErrorResponse "Internal server error"
// @Router       /api/products/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request, id string) {
	var req upsertRequest
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
	product := &domain.Product{
		ID:    id,
		Name:  req.Name,
		Price: *req.Price,
		Stock: *req.Stock,
	}
	if err := h.service.Update(r.Context(), product); err != nil {
		switch err {
		case ErrNotFound:
			code := apperror.ErrProductNotFound
			rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, "Product not found", nil))
			return
		default:
			code := apperror.ErrInternal
			rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, err.Error(), nil))
			return
		}
	}
	rest.JSON(w, http.StatusOK, rest.Success(toResponse(product), "Product updated successfully"))
}

// Delete godoc
// @Summary      Delete a product
// @Description  Remove a product from the database by its ID
// @Tags         products
// @Produce      json
// @Param        id   path      string  true  "Product ID" example(550e8400-e29b-41d4-a716-446655440000)
// @Success      200  {object}  product.BaseSuccessResponse "Success deleted"
// @Failure      404  {object}  product.ProductNotFoundResponse "Product Not Found"
// @Failure      500  {object}  product.InternalServerErrorResponse "Internal server error"
// @Router       /api/products/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.Delete(r.Context(), id); err != nil {
		switch err {
		case ErrNotFound:
			code := apperror.ErrProductNotFound
			rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, "Product not found", nil))
			return
		default:
			code := apperror.ErrInternal
			rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, err.Error(), nil))
			return
		}
	}
	rest.JSON(w, http.StatusOK, rest.Success(nil, "Product deleted successfully"))
}

func toResponse(p *domain.Product) productResponse {
	return productResponse{
		ID:    p.ID,
		Name:  p.Name,
		Price: p.Price,
		Stock: p.Stock,
	}
}
