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
	minio          *m.Client
	fileRepository repositories.FileRepository
}

func NewGetFileUseCase(
	minio *m.Client,
	fileRepository repositories.FileRepository,
) GetFileUseCase {
	return &getFile{
		minio:          minio,
		fileRepository: fileRepository,
	}
}

func (g *getFile) GetFile(ctx context.Context, accountId string, req requests.File) (responses.File, error) {
	file, err := g.fileRepository.SelectFile(ctx, accountId, req.Category.Name, req.Folder.Name, req.Name)
	if err != nil {
		return responses.File{}, fmt.Errorf("failed to get file from database: %w", err)
	}

	fileUrl, err := g.minio.MinioClient.PresignedGetObject(
		ctx,
		g.minio.BucketName,
		file.Path,
		24*time.Hour,
		nil,
	)

	if err != nil {
		return responses.File{}, fmt.Errorf("failed to get file url: %w", err)
	}

	return responses.File{
		Name:      file.Name,
		Path:      fileUrl.String(),
		Size:      file.Size,
		Type:      file.Type,
		CreatedAt: file.CreatedAt,
	}, nil
}
