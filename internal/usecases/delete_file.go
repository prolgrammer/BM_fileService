package usecases

import (
	"app/controllers/requests"
	m "app/infrastructure/minio"
	"app/internal/repositories"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
)

type DeleteFileUseCase interface {
	DeleteFile(ctx context.Context, accountId string, req requests.File) error
}

type deleteFileUseCase struct {
	minio          *m.Client
	fileRepository repositories.FileRepository
}

func NewDeleteFileUseCase(
	minio *m.Client,
	fileRepository repositories.FileRepository,
) DeleteFileUseCase {
	return &deleteFileUseCase{
		minio:          minio,
		fileRepository: fileRepository,
	}
}

func (d *deleteFileUseCase) DeleteFile(ctx context.Context, accountId string, req requests.File) error {
	file, err := d.fileRepository.SelectFile(ctx, accountId, req.Category.Name, req.Folder.Name, req.Name)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	err = d.minio.MinioClient.RemoveObject(ctx, d.minio.BucketName, file.Path, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from minio: %w", err)
	}

	err = d.fileRepository.DeleteFile(ctx, accountId, req.Category.Name, req.Folder.Name, req.Name)
	if err != nil {
		return fmt.Errorf("failed to delete file from db: %w", err)
	}

	return nil
}
