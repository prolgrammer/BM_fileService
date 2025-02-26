package usecases

import (
	"app/controllers/requests"
	"app/internal/repositories"
	"context"
)

type deleteCategory struct {
	categoryRepository repositories.CategoryRepository
}

type DeleteCategoryUseCase interface {
	DeleteCategory(ctx context.Context, accountId string, req requests.Category) error
}

func NewDeleteCategoryUseCase(
	categoryRepository repositories.CategoryRepository,
) DeleteCategoryUseCase {

	return &deleteCategory{
		categoryRepository: categoryRepository,
	}
}

func (uc *deleteCategory) DeleteCategory(ctx context.Context, accountId string, req requests.Category) error {
	return uc.categoryRepository.DeleteCategory(ctx, accountId, req.Name)
}
