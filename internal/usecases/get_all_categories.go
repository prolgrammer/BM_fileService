package usecases

import (
	"app/controllers/responses"
	"app/internal/repositories"
	"context"
)

type getAllCategory struct {
	categoryRepository repositories.CategoryRepository
}

type GetAllCategoryUseCase interface {
	GetAllCategory(ctx context.Context, accountId string) ([]responses.Category, error)
}

func NewGetAllCategory(categoryRepository repositories.CategoryRepository) GetAllCategoryUseCase {
	return &getAllCategory{
		categoryRepository: categoryRepository,
	}
}

func (uc *getAllCategory) GetAllCategory(ctx context.Context, accountId string) ([]responses.Category, error) {
	categories, err := uc.categoryRepository.SelectAllCategories(ctx, accountId)
	if err != nil {
		return nil, err
	}

	respCategories := make([]responses.Category, len(categories))
	for i, category := range categories {
		respCategories[i] = responses.NewCategory(category.Name, accountId, category.Folders)
	}

	return respCategories, nil
}
