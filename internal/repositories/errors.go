package repositories

import "errors"

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrFolderNotFound   = errors.New("folder not found")
	ErrFileNotFound     = errors.New("file not found")
)
