package auth

import (
	"context"

	"github.com/dratum/auth/internal/model"
)

func (s *serv) Create(ctx context.Context, fields *model.User) (int64, error) {
	id, err := s.userRepository.Create(ctx, fields)
	if err != nil {
		return 0, err
	}

	return id, nil
}
