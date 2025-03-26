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
	minio          *m.Client
	fileRepository repositories.FileRepository
}

type GetAllFilesUseCase interface {
	GetAllFiles(ctx context.Context, accountId string, req requests.File) ([]responses.File, error)
}

func NewGetAllFilesUseCase(
	minio *m.Client,
	fileRepository repositories.FileRepository,
) GetAllFilesUseCase {
	return &getAllFiles{
		minio:          minio,
		fileRepository: fileRepository}
}

func (g *getAllFiles) GetAllFiles(ctx context.Context, accountId string, req requests.File) ([]responses.File, error) {
	files, err := g.fileRepository.SelectFiles(ctx, accountId, req.Category.Name, req.Folder.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get files list: %w", err)
	}

	result := make([]responses.File, len(files))
	for _, file := range files {
		fileUrl, err := g.minio.MinioClient.PresignedGetObject(
			ctx,
			g.minio.BucketName,
			file.Path,
			24*time.Hour,
			nil,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to get file url: %w", err)
		}

		result = append(result, responses.File{
			Name:      file.Name,
			Path:      fileUrl.String(),
			Size:      file.Size,
			Type:      file.Type,
			CreatedAt: file.CreatedAt,
		})
	}

	return result, nil
}
