package usecases

import (
	"app/controllers/requests"
	"app/internal/repositories"
	"context"
)

type createFolderUseCase struct {
	categoryRepository repositories.CategoryRepository
	folderRepository   repositories.FolderRepository
}

type CreateFolderUseCase interface {
	CreateFolder(ctx context.Context, accountId string, request requests.Folder) error
}

func NewCreateFolderUseCase(categoryRepository repositories.CategoryRepository, folderRepository repositories.FolderRepository) CreateFolderUseCase {
	return &createFolderUseCase{
		categoryRepository: categoryRepository,
		folderRepository:   folderRepository}
}

func (c *createFolderUseCase) CreateFolder(ctx context.Context, accountId string, request requests.Folder) error {
	category, err := c.categoryRepository.SelectCategory(ctx, accountId, request.Category.Name)
	if err != nil {
		return err
	}

	for _, folder := range category.Folders {
		if folder.Name == request.Name {
			return ErrFolderAlreadyExists
		}
	}

	return c.folderRepository.CreateFolder(ctx, accountId, request.Category.Name, request.Name)
}
