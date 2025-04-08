package usecases

import (
	"app/controllers/requests"
	m "app/infrastructure/minio"
	"app/internal/repositories"
	"context"
	"fmt"
)

type deleteCategory struct {
	minio              *m.Client
	categoryRepository repositories.CategoryRepository
	folderRepository   repositories.FolderRepository
	fileRepository     repositories.FileRepository
}

type DeleteCategoryUseCase interface {
	DeleteCategory(ctx context.Context, accountId string, req requests.Category) error
}

func NewDeleteCategoryUseCase(
	minio *m.Client,
	categoryRepository repositories.CategoryRepository,
	folderRepository repositories.FolderRepository,
	fileRepository repositories.FileRepository,
) DeleteCategoryUseCase {

	return &deleteCategory{
		minio:              minio,
		categoryRepository: categoryRepository,
		folderRepository:   folderRepository,
		fileRepository:     fileRepository,
	}
}

func (uc *deleteCategory) DeleteCategory(ctx context.Context, accountId string, req requests.Category) error {
	category, err := uc.categoryRepository.SelectCategory(ctx, accountId, req.Name)
	if err != nil {
		return fmt.Errorf("failed to get category: %s: %w", req.Name, err)
	}

	folders, err := uc.folderRepository.SelectFolders(ctx, accountId, req.Name)
	if err != nil {
		return fmt.Errorf("failed to get category folders: %w", err)
	}

	for _, folder := range folders {
		files, err := uc.fileRepository.SelectFiles(ctx, category.Id, folder.Name)
		if err != nil {
			return fmt.Errorf("failed to get category files: %w", err)
		}

		for _, file := range files {
			err = cleanupFile(ctx, uc.minio.MinioClient, uc.minio.BucketName, uc.fileRepository,
				accountId, category.Id, folder.Name,
				&file)

			if err != nil {
				return err
			}
		}
	}

	if err := uc.categoryRepository.DeleteCategory(ctx, accountId, req.Name); err != nil {
		return fmt.Errorf("failed to delete files from category db: %w", err)
	}

	fmt.Printf("Category %s deleted for account %s\n", req.Name, accountId)
	return nil
}
