package usecases

import (
	"app/controllers/requests"
	"app/internal/repositories"
	"context"
)

type deleteFolderUseCase struct {
	categoryRepository repositories.CategoryRepository
	folderRepository   repositories.FolderRepository
}

type DeleteFolderUseCase interface {
	DeleteFolder(ctx context.Context, accountId string, request requests.Folder) error
}

func NewDeleteFolderUseCase(categoryRepository repositories.CategoryRepository, folderRepository repositories.FolderRepository) DeleteFolderUseCase {
	return deleteFolderUseCase{
		categoryRepository: categoryRepository,
		folderRepository:   folderRepository,
	}
}

func (uc deleteFolderUseCase) DeleteFolder(ctx context.Context, accountId string, request requests.Folder) error {
	_, err := uc.categoryRepository.SelectCategory(ctx, accountId, request.Category.Name)
	if err != nil {
		return err
	}

	return uc.folderRepository.DeleteFolder(ctx, accountId, request.Category.Name, request.Name)
}
