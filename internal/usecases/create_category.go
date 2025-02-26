package usecases

import (
	"app/controllers/requests"
	"app/internal/repositories"
	"context"
)

type createCategory struct {
	categoryRepository repositories.CategoryRepository
}

type CreateCategoryUseCase interface {
	CreateCategory(ctx context.Context, accountId string, req requests.Category) error
}

func NewCreateCategoryUseCase(
	categoryRepository repositories.CategoryRepository) CreateCategoryUseCase {
	return &createCategory{
		categoryRepository: categoryRepository,
	}
}
func (uc *createCategory) CreateCategory(ctx context.Context, accountId string, req requests.Category) error {
	_, err := uc.categoryRepository.SelectCategory(ctx, accountId, req.Name)
	if err != nil {
		return ErrCategoryAlreadyExists
	}

	return uc.categoryRepository.CreateCategory(ctx, accountId, req.Name)
}
