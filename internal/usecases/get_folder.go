package usecases

import (
	"app/controllers/requests"
	"app/controllers/responses"
	"app/internal/repositories"
	"context"
)

type selectFolderUseCase struct {
	categoryRepository repositories.CategoryRepository
	folderRepository   repositories.FolderRepository
}

type SelectFolderUseCase interface {
	SelectFolder(ctx context.Context, accountId string, request requests.Folder) (responses.Folder, error)
}

func NewSelectFolderUseCase(categoryRepository repositories.CategoryRepository, folderRepository repositories.FolderRepository) SelectFolderUseCase {
	return selectFolderUseCase{
		categoryRepository: categoryRepository,
		folderRepository:   folderRepository,
	}
}

func (s selectFolderUseCase) SelectFolder(ctx context.Context, accountId string, request requests.Folder) (responses.Folder, error) {
	_, err := s.categoryRepository.SelectCategory(ctx, accountId, request.Category.Name)
	if err != nil {
		return responses.Folder{}, err
	}

	folder, err := s.folderRepository.SelectFolder(ctx, accountId, request.Category.Name, request.Name)
	if err != nil {
		return responses.Folder{}, err
	}

	return responses.NewFolder(folder.Name, folder.Files), nil
}
