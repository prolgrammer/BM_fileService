package usecases

import (
	"app/internal/entities"
	"app/internal/repositories"
	"context"
)

type getAllCategory struct {
	categoryRepository repositories.CategoryRepository
}

type GetAllCategoryUseCase interface {
	GetAllCategory(ctx context.Context, accountId string) ([]entities.Category, error)
}

func NewGetAllCategory(categoryRepository repositories.CategoryRepository) GetAllCategoryUseCase {
	return &getAllCategory{
		categoryRepository: categoryRepository,
	}
}

func (uc *getAllCategory) GetAllCategory(ctx context.Context, accountId string) ([]entities.Category, error) {
	return uc.categoryRepository.SelectAllCategories(ctx, accountId)
}
