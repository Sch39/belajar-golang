package category

import (
	"errors"

	"sch.dev/my-kasir-gw/internal/domain"
	"sch.dev/my-kasir-gw/internal/pkg/apperror"
	"sch.dev/my-kasir-gw/internal/storage/repository"
)

func mapRepoError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, repository.ErrNotFound):
		return ErrCategoryNotFound

	case errors.Is(err, repository.ErrAlreadyExists):
		return ErrConflict

	default:
		return err

	}
}

func mapServiceError(err error) apperror.ErrorCode {
	switch err {
	case ErrCategoryNotFound:
		return apperror.ErrCategoryNotFound
	default:
		return apperror.ErrInternal
	}
}

func mapToCategoryOutput(c domain.Category) *CategoryOutput {
	return &CategoryOutput{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func mapToCategoryResponse(category CategoryOutput) *categoryResponse {
	return &categoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}
