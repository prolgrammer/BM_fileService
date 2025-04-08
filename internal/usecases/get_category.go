package usecases

import (
	"app/controllers/requests"
	"app/controllers/responses"
	"app/internal/repositories"
	"context"
	"errors"
)

type getCategory struct {
	categoryRepository repositories.CategoryRepository
}

type GetCategoryUseCase interface {
	GetCategory(ctx context.Context, accountId string, req requests.Category) (responses.Category, error)
}

func NewGetCategory(categoryRepository repositories.CategoryRepository) GetCategoryUseCase {
	return &getCategory{
		categoryRepository: categoryRepository,
	}
}

func (uc *getCategory) GetCategory(ctx context.Context, accountId string, req requests.Category) (responses.Category, error) {
	category, err := uc.categoryRepository.SelectCategory(ctx, accountId, req.Name)
	if err != nil {
		if errors.Is(err, repositories.ErrCategoryNotFound) {
			return responses.Category{}, ErrCategoryNotFound
		}

		return responses.Category{}, err
	}

	return responses.NewCategory(category.Name, accountId, category.Folders), nil
}
