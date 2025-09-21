package repository

import (
	"context"

	"github.com/dratum/auth/pkg/auth_v1"
)

type UserRepository interface {
	Get(ctx context.Context, id int64) (*auth_v1.GetResponse, error)
	Create(ctx context.Context, fields *auth_v1.CreateRequest) (int64, error)
	// Update(ctx context.Context, id int64) error
	// Delete(ctx context.Context, id int64) error
}
