package usecases

import "errors"

var (
	ErrCategoryAlreadyExists = errors.New("category already exists")
	ErrCategoryNotFound      = errors.New("category not found")

	ErrFolderAlreadyExists = errors.New("folder already exists")
)
