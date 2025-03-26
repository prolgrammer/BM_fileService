package usecases

import (
	"app/controllers/requests"
	"app/controllers/responses"
	"app/internal/repositories"
	"context"
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
		return nil, err
	}

	folders, err := s.folderRepository.SelectFolders(ctx, accountId, request.Name)
	if err != nil {
		return nil, err
	}

	foldersResponse := make([]responses.Folder, len(folders))
	for i, folder := range folders {
		foldersResponse[i] = responses.NewFolder(folder.Name, folder.Files)
	}

	return foldersResponse, nil
}
