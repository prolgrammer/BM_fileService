package usecases

import (
	"app/controllers/requests"
	m "app/infrastructure/minio"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
)

type LoadFileUseCases interface {
	LoadFile(ctx context.Context, accountId string, req requests.LoadFile) error
}

type loadFileUseCase struct {
	minio *m.Client
}

func NewLoadFileUseCase(minio *m.Client) LoadFileUseCases {
	return &loadFileUseCase{
		minio: minio,
	}
}

func (l loadFileUseCase) LoadFile(ctx context.Context, accountId string, req requests.LoadFile) error {
	if len(req.Files) == 0 {
		return fmt.Errorf("no files provided")
	}

	for _, fileHeader := range req.Files {
		file, err := fileHeader.Open()
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", fileHeader.Filename, err)
		}
		defer file.Close()

		objectName := fmt.Sprintf("%s/%s/%s", accountId, req.Folder, fileHeader.Filename)

		_, err = l.minio.MinioClient.PutObject(ctx, l.minio.BucketName, objectName, file, fileHeader.Size, minio.PutObjectOptions{})
		if err != nil {
			return fmt.Errorf("failed to upload file %s to minio: %w", fileHeader.Filename, err)
		}

		fmt.Printf("uploaded file %s to minio\n", fileHeader.Filename)
	}

	return nil
}
