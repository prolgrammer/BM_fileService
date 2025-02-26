package usecases

import (
	"app/controllers/requests"
	"app/internal/entities"
	"app/internal/repositories"
	"context"
	"errors"
)

type getCategory struct {
	categoryRepository repositories.CategoryRepository
}

type GetCategoryUseCase interface {
	GetCategory(ctx context.Context, accountId string, req requests.Category) (entities.Category, error)
}

func NewGetCategory(categoryRepository repositories.CategoryRepository) GetCategoryUseCase {
	return &getCategory{
		categoryRepository: categoryRepository,
	}
}

func (uc *getCategory) GetCategory(ctx context.Context, accountId string, req requests.Category) (entities.Category, error) {
	category, err := uc.categoryRepository.SelectCategory(ctx, accountId, req.Name)
	if err != nil {
		if errors.Is(err, repositories.ErrCategoryNotFound) {
			return entities.Category{}, ErrCategoryNotFound
		}

		return entities.Category{}, err
	}
	return category, nil
}
