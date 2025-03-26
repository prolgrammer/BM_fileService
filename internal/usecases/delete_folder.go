package usecases

import (
	"app/controllers/requests"
	m "app/infrastructure/minio"
	"app/internal/repositories"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
)

type deleteFolderUseCase struct {
	minio            *m.Client
	folderRepository repositories.FolderRepository
	fileRepository   repositories.FileRepository
}

type DeleteFolderUseCase interface {
	DeleteFolder(ctx context.Context, accountId string, request requests.Folder) error
}

func NewDeleteFolderUseCase(
	minio *m.Client,
	folderRepository repositories.FolderRepository,
	fileRepository repositories.FileRepository) DeleteFolderUseCase {
	return &deleteFolderUseCase{
		minio:            minio,
		folderRepository: folderRepository,
		fileRepository:   fileRepository,
	}
}

func (uc *deleteFolderUseCase) DeleteFolder(ctx context.Context, accountId string, req requests.Folder) error {
	files, err := uc.fileRepository.SelectFiles(ctx, accountId, req.Category.Name, req.Name)
	if err != nil {
		return fmt.Errorf("failed to get folder files: %w", err)
	}

	for _, file := range files {
		err = uc.minio.MinioClient.RemoveObject(ctx, uc.minio.BucketName, file.Path, minio.RemoveObjectOptions{})
		if err != nil {
			return fmt.Errorf("failed to delete file %s from minio: %w", file.Name, err)
		}
	}

	err = uc.folderRepository.DeleteFolder(ctx, accountId, req.Category.Name, req.Name)
	if err != nil {
		return fmt.Errorf("failed to delete folder from db: %w", err)
	}

	return nil
}
