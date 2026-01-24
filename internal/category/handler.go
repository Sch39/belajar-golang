// internal/category/handler.go
package category

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

func NewHandler(s Service) *Handler{
	return &Handler{
		service: s,
	}
}

// Create godoc
// @Summary      Create a new category
// @Description  Save a new category to database
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        category  body      upsertRequest  true  "Category Request Body"
// @Success      201      {object}  category.CategorySuccessResponse
// @Failure      400      {object}  category.InvalidBodyResponse "Invalid JSON body"
// @Failure      422      {object}  category.ValidateErrorResponse "Validation error"
// @Failure 500 {object} category.InternalServerErrorResponse "Internal server error"
// @Router       /api/categories [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request){
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

	category := &domain.Category{
		Name: req.Name,
		Description: req.Description,
	}

	if err := h.service.Add(r.Context(), category); err != nil {
		code := apperror.ErrInternal
		rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, err.Error(), nil))
		return
	}

	rest.JSON(w, http.StatusCreated, rest.Success( toResponse(category), "Category created successfully"))
}

// GetAll godoc
// @Summary      Get all categories
// @Description  Retrieve a list of all categories
// @Tags         categories
// @Produce      json
// @Success      200      {object}  category.CategoryListSuccessResponse
// @Failure      500      {object}  category.InternalServerErrorResponse "Internal server error"
// @Router       /api/categories [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request){
	categories, err := h.service.GetAll(r.Context())
	if err != nil {
		code := apperror.ErrInternal
		rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, err.Error(), nil))
		return
	}

	resp := []categoryResponse{}

	for _, category := range categories {
		resp = append(resp, toResponse(&category))
	}

	rest.JSON(w, http.StatusOK, rest.Success( resp))
}

// GetByID godoc
// @Summary      Get a category by ID
// @Description  Get detailed information about a specific category
// @Tags         categories
// @Produce      json
// @Param        id   path      string  true  "Category ID" example(550e8400-e29b-41d4-a716-446655440000)
// @Success      200  {object}  category.CategorySuccessResponse
// @Failure      404  {object}  category.CategoryNotFoundResponse "Category Not Found"
// @Failure      500  {object}  category.InternalServerErrorResponse "Internal server error"
// @Router       /api/categories/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request, id string){
	category, err := h.service.GetByID(r.Context(), id)

	if err != nil {
		switch err {
		case ErrNotFound:
			code := apperror.ErrCategoryNotFound
			rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, "Category not found", nil))
			return
		default:
			code := apperror.ErrInternal
			rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, err.Error(), nil))
			return
		}
	}
		
	rest.JSON(w, http.StatusOK, rest.Success( toResponse(category)))
}

// Update godoc
// @Summary      Update an existing category
// @Description  Update category details by its ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id       path      string         true  "Category ID" example(550e8400-e29b-41d4-a716-446655440000)
// @Param        category  body      upsertRequest  true  "Updated Category Data"
// @Success      200      {object}  category.CategorySuccessResponse
// @Failure      400      {object}  category.InvalidBodyResponse "Invalid JSON body"
// @Failure      404      {object}  category.CategoryNotFoundResponse "Category Not Found"
// @Failure      422      {object}  category.ValidateErrorResponse "Validation error"
// @Failure 500 {object} category.InternalServerErrorResponse "Internal server error"
// @Router       /api/categories/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request, id string){
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

	category := &domain.Category{
		ID: id,
		Name: req.Name,
		Description: req.Description,
	}

	if err := h.service.Update(r.Context(), category); err != nil {
		switch err {
		case ErrNotFound:
			code := apperror.ErrCategoryNotFound
			rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, "Category not found", nil))
			return
		default:
			code := apperror.ErrInternal
			rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, err.Error(), nil))
			return
		}
	}
	rest.JSON(w, http.StatusOK, rest.Success(toResponse(category), "Category updated successfully"))
}

// Delete godoc
// @Summary      Delete a category
// @Description  Remove a category from the database by its ID
// @Tags         categories
// @Produce      json
// @Param        id   path      string  true  "Category ID" example(550e8400-e29b-41d4-a716-446655440000)
// @Success      200  {object}  category.BaseSuccessResponse "Success deleted"
// @Failure      404  {object}  category.CategoryNotFoundResponse "Category Not Found"
// @Failure      500  {object}  category.InternalServerErrorResponse "Internal server error"
// @Router       /api/categories/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, id string){
	if err := h.service.Delete(r.Context(), id); err != nil {
		switch err {
		case ErrNotFound:
			code := apperror.ErrCategoryNotFound
			rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, "Category not found", nil))
			return
		default:
			code := apperror.ErrInternal
			rest.JSON(w, code.ToHttpStatus(), rest.Fail(code, err.Error(), nil))
			return
		}
	}
	rest.JSON(w, http.StatusOK, rest.Success(nil, "Category deleted successfully"))
}

func toResponse(category *domain.Category) categoryResponse {
	return categoryResponse{
		ID: category.ID,
		Name: category.Name,
		Description: category.Description,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}