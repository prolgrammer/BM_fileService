package usecases

import (
	"app/controllers/requests"
	"app/internal/repositories"
	"context"
)

type checkFolderExistUseCase struct {
	categoryRepository repositories.CategoryRepository
	folderRepository   repositories.FolderRepository
}

type CheckFolderExistUseCase interface {
	CheckFolderExist(ctx context.Context, accountId string, request requests.Folder) (bool, error)
}

func NewCheckFolderExistUseCase(categoryRepository repositories.CategoryRepository, folderRepository repositories.FolderRepository) CheckFolderExistUseCase {
	return &checkFolderExistUseCase{
		categoryRepository: categoryRepository,
		folderRepository:   folderRepository,
	}
}

func (c *checkFolderExistUseCase) CheckFolderExist(ctx context.Context, accountId string, request requests.Folder) (bool, error) {
	_, err := c.categoryRepository.SelectCategory(ctx, accountId, request.Category.Name)
	if err != nil {
		return false, err
	}

	return c.folderRepository.CheckFolderExists(ctx, accountId, request.Category.Name, request.Name)
}
