package product

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
		return ErrProductNotFound

	case errors.Is(err, repository.ErrAlreadyExists):
		return ErrConflict

	default:
		return err

	}
}

func mapServiceError(err error) apperror.ErrorCode {
	switch err {
	case ErrProductNotFound:
		return apperror.ErrProductNotFound
	case ErrCategoryNotFound:
		return apperror.ErrCategoryNotFound
	default:
		return apperror.ErrInternal
	}
}

func mapToProductOutput(
	product domain.Product,
) *ProductOutput {
	return &ProductOutput{
		ID:         product.ID,
		Name:       product.Name,
		Price:      product.Price,
		Stock:      product.Stock,
		CategoryID: product.CategoryID,

		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}

func mapToProductDetailOutput(
	p domain.Product,
	c domain.Category,
) *ProductDetailOutput {
	return &ProductDetailOutput{
		ProductOutput: *mapToProductOutput(p),
		Category: &CategoryOutput{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description,
		},
	}
}

func mapToProductResponse(
	productOut ProductOutput,
) *productResponse {
	return &productResponse{
		ID:         productOut.ID,
		Name:       productOut.Name,
		Price:      productOut.Price,
		Stock:      productOut.Stock,
		CategoryID: productOut.CategoryID,

		CreatedAt: productOut.CreatedAt,
		UpdatedAt: productOut.UpdatedAt,
	}
}

func mapToProductDetailResponse(
	productDetailOut ProductDetailOutput,
) *productDetailResponse {
	return &productDetailResponse{
		productResponse: *mapToProductResponse(productDetailOut.ProductOutput),
		Category: &categoryResponse{
			ID:          productDetailOut.Category.ID,
			Name:        productDetailOut.Category.Name,
			Description: productDetailOut.Category.Description,
		},
	}
}
