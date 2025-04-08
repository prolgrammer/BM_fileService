package usecases

import (
	"app/controllers/requests"
	m "app/infrastructure/minio"
	"app/internal/entities"
	"app/internal/repositories"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
)

type DeleteFileUseCase interface {
	DeleteFile(ctx context.Context, accountId string, req requests.File) error
}

type deleteFileUseCase struct {
	minio              *m.Client
	categoryRepository repositories.CategoryRepository
	fileRepository     repositories.FileRepository
}

func NewDeleteFileUseCase(
	minio *m.Client,
	categoryRepository repositories.CategoryRepository,
	fileRepository repositories.FileRepository,
) DeleteFileUseCase {
	return &deleteFileUseCase{
		minio:              minio,
		categoryRepository: categoryRepository,
		fileRepository:     fileRepository,
	}
}

func (d *deleteFileUseCase) DeleteFile(ctx context.Context, accountId string, req requests.File) error {
	category, err := d.categoryRepository.SelectCategory(ctx, accountId, req.Category.Name)
	if err != nil {
		return fmt.Errorf("failed to get category %s: %w", req.Category.Name, err)
	}

	file, err := d.fileRepository.SelectFile(ctx, category.Id, req.Folder.Name, req.Name)
	if err != nil {
		return fmt.Errorf("failed to get file: %w", err)
	}

	folderRemoved := false
	for _, fileCategory := range file.Categories {
		if fileCategory.CategoryId == category.Id {
			for _, folder := range fileCategory.Folders {
				if folder.Name == req.Folder.Name {
					folderRemoved = true
					break
				}
			}
			break
		}
	}

	if !folderRemoved {
		return fmt.Errorf("folder %s not found in category %s for file %s", req.Folder.Name, req.Category.Name, req.Name)
	}

	err = cleanupFile(ctx, d.minio.MinioClient, d.minio.BucketName, d.fileRepository,
		accountId, category.Id, req.Folder.Name,
		&file)

	return nil
}

func cleanupFile(ctx context.Context,
	minioClient *minio.Client,
	bucketName string,
	fileRepository repositories.FileRepository,
	accountId, categoryId, folderName string,
	file *entities.File) error {
	updatedCategories := make([]entities.FileCategory, 0, len(file.Categories))
	for _, fileCategory := range file.Categories {
		if fileCategory.CategoryId == categoryId {
			updatedFolders := make([]entities.Folder, 0, len(fileCategory.Folders))
			for _, folder := range fileCategory.Folders {
				if folder.Name == folderName {
					updatedFolders = append(updatedFolders, folder)
				}
			}
			if len(updatedFolders) > 0 {
				fileCategory.Folders = updatedFolders
				updatedCategories = append(updatedCategories, fileCategory)
			}
			continue
		}
		updatedCategories = append(updatedCategories, fileCategory)
	}

	file.Categories = updatedCategories
	err := fileRepository.UpdateFile(ctx, *file)
	if err != nil {
		return fmt.Errorf("failed to update file %s in database: %w", file.Name, err)
	}

	if len(file.Categories) == 0 {
		path := fmt.Sprintf("%s/%s/%s", accountId, file.Version, file.Name)

		err := minioClient.RemoveObject(ctx, bucketName, path, minio.RemoveObjectOptions{})
		if err != nil {
			return fmt.Errorf("failed to delete file %s from MinIO: %w", file.Name, err)
		}

		err = fileRepository.DeleteFile(ctx, categoryId, folderName, file.Name)
		if err != nil {
			return fmt.Errorf("failed to delete file %s from database: %w", file.Name, err)
		}
		fmt.Printf("File %s completely removed as it has no remaining categories\n", file.Name)
		return nil
	}
	fmt.Printf("File %s updated, folder %s removed from category %s\n", file.Name, folderName, categoryId)

	return nil
}
