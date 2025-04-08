package usecases

import "errors"

var (
	ErrAccountAlreadyExists = errors.New("account already exists")

	ErrCategoryAlreadyExists = errors.New("category already exists")
	ErrCategoryNotFound      = errors.New("category not found")

	ErrFolderAlreadyExists = errors.New("folder already exists")
)
