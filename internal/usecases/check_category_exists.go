package usecases

import (
	"app/controllers/requests"
	"app/internal/repositories"
	"context"
)

type checkCategoryExistUseCase struct {
	categoryRepository repositories.CategoryRepository
}

type CheckCategoryExistUseCase interface {
	CheckCategoryExist(ctx context.Context, accountId string, request requests.Category) (bool, error)
}

func NewCheckCategoryExist(categoryRepository repositories.CategoryRepository) CheckCategoryExistUseCase {
	return &checkCategoryExistUseCase{
		categoryRepository: categoryRepository,
	}
}

func (c *checkCategoryExistUseCase) CheckCategoryExist(ctx context.Context, accountId string, request requests.Category) (bool, error) {
	return c.categoryRepository.CheckCategoryExists(ctx, accountId, request.Name)
}
