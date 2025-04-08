package usecases

import (
	"app/controllers/requests"
	m "app/infrastructure/minio"
	"app/internal/repositories"
	"context"
	"fmt"
)

type deleteFolderUseCase struct {
	minio              *m.Client
	categoryRepository repositories.CategoryRepository
	folderRepository   repositories.FolderRepository
	fileRepository     repositories.FileRepository
}

type DeleteFolderUseCase interface {
	DeleteFolder(ctx context.Context, accountId string, request requests.Folder) error
}

func NewDeleteFolderUseCase(
	minio *m.Client,
	categoryRepository repositories.CategoryRepository,
	folderRepository repositories.FolderRepository,
	fileRepository repositories.FileRepository) DeleteFolderUseCase {
	return &deleteFolderUseCase{
		minio:              minio,
		categoryRepository: categoryRepository,
		folderRepository:   folderRepository,
		fileRepository:     fileRepository,
	}
}

func (uc *deleteFolderUseCase) DeleteFolder(ctx context.Context, accountId string, req requests.Folder) error {
	category, err := uc.categoryRepository.SelectCategory(ctx, accountId, req.Category.Name)
	if err != nil {
		return fmt.Errorf("failed to get category %s: %w", req.Category.Name, err)
	}

	files, err := uc.fileRepository.SelectFiles(ctx, category.Id, req.Name)
	if err != nil {
		return fmt.Errorf("failed to get folder files: %w", err)
	}

	for _, file := range files {
		err := cleanupFile(ctx, uc.minio.MinioClient, uc.minio.BucketName, uc.fileRepository,
			accountId, category.Id, req.Name,
			&file)

		if err != nil {
			return err
		}
	}

	err = uc.folderRepository.DeleteFolder(ctx, accountId, req.Category.Name, req.Name)
	if err != nil {
		return fmt.Errorf("failed to delete folder from db: %w", err)
	}

	return nil
}
