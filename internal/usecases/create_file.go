package usecases

import (
	"app/controllers/requests"
	m "app/infrastructure/minio"
	"app/internal/entities"
	"app/internal/repositories"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"time"
)

type CreateFileUseCases interface {
	CreateFile(ctx context.Context, accountId string, req requests.CreateFile) error
}

type createFileUseCase struct {
	minio              *m.Client
	categoryRepository repositories.CategoryRepository
	folderRepository   repositories.FolderRepository
	fileRepository     repositories.FileRepository
}

func NewCreateFileUseCase(minio *m.Client, categoryRepository repositories.CategoryRepository, folderRepository repositories.FolderRepository, fileRepository repositories.FileRepository) CreateFileUseCases {
	return &createFileUseCase{
		minio:              minio,
		categoryRepository: categoryRepository,
		folderRepository:   folderRepository,
		fileRepository:     fileRepository,
	}
}

func (uc *createFileUseCase) CreateFile(ctx context.Context, accountId string, req requests.CreateFile) error {
	if len(req.Files) == 0 {
		return fmt.Errorf("no files prodied")
	}

	category, err := uc.categoryRepository.SelectCategory(ctx, accountId, req.Category.Name)
	if err != nil {
		return err
	}

	exists, err := uc.folderRepository.CheckFolderExists(ctx, accountId, req.Category.Name, req.Folder.Name)
	if !exists {
		return repositories.ErrFolderNotFound
	}
	if err != nil {
		return err
	}

	for _, fileHeader := range req.Files {
		file, err := fileHeader.Open()
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", fileHeader.Filename, err)
		}
		defer file.Close()

		objectName := fmt.Sprintf("%s/%s/%s", accountId, req.Version, fileHeader.Filename)

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

		exists, err := uc.fileRepository.CheckFileExistsByNameAndVersion(ctx, fileHeader.Filename, req.Version)
		if err != nil {
			_ = uc.minio.MinioClient.RemoveObject(ctx, uc.minio.BucketName, objectName, minio.RemoveObjectOptions{})
			return fmt.Errorf("failed to check file existence for %s: %w", fileHeader.Filename, err)
		}

		if exists {
			fileEntity, err := uc.fileRepository.SelectFileByNameAndVersion(ctx, category.Id, fileHeader.Filename, req.Version)
			if err != nil {
				_ = uc.minio.MinioClient.RemoveObject(ctx, uc.minio.BucketName, objectName, minio.RemoveObjectOptions{})
				return fmt.Errorf("failed to select file %s: %w", fileHeader.Filename, err)
			}

			categoryExists := false
			for i, cat := range fileEntity.Categories {
				if cat.CategoryId == category.Id {
					folderExists := false
					for _, folder := range cat.Folders {
						if folder.Name == req.Folder.Name {
							folderExists = true
							break
						}
					}
					if !folderExists {
						fileEntity.Categories[i].Folders = append(fileEntity.Categories[i].Folders, entities.CreateFolder(req.Folder.Name))
					}
					categoryExists = true
					break
				}
			}

			if !categoryExists {
				fileEntity.Categories = append(fileEntity.Categories, entities.FileCategory{
					CategoryId: category.Id,
					Folders:    []entities.Folder{{Name: req.Folder.Name}},
				})
			}

			err = uc.fileRepository.UpdateFile(ctx, fileEntity)
			if err != nil {
				_ = uc.minio.MinioClient.RemoveObject(ctx, uc.minio.BucketName, objectName, minio.RemoveObjectOptions{})
				return fmt.Errorf("failed to update file %s in database: %w", fileHeader.Filename, err)
			}
		} else {
			fileEntity := entities.File{
				Id:      uuid.New().String(),
				Name:    fileHeader.Filename,
				Size:    int(fileHeader.Size),
				Type:    fileHeader.Header.Get("Content-Type"),
				Version: req.Version,
				Categories: []entities.FileCategory{
					{
						CategoryId: category.Id,
						Folders: []entities.Folder{
							{
								req.Folder.Name,
							},
						},
					},
				},
				CreatedAt: time.Now(),
			}

			err = uc.fileRepository.CreateFile(
				ctx,
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
	}

	return nil
}
