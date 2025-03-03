package usecases

import (
	"app/controllers/requests"
	"app/internal/entities"
	"app/internal/repositories"
	"context"
)

type selectFoldersUseCase struct {
	categoryRepository repositories.CategoryRepository
	folderRepository   repositories.FolderRepository
}

type SelectFoldersUseCase interface {
	SelectFolders(ctx context.Context, accountId string, request requests.Category) ([]entities.Folder, error)
}

func NewSelectFoldersUseCase(categoryRepository repositories.CategoryRepository, folderRepository repositories.FolderRepository) SelectFoldersUseCase {
	return selectFoldersUseCase{
		categoryRepository: categoryRepository,
		folderRepository:   folderRepository,
	}
}

func (s selectFoldersUseCase) SelectFolders(ctx context.Context, accountId string, request requests.Category) ([]entities.Folder, error) {
	_, err := s.categoryRepository.SelectCategory(ctx, accountId, request.Name)
	if err != nil {
		return nil, err
	}

	return s.folderRepository.SelectFolders(ctx, accountId, request.Name)
}
