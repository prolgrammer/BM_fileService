package usecases

import (
	"app/controllers/requests"
	m "app/infrastructure/minio"
	"app/internal/entities"
	"app/internal/repositories"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"time"
)

type CreateFileUseCases interface {
	CreateFile(ctx context.Context, accountId string, req requests.CreateFile) error
}

type createFileUseCase struct {
	minio          *m.Client
	fileRepository repositories.FileRepository
}

func NewCreateFileUseCase(minio *m.Client, fileRepository repositories.FileRepository) CreateFileUseCases {
	return &createFileUseCase{
		minio,
		fileRepository,
	}
}

func (uc *createFileUseCase) CreateFile(ctx context.Context, accountId string, req requests.CreateFile) error {
	if len(req.Files) == 0 {
		return fmt.Errorf("no files prodied")
	}

	for _, fileHeader := range req.Files {
		file, err := fileHeader.Open()
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", fileHeader.Filename, err)
		}
		defer file.Close()

		objectName := fmt.Sprintf("%s/%s/%s/%s", accountId, req.Category.Name, req.Folder.Name, fileHeader.Filename)

		_, err = uc.minio.MinioClient.PutObject(
			ctx,
			uc.minio.BucketName,
			objectName,
			file,
			fileHeader.Size,
			minio.PutObjectOptions{})
		if err != nil {
			return fmt.Errorf("failed to upload file %s: %w", fileHeader.Filename, err)
		}

		fmt.Printf("file uploaded successfully to minio: %s\n", fileHeader.Filename)

		fileEntity := entities.File{
			Name:      fileHeader.Filename,
			Path:      objectName,
			Size:      int(fileHeader.Size),
			Type:      fileHeader.Header.Get("Content-Type"),
			CreatedAt: time.Now(),
		}

		err = uc.fileRepository.CreateFile(
			ctx,
			accountId,
			req.Category.Name,
			req.Folder.Name,
			fileEntity,
		)
		if err != nil {
			_ = uc.minio.MinioClient.RemoveObject(
				ctx,
				uc.minio.BucketName,
				objectName,
				minio.RemoveObjectOptions{})
			return fmt.Errorf("failed to upload file %s: %w into database", fileHeader.Filename, err)
		}
	}

	return nil
}
