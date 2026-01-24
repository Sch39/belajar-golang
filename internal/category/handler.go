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