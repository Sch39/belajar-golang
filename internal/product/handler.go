// internal/handler/handler.go
package product

import (
	"encoding/json"
	"net/http"
	"strconv"

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

	input := UpsertProductInput{
		Name:       req.Name,
		Price:      *req.Price,
		Stock:      *req.Stock,
		CategoryID: req.CategoryID,
	}

	result, err := h.service.Add(r.Context(), input)
	if err != nil {
		code := mapServiceError(err)
		rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, err.Error(), nil))
		return
	}

	rest.JSON(w, http.StatusCreated, rest.Success(mapToProductResponse(*result), "Product created successfully"))
}

// GetAll godoc
// @Summary      Get all products
// @Description  Retrieve a list of all products with pagination
// @Tags         products
// @Produce      json
// @Param        query    query     string  false  "Search query by name"
// @Param        page     query     int     false  "Page number (default 1)"
// @Param        limit    query     int     false  "Items per page (default 10)"
// @Success      200      {object}  product.ProductsSuccessResponse
// @Failure      500      {object}  product.InternalServerErrorResponse "Internal server error"
// @Router       /api/products [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	query := r.URL.Query().Get("query")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	products, total, err := h.service.GetAll(r.Context(), query, page, limit)
	if err != nil {
		code := mapServiceError(err)
		rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, err.Error(), nil))
		return
	}

	resp := []productResponse{}
	for _, p := range products {
		resp = append(resp, *mapToProductResponse(p))
	}

	totalPages := (total + limit - 1) / limit
	pagination := rest.Pagination{
		Total:       total,
		TotalPage:   totalPages,
		CurrentPage: page,
		Limit:       limit,
	}

	rest.JSON(w, http.StatusOK, rest.SuccessWithPagination(resp, pagination))
}

// GetByID godoc
// @Summary      Get a product by ID
// @Description  Get detailed information about a specific product
// @Tags         products
// @Produce      json
// @Param        id   path      string  true  "Product ID" example(550e8400-e29b-41d4-a716-446655440000)
// @Success      200  {object}  product.ProductDetailSuccessResponse
// @Failure      404  {object}  product.ProductNotFoundResponse "Product Not Found"
// @Failure      500  {object}  product.InternalServerErrorResponse "Internal server error"
// @Router       /api/products/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request, id string) {
	product, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case ErrProductNotFound:
			code := apperror.ErrProductNotFound
			rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, "Product not found", nil))
			return
		default:
			code := mapServiceError(err)
			rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, err.Error(), nil))
			return
		}
	}
	rest.JSON(w, http.StatusOK, rest.Success(mapToProductDetailResponse(*product)))
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
	input := UpsertProductInput{
		Name:       req.Name,
		Price:      *req.Price,
		Stock:      *req.Stock,
		CategoryID: req.CategoryID,
	}
	result, err := h.service.Update(r.Context(), id, input)
	if err != nil {
		switch err {
		case ErrProductNotFound:
			code := apperror.ErrProductNotFound
			rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, "Product not found", nil))
			return
		default:
			code := mapServiceError(err)
			rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, err.Error(), nil))
			return
		}
	}
	rest.JSON(w, http.StatusOK, rest.Success(mapToProductResponse(*result), "Product updated successfully"))
}

// Delete godoc
// @Summary      Delete a product
// @Description  Remove a product from the database by its ID
// @Tags         products
// @Produce      json
// @Param        id   path      string  true  "Product ID" example(550e8400-e29b-41d4-a716-446655440000)
// @Success      204  "Product deleted"
// @Failure      404  {object}  product.ProductNotFoundResponse "Product Not Found"
// @Failure      500  {object}  product.InternalServerErrorResponse "Internal server error"
// @Router       /api/products/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.Delete(r.Context(), id); err != nil {
		switch err {
		case ErrProductNotFound:
			code := apperror.ErrProductNotFound
			rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, "Product not found", nil))
			return
		default:
			code := mapServiceError(err)
			rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, err.Error(), nil))
			return
		}
	}
	rest.JSON(w, http.StatusNoContent, rest.Success(nil))
}
