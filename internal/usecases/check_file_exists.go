package usecases

import (
	"app/controllers/requests"
	"app/internal/repositories"
	"context"
)

type checkFileExistsUseCase struct {
	categoryRepository repositories.CategoryRepository
	folderRepository   repositories.FolderRepository
	fileRepository     repositories.FileRepository
}

type CheckFileExistsUseCase interface {
	CheckFileExists(ctx context.Context, accountId string, request requests.File) (bool, error)
}

func NewCheckFileExistsUseCase(
	categoryRepository repositories.CategoryRepository,
	folderRepository repositories.FolderRepository,
	fileRepository repositories.FileRepository) CheckFileExistsUseCase {
	return &checkFileExistsUseCase{
		categoryRepository: categoryRepository,
		folderRepository:   folderRepository,
		fileRepository:     fileRepository,
	}
}

func (c *checkFileExistsUseCase) CheckFileExists(ctx context.Context, accountId string, request requests.File) (bool, error) {
	_, err := c.categoryRepository.SelectCategory(ctx, accountId, request.Category.Name)
	if err != nil {
		return false, err
	}

	folderExists, err := c.folderRepository.CheckFolderExists(ctx, accountId, request.Category.Name, request.Folder.Name)
	if err != nil || !folderExists {
		return false, err
	}

	return c.fileRepository.CheckFileExists(ctx, accountId, request.Category.Name, request.Folder.Name, request.Name)

}
