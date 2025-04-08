package repositories

import (
	"app/internal/entities"
	"context"
)

type (
	AccountRepository interface {
		CreateAccount(ctx context.Context, userId string) error
		SelectAccount(ctx context.Context, userId string) (entities.Account, error)
		CheckAccountExists(ctx context.Context, userId string) (bool, error)
	}

	CategoryRepository interface {
		CreateCategory(ctx context.Context, userId, category string) error
		SelectCategory(ctx context.Context, userId, category string) (entities.Category, error)
		SelectAllCategories(ctx context.Context, userId string) ([]entities.Category, error)
		DeleteCategory(ctx context.Context, userId, category string) error
		CheckCategoryExists(ctx context.Context, userId, category string) (bool, error)
	}

	FolderRepository interface {
		CreateFolder(ctx context.Context, userId, category, folderName string) error
		SelectFolder(ctx context.Context, userId, category, folder string) (entities.Folder, error)
		SelectFolders(ctx context.Context, userId, category string) ([]entities.Folder, error)
		DeleteFolder(ctx context.Context, userId, category, folder string) error
		CheckFolderExists(ctx context.Context, userId, category, folder string) (bool, error)
	}

	FileRepository interface {
		CreateFile(ctx context.Context, data entities.File) error
		SelectFile(ctx context.Context, categoryId, folderName, fileName string) (entities.File, error)
		SelectFiles(ctx context.Context, categoryId, folderName string) ([]entities.File, error)
		UpdateFile(ctx context.Context, file entities.File) error
		DeleteFile(ctx context.Context, categoryId, folderName, fileName string) error
		CheckFileExists(ctx context.Context, categoryId, folderName, fileName string) (bool, error)
	}
)
