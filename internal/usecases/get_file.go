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

type GetFileUseCase interface {
	GetFile(ctx context.Context, accountId string, req requests.File) (responses.File, error)
}

type getFile struct {
	minio              *m.Client
	categoryRepository repositories.CategoryRepository
	folderRepository   repositories.FolderRepository
	fileRepository     repositories.FileRepository
}

func NewGetFileUseCase(
	minio *m.Client,
	categoryRepository repositories.CategoryRepository,
	folderRepository repositories.FolderRepository,
	fileRepository repositories.FileRepository,
) GetFileUseCase {
	return &getFile{
		minio:              minio,
		categoryRepository: categoryRepository,
		folderRepository:   folderRepository,
		fileRepository:     fileRepository,
	}
}

func (g *getFile) GetFile(ctx context.Context, accountId string, req requests.File) (responses.File, error) {
	category, err := g.categoryRepository.SelectCategory(ctx, accountId, req.Category.Name)
	if err != nil {
		return responses.File{}, fmt.Errorf("category not found in database: %w", err)
	}

	exists, err := g.folderRepository.CheckFolderExists(ctx, accountId, req.Category.Name, req.Folder.Name)
	if !exists {
		return responses.File{}, repositories.ErrFolderNotFound
	}
	if err != nil {
		return responses.File{}, err
	}

	file, err := g.fileRepository.SelectFile(ctx, category.Id, req.Folder.Name, req.Name)
	if err != nil {
		return responses.File{}, fmt.Errorf("failed to get file from database: %w", err)
	}

	objectName := fmt.Sprintf("%s/%s/%s", accountId, file.Version, file.Name)

	fileUrl, err := g.minio.MinioClient.PresignedGetObject(
		ctx,
		g.minio.BucketName,
		objectName,
		24*time.Hour,
		nil,
	)

	if err != nil {
		return responses.File{}, fmt.Errorf("failed to get file url: %w", err)
	}

	return responses.File{
		Name:        file.Name,
		Description: file.Description,
		Size:        file.Size,
		Type:        file.Type,
		Version:     file.Version,
		CreatedAt:   file.CreatedAt,
		URL:         fileUrl.String(),
	}, nil
}
