package usecases

import (
	"app/controllers/requests"
	"app/internal/repositories"
	"context"
	"fmt"
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

func (uc *createCategory) CreateCategory(ctx context.Context, accountId string, req requests.Category) error { //TODO добавить проверку на наличие аккаунта
	exists, err := uc.categoryRepository.CheckCategoryExists(ctx, accountId, req.Name)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	if exists {
		return ErrCategoryAlreadyExists
	}

	return uc.categoryRepository.CreateCategory(ctx, accountId, req.Name)
}
