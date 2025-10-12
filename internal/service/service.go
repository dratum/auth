package service

import (
	"context"

	"github.com/dratum/auth/internal/model"
)

type AuthService interface {
	Create(ctx context.Context, fields *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, fields *model.UserUpdate) error
	Delete(ctx context.Context, id int64) error
}
