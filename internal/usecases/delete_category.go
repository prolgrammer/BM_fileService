package usecases

import (
	"app/controllers/requests"
	m "app/infrastructure/minio"
	"app/internal/repositories"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
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

func (uc *deleteCategory) DeleteCategory(ctx context.Context, accountId string, req requests.Category) error { //TODO
	folders, err := uc.folderRepository.SelectFolders(ctx, accountId, req.Name)
	if err != nil {
		return fmt.Errorf("failed to get category folders: %w", err)
	}

	var filePaths []string
	for _, folder := range folders {
		files := folder.Files
		for _, file := range files {
			filePaths = append(filePaths, file.Path)
		}
	}

	if len(filePaths) > 0 {
		objectCh := make(chan minio.ObjectInfo)

		go func() {
			for _, filePath := range filePaths {
				objectCh <- minio.ObjectInfo{Key: filePath}
			}
		}()

		for err := range uc.minio.MinioClient.RemoveObjects(ctx, uc.minio.BucketName, objectCh, minio.RemoveObjectsOptions{}) {
			return err.Err
		}
		return nil
	}

	if err := uc.categoryRepository.DeleteCategory(ctx, accountId, req.Name); err != nil {
		return fmt.Errorf("failed to delete files from category db: %w", err)
	}

	return nil
}
