package usecases

import (
	"app/controllers/requests"
	"app/controllers/responses"
	m "app/infrastructure/minio"
	"app/internal/repositories"
	"context"
	"fmt"
	"time"
)

type getAllFiles struct {
	minio              *m.Client
	categoryRepository repositories.CategoryRepository
	folderRepository   repositories.FolderRepository
	fileRepository     repositories.FileRepository
}

type GetAllFilesUseCase interface {
	GetAllFiles(ctx context.Context, accountId string, req requests.Folder) ([]responses.File, error)
}

func NewGetAllFilesUseCase(
	minio *m.Client,
	categoryRepository repositories.CategoryRepository,
	folderRepository repositories.FolderRepository,
	fileRepository repositories.FileRepository,
) GetAllFilesUseCase {
	return &getAllFiles{
		minio:              minio,
		categoryRepository: categoryRepository,
		folderRepository:   folderRepository,
		fileRepository:     fileRepository}
}

func (g *getAllFiles) GetAllFiles(ctx context.Context, accountId string, req requests.Folder) ([]responses.File, error) {
	category, err := g.categoryRepository.SelectCategory(ctx, accountId, req.Category.Name)
	if err != nil {
		return nil, fmt.Errorf("category not found in database: %w", err)
	}

	exists, err := g.folderRepository.CheckFolderExists(ctx, accountId, req.Category.Name, req.Name)
	if !exists {
		return nil, repositories.ErrFolderNotFound
	}
	if err != nil {
		return nil, err
	}

	files, err := g.fileRepository.SelectFiles(ctx, category.Id, req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get files list: %w", err)
	}

	result := make([]responses.File, 0, len(files))
	for _, file := range files {
		objectName := fmt.Sprintf("%s/%s/%s", accountId, file.Version, file.Name)
		fileUrl, err := g.minio.MinioClient.PresignedGetObject(
			ctx,
			g.minio.BucketName,
			objectName,
			24*time.Hour,
			nil,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to get file url: %w", err)
		}

		result = append(result, responses.File{
			Name:        file.Name,
			Description: file.Description,
			Size:        file.Size,
			Type:        file.Type,
			Version:     file.Version,
			CreatedAt:   file.CreatedAt,
			URL:         fileUrl.String(),
		})
	}

	return result, nil
}
