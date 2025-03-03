package usecases

import (
	"app/controllers/requests"
	"app/internal/entities"
	"app/internal/repositories"
	"context"
)

type selectFolderUseCase struct {
	categoryRepository repositories.CategoryRepository
	folderRepository   repositories.FolderRepository
}

type SelectFolderUseCase interface {
	SelectFolder(ctx context.Context, accountId string, request requests.Folder) (entities.Folder, error)
}

func NewSelectFolderUseCase(categoryRepository repositories.CategoryRepository, folderRepository repositories.FolderRepository) SelectFolderUseCase {
	return selectFolderUseCase{
		categoryRepository: categoryRepository,
		folderRepository:   folderRepository,
	}
}

func (s selectFolderUseCase) SelectFolder(ctx context.Context, accountId string, request requests.Folder) (entities.Folder, error) {
	_, err := s.categoryRepository.SelectCategory(ctx, accountId, request.Category.Name)
	if err != nil {
		return entities.Folder{}, err
	}

	return s.folderRepository.SelectFolder(ctx, accountId, request.Category.Name, request.Name)
}
