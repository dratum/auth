package repository

import (
	"context"

	"github.com/dratum/auth/internal/model"
)

type UserRepository interface {
	Get(ctx context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, fields *model.User) (int64, error)
	Update(ctx context.Context, fields *model.UserUpdate) error
	Delete(ctx context.Context, id int64) error
}
