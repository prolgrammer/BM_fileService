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
	category, err := c.categoryRepository.SelectCategory(ctx, accountId, request.Category.Name)
	if err != nil {
		return false, err
	}

	folderExist := false
	for _, folder := range category.Folders {
		if folder.Name == request.Folder.Name {
			folderExist = true
			break
		}
	}

	if !folderExist {
		return false, repositories.ErrFolderNotFound
	}

	return c.fileRepository.CheckFileExists(ctx, category.Id, request.Folder.Name, request.Name)
}
