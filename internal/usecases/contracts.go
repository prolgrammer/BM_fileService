package usecases

import "context"

type (
	Repository interface {
		Insert(ctx context.Context, data interface{}) error
		SelectById(ctx context.Context, id string) error
		UpdateById(ctx context.Context, id string, data interface{}) error
		DeleteById(ctx context.Context, id string) error

		SelectByCategoryId(ctx context.Context, categoryId string) error
	}
)
