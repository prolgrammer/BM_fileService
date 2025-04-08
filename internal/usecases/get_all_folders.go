package usecases

import (
	"app/controllers/requests"
	"app/controllers/responses"
	"app/internal/repositories"
	"context"
	"errors"
	"fmt"
)

type selectFoldersUseCase struct {
	categoryRepository repositories.CategoryRepository
	folderRepository   repositories.FolderRepository
}

type SelectFoldersUseCase interface {
	SelectFolders(ctx context.Context, accountId string, request requests.Category) ([]responses.Folder, error)
}

func NewSelectFoldersUseCase(categoryRepository repositories.CategoryRepository, folderRepository repositories.FolderRepository) SelectFoldersUseCase {
	return &selectFoldersUseCase{
		categoryRepository: categoryRepository,
		folderRepository:   folderRepository,
	}
}

func (s *selectFoldersUseCase) SelectFolders(ctx context.Context, accountId string, request requests.Category) ([]responses.Folder, error) {
	_, err := s.categoryRepository.SelectCategory(ctx, accountId, request.Name)
	if err != nil {
		if errors.Is(err, repositories.ErrCategoryNotFound) {
			return nil, fmt.Errorf("category '%s' not found: %w", request.Name, err)
		}
		return nil, fmt.Errorf("failed to verify category: %w", err)
	}

	folders, err := s.folderRepository.SelectFolders(ctx, accountId, request.Name)
	if err != nil {
		return nil, err
	}

	foldersResponse := make([]responses.Folder, 0, len(folders))
	for _, folder := range folders {
		foldersResponse = append(foldersResponse, responses.NewFolder(folder.Name))
	}

	return foldersResponse, nil
}
